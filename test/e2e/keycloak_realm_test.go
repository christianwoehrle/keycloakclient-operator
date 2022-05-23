package e2e

import (
	"testing"

	keycloakv1alpha1 "github.com/keycloak/keycloak-operator/pkg/apis/keycloak/v1alpha1"
	"github.com/operator-framework/operator-sdk/pkg/test"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

const (
	realmName = "test-realm"
)

func NewKeycloakRealmsCRDTestStruct() *CRDTestStruct {
	return &CRDTestStruct{
		prepareEnvironmentSteps: []environmentInitializationStep{
			prepareUnmanagedKeycloaksCR,
		},
		testSteps: map[string]deployedOperatorTestStep{
			"unmanagedKeycloakRealmTest": {
				testFunction: keycloakUnmanagedRealmTest,
			},
		},
	}
}

func getKeycloakRealmCR(namespace string) *keycloakv1alpha1.KeycloakRealm {
	return &keycloakv1alpha1.KeycloakRealm{
		ObjectMeta: metav1.ObjectMeta{
			Name:      testKeycloakRealmCRName,
			Namespace: namespace,
			Labels:    CreateLabel(namespace),
		},
		Spec: keycloakv1alpha1.KeycloakRealmSpec{
			InstanceSelector: &metav1.LabelSelector{
				MatchLabels: CreateLabel(namespace),
			},
			Realm: &keycloakv1alpha1.KeycloakAPIRealm{
				ID:      realmName,
				Realm:   realmName,
				Enabled: true,
			},
		},
	}
}

func keycloakRealmBasicTest(t *testing.T, framework *test.Framework, ctx *test.Context, namespace string) error {
	return WaitForRealmToBeReady(t, framework, namespace)
}

func keycloakUnmanagedRealmTest(t *testing.T, framework *test.Framework, ctx *test.Context, namespace string) error {
	keycloakRealmCR := getKeycloakRealmCR(namespace)
	keycloakRealmCR.Spec.Unmanaged = true

	err := Create(framework, keycloakRealmCR, ctx)
	if err != nil {
		return err
	}

	err = WaitForRealmToBeReady(t, framework, namespace)
	if err != nil {
		return err
	}

	return err
}

func prepareExternalKeycloakRealmCR(t *testing.T, f *test.Framework, ctx *test.Context, namespace string) error {
	keycloakRealmCR := getKeycloakRealmCR(namespace)
	keycloakRealmCR.Spec.Unmanaged = true

	err := Create(f, keycloakRealmCR, ctx)

	if err != nil {
		return err
	}

	err = WaitForRealmToBeReady(t, f, namespace)
	if err != nil {
		return err
	}

	return err
}
