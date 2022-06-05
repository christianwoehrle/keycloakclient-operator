package model

// Constants for a community Keycloak installation
const (
	ApplicationName                  = "keycloak"
	AdminUsernameProperty            = "ADMIN_USERNAME"
	AdminPasswordProperty            = "ADMIN_PASSWORD"
	ServingCertSecretName            = "sso-x509-https-secret" // nolint
	ClientSecretName                 = ApplicationName + "-client-secret"
	ClientSecretClientIDProperty     = "CLIENT_ID"
	ClientSecretClientSecretProperty = "CLIENT_SECRET"
)

var PodLabels = map[string]string{}
