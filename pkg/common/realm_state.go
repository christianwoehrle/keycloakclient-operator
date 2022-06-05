package common

import (
	"context"

	kc "github.com/christianwoehrle/keycloakclient-operator/pkg/apis/keycloak/v1alpha1"
	v1 "k8s.io/api/core/v1"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

type RealmState struct {
	Realm            *kc.KeycloakRealm
	RealmUserSecrets map[string]*v1.Secret
	Context          context.Context
	Keycloak         *kc.Keycloak
}

func NewRealmState(context context.Context, keycloak kc.Keycloak) *RealmState {
	return &RealmState{
		Context:  context,
		Keycloak: &keycloak,
	}
}

func (i *RealmState) Read(cr *kc.KeycloakRealm, realmClient KeycloakInterface, controllerClient client.Client) error {
	realm, err := realmClient.GetRealm(cr.Spec.Realm.Realm)
	if err != nil {
		i.Realm = nil
		return err
	}

	i.Realm = realm
	if realm == nil {
		return nil
	}

	return nil
}
