package e2e

import (
	"context"
	"crypto/tls"
	"net/http"
	"testing"

	"k8s.io/client-go/kubernetes"

	v1 "k8s.io/api/core/v1"

	"github.com/pkg/errors"

	"github.com/stretchr/testify/assert"

	keycloakv1alpha1 "github.com/christianwoehrle/keycloakclient-operator/pkg/apis/keycloak/v1alpha1"
	framework "github.com/operator-framework/operator-sdk/pkg/test"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

const podName = "keycloak-0"
const extraLabelName = "extra"
const extraLabelValue = "value"

func NewUnmanagedKeycloaksCRDTestStruct() *CRDTestStruct {
	return &CRDTestStruct{
		prepareEnvironmentSteps: []environmentInitializationStep{
			prepareUnmanagedKeycloaksCR,
		},
		testSteps: map[string]deployedOperatorTestStep{
			"keycloakUnmanagedDeploymentTest": {testFunction: keycloakUnmanagedDeploymentTest},
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
	keycloak.Name = testKeycloakUnmanagedCRName
	keycloak.Spec.Unmanaged = true
	return keycloak
}

func getExternalKeycloakCR(namespace string, url string) *keycloakv1alpha1.Keycloak {
	keycloak := getUnmanagedKeycloakCR(namespace)
	keycloak.Name = testKeycloakExternalCRName
	keycloak.Labels = CreateExternalLabel(namespace)
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
			Name:      "credential-" + testKeycloakExternalCRName,
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

	err = WaitForKeycloakToBeReady(t, f, namespace, testKeycloakUnmanagedCRName)
	if err != nil {
		return err
	}

	return err
}

func prepareExternalKeycloaksCR(t *testing.T, f *framework.Framework, ctx *framework.Context, namespace string) error {
	keycloakCR := getDeployedKeycloakCR(f, namespace)
	keycloakURL := keycloakCR.Status.ExternalURL

	secret, err := getExternalKeycloakSecret(f, namespace)
	if err != nil {
		return err
	}

	err = Create(f, secret, ctx)
	if err != nil {
		return err
	}

	externalKeycloakCR := getExternalKeycloakCR(namespace, keycloakURL)
	err = Create(f, externalKeycloakCR, ctx)
	if err != nil {
		return err
	}

	err = WaitForKeycloakToBeReady(t, f, namespace, testKeycloakExternalCRName)
	if err != nil {
		return err
	}

	return err
}

func keycloakDeploymentTest(t *testing.T, f *framework.Framework, ctx *framework.Context, namespace string) error {
	keycloakCR := getDeployedKeycloakCR(f, namespace)
	assert.NotEmpty(t, keycloakCR.Status.ExternalURL)

	err := WaitForKeycloakToBeReady(t, f, namespace, testKeycloakCRName)
	if err != nil {
		return err
	}

	keycloakURL := keycloakCR.Status.ExternalURL

	// Skipping TLS verification is actually part of the test. In Kubernetes, if there's no signing
	// manager installed, Keycloak will generate its own, self-signed cert. Of course
	// we don't have a matching truststore for it, hence we need to skip TLS verification.
	http.DefaultTransport.(*http.Transport).TLSClientConfig = &tls.Config{InsecureSkipVerify: true} //nolint

	err = WaitForSuccessResponse(t, f, keycloakURL+"/auth")
	if err != nil {
		return err
	}

	metricsBody, err := GetSuccessfulResponseBody(keycloakURL + "/auth/realms/master/metrics")
	if err != nil {
		return err
	}

	masterRealmBody, err := GetSuccessfulResponseBody(keycloakURL + "/auth/realms/master")
	if err != nil {
		return err
	}

	// there should be a redirect/rewrite from the metrics endpoint to master realm
	assert.Equal(t, masterRealmBody, metricsBody)

	return err
}

func keycloakUnmanagedDeploymentTest(t *testing.T, f *framework.Framework, ctx *framework.Context, namespace string) error {
	keycloakCR := getDeployedKeycloakCR(f, namespace)
	assert.Empty(t, keycloakCR.Status.ExternalURL)

	err := WaitForCondition(t, f.KubeClient, func(t *testing.T, c kubernetes.Interface) error {
		sts, err := f.KubeClient.AppsV1().StatefulSets(namespace).List(context.TODO(), metav1.ListOptions{})
		if err != nil {
			return errors.Errorf("list StatefulSet failed, ignoring for %v: %v", pollRetryInterval, err)
		}
		if len(sts.Items) != 1 {
			return nil
		}
		return errors.Errorf("should find one Statefulset, as the cluster has been prepared with a keycloak isntallation")
	})
	return err
}

func keycloakDeploymentDefaultImagePullPolicyTest(t *testing.T, f *framework.Framework, ctx *framework.Context, namespace string) error {
	// check that the default imagePolicy is set to Always, even though not defined by the user
	keycloakPod := v1.Pod{}
	err := GetNamespacedObject(f, namespace, podName, &keycloakPod)
	assert.Contains(t, keycloakPod.Spec.Containers[0].ImagePullPolicy, v1.PullAlways)
	return err
}
