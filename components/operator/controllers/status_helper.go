package controllers

import (
	"context"
	"fmt"

	cachev1beta1 "github.com/amanhigh/go-fun/components/operator/api/v1beta1"
	"k8s.io/apimachinery/pkg/api/meta"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

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

func (s *statusHelperImpl) InitializeStatus(ctx context.Context, memcached *cachev1beta1.Memcached) error {
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

func (s *statusHelperImpl) UpdateSuccessStatus(ctx context.Context, memcached *cachev1beta1.Memcached, size int32) error {
	meta.SetStatusCondition(&memcached.Status.Conditions, metav1.Condition{
		Type:    typeAvailableMemcached,
		Status:  metav1.ConditionTrue,
		Reason:  "Reconciling",
		Message: fmt.Sprintf("Deployment for custom resource (%s) with %d replicas created successfully", memcached.Name, size),
	})

	return s.controller.Status().Update(ctx, memcached)
}

func (s *statusHelperImpl) UpdateStatusWithError(ctx context.Context, memcached *cachev1beta1.Memcached, message string, err error) error {
	meta.SetStatusCondition(&memcached.Status.Conditions, metav1.Condition{
		Type:    typeAvailableMemcached,
		Status:  metav1.ConditionFalse,
		Reason:  "Reconciling",
		Message: fmt.Sprintf("%s: %s", message, err),
	})

	return s.controller.Status().Update(ctx, memcached)
}

func (s *statusHelperImpl) UpdateDegradedStatus(ctx context.Context, memcached *cachev1beta1.Memcached, message string) error {
	meta.SetStatusCondition(&memcached.Status.Conditions, metav1.Condition{
		Type:    typeDegradedMemcached,
		Status:  metav1.ConditionUnknown,
		Reason:  "Finalizing",
		Message: message,
	})

	return s.controller.Status().Update(ctx, memcached)
}
