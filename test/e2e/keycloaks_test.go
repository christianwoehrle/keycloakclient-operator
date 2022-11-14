package e2e

import (
	"context"
	"testing"

	apiErrors "k8s.io/apimachinery/pkg/api/errors"

	"k8s.io/client-go/kubernetes"

	v1 "k8s.io/api/core/v1"

	"github.com/pkg/errors"

	"github.com/stretchr/testify/assert"

	keycloakv1alpha1 "github.com/christianwoehrle/keycloakclient-operator/pkg/apis/keycloak/v1alpha1"
	framework "github.com/operator-framework/operator-sdk/pkg/test"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func NewExternalKeycloaksCRDTestStruct() *CRDTestStruct {
	return &CRDTestStruct{
		prepareEnvironmentSteps: []environmentInitializationStep{
			prepareExternalKeycloaksCR,
		},
		testSteps: map[string]deployedOperatorTestStep{
			"keycloakExternalDeploymentTest": {testFunction: keycloakExternalDeploymentTest},
		},
	}
}

func getKeycloakCR(namespace string) *keycloakv1alpha1.Keycloak {
	return &keycloakv1alpha1.Keycloak{
		ObjectMeta: metav1.ObjectMeta{
			Name:      testKeycloakCRName,
			Namespace: namespace,
			Labels:    CreateLabel(namespace),
		},
		Spec: keycloakv1alpha1.KeycloakSpec{},
	}
}

func getUnmanagedKeycloakCR(namespace string) *keycloakv1alpha1.Keycloak {
	keycloak := getKeycloakCR(namespace)
	keycloak.Name = testKeycloakCRName
	keycloak.Spec.Unmanaged = true
	return keycloak
}

func getExternalKeycloakCR(namespace string, url string) *keycloakv1alpha1.Keycloak {
	keycloak := getUnmanagedKeycloakCR(namespace)
	keycloak.Name = testKeycloakCRName
	keycloak.Labels = CreateLabel(namespace)
	keycloak.Spec.External.Enabled = true
	keycloak.Spec.External.URL = url
	return keycloak
}

func getDeployedKeycloakCR(f *framework.Framework, namespace string) keycloakv1alpha1.Keycloak {
	keycloakCR := keycloakv1alpha1.Keycloak{}
	_ = GetNamespacedObject(f, namespace, testKeycloakCRName, &keycloakCR)
	return keycloakCR
}

func getExternalKeycloakSecret(f *framework.Framework, namespace string) (*v1.Secret, error) {
	secret, err := f.KubeClient.CoreV1().Secrets(namespace).Get(context.TODO(), "credential-"+testKeycloakCRName, metav1.GetOptions{})

	if err != nil {
		return nil, err
	}

	return &v1.Secret{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "credential-" + testKeycloakCRName,
			Namespace: namespace,
		},
		Data:       secret.Data,
		StringData: secret.StringData,
		Type:       secret.Type,
	}, nil
}

func prepareUnmanagedKeycloaksCR(t *testing.T, f *framework.Framework, ctx *framework.Context, namespace string) error {
	keycloakCR := getUnmanagedKeycloakCR(namespace)
	err := Create(f, keycloakCR, ctx)
	if err != nil {
		return err
	}

	err = WaitForKeycloakToBeReady(t, f, namespace, testKeycloakCRName)
	if err != nil {
		return err
	}

	return err
}

func prepareExternalKeycloaksCR(t *testing.T, f *framework.Framework, ctx *framework.Context, namespace string) error {
	keycloakURL := "http://keycloak.local:80"

	secret, err := getExternalKeycloakSecret(f, namespace)
	if err != nil && !apiErrors.IsNotFound(err) {
		return err
	}

	if err != nil && !apiErrors.IsNotFound(err) {
		secret = &v1.Secret{
			ObjectMeta: metav1.ObjectMeta{
				Name:      "credential-" + testKeycloakCRName,
				Namespace: namespace,
			},
			StringData: map[string]string{
				"user":     "admin",
				"password": "admin",
			},
			Type: v1.SecretTypeOpaque,
		}

		err = Create(f, secret, ctx)
		if err != nil {
			return err
		}
	}

	externalKeycloakCR := getExternalKeycloakCR(namespace, keycloakURL)

	err = Create(f, externalKeycloakCR, ctx)
	if err != nil && !apiErrors.IsAlreadyExists(err) {
		return err
	}

	err = WaitForKeycloakToBeReady(t, f, namespace, testKeycloakCRName)
	if err != nil {
		return err
	}

	return err
}

func keycloakExternalDeploymentTest(t *testing.T, f *framework.Framework, ctx *framework.Context, namespace string) error {
	keycloakCR := getDeployedKeycloakCR(f, namespace)
	assert.NotEmpty(t, keycloakCR.Status.ExternalURL)

	err := WaitForCondition(t, f.KubeClient, func(t *testing.T, c kubernetes.Interface) error {
		sts, err := f.KubeClient.AppsV1().StatefulSets(namespace).List(context.TODO(), metav1.ListOptions{})

		if err != nil {
			return errors.Errorf("list StatefulSet failed, ignoring for %v: %v", pollRetryInterval, err)
		}
		if len(sts.Items) == 1 {
			return nil
		}
		return errors.Errorf("should find one Statefulset, as the cluster has been prepared with a keycloak installation")
	})
	return err
}
