package e2e

import (
	"testing"

	"github.com/operator-framework/operator-sdk/pkg/test"
)

const (
	userID = "test-user"
)

func NewKeycloakUserCRDTestStruct() *CRDTestStruct {
	return &CRDTestStruct{
		prepareEnvironmentSteps: []environmentInitializationStep{
			prepareKeycloaksCR,
			prepareKeycloakRealmCR,
		},
		testSteps: map[string]deployedOperatorTestStep{
			"keycloakUserBasicTest": {
				prepareTestEnvironmentSteps: []environmentInitializationStep{
					prepareKeycloakUserCR,
				},
				testFunction: keycloakUserBasicTest,
			},
		},
	}
}

func prepareKeycloakUserCR(t *testing.T, framework *test.Framework, ctx *test.Context, namespace string) error {
	keycloakUserCR := getKeycloakUserCR(namespace)
	return Create(framework, keycloakUserCR, ctx)
}

func keycloakUserBasicTest(t *testing.T, framework *test.Framework, ctx *test.Context, namespace string) error {
	return WaitForUserToBeReady(t, framework, namespace)
}
