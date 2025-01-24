package controllers

import (
	"context"
	"fmt"

	cachev1beta1 "github.com/amanhigh/go-fun/components/operator/api/v1beta1"
	appsv1 "k8s.io/api/apps/v1"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/types"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
	"sigs.k8s.io/controller-runtime/pkg/log"
)

/*
ReconciliationHelper provides core reconciliation operations for the Memcached controller.

It implements the core reconciliation loop which aims to move the current state of
the cluster closer to the desired state. The reconciliation must be idempotent.
*/
type ReconciliationHelper interface {
	FetchMemcachedInstance(ctx context.Context, req ctrl.Request) (*cachev1beta1.Memcached, error)
	HandleFinalizers(ctx context.Context, memcached *cachev1beta1.Memcached) (ctrl.Result, error)
	ReconcileDeployment(ctx context.Context, memcached *cachev1beta1.Memcached) (ctrl.Result, error)
}

type reconciliationHelperImpl struct {
	controller   *MemcachedReconciler
	statusHelper StatusHelper
}

func NewReconciliationHelper(statusHelper StatusHelper, controller *MemcachedReconciler) ReconciliationHelper {
	return &reconciliationHelperImpl{
		controller:   controller,
		statusHelper: statusHelper,
	}
}

// FetchMemcachedInstance retrieves the Memcached instance from the cluster
// If the resource is not found, it means it was deleted or not created
// in which case we stop the reconciliation
func (r *reconciliationHelperImpl) FetchMemcachedInstance(ctx context.Context, req ctrl.Request) (*cachev1beta1.Memcached, error) {
	log := log.FromContext(ctx)
	memcached := &cachev1beta1.Memcached{}

	if err := r.controller.Get(ctx, req.NamespacedName, memcached); err != nil {
		if apierrors.IsNotFound(err) {
			// If the custom resource is not found then, it usually means that it was deleted or not created
			// In this way, we will stop the reconciliation
			log.Info("memcached resource not found. Ignoring since object must be deleted")
			return nil, nil
		}
		log.Error(err, "Failed to get memcached")
		return nil, err
	}
	return memcached, nil
}

// HandleFinalizers manages finalizer operations for the Memcached resource
// Finalizers allow controllers to implement cleanup tasks before an object is deleted
// More info: https://kubernetes.io/docs/concepts/overview/working-with-objects/finalizers/
func (r *reconciliationHelperImpl) HandleFinalizers(ctx context.Context, memcached *cachev1beta1.Memcached) (ctrl.Result, error) {
	log := log.FromContext(ctx)

	// Check if the Memcached instance is marked to be deleted, which is
	// indicated by the deletion timestamp being set.
	if memcached.GetDeletionTimestamp() != nil {
		if controllerutil.ContainsFinalizer(memcached, memcachedFinalizer) {
			log.Info("Performing Finalizer Operations for Memcached before delete CR")

			// Note: It is not recommended to use finalizers with the purpose of delete resources which are
			// created and managed in the reconciliation. These ones, such as the Deployment created on this reconcile,
			// are defined as dependent of the custom resource. See that we use the method ctrl.SetControllerReference.
			// to set the ownerRef which means that the Deployment will be deleted by the Kubernetes API.
			// More info: https://kubernetes.io/docs/tasks/administer-cluster/use-cascading-deletion/

			// Update status to indicate deletion
			if err := r.statusHelper.UpdateDegradedStatus(ctx, memcached,
				fmt.Sprintf("Performing finalizer operations for the custom resource: %s", memcached.Name)); err != nil {
				return ctrl.Result{}, err
			}

			// Record event for deletion
			r.controller.Recorder.Event(memcached, "Warning", "Deleting",
				fmt.Sprintf("Custom Resource %s is being deleted from the namespace %s",
					memcached.Name, memcached.Namespace))

			// Remove finalizer
			if ok := controllerutil.RemoveFinalizer(memcached, memcachedFinalizer); !ok {
				log.Error(nil, "Failed to remove finalizer")
				return ctrl.Result{Requeue: true}, nil
			}

			if err := r.controller.Update(ctx, memcached); err != nil {
				log.Error(err, "Failed to remove finalizer")
				return ctrl.Result{}, err
			}
		}
		return ctrl.Result{}, nil
	}

	// Add finalizer if it doesn't exist
	// This will ensure our cleanup operations are performed when the resource is deleted
	if !controllerutil.ContainsFinalizer(memcached, memcachedFinalizer) {
		log.Info("Adding Finalizer for Memcached")
		if ok := controllerutil.AddFinalizer(memcached, memcachedFinalizer); !ok {
			log.Error(nil, "Failed to add finalizer")
			return ctrl.Result{Requeue: true}, nil
		}

		if err := r.controller.Update(ctx, memcached); err != nil {
			log.Error(err, "Failed to update CR to add finalizer")
			return ctrl.Result{}, err
		}
	}

	return ctrl.Result{}, nil
}

// ReconcileDeployment ensures the deployment exists and matches the desired state
// It handles creation and updates of the deployment based on the Memcached spec
// The reconciliation process implements the following cases:
// - Create Deployment if it doesn't exist
// - Update Deployment if replicas don't match
// - Update status based on reconciliation results
func (r *reconciliationHelperImpl) ReconcileDeployment(ctx context.Context, memcached *cachev1beta1.Memcached) (ctrl.Result, error) {
	log := log.FromContext(ctx)

	// Skip deployment logic if resource is being deleted
	if memcached.GetDeletionTimestamp() != nil {
		log.Info("Resource is being deleted, skipping deployment reconciliation")
		return ctrl.Result{}, nil
	}

	dep := &appsv1.Deployment{}

	// Check if deployment exists
	err := r.controller.Get(ctx, types.NamespacedName{
		Name:      memcached.Name,
		Namespace: memcached.Namespace,
	}, dep)

	if err != nil && apierrors.IsNotFound(err) {
		// Attempt creation with validation
		return r.controller.deployHelper.ValidateAndCreateDeployment(ctx, memcached)
	} else if err != nil {
		log.Error(err, "Failed to get Deployment")
		return ctrl.Result{}, err
	}

	// Update if size doesn't match
	// The CRD API is defining that the Memcached type, have a MemcachedSpec.Size field
	// to set the quantity of Deployment instances is the desired state on the cluster.
	size := memcached.Spec.Size
	if *dep.Spec.Replicas != size {
		dep.Spec.Replicas = &size
		if err := r.controller.Update(ctx, dep); err != nil {
			log.Error(err, "Failed to update Deployment",
				"Deployment.Namespace", dep.Namespace, "Deployment.Name", dep.Name)
			return ctrl.Result{}, err
		}
		return ctrl.Result{Requeue: true}, nil
	}

	// Update success status
	if err := r.statusHelper.UpdateSuccessStatus(ctx, memcached, size); err != nil {
		return ctrl.Result{}, err
	}

	return ctrl.Result{}, nil
}
