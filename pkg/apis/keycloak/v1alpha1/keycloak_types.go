package v1alpha1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// KeycloakSpec defines the desired state of Keycloak.
// +k8s:openapi-gen=true
type KeycloakSpec struct {
	// When set to true, this Keycloak will be marked as unmanaged and will not be managed by this operator.
	// It can then be used for targeting purposes.
	// +optional
	// +kubebuilder:default:=true
	Unmanaged bool `json:"unmanaged,omitempty"`
	// Contains configuration for external Keycloak instances. Unmanaged needs to be set to true to use this.
	// +optional
	External KeycloakExternal `json:"external"`
}

type KeycloakExternal struct {
	// If set to true, this Keycloak will be treated as an external instance.
	// The unmanaged field also needs to be set to true if this field is true.
	Enabled bool `json:"enabled,omitempty"`
	// The URL to use for the keycloak admin API. Needs to be set if external is true.
	// +optional
	URL string `json:"url,omitempty"`
}

// KeycloakStatus defines the observed state of Keycloak.
// +k8s:openapi-gen=true
type KeycloakStatus struct {
	// Current phase of the operator.
	Phase StatusPhase `json:"phase"`
	// Human-readable message indicating details about current operator phase or error.
	Message string `json:"message"`
	// True if all resources are in a ready state and all work is done.
	Ready bool `json:"ready"`
	// A map of all the secondary resources types and names created for this CR. e.g "Deployment": [ "DeploymentName1", "DeploymentName2" ].
	SecondaryResources map[string][]string `json:"secondaryResources,omitempty"`
	// Version of Keycloak or RHSSO running on the cluster.
	Version string `json:"version"`
	// External URL for accessing Keycloak instance from outside the cluster. Is identical to external.URL if it's specified, otherwise is computed (e.g. from Ingress).
	ExternalURL string `json:"externalURL,omitempty"`
	// The secret where the admin credentials are to be found.
	CredentialSecret string `json:"credentialSecret"`
}

type StatusPhase string

var (
	NoPhase           StatusPhase
	PhaseReconciling  StatusPhase = "reconciling"
	PhaseFailing      StatusPhase = "failing"
	PhaseInitialising StatusPhase = "initialising"
)

// Keycloak is the Schema for the keycloaks API.
// +genclient
// +k8s:openapi-gen=true
// +kubebuilder:subresource:status
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object
type Keycloak struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   KeycloakSpec   `json:"spec,omitempty"`
	Status KeycloakStatus `json:"status,omitempty"`
}

// KeycloakList contains a list of Keycloak.
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object
type KeycloakList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []Keycloak `json:"items"`
}

func init() {
	SchemeBuilder.Register(&Keycloak{}, &KeycloakList{})
}

func (i *Keycloak) UpdateStatusSecondaryResources(kind string, resourceName string) {
	i.Status.SecondaryResources = UpdateStatusSecondaryResources(i.Status.SecondaryResources, kind, resourceName)
}
