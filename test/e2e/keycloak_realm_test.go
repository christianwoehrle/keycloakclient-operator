package e2e

import (
	"testing"

	keycloakv1alpha1 "github.com/christianwoehrle/keycloakclient-operator/pkg/apis/keycloak/v1alpha1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"github.com/operator-framework/operator-sdk/pkg/test"
)

const (
	realmName                  = "test-realm"
	testOperatorIDPDisplayName = "Test Operator IDP"
)

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
			Unmanaged: true,
			Realm: &keycloakv1alpha1.KeycloakAPIRealm{
				Enabled: true,
				Realm:   "test-realm",
				ID:      "test-realm",
			},
		},
	}
}

func NewKeycloakRealmsCRDTestStruct() *CRDTestStruct {
	return &CRDTestStruct{
		prepareEnvironmentSteps: []environmentInitializationStep{
			prepareUnmanagedKeycloaksCR,
		},
		testSteps: map[string]deployedOperatorTestStep{
			"keycloakRealmBasicTest": {
				prepareTestEnvironmentSteps: []environmentInitializationStep{
					prepareKeycloakRealmCR,
				},
				testFunction: keycloakRealmBasicTest,
			},
			"unmanagedKeycloakRealmTest": {
				testFunction: keycloakUnmanagedRealmTest,
			},
		},
	}
}

func prepareKeycloakRealmCR(t *testing.T, framework *test.Framework, ctx *test.Context, namespace string) error {
	keycloakRealmCR := getKeycloakRealmCR(namespace)
	return Create(framework, keycloakRealmCR, ctx)
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
