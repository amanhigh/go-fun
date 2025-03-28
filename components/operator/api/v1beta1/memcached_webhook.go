/*
Copyright 2023.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package v1beta1

import (
	"fmt"

	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	logf "sigs.k8s.io/controller-runtime/pkg/log"
	"sigs.k8s.io/controller-runtime/pkg/webhook"
)

const (
	minMemcachedPort = 8000
)

// log is for logging in this package.
var memcachedlog = logf.Log.WithName("memcached-resource")

func (r *Memcached) SetupWebhookWithManager(mgr ctrl.Manager) error {
	return ctrl.NewWebhookManagedBy(mgr).
		For(r).
		Complete()
}

// nolint:lll
// +kubebuilder:webhook:path=/mutate-cache-aman-com-v1beta1-memcached,mutating=true,failurePolicy=fail,sideEffects=None,groups=cache.aman.com,resources=memcacheds,verbs=create;update,versions=v1beta1,name=mmemcached.kb.io,admissionReviewVersions=v1

var _ webhook.Defaulter = &Memcached{}

// Default implements webhook.Defaulter so a webhook will be registered for the type
func (r *Memcached) Default() {
	memcachedlog.Info("Defaulting", "name", r.Name, "size", r.Spec.Size)
	if r.Spec.Size < 0 {
		r.Spec.Size = 1
		memcachedlog.Info("Detected Negative Size Defaulting to 1")
	}
}

// Change verbs to "verbs=create;update;delete" if you want to enable deletion validation.
// nolint:lll
// +kubebuilder:webhook:path=/validate-cache-aman-com-v1beta1-memcached,mutating=false,failurePolicy=fail,sideEffects=None,groups=cache.aman.com,resources=memcacheds,verbs=create;update;delete,versions=v1beta1,name=vmemcached.kb.io,admissionReviewVersions=v1

var _ webhook.Validator = &Memcached{}

// ValidateCreate implements webhook.Validator so a webhook will be registered for the type
func (r *Memcached) ValidateCreate() (err error) {
	memcachedlog.Info("validate create", "name", r.Name, "port", r.Spec.ContainerPort)

	// Verify Container Port is in Right Range
	if r.Spec.ContainerPort < minMemcachedPort {
		err = fmt.Errorf("Memcached Port %d should be between %d and 10000", r.Spec.ContainerPort, minMemcachedPort)
	}

	return
}

// ValidateUpdate implements webhook.Validator so a webhook will be registered for the type
func (r *Memcached) ValidateUpdate(_ runtime.Object) (err error) {
	memcachedlog.Info("validate update", "name", r.Name, "port", r.Spec.ContainerPort)

	// Verify Container Port is in Right Range
	if r.Spec.ContainerPort < minMemcachedPort {
		err = fmt.Errorf("Memcached Port %d should be between %d and 10000", r.Spec.ContainerPort, minMemcachedPort)
	}
	return
}

// ValidateDelete implements webhook.Validator so a webhook will be registered for the type
func (r *Memcached) ValidateDelete() (err error) {
	memcachedlog.Info("validate delete", "name", r.Name)
	if r.Labels["type"] == "critical" {
		err = fmt.Errorf("Cannot delete pod %s because it is marked critical", r.Name)
	}
	return
}
