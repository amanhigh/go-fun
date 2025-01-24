package controllers

import (
	"context"

	cachev1beta1 "github.com/amanhigh/go-fun/components/operator/api/v1beta1"
	ctrl "sigs.k8s.io/controller-runtime"
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
