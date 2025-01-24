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

package controllers

import (
	"context"

	appsv1 "k8s.io/api/apps/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/client-go/tools/record"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/log"

	cachev1beta1 "github.com/amanhigh/go-fun/components/operator/api/v1beta1"
)

const memcachedFinalizer = "cache.aman.com/finalizer"

// Definitions to manage status conditions
const (
	// typeAvailableMemcached represents the status of the Deployment reconciliation
	typeAvailableMemcached = "Available"
	// typeDegradedMemcached represents the status used when the custom resource is deleted and the finalizer operations are must to occur.
	typeDegradedMemcached = "Degraded"
)

// MemcachedReconciler reconciles a Memcached object
type MemcachedReconciler struct {
	client.Client
	Scheme   *runtime.Scheme
	Recorder record.EventRecorder

	statusHelper    StatusHelper
	reconcileHelper ReconciliationHelper
	deployHelper    DeploymentHelper
}

// The following markers are used to generate the rules permissions (RBAC) on config/rbac using controller-gen
// when the command <make manifests> is executed.
// To know more about markers see: https://book.kubebuilder.io/reference/markers.html

//+kubebuilder:rbac:groups=cache.aman.com,resources=memcacheds,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=cache.aman.com,resources=memcacheds/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=cache.aman.com,resources=memcacheds/finalizers,verbs=update
//+kubebuilder:rbac:groups=core,resources=events,verbs=create;patch
//+kubebuilder:rbac:groups=apps,resources=deployments,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=core,resources=pods,verbs=get;list;watch

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
//
// It is essential for the controller's reconciliation loop to be idempotent. By following the Operator
// pattern you will create Controllers which provide a reconcile function
// responsible for synchronizing resources until the desired state is reached on the cluster.
// Breaking this recommendation goes against the design principles of controller-runtime.
// and may lead to unforeseen consequences such as resources becoming stuck and requiring manual intervention.
// For further info:
// - About Operator Pattern: https://kubernetes.io/docs/concepts/extend-kubernetes/operator/
// - About Controllers: https://kubernetes.io/docs/concepts/architecture/controller/
// - https://pkg.go.dev/sigs.k8s.io/controller-runtime@v0.13.0/pkg/reconcile
func (r *MemcachedReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	log := log.FromContext(ctx)
	log.Info("Reconciling Memcached")

	// Fetch Instance
	memcached, err := r.reconcileHelper.FetchMemcachedInstance(ctx, req)
	if memcached == nil || err != nil {
		return ctrl.Result{}, err
	}

	// Initialize Status
	if err := r.statusHelper.InitializeStatus(ctx, memcached); err != nil {
		return ctrl.Result{}, err
	}

	// Handle Finalizers
	if result, err := r.reconcileHelper.HandleFinalizers(ctx, memcached); err != nil || result.Requeue {
		return result, err
	}

	// Reconcile Deployment
	return r.reconcileHelper.ReconcileDeployment(ctx, memcached)
}

// SetupWithManager sets up the controller with the Manager.
// Note that the Deployment will be also watched in order to ensure its
// desirable state on the cluster
func (r *MemcachedReconciler) SetupWithManager(mgr ctrl.Manager) error {
	// Initialize all helpers
	r.statusHelper = NewStatusHelper(r)
	r.reconcileHelper = NewReconciliationHelper(r.statusHelper, r)
	r.deployHelper = NewDeploymentHelper(r)

	return ctrl.NewControllerManagedBy(mgr).
		For(&cachev1beta1.Memcached{}).
		//Inform Reconciler when any change happens in Owned Resources including deletion.
		Owns(&appsv1.Deployment{}).
		Complete(r)
}

/**
Reconcile Result
---------
In the context of Kubernetes controllers, the Reconcile function typically returns a value of type ctrl.Result. This value represents the result of the reconciliation process, and can have the following possible values:

ctrl.Result{}: Indicates that the reconciliation was successful, and that no further action is required at this time.

ctrl.Result{Requeue: true}: Indicates that the reconciliation was not successful, and that the controller should attempt to reconcile the resource again after a short delay (usually a few seconds). This is useful when the controller is waiting for some external resource to become available before proceeding with the reconciliation.

ctrl.Result{RequeueAfter: time.Second * 30}: Indicates that the reconciliation was not successful, and that the controller should attempt to reconcile the resource again after a specified delay (in this case, 30 seconds). This is useful when the controller is waiting for some long-running process to complete before proceeding with the reconciliation.

ctrl.Result{RequeueAfter: -1}: Indicates that the reconciliation was not successful, and that the controller should attempt to reconcile the resource again as soon as possible (i.e., without any delay). This is useful when the controller needs to retry the reconciliation immediately, without waiting for any external events.

Overall, the ctrl.Result type provides a flexible way for controllers to communicate the results of their reconciliation process to the Kubernetes API server. By using different combinations of Requeue and RequeueAfter values, controllers can implement a wide variety of reconciliation strategies, depending on the specific requirements of the resource being managed.

**/
