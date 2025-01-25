package controllers

import (
	"context"
	"fmt"

	cachev1beta1 "github.com/amanhigh/go-fun/components/operator/api/v1beta1"
	"k8s.io/apimachinery/pkg/api/meta"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

/*
StatusHelper manages all status-related operations for the Memcached controller.

It handles setting and updating various status conditions that reflect the state
of the Memcached custom resource and its associated components. Status conditions
are crucial for providing visibility into the operational state of the resource.

Key responsibilities:
- Initializing status conditions
- Updating success/failure states
- Managing degraded states during deletion
*/
type StatusHelper interface {
	InitializeStatus(ctx context.Context, memcached *cachev1beta1.Memcached) error
	UpdateSuccessStatus(ctx context.Context, memcached *cachev1beta1.Memcached, size int32) error
	UpdateStatusWithError(ctx context.Context, memcached *cachev1beta1.Memcached, message string, err error) error
	UpdateDegradedStatus(ctx context.Context, memcached *cachev1beta1.Memcached, message string) error
}

type statusHelperImpl struct {
	controller *MemcachedReconciler
}

func NewStatusHelper(controller *MemcachedReconciler) StatusHelper {
	return &statusHelperImpl{
		controller: controller,
	}
}

// InitializeStatus sets up the initial status conditions for a newly created Memcached resource
// If no status conditions exist, it creates an initial "Unknown" state
// This helps track the resource's state from the beginning of its lifecycle
func (s *statusHelperImpl) InitializeStatus(ctx context.Context, memcached *cachev1beta1.Memcached) error {
	// Let's just set the status as Unknown when no status are available
	if len(memcached.Status.Conditions) == 0 {
		meta.SetStatusCondition(&memcached.Status.Conditions, metav1.Condition{
			Type:    typeAvailableMemcached,
			Status:  metav1.ConditionUnknown,
			Reason:  "Reconciling",
			Message: "Starting reconciliation",
		})

		return s.controller.Status().Update(ctx, memcached)
	}
	return nil
}

// UpdateSuccessStatus updates the status to reflect successful deployment
// This is called when the deployment has been created/updated successfully
// and the desired number of replicas are running
func (s *statusHelperImpl) UpdateSuccessStatus(ctx context.Context, memcached *cachev1beta1.Memcached, size int32) error {
	// The following implementation will update the status
	meta.SetStatusCondition(&memcached.Status.Conditions, metav1.Condition{
		Type:    typeAvailableMemcached,
		Status:  metav1.ConditionTrue,
		Reason:  "Reconciling",
		Message: fmt.Sprintf("Deployment for custom resource (%s) with %d replicas created successfully", memcached.Name, size),
	})

	return s.controller.Status().Update(ctx, memcached)
}

// UpdateStatusWithError updates status when an error occurs during reconciliation
// This provides visibility into failures and helps diagnose issues
// The message and error are combined to give detailed information about the failure
func (s *statusHelperImpl) UpdateStatusWithError(ctx context.Context, memcached *cachev1beta1.Memcached, message string, err error) error {
	// Update status with error information
	meta.SetStatusCondition(&memcached.Status.Conditions, metav1.Condition{
		Type:    typeAvailableMemcached,
		Status:  metav1.ConditionFalse,
		Reason:  "Reconciling",
		Message: fmt.Sprintf("%s: %s", message, err),
	})

	return s.controller.Status().Update(ctx, memcached)
}

// UpdateDegradedStatus sets the status to a degraded state, typically used during deletion
// This is important for tracking the resource's state during cleanup operations
// More info: https://kubernetes.io/docs/concepts/workloads/controllers/garbage-collection/
func (s *statusHelperImpl) UpdateDegradedStatus(ctx context.Context, memcached *cachev1beta1.Memcached, message string) error {
	// Let's add here a status "Degraded" to define that this resource begin its process to be terminated
	meta.SetStatusCondition(&memcached.Status.Conditions, metav1.Condition{
		Type:    typeDegradedMemcached,
		Status:  metav1.ConditionUnknown,
		Reason:  "Finalizing",
		Message: message,
	})

	return s.controller.Status().Update(ctx, memcached)
}
