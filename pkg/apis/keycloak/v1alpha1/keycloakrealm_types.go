package v1alpha1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// KeycloakRealmSpec defines the desired state of KeycloakRealm.
// +k8s:openapi-gen=true
type KeycloakRealmSpec struct {
	// When set to true, this KeycloakRealm will be marked as unmanaged and not be managed by this operator.
	// It can then be used for targeting purposes.
	// +optional
	Unmanaged bool `json:"unmanaged,omitempty"`
	// Selector for looking up Keycloak Custom Resources.
	// +kubebuilder:validation:Required
	InstanceSelector *metav1.LabelSelector `json:"instanceSelector,omitempty"`
	// Keycloak Realm REST object.
	// +kubebuilder:validation:Required
	Realm *KeycloakAPIRealm `json:"realm"`
}

type KeycloakAPIRealm struct {
	// +kubebuilder:validation:Required
	// +optional
	ID string `json:"id,omitempty"`
	// Realm name.
	// +kubebuilder:validation:Required
	Realm string `json:"realm"`
	// Realm enabled flag.
	// +optional
	Enabled bool `json:"enabled"`
	// Client scopes
	// +optional
	ClientScopes []KeycloakClientScope `json:"clientScopes,omitempty"`

	// Default role
	// +optional
	DefaultRole *RoleRepresentation `json:"defaultRole,omitempty"`
}

type KeycloakClientScope struct {
	// +optional
	Attributes map[string]string `json:"attributes,omitempty"`
	// +optional
	Description string `json:"description,omitempty"`
	// +optional
	ID string `json:"id,omitempty"`
	// +optional
	Name string `json:"name,omitempty"`
	// +optional
	Protocol string `json:"protocol,omitempty"`
	// Protocol Mappers.
	// +optional
	ProtocolMappers []KeycloakProtocolMapper `json:"protocolMappers,omitempty"`
}

type RoleRepresentationArray []RoleRepresentation

// https://www.keycloak.org/docs-api/11.0/rest-api/index.html#_rolesrepresentation
type RolesRepresentation struct {
	// Client Roles
	// +optional
	Client map[string]RoleRepresentationArray `json:"client,omitempty"`

	// Realm Roles
	// +optional
	Realm []RoleRepresentation `json:"realm,omitempty"`
}

// https://www.keycloak.org/docs-api/11.0/rest-api/index.html#_rolerepresentation
type RoleRepresentation struct {
	// Role Attributes
	// +optional
	Attributes map[string][]string `json:"attributes,omitempty"`

	// Client Role
	// +optional
	ClientRole *bool `json:"clientRole,omitempty"`

	// Composite
	// +optional
	Composite *bool `json:"composite,omitempty"`

	// Composites
	// +optional
	Composites *RoleRepresentationComposites `json:"composites,omitempty"`

	// Container Id
	// +optional
	ContainerID string `json:"containerId,omitempty"`

	// Description
	// +optional
	Description string `json:"description,omitempty"`

	// Id
	// +optional
	ID string `json:"id,omitempty"`

	// Name
	Name string `json:"name"`
}

type ScopeMappingRepresentationArray []ScopeMappingRepresentation

// https://www.keycloak.org/docs-api/11.0/rest-api/index.html#_scopemappingrepresentation
type ScopeMappingRepresentation struct {
	// Client
	// +optional
	Client string `json:"client,omitempty"`

	// Client Scope
	// +optional
	ClientScope string `json:"clientScope,omitempty"`

	// Roles
	// +optional
	Roles []string `json:"roles,omitempty"`

	// Self
	// +optional
	Self string `json:"self,omitempty"`
}

// https://www.keycloak.org/docs-api/11.0/rest-api/index.html#_rolerepresentation-composites
type RoleRepresentationComposites struct {
	// Map client => []role
	// +optional
	Client map[string][]string `json:"client,omitempty"`

	// Realm roles
	// +optional
	Realm []string `json:"realm,omitempty"`
}

type KeycloakUserRole struct {
	ID          string `json:"id,omitempty"`
	Name        string `json:"name,omitempty"`
	Description string `json:"description,omitempty"`
	Composite   bool   `json:"composite,omitempty"`
	ClientRole  bool   `json:"clientRole,omitempty"`
	ContainerID string `json:"containerId,omitempty"`
}

type TokenResponse struct {
	// Token Response Access Token.
	// +optional
	AccessToken string `json:"access_token"`
	// Token Response Expired In setting.
	// +optional
	ExpiresIn int `json:"expires_in"`
	// Token Response Refresh Expires In setting.
	// +optional
	RefreshExpiresIn int `json:"refresh_expires_in"`
	// Token Response Refresh Token.
	// +optional
	RefreshToken string `json:"refresh_token"`
	// Token Response Token Type.
	// +optional
	TokenType string `json:"token_type"`
	// Token Response Not Before Policy setting.
	// +optional
	NotBeforePolicy int `json:"not-before-policy"`
	// Token Response Session State.
	// +optional
	SessionState string `json:"session_state"`
	// Token Response Error.
	// +optional
	Error string `json:"error"`
	// Token Response Error Description.
	// +optional
	ErrorDescription string `json:"error_description"`
}

// KeycloakRealmStatus defines the observed state of KeycloakRealm
// +k8s:openapi-gen=true
type KeycloakRealmStatus struct {
	// Current phase of the operator.
	Phase StatusPhase `json:"phase"`
	// Human-readable message indicating details about current operator phase or error.
	Message string `json:"message"`
	// True if all resources are in a ready state and all work is done.
	Ready bool `json:"ready"`
	// A map of all the secondary resources types and names created for this CR. e.g "Deployment": [ "DeploymentName1", "DeploymentName2" ]
	SecondaryResources map[string][]string `json:"secondaryResources,omitempty"`
	// TODO
	LoginURL string `json:"loginURL"`
}

// KeycloakRealm is the Schema for the keycloakrealms API
// +genclient
// +k8s:openapi-gen=true
// +kubebuilder:subresource:status
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object
type KeycloakRealm struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   KeycloakRealmSpec   `json:"spec,omitempty"`
	Status KeycloakRealmStatus `json:"status,omitempty"`
}

// KeycloakRealmList contains a list of KeycloakRealm
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object
type KeycloakRealmList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []KeycloakRealm `json:"items"`
}

func init() {
	SchemeBuilder.Register(&KeycloakRealm{}, &KeycloakRealmList{})
}

func (i *KeycloakRealm) UpdateStatusSecondaryResources(kind string, resourceName string) {
	i.Status.SecondaryResources = UpdateStatusSecondaryResources(i.Status.SecondaryResources, kind, resourceName)
}
