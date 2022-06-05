package keycloak

import (
	"testing"

	"github.com/christianwoehrle/keycloakclient-operator/pkg/apis/keycloak/v1alpha1"
	"github.com/christianwoehrle/keycloakclient-operator/pkg/model"
	"github.com/stretchr/testify/assert"
)

func TestKeycloakReconciler_Test_Recreate_Credentials_When_Missig(t *testing.T) {
	// given
	cr := &v1alpha1.Keycloak{}
	secret := model.KeycloakAdminSecret(cr)

	// when
	secret.Data[model.AdminUsernameProperty] = nil
	secret.Data[model.AdminPasswordProperty] = nil
	secret = model.KeycloakAdminSecretReconciled(cr, secret)

	// then
	assert.NotEmpty(t, secret.Data[model.AdminUsernameProperty])
	assert.NotEmpty(t, secret.Data[model.AdminPasswordProperty])
}

func TestKeycloakReconciler_Test_Recreate_Does_Not_Update_Existing_Credentials(t *testing.T) {
	// given
	cr := &v1alpha1.Keycloak{}
	secret := model.KeycloakAdminSecret(cr)

	// when
	username := secret.Data[model.AdminUsernameProperty]
	password := secret.Data[model.AdminPasswordProperty]
	secret = model.KeycloakAdminSecretReconciled(cr, secret)

	// then
	assert.Equal(t, username, secret.Data[model.AdminUsernameProperty])
	assert.Equal(t, password, secret.Data[model.AdminPasswordProperty])
}
