package controllers

import (
    "context"
    cachev1beta1 "github.com/amanhigh/go-fun/components/operator/api/v1beta1"
    appsv1 "k8s.io/api/apps/v1"
    corev1 "k8s.io/api/core/v1"
    ctrl "sigs.k8s.io/controller-runtime"
)

/* DeploymentHelper handles all deployment related operations
for the Memcached controller including creation, updates and
security configuration. */
type DeploymentHelper interface {
    CreateNewDeployment(ctx context.Context, memcached *cachev1beta1.Memcached) (ctrl.Result, error)
    GenerateDeploymentSpec(memcached *cachev1beta1.Memcached) (*appsv1.DeploymentSpec, error)
    GeneratePodSpec(memcached *cachev1beta1.Memcached, image string) (*corev1.PodSpec, error)
    GenerateSecurityContext() *corev1.SecurityContext
}

type deploymentHelperImpl struct {
    controller *MemcachedReconciler
}

func NewDeploymentHelper(controller *MemcachedReconciler) DeploymentHelper {
    return &deploymentHelperImpl{
        controller: controller,
    }
}