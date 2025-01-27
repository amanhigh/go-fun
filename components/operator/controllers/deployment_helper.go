package controllers

import (
	"context"
	"fmt"
	"os"
	"strings"
	"time"

	cachev1beta1 "github.com/amanhigh/go-fun/components/operator/api/v1beta1"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/log"
)

/*
	DeploymentHelper handles all deployment related operations

for the Memcached controller including creation, updates and
security configuration.
*/
type DeploymentHelper interface {
	ValidateAndCreateDeployment(ctx context.Context, memcached *cachev1beta1.Memcached) (ctrl.Result, error)
	CreateNewDeployment(ctx context.Context, memcached *cachev1beta1.Memcached) (ctrl.Result, error)
	GenerateDeploymentSpec(memcached *cachev1beta1.Memcached) (*appsv1.DeploymentSpec, error)
	GeneratePodSpec(memcached *cachev1beta1.Memcached, image string) (*corev1.PodSpec, error)
	GenerateSecurityContext() *corev1.SecurityContext
	GetLabels(name string, image string) map[string]string
}

type deploymentHelperImpl struct {
	controller *MemcachedReconciler
}

func NewDeploymentHelper(controller *MemcachedReconciler) DeploymentHelper {
	return &deploymentHelperImpl{
		controller: controller,
	}
}

func (d *deploymentHelperImpl) ValidateAndCreateDeployment(ctx context.Context, memcached *cachev1beta1.Memcached) (ctrl.Result, error) {
	log := log.FromContext(ctx)

	// Check image availability first - let error propagate up
	if _, err := d.getMemcachedImage(); err != nil {
		log.Error(err, "Failed to get Memcached image")
		return ctrl.Result{}, err
	}

	// Now proceed with regular deployment
	return d.CreateNewDeployment(ctx, memcached)
}

func (d *deploymentHelperImpl) CreateNewDeployment(ctx context.Context, memcached *cachev1beta1.Memcached) (ctrl.Result, error) {
	log := log.FromContext(ctx)

	// Generate deployment spec
	deploymentSpec, err := d.GenerateDeploymentSpec(memcached)
	if err != nil {
		log.Error(err, "Failed to generate deployment spec")
		return ctrl.Result{}, err
	}

	// Create deployment object
	dep := &appsv1.Deployment{
		ObjectMeta: metav1.ObjectMeta{
			Name:      memcached.Name,
			Namespace: memcached.Namespace,
		},
		Spec: *deploymentSpec,
	}

	// Set controller reference
	if err := ctrl.SetControllerReference(memcached, dep, d.controller.Scheme); err != nil {
		log.Error(err, "Failed to set controller reference")
		return ctrl.Result{}, err
	}

	// Create deployment
	log.Info("Creating a new Deployment",
		"Deployment.Namespace", dep.Namespace, "Deployment.Name", dep.Name)
	if err := d.controller.Create(ctx, dep); err != nil {
		log.Error(err, "Failed to create deployment")
		return ctrl.Result{}, err
	}

	return ctrl.Result{RequeueAfter: time.Minute}, nil
}

func (d *deploymentHelperImpl) GenerateDeploymentSpec(memcached *cachev1beta1.Memcached) (*appsv1.DeploymentSpec, error) {
	// Get image and labels
	image, err := d.getMemcachedImage()
	if err != nil {
		return nil, err
	}
	ls := d.GetLabels(memcached.Name, image)

	// Generate pod spec
	podSpec, err := d.GeneratePodSpec(memcached, image)
	if err != nil {
		return nil, err
	}

	replicas := memcached.Spec.Size
	return &appsv1.DeploymentSpec{
		Replicas: &replicas,
		Selector: &metav1.LabelSelector{
			MatchLabels: ls,
		},
		Template: corev1.PodTemplateSpec{
			ObjectMeta: metav1.ObjectMeta{
				Labels: ls,
			},
			Spec: *podSpec,
		},
	}, nil
}

func (d *deploymentHelperImpl) GeneratePodSpec(memcached *cachev1beta1.Memcached, image string) (*corev1.PodSpec, error) {
	return &corev1.PodSpec{
		SecurityContext: &corev1.PodSecurityContext{
			RunAsNonRoot: &[]bool{true}[0],
			SeccompProfile: &corev1.SeccompProfile{
				Type: corev1.SeccompProfileTypeRuntimeDefault,
			},
		},
		Containers: []corev1.Container{
			{
				Image:           image,
				Name:            "memcached",
				ImagePullPolicy: corev1.PullIfNotPresent,
				SecurityContext: d.GenerateSecurityContext(),
				Command:         []string{"memcached", "-m=64", "modern", "-v"},
			},
			{
				Image:           memcached.Spec.SidecarImage,
				Name:            "sidecar",
				ImagePullPolicy: corev1.PullIfNotPresent,
				Command:         []string{"sleep", "infinity"},
				SecurityContext: d.GenerateSecurityContext(),
			},
		},
	}, nil
}

func (d *deploymentHelperImpl) GenerateSecurityContext() *corev1.SecurityContext {
	return &corev1.SecurityContext{
		RunAsNonRoot:             &[]bool{true}[0],
		RunAsUser:                &[]int64{1001}[0],
		AllowPrivilegeEscalation: &[]bool{false}[0],
		Capabilities: &corev1.Capabilities{
			Drop: []corev1.Capability{"ALL"},
		},
	}
}

// Helper methods
func (d *deploymentHelperImpl) getMemcachedImage() (string, error) {
	var imageEnvVar = "MEMCACHED_IMAGE"
	image, found := os.LookupEnv(imageEnvVar)
	if !found {
		return "", fmt.Errorf("Unable to find %s environment variable with the image", imageEnvVar)
	}
	return image, nil
}

func (d *deploymentHelperImpl) GetLabels(name, image string) map[string]string {
	imageTag := "latest"
	if image != "" {
		parts := strings.Split(image, ":")
		if len(parts) > 1 {
			imageTag = parts[1]
		}
	}

	return map[string]string{
		"app.kubernetes.io/name":       "Memcached",
		"app.kubernetes.io/instance":   name,
		"app.kubernetes.io/version":    imageTag,
		"app.kubernetes.io/part-of":    "operator",
		"app.kubernetes.io/created-by": "controller-manager",
	}
}
