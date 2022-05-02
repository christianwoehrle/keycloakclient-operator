package v1alpha1

type KeycloakAPIUser struct {
	// User ID.
	ID string `json:"id,omitempty"`
	// User Name.
	UserName string `json:"username,omitempty"`
	// First Name.
	FirstName string `json:"firstName,omitempty"`
	// Last Name.
	LastName string `json:"lastName,omitempty"`
	// Email.
	Email string `json:"email,omitempty"`
	// True if email has already been verified.
	EmailVerified bool `json:"emailVerified,omitempty"`
	// User enabled flag.
	Enabled bool `json:"enabled,omitempty"`
	// A set of Realm Roles.
	RealmRoles []string `json:"realmRoles,omitempty"`
	// A set of Client Roles.
	ClientRoles map[string][]string `json:"clientRoles,omitempty"`
	// A set of Required Actions.
	RequiredActions []string `json:"requiredActions,omitempty"`
	// A set of Groups.
	Groups []string `json:"groups,omitempty"`
	// A set of Federated Identities.
	FederatedIdentities []FederatedIdentity `json:"federatedIdentities,omitempty"`
	// A set of Credentials.
	Credentials []KeycloakCredential `json:"credentials,omitempty"`
	// A set of Attributes.
	Attributes map[string][]string `json:"attributes,omitempty"`
}

type KeycloakCredential struct {
	// Credential Type.
	Type string `json:"type,omitempty"`
	// Credential Value.
	Value string `json:"value,omitempty"`
	// True if this credential object is temporary.
	Temporary bool `json:"temporary,omitempty"`
}

type FederatedIdentity struct {
	// Federated Identity Provider.
	IdentityProvider string `json:"identityProvider,omitempty"`
	// Federated Identity User ID.
	UserID string `json:"userId,omitempty"`
	// Federated Identity User Name.
	UserName string `json:"userName,omitempty"`
}
