package common

import (
	"context"
	"time"

	kc "github.com/christianwoehrle/keycloakclient-operator/pkg/apis/keycloak/v1alpha1"
	"github.com/christianwoehrle/keycloakclient-operator/pkg/model"
	v12 "k8s.io/api/apps/v1"
	v1 "k8s.io/api/core/v1"
	apiErrors "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/api/meta"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

// BackupTime is used for generating a unique Backup job name
var BackupTime string

func init() {
	BackupTime = time.Now().Format("20060102-150405")
}

// The desired cluster state is defined by a list of actions that have to be run to
// get from the current state to the desired state
type DesiredClusterState []ClusterAction

func (d *DesiredClusterState) AddAction(action ClusterAction) DesiredClusterState {
	if action != nil {
		*d = append(*d, action)
	}
	return *d
}

func (d *DesiredClusterState) AddActions(actions []ClusterAction) DesiredClusterState {
	for _, action := range actions {
		*d = d.AddAction(action)
	}
	return *d
}

type ClusterState struct {
	KeycloakDeployment  *v12.StatefulSet
	KeycloakAdminSecret *v1.Secret
}

func NewClusterState() *ClusterState {
	return &ClusterState{}
}

func (i *ClusterState) Read(context context.Context, cr *kc.Keycloak, controllerClient client.Client) error {
	err := i.readKeycloakAdminSecretCurrentState(context, cr, controllerClient)
	if err != nil {
		return err
	}

	// Read other things
	return nil
}

func (i *ClusterState) readKeycloakAdminSecretCurrentState(context context.Context, cr *kc.Keycloak, controllerClient client.Client) error {
	keycloakAdminSecret := model.KeycloakAdminSecret(cr)
	keycloakAdminSecretSelector := model.KeycloakAdminSecretSelector(cr)

	err := controllerClient.Get(context, keycloakAdminSecretSelector, keycloakAdminSecret)

	if err != nil {
		// If the resource type doesn't exist on the cluster or does exist but is not found
		if meta.IsNoMatchError(err) || apiErrors.IsNotFound(err) {
			i.KeycloakAdminSecret = nil
		} else {
			return err
		}
	} else {
		i.KeycloakAdminSecret = keycloakAdminSecret.DeepCopy()
		cr.UpdateStatusSecondaryResources(i.KeycloakAdminSecret.Kind, i.KeycloakAdminSecret.Name)
	}
	return nil
}

func (i *ClusterState) IsResourcesReady(cr *kc.Keycloak) (bool, error) {
	if cr.Spec.Unmanaged {
		return true, nil
	}

	return true, nil
}
