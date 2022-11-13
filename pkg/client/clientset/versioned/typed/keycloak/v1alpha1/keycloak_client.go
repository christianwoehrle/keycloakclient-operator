// Code generated by client-gen. DO NOT EDIT.

package v1alpha1

import (
	v1alpha1 "github.com/christianwoehrle/keycloakclient-operator/pkg/apis/keycloak/v1alpha1"
	"github.com/christianwoehrle/keycloakclient-operator/pkg/client/clientset/versioned/scheme"
	rest "k8s.io/client-go/rest"
)

type KeycloakV1alpha1Interface interface {
	RESTClient() rest.Interface
	KeycloaksGetter
	KeycloakClientsGetter
	KeycloakRealmsGetter
}

// KeycloakV1alpha1Client is used to interact with features provided by the keycloak.org group.
type KeycloakV1alpha1Client struct {
	restClient rest.Interface
}

func (c *KeycloakV1alpha1Client) Keycloaks(namespace string) KeycloakInterface {
	return newKeycloaks(c, namespace)
}

func (c *KeycloakV1alpha1Client) KeycloakClients(namespace string) KeycloakClientInterface {
	return newKeycloakClients(c, namespace)
}

func (c *KeycloakV1alpha1Client) KeycloakRealms(namespace string) KeycloakRealmInterface {
	return newKeycloakRealms(c, namespace)
}

// NewForConfig creates a new KeycloakV1alpha1Client for the given config.
func NewForConfig(c *rest.Config) (*KeycloakV1alpha1Client, error) {
	config := *c
	if err := setConfigDefaults(&config); err != nil {
		return nil, err
	}
	client, err := rest.RESTClientFor(&config)
	if err != nil {
		return nil, err
	}
	return &KeycloakV1alpha1Client{client}, nil
}

// NewForConfigOrDie creates a new KeycloakV1alpha1Client for the given config and
// panics if there is an error in the config.
func NewForConfigOrDie(c *rest.Config) *KeycloakV1alpha1Client {
	client, err := NewForConfig(c)
	if err != nil {
		panic(err)
	}
	return client
}

// New creates a new KeycloakV1alpha1Client for the given RESTClient.
func New(c rest.Interface) *KeycloakV1alpha1Client {
	return &KeycloakV1alpha1Client{c}
}

func setConfigDefaults(config *rest.Config) error {
	gv := v1alpha1.SchemeGroupVersion
	config.GroupVersion = &gv
	config.APIPath = "/apis"
	config.NegotiatedSerializer = scheme.Codecs.WithoutConversion()

	if config.UserAgent == "" {
		config.UserAgent = rest.DefaultKubernetesUserAgent()
	}

	return nil
}

// RESTClient returns a RESTClient that is used to communicate
// with API server by this client implementation.
func (c *KeycloakV1alpha1Client) RESTClient() rest.Interface {
	if c == nil {
		return nil
	}
	return c.restClient
}
