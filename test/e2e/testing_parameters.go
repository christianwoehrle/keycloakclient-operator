package e2e

import "flag"

var (
	isProductBuild bool
)

func init() {
	flag.BoolVar(&isProductBuild, "product", false, "Using Keycloak")
}

const keycloakProfile = "keycloak"

func currentProfile() string {
	return keycloakProfile
}
