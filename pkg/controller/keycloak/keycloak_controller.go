package keycloak

import (
	"context"
	"time"

	"github.com/christianwoehrle/keycloakclient-operator/version"

	"k8s.io/client-go/tools/record"
	"sigs.k8s.io/controller-runtime/pkg/controller"

	corev1 "k8s.io/api/core/v1"

	"github.com/christianwoehrle/keycloakclient-operator/pkg/apis/keycloak/v1alpha1"
	kc "github.com/christianwoehrle/keycloakclient-operator/pkg/apis/keycloak/v1alpha1"
	keycloakv1alpha1 "github.com/christianwoehrle/keycloakclient-operator/pkg/apis/keycloak/v1alpha1"
	"github.com/christianwoehrle/keycloakclient-operator/pkg/common"
	"github.com/pkg/errors"

	kubeerrors "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/handler"
	logf "sigs.k8s.io/controller-runtime/pkg/log"
	"sigs.k8s.io/controller-runtime/pkg/manager"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
	"sigs.k8s.io/controller-runtime/pkg/source"
)

var log = logf.Log.WithName("controller_keycloak")

const (
	RequeueDelay      = 150 * time.Second
	RequeueDelayError = 60 * time.Second
	ControllerName    = "keycloak-controller"
)

// Add creates a new Keycloak Controller and adds it to the Manager. The Manager will set fields on the Controller
// and Start it when the Manager is Started.
func Add(mgr manager.Manager) error {
	return add(mgr, newReconciler(mgr))
}

// newReconciler returns a new reconcile.Reconciler
func newReconciler(mgr manager.Manager) reconcile.Reconciler {
	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)
	client := mgr.GetClient()

	return &ReconcileKeycloak{
		client:   client,
		scheme:   mgr.GetScheme(),
		context:  ctx,
		cancel:   cancel,
		recorder: mgr.GetEventRecorderFor(ControllerName),
	}
}

// add adds a new Controller to mgr with r as the reconcile.Reconciler
func add(mgr manager.Manager, r reconcile.Reconciler) error {
	// Create a new controller
	c, err := controller.New(ControllerName, mgr, controller.Options{Reconciler: r})
	if err != nil {
		return err
	}

	// Watch for changes to primary resource Keycloak
	err = c.Watch(&source.Kind{Type: &keycloakv1alpha1.Keycloak{}}, &handler.EnqueueRequestForObject{})
	if err != nil {
		return err
	}

	if err := common.WatchSecondaryResource(c, ControllerName, common.SecretKind, &corev1.Secret{}, &kc.Keycloak{}); err != nil {
		return err
	}

	return nil
}

// blank assignment to verify that ReconcileKeycloak implements reconcile.Reconciler
var _ reconcile.Reconciler = &ReconcileKeycloak{}

// ReconcileKeycloak reconciles a Keycloak object
type ReconcileKeycloak struct {
	// This client, initialized using mgr.Client() above, is a split client
	// that reads objects from the cache and writes to the apiserver
	client   client.Client
	scheme   *runtime.Scheme
	context  context.Context
	cancel   context.CancelFunc
	recorder record.EventRecorder
}

func (r *ReconcileKeycloak) Reconcile(request reconcile.Request) (reconcile.Result, error) {
	reqLogger := log.WithValues("Request.Namespace", request.Namespace, "Request.Name", request.Name)
	reqLogger.Info("Reconciling Keycloak")

	// Fetch the Keycloak instance
	instance := &keycloakv1alpha1.Keycloak{}

	err := r.client.Get(r.context, request.NamespacedName, instance)
	if err != nil {
		if kubeerrors.IsNotFound(err) {
			// Request object not found, could have been deleted after reconcile request.
			// Owned objects are automatically garbage collected. For additional cleanup logic use finalizers.
			// Return and don't requeue
			return reconcile.Result{}, nil
		}
		// Error reading the object - requeue the request.
		return reconcile.Result{}, err
	}
	currentState := common.NewClusterState()

	if instance.Spec.Unmanaged {
		return r.ManageSuccess(instance, currentState)
	}

	if instance.Spec.External.Enabled {
		return r.ManageError(instance, errors.Errorf("if external.enabled is true, unmanaged also needs to be true"))
	}

	// Read current state
	err = currentState.Read(r.context, instance, r.client)
	if err != nil {
		return r.ManageError(instance, err)
	}

	// Get Action to reconcile current state into desired state
	reconciler := NewKeycloakReconciler()
	desiredState := reconciler.Reconcile(currentState, instance)

	// Run the actions to reach the desired state
	actionRunner := common.NewClusterActionRunner(r.context, r.client, r.scheme, instance)
	err = actionRunner.RunAll(desiredState)
	if err != nil {
		return r.ManageError(instance, err)
	}

	return r.ManageSuccess(instance, currentState)
}

func (r *ReconcileKeycloak) ManageError(instance *v1alpha1.Keycloak, issue error) (reconcile.Result, error) {
	r.recorder.Event(instance, "Warning", "ProcessingError", issue.Error())

	instance.Status.Message = issue.Error()
	instance.Status.Ready = false
	instance.Status.Phase = v1alpha1.PhaseFailing

	r.setVersion(instance)

	err := r.client.Status().Update(r.context, instance)
	if err != nil {
		log.Error(err, "unable to update status")
	}

	return reconcile.Result{
		RequeueAfter: RequeueDelayError,
		Requeue:      true,
	}, nil
}

func (r *ReconcileKeycloak) ManageSuccess(instance *v1alpha1.Keycloak, currentState *common.ClusterState) (reconcile.Result, error) {
	// Check if the resources are ready
	resourcesReady, err := currentState.IsResourcesReady(instance)
	if err != nil {
		return r.ManageError(instance, err)
	}

	instance.Status.Ready = resourcesReady
	instance.Status.Message = ""

	// If resources are ready and we have not errored before now, we are in a reconciling phase
	if resourcesReady {
		instance.Status.Phase = v1alpha1.PhaseReconciling
	} else {
		instance.Status.Phase = v1alpha1.PhaseInitialising
	}

	if instance.Spec.External.URL != "" { //nolint
		instance.Status.ExternalURL = instance.Spec.External.URL
	}

	// Let the clients know where the admin credentials are stored
	if currentState.KeycloakAdminSecret != nil {
		instance.Status.CredentialSecret = currentState.KeycloakAdminSecret.Name
	}

	r.setVersion(instance)

	err = r.client.Status().Update(r.context, instance)
	if err != nil {
		log.Error(err, "unable to update status")
		return reconcile.Result{
			RequeueAfter: RequeueDelayError,
			Requeue:      true,
		}, nil
	}

	log.Info("desired cluster state met")
	return reconcile.Result{RequeueAfter: RequeueDelay}, nil
}

func (r *ReconcileKeycloak) setVersion(instance *v1alpha1.Keycloak) {
	instance.Status.Version = version.Version
}
