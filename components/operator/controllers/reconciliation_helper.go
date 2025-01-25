package controllers

import (
	"context"
	"fmt"
	"time"

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
	ReconcileDeployment(ctx context.Context, memcached *cachev1beta1.Memcached) (ctrl.Result, error)
	ExecuteFinalizer(ctx context.Context, memcached *cachev1beta1.Memcached) (ctrl.Result, error)
	AddFinalizer(ctx context.Context, memcached *cachev1beta1.Memcached) (ctrl.Result, error)
}

type reconciliationHelperImpl struct {
	controller   *MemcachedReconciler
	statusHelper StatusHelper
	deployHelper DeploymentHelper // Add new field
}

func NewReconciliationHelper(statusHelper StatusHelper, deployHelper DeploymentHelper, controller *MemcachedReconciler) ReconciliationHelper {
	return &reconciliationHelperImpl{
		controller:   controller,
		statusHelper: statusHelper,
		deployHelper: deployHelper,
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

// ExecuteFinalizer performs finalizer operations when resource is being deleted
// It updates status, records events, and removes the finalizer
func (r *reconciliationHelperImpl) ExecuteFinalizer(
	ctx context.Context,
	memcached *cachev1beta1.Memcached,
) (ctrl.Result, error) {
	log := log.FromContext(ctx)
	if !controllerutil.ContainsFinalizer(memcached, memcachedFinalizer) {
		return ctrl.Result{}, nil
	}

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

	return ctrl.Result{}, nil
}

// AddFinalizer adds finalizer if it doesn't exist
// This ensures cleanup operations are performed when the resource is deleted
func (r *reconciliationHelperImpl) AddFinalizer(
	ctx context.Context,
	memcached *cachev1beta1.Memcached,
) (ctrl.Result, error) {
	log := log.FromContext(ctx)
	if controllerutil.ContainsFinalizer(memcached, memcachedFinalizer) {
		return ctrl.Result{}, nil
	}

	log.Info("Adding Finalizer for Memcached")
	if ok := controllerutil.AddFinalizer(memcached, memcachedFinalizer); !ok {
		log.Error(nil, "Failed to add finalizer")
		return ctrl.Result{Requeue: true}, nil
	}

	if err := r.controller.Update(ctx, memcached); err != nil {
		log.Error(err, "Failed to update CR to add finalizer")
		return ctrl.Result{}, err
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
	dep := &appsv1.Deployment{}
	// Handle deployment creation and capture result
	createResult, err := r.handleDeploymentCreation(ctx, memcached, dep)
	if err != nil {
		return createResult, err
	}

	// If creation requires requeue, return immediately
	if createResult.Requeue || createResult.RequeueAfter > 0 {
		return createResult, nil
	}

	return r.handleDeploymentUpdate(ctx, memcached, dep)
}

// handleDeploymentCreation handles deployment existence check and creation
// Returns ctrl.Result and error if deployment needs to be created or if there's an error
func (r *reconciliationHelperImpl) handleDeploymentCreation(
	ctx context.Context,
	memcached *cachev1beta1.Memcached,
	dep *appsv1.Deployment,
) (ctrl.Result, error) {
	log := log.FromContext(ctx)

	// Check if deployment exists
	err := r.controller.Get(ctx, types.NamespacedName{
		Name:      memcached.Name,
		Namespace: memcached.Namespace,
	}, dep)

	if err != nil && apierrors.IsNotFound(err) {
		result, err := r.deployHelper.ValidateAndCreateDeployment(ctx, memcached)
		if err != nil {
			return result, err
		}
		// Requeue for deployment creation
		return ctrl.Result{RequeueAfter: time.Minute}, nil
	} else if err != nil {
		log.Error(err, "Failed to get Deployment")
		return ctrl.Result{}, err
	}

	return ctrl.Result{}, nil
}

// handleDeploymentUpdate handles deployment update and status management
func (r *reconciliationHelperImpl) handleDeploymentUpdate(
	ctx context.Context,
	memcached *cachev1beta1.Memcached,
	dep *appsv1.Deployment,
) (ctrl.Result, error) {
	log := log.FromContext(ctx)
	size := memcached.Spec.Size

	// Guard against nil deployment or replicas
	if dep == nil || dep.Spec.Replicas == nil {
		log.Info("Deployment not yet initialized, requeueing")
		return ctrl.Result{Requeue: true}, nil
	}

	// Update if size doesn't match
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
