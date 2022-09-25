package keycloakrealm

import (
	"context"
	"fmt"
	"time"

	"github.com/christianwoehrle/keycloakclient-operator/pkg/apis/keycloak/v1alpha1"
	kc "github.com/christianwoehrle/keycloakclient-operator/pkg/apis/keycloak/v1alpha1"
	corev1 "k8s.io/api/core/v1"
	kubeerrors "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/client-go/tools/record"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller"
	"sigs.k8s.io/controller-runtime/pkg/handler"
	logf "sigs.k8s.io/controller-runtime/pkg/log"
	"sigs.k8s.io/controller-runtime/pkg/manager"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
	"sigs.k8s.io/controller-runtime/pkg/source"
)

const (
	RealmFinalizer    = "realm.cleanup"
	RequeueDelayError = 60 * time.Second
	ControllerName    = "controller_keycloakrealm"
)

var log = logf.Log.WithName(ControllerName)

/**
* USER ACTION REQUIRED: This is a scaffold file intended for the user to modify with their own Controller
* business logic.  Delete these comments after modifying this file.*
 */

// Add creates a new KeycloakRealm Controller and adds it to the Manager. The Manager will set fields on the Controller
// and Start it when the Manager is Started.
func Add(mgr manager.Manager) error {
	return add(mgr, newReconciler(mgr))
}

// newReconciler returns a new reconcile.Reconciler
func newReconciler(mgr manager.Manager) reconcile.Reconciler {
	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)

	return &ReconcileKeycloakRealm{
		client:   mgr.GetClient(),
		scheme:   mgr.GetScheme(),
		cancel:   cancel,
		context:  ctx,
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

	// Watch for changes to primary resource KeycloakRealm
	err = c.Watch(&source.Kind{Type: &kc.KeycloakRealm{}}, &handler.EnqueueRequestForObject{})
	if err != nil {
		return err
	}

	// Make sure to watch the credential secrets
	err = c.Watch(&source.Kind{Type: &corev1.Secret{}}, &handler.EnqueueRequestForOwner{
		IsController: true,
		OwnerType:    &kc.KeycloakRealm{},
	})
	if err != nil {
		return err
	}

	return nil
}

// blank assignment to verify that ReconcileKeycloakRealm implements reconcile.Reconciler
var _ reconcile.Reconciler = &ReconcileKeycloakRealm{}

// ReconcileKeycloakRealm reconciles a KeycloakRealm object
type ReconcileKeycloakRealm struct {
	// This client, initialized using mgr.Client() above, is a split client
	// that reads objects from the cache and writes to the apiserver
	client   client.Client
	scheme   *runtime.Scheme
	context  context.Context
	cancel   context.CancelFunc
	recorder record.EventRecorder
}

// Reconcile reads that state of the cluster for a KeycloakRealm object and makes changes based on the state read
// and what is in the KeycloakRealm.Spec
func (r *ReconcileKeycloakRealm) Reconcile(request reconcile.Request) (reconcile.Result, error) {
	reqLogger := log.WithValues("Request.Namespace", request.Namespace, "Request.Name", request.Name)
	reqLogger.Info("Reconciling KeycloakRealm")

	// Fetch the KeycloakRealm instance
	instance := &kc.KeycloakRealm{}
	err := r.client.Get(context.TODO(), request.NamespacedName, instance)
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

	if !instance.Spec.Unmanaged {
		log.Info(fmt.Sprintf("ignore unmanaged==true flag in realm %v/%v handle it as unmanaged", instance.Namespace, instance.Name))
	}

	return reconcile.Result{Requeue: false}, r.manageSuccess(instance, instance.DeletionTimestamp != nil)

}

func (r *ReconcileKeycloakRealm) manageSuccess(realm *kc.KeycloakRealm, deleted bool) error {
	realm.Status.Ready = true
	realm.Status.Message = ""
	realm.Status.Phase = v1alpha1.PhaseReconciling

	err := r.client.Status().Update(r.context, realm)
	if err != nil {
		log.Error(err, "unable to update status")
	}

	// Finalizer already set?
	finalizerExists := false
	for _, finalizer := range realm.Finalizers {
		if finalizer == RealmFinalizer {
			finalizerExists = true
			break
		}
	}

	// Resource created and finalizer exists: nothing to do
	if !deleted && finalizerExists {
		return nil
	}

	// Resource created and finalizer does not exist: add finalizer
	if !deleted && !finalizerExists {
		realm.Finalizers = append(realm.Finalizers, RealmFinalizer)
		log.Info(fmt.Sprintf("added finalizer to keycloak realm %v/%v",
			realm.Namespace,
			realm.Spec.Realm.Realm))

		return r.client.Update(r.context, realm)
	}

	// Otherwise remove the finalizer
	newFinalizers := []string{}
	for _, finalizer := range realm.Finalizers {
		if finalizer == RealmFinalizer {
			log.Info(fmt.Sprintf("removed finalizer from keycloak realm %v/%v",
				realm.Namespace,
				realm.Spec.Realm.Realm))

			continue
		}
		newFinalizers = append(newFinalizers, finalizer)
	}

	realm.Finalizers = newFinalizers
	return r.client.Update(r.context, realm)
}
