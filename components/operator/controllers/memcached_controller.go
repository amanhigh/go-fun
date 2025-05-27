package controllers

import (
	"context"
	"fmt"

	appsv1 "k8s.io/api/apps/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/client-go/tools/record"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/log"

	cachev1beta1 "github.com/amanhigh/go-fun/components/operator/api/v1beta1"
)

// Constants used for defining the state and finalizers of the Memcached resource
const memcachedFinalizer = "cache.aman.com/finalizer"

// Definitions to manage status conditions
const (
	// typeAvailableMemcached represents the status of the Deployment reconciliation
	typeAvailableMemcached = "Available"
	// typeDegradedMemcached represents the status used when the custom resource is deleted and the finalizer operations are must to occur.
	typeDegradedMemcached = "Degraded"
)

// MemcachedReconciler reconciles a Memcached object
// This is the primary controller struct that handles the reconciliation of Memcached resources
type MemcachedReconciler struct {
	client.Client
	Scheme   *runtime.Scheme
	Recorder record.EventRecorder

	// Helper interfaces for modular functionality
	statusHelper    StatusHelper
	reconcileHelper ReconciliationHelper
	deployHelper    DeploymentHelper
}

// The following markers are used to generate the rules permissions (RBAC) on config/rbac using controller-gen
// when the command <make manifests> is executed.
// To know more about markers see: https://book.kubebuilder.io/reference/markers.html

// +kubebuilder:rbac:groups=cache.aman.com,resources=memcacheds,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=cache.aman.com,resources=memcacheds/status,verbs=get;update;patch
// +kubebuilder:rbac:groups=cache.aman.com,resources=memcacheds/finalizers,verbs=update
// +kubebuilder:rbac:groups=core,resources=events,verbs=create;patch
// +kubebuilder:rbac:groups=apps,resources=deployments,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=core,resources=pods,verbs=get;list;watch

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
//
// It is essential for the controller's reconciliation loop to be idempotent. By following the Operator
// pattern you will create Controllers which provide a reconcile function
// responsible for synchronizing resources until the desired state is reached on the cluster.
// Breaking this recommendation goes against the design principles of controller-runtime.
// and may lead to unforeseen consequences such as resources becoming stuck and requiring manual intervention.
//
// For further info:
// - About Operator Pattern: https://kubernetes.io/docs/concepts/extend-kubernetes/operator/
// - About Controllers: https://kubernetes.io/docs/concepts/architecture/controller/
// - https://pkg.go.dev/sigs.k8s.io/controller-runtime@v0.13.0/pkg/reconcile
func (r *MemcachedReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	log := log.FromContext(ctx)
	log.Info("Reconciling Memcached")

	// Fetch Instance - If not found, it means the Custom Resource was deleted
	memcached, err := r.reconcileHelper.FetchMemcachedInstance(ctx, req)
	var result ctrl.Result // Declare result variable
	if memcached == nil || err != nil {
		return ctrl.Result{}, fmt.Errorf("failed to fetch Memcached instance: %w", err)
	}

	// Initialize Status - Sets up initial conditions if none exist
	if err = r.statusHelper.InitializeStatus(ctx, memcached); err != nil {
		return ctrl.Result{}, fmt.Errorf("failed to initialize Memcached status: %w", err)
	}

	// Check if Resource is being deleted
	if memcached.GetDeletionTimestamp() != nil {
		log.Info("Resource is being deleted, performing finalizer operations")
		result, err = r.reconcileHelper.ExecuteFinalizer(ctx, memcached)
		if err != nil {
			return result, fmt.Errorf("failed to execute finalizer: %w", err)
		}
		return result, nil
	}

	// Add Finalizer if it doesn't exist
	if result, err = r.reconcileHelper.AddFinalizer(ctx, memcached); err != nil {
		return result, fmt.Errorf("failed to add finalizer: %w", err)
	}

	// Reconcile Deployment - Creates or updates the deployment to match desired state
	result, err = r.reconcileHelper.ReconcileDeployment(ctx, memcached)
	if err != nil {
		return result, fmt.Errorf("failed to reconcile deployment: %w", err)
	}
	return result, nil
}

// SetupWithManager sets up the controller with the Manager.
// Note that the Deployment will be also watched in order to ensure its
// desirable state on the cluster
func (r *MemcachedReconciler) SetupWithManager(mgr ctrl.Manager) error {
	// Initialize all helpers
	r.statusHelper = NewStatusHelper(r)
	r.deployHelper = NewDeploymentHelper(r)
	r.reconcileHelper = NewReconciliationHelper(r.statusHelper, r.deployHelper, r)

	err := ctrl.NewControllerManagedBy(mgr).
		For(&cachev1beta1.Memcached{}).
		// Inform Reconciler when any change happens in Owned Resources including deletion
		Owns(&appsv1.Deployment{}).
		Complete(r)
	if err != nil {
		return fmt.Errorf("failed to complete controller setup: %w", err)
	}
	return nil
}

/*
Reconcile Result Types
---------------------
The Reconcile function returns ctrl.Result to communicate the outcome of reconciliation:

1. ctrl.Result{}:
   - Indicates successful reconciliation
   - No further action needed

2. ctrl.Result{Requeue: true}:
   - Indicates reconciliation should be retried after a short delay
   - Useful when waiting for external resources
   - Default retry occurs after a few seconds

3. ctrl.Result{RequeueAfter: time.Second * 30}:
   - Schedules next reconciliation after specified duration
   - Useful for periodic checks or waiting for long operations
   - Example: Checking deployment status after 30 seconds

4. ctrl.Result{RequeueAfter: -1}:
   - Triggers immediate reconciliation
   - Useful when immediate retry is needed without delay

The ctrl.Result pattern enables flexible reconciliation strategies based on the
resource's needs and state.
*/
