package e2e

import (
	"testing"

	"github.com/operator-framework/operator-sdk/pkg/test"
)

const (
	realmName                  = "test-realm"
	testOperatorIDPDisplayName = "Test Operator IDP"
)

func NewKeycloakRealmsCRDTestStruct() *CRDTestStruct {
	return &CRDTestStruct{
		prepareEnvironmentSteps: []environmentInitializationStep{
			prepareKeycloaksCR,
		},
		testSteps: map[string]deployedOperatorTestStep{
			"keycloakRealmBasicTest": {
				prepareTestEnvironmentSteps: []environmentInitializationStep{
					prepareKeycloakRealmCR,
				},
				testFunction: keycloakRealmBasicTest,
			},
			"keycloakRealmWithEventsTest": {
				testFunction: keycloakRealmWithEventsTest,
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

func keycloakRealmWithEventsTest(t *testing.T, framework *test.Framework, ctx *test.Context, namespace string) error {
	keycloakRealmCR := getKeycloakRealmCR(namespace)

	keycloakRealmCR.Spec.Realm.EventsEnabled = &[]bool{true}[0]
	keycloakRealmCR.Spec.Realm.EnabledEventTypes = []string{"SEND_RESET_PASSWORD", "LOGIN_ERROR"}

	err := Create(framework, keycloakRealmCR, ctx)
	if err != nil {
		return err
	}

	return WaitForRealmToBeReady(t, framework, namespace)
}
