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
	ReconciliationHelper provides core reconciliation operations

for the Memcached controller. It handles fetching instances,
status management and finalizer operations.
*/
type ReconciliationHelper interface {
	FetchMemcachedInstance(ctx context.Context, req ctrl.Request) (*cachev1beta1.Memcached, error)
	InitializeStatus(ctx context.Context, memcached *cachev1beta1.Memcached) error
	HandleFinalizers(ctx context.Context, memcached *cachev1beta1.Memcached) (ctrl.Result, error)
	ReconcileDeployment(ctx context.Context, memcached *cachev1beta1.Memcached) (ctrl.Result, error)
}

type reconciliationHelperImpl struct {
	controller *MemcachedReconciler
}

func NewReconciliationHelper(controller *MemcachedReconciler) ReconciliationHelper {
	return &reconciliationHelperImpl{
		controller: controller,
	}
}

func (r *reconciliationHelperImpl) FetchMemcachedInstance(ctx context.Context, req ctrl.Request) (*cachev1beta1.Memcached, error) {
	log := log.FromContext(ctx)
	memcached := &cachev1beta1.Memcached{}

	if err := r.controller.Get(ctx, req.NamespacedName, memcached); err != nil {
		if apierrors.IsNotFound(err) {
			log.Info("memcached resource not found. Ignoring since object must be deleted")
			return nil, nil
		}
		log.Error(err, "Failed to get memcached")
		return nil, err
	}
	return memcached, nil
}

func (r *reconciliationHelperImpl) HandleFinalizers(ctx context.Context, memcached *cachev1beta1.Memcached) (ctrl.Result, error) {
	log := log.FromContext(ctx)

	// Check if object is being deleted
	if memcached.GetDeletionTimestamp() != nil {
		if controllerutil.ContainsFinalizer(memcached, memcachedFinalizer) {
			// Begin finalizer operations
			log.Info("Performing Finalizer Operations for Memcached before delete CR")

			// Update status to indicate deletion
			if err := r.controller.statusHelper.UpdateDegradedStatus(ctx, memcached,
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

func (r *reconciliationHelperImpl) ReconcileDeployment(ctx context.Context, memcached *cachev1beta1.Memcached) (ctrl.Result, error) {
	log := log.FromContext(ctx)
	dep := &appsv1.Deployment{}

	// Check if deployment exists
	err := r.controller.Get(ctx, types.NamespacedName{
		Name:      memcached.Name,
		Namespace: memcached.Namespace,
	}, dep)

	if err != nil && apierrors.IsNotFound(err) {
		return r.controller.deployHelper.CreateNewDeployment(ctx, memcached)
	} else if err != nil {
		log.Error(err, "Failed to get Deployment")
		return ctrl.Result{}, err
	}

	// Update if size doesn't match
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
	if err := r.controller.statusHelper.UpdateSuccessStatus(ctx, memcached, size); err != nil {
		return ctrl.Result{}, err
	}

	return ctrl.Result{}, nil
}
