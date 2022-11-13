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
	"github.com/christianwoehrle/keycloakclient-operator/pkg/model"
	framework "github.com/operator-framework/operator-sdk/pkg/test"
	v1apps "k8s.io/api/apps/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

const podName = "keycloak-0"
const extraLabelName = "extra"
const extraLabelValue = "value"

func NewKeycloaksCRDTestStruct() *CRDTestStruct {
	return &CRDTestStruct{
		prepareEnvironmentSteps: []environmentInitializationStep{
			prepareKeycloaksCRWithExtension,
		},
		testSteps: map[string]deployedOperatorTestStep{
			"keycloakDeploymentTest": {testFunction: keycloakDeploymentTest},
		},
	}
}

func NewKeycloaksWithLabelsCRDTestStruct() *CRDTestStruct {
	return &CRDTestStruct{
		prepareEnvironmentSteps: []environmentInitializationStep{
			prepareKeycloaksCRWithPodLabels,
		},
		testSteps: map[string]deployedOperatorTestStep{
			"keycloakWithPodLabelsDeploymentTest": {testFunction: keycloakDeploymentWithLabelsTest},
		},
	}
}

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

func NewKeycloaksSSLTestStruct() *CRDTestStruct {
	return &CRDTestStruct{
		prepareEnvironmentSteps: []environmentInitializationStep{
			prepareKeycloaksSSLWithDB,
		},
		testSteps: map[string]deployedOperatorTestStep{
			"keycloakSSLDBTest": {testFunction: keycloakSSLDBTest},
		},
	}
}

func NewKeycloaksWithDefaultImagePullPolicyTestStruct() *CRDTestStruct {
	return &CRDTestStruct{
		prepareEnvironmentSteps: []environmentInitializationStep{
			prepareKeycloaksCR,
		},
		testSteps: map[string]deployedOperatorTestStep{
			"keycloakDeploymentDefaultImagePullPolicyTest": {testFunction: keycloakDeploymentDefaultImagePullPolicyTest},
		},
	}
}

func NewKeycloakStatefulSetSelectorTestStruct() *CRDTestStruct {
	return &CRDTestStruct{
		testSteps: map[string]deployedOperatorTestStep{
			"keycloakStatefulSetSelectorTest": {testFunction: keycloakStatefulSetSelectorTest},
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
		Spec: keycloakv1alpha1.KeycloakSpec{
			Instances:      1,
			ExternalAccess: keycloakv1alpha1.KeycloakExternalAccess{Enabled: true},
			Profile:        currentProfile(),
		},
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

func prepareKeycloaksCR(t *testing.T, f *framework.Framework, ctx *framework.Context, namespace string) error {
	return deployKeycloaksCR(t, f, ctx, namespace, getKeycloakCR(namespace))
}

func prepareKeycloaksCRWithPodLabels(t *testing.T, f *framework.Framework, ctx *framework.Context, namespace string) error {
	keycloakCR := getKeycloakCR(namespace)
	keycloakCR.Spec.KeycloakDeploymentSpec.PodLabels = map[string]string{"cr.first.label": "first.value", "cr.second.label": "second.value"}
	return deployKeycloaksCR(t, f, ctx, namespace, keycloakCR)
}

func deployKeycloaksCR(t *testing.T, f *framework.Framework, ctx *framework.Context, namespace string, keycloakCR *keycloakv1alpha1.Keycloak) error {
	err := doWorkaroundIfNecessary(f, ctx, namespace)
	if err != nil {
		return err
	}

	err = Create(f, keycloakCR, ctx)
	if err != nil {
		return err
	}

	err = WaitForStatefulSetReplicasReady(t, f.KubeClient, model.ApplicationName, namespace)
	if err != nil {
		return err
	}

	return err
}

func prepareUnmanagedKeycloaksCR(t *testing.T, f *framework.Framework, ctx *framework.Context, namespace string) error {
	err := doWorkaroundIfNecessary(f, ctx, namespace)
	if err != nil {
		return err
	}

	keycloakCR := getUnmanagedKeycloakCR(namespace)
	err = Create(f, keycloakCR, ctx)
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
	assert.NotEmpty(t, keycloakCR.Status.InternalURL)
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
func keycloakDeploymentWithLabelsTest(t *testing.T, f *framework.Framework, ctx *framework.Context, namespace string) error {
	// check that the creation labels are present
	keycloakPod := v1.Pod{}
	_ = GetNamespacedObject(f, namespace, podName, &keycloakPod)
	assert.Contains(t, keycloakPod.Labels, "cr.first.label")
	assert.Contains(t, keycloakPod.Labels, "cr.second.label")

	//add runtime labels to the pod (as if it was the existing labels from previous installation)
	keycloakPod.ObjectMeta.Labels["pod.label.one"] = "value1"
	keycloakPod.ObjectMeta.Labels["pod.label.two"] = "value2"
	err := Update(f, &keycloakPod)
	if err != nil {
		return err
	}

	//modify the CR adding labels, to see ifthe reconcile process also adds the labels
	keycloakCR := getDeployedKeycloakCR(f, namespace)
	newlabels := map[string]string{"cr-reconc.label.one": "value1", "cr-reconc.label.two": "value1"}
	keycloakCR.Spec.KeycloakDeploymentSpec.PodLabels = model.AddPodLabels(&keycloakCR, newlabels)
	err = Update(f, &keycloakCR)
	if err != nil {
		return err
	}

	// we need to wait for the reconciliation
	err = WaitForPodHavingLabels(t, f.KubeClient, podName, namespace, keycloakCR.Spec.KeycloakDeploymentSpec.PodLabels)
	if err != nil {
		return err
	}

	// assert that runtime  labels added directly to the pod are still there
	// assert that new labels added to the CR are also present in the pod
	_ = GetNamespacedObject(f, namespace, podName, &keycloakPod)
	// Labels set in the CR on the creation
	assert.Contains(t, keycloakPod.Labels, "cr.first.label")
	assert.Contains(t, keycloakPod.Labels, "cr.second.label")
	// Labels in the pod set by the user
	assert.Contains(t, keycloakPod.Labels, "pod.label.one")
	assert.Contains(t, keycloakPod.Labels, "pod.label.two")
	// Labels added to the CR during runtime
	assert.Contains(t, keycloakPod.Labels, "cr-reconc.label.one")
	assert.Contains(t, keycloakPod.Labels, "cr-reconc.label.two")

	return nil
}

func keycloakUnmanagedDeploymentTest(t *testing.T, f *framework.Framework, ctx *framework.Context, namespace string) error {
	keycloakCR := getDeployedKeycloakCR(f, namespace)
	assert.Empty(t, keycloakCR.Status.InternalURL)
	assert.Empty(t, keycloakCR.Status.ExternalURL)

	err := WaitForCondition(t, f.KubeClient, func(t *testing.T, c kubernetes.Interface) error {
		sts, err := f.KubeClient.AppsV1().StatefulSets(namespace).List(context.TODO(), metav1.ListOptions{})
		if err != nil {
			return errors.Errorf("list StatefulSet failed, ignoring for %v: %v", pollRetryInterval, err)
		}
		if len(sts.Items) == 0 {
			return nil
		}
		return errors.Errorf("found Statefulsets, this shouldn't be the case")
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

func keycloakStatefulSetSelectorTest(t *testing.T, f *framework.Framework, ctx *framework.Context, namespace string) error {
	// We can't modify existing selector on StatefulSet.
	// The selector might be wrongly set by e.g. RH-SSO 7.5.2.
	// In such case, we need to recreate the StatefulSet

	keycloakCR := getKeycloakCR(namespace)

	var kcDeployment *v1apps.StatefulSet
	if currentProfile() == keycloakProfile {
		kcDeployment = model.KeycloakDeployment(keycloakCR, nil, nil)
	} else {
		kcDeployment = model.RHSSODeployment(keycloakCR, nil, nil)
	}

	// Add extra selectors/labels which should be removed by the Operator
	kcDeployment.Spec.Selector.MatchLabels[extraLabelName] = extraLabelValue
	kcDeployment.Spec.Template.Labels[extraLabelName] = extraLabelValue

	// Let's create a faulty StatefulSet before Operator does it
	err := Create(f, kcDeployment, ctx)
	if err != nil {
		return err
	}

	// Operator should take over, recreate the SS and all should be fine
	err = deployKeycloaksCR(t, f, ctx, namespace, keycloakCR)
	if err != nil {
		return err
	}

	err = WaitForCondition(t, f.KubeClient, func(t *testing.T, c kubernetes.Interface) error {
		foundSS := &v1apps.StatefulSet{}
		err = GetNamespacedObject(f, namespace, kcDeployment.Name, foundSS)
		if err != nil {
			return err
		}
		if _, ok := foundSS.Spec.Selector.MatchLabels[extraLabelName]; ok {
			return errors.Errorf("Bad Selector not removed")
		}
		return nil
	})
	return err
}
