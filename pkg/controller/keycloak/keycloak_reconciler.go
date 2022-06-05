package keycloak

import (
	kc "github.com/christianwoehrle/keycloakclient-operator/pkg/apis/keycloak/v1alpha1"
	"github.com/christianwoehrle/keycloakclient-operator/pkg/common"
	"github.com/christianwoehrle/keycloakclient-operator/pkg/model"
)

type Reconciler interface {
	Reconcile(clusterState *common.ClusterState, cr *kc.Keycloak) (common.DesiredClusterState, error)
}

type KeycloakReconciler struct { // nolint
}

func NewKeycloakReconciler() *KeycloakReconciler {
	return &KeycloakReconciler{}
}

func (i *KeycloakReconciler) Reconcile(clusterState *common.ClusterState, cr *kc.Keycloak) common.DesiredClusterState {
	desired := common.DesiredClusterState{}

	desired = desired.AddAction(i.GetKeycloakAdminSecretDesiredState(clusterState, cr))

	return desired
}

func (i *KeycloakReconciler) GetKeycloakAdminSecretDesiredState(clusterState *common.ClusterState, cr *kc.Keycloak) common.ClusterAction {
	keycloakAdminSecret := model.KeycloakAdminSecret(cr)

	if clusterState.KeycloakAdminSecret == nil {
		return common.GenericCreateAction{
			Ref: keycloakAdminSecret,
			Msg: "Create Keycloak admin secret",
		}
	}
	return common.GenericUpdateAction{
		Ref: model.KeycloakAdminSecretReconciled(cr, clusterState.KeycloakAdminSecret),
		Msg: "Update Keycloak admin secret",
	}
}
