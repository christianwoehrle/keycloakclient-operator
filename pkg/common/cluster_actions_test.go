package common
/*
import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/christianwoehrle/keycloakclient-operator/pkg/apis/keycloak/v1alpha1"
	"github.com/stretchr/testify/assert"
)

func TestClusterAction_CreateKeycloakCLient(t *testing.T) {
	// given
	realm := getDummyRealm()
	keycloakClient := getDummyClient()
	var method string
	var ID string

	ctx := context.TODO()
	handler := http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		assert.Equal(t, fmt.Sprintf(ClientCreatePath, realm.Spec.Realm.Realm), req.URL.Path)
		method = req.Method

		var kcc v1alpha1.KeycloakAPIClient

		err := json.NewDecoder(req.Body).Decode(&kcc)

		if err != nil {
			fmt.Println("Err:", err)
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		ID = kcc.ID
		w.Header().Set("location", "https://keycloak.tue.private.dwpbank.io/auth/admin/realms/dwpbank/clients/4710701c-3f81-4b96-8ba0-f6b73651fbca")
		w.WriteHeader(409)

	})
	server := httptest.NewServer(handler)
	defer server.Close()

	client := Client{
		requester: server.Client(),
		URL:       server.URL,
		token:     "dummy",
	}

	createClientAction := CreateClientAction{
		Ref:   keycloakClient,
		Msg:   "Create Keycloak Client",
		Realm: realm.Name,
	}

	clientState := NewClientState(ctx, realm.DeepCopy(), *getDummyKeycloak())
	// client muss ein KeycloakInterface sein.
	err := clientState.Read(ctx, keycloakClient, client, apiServerClient)
	assert.NoError(t, err)

	// Figure out the actions to keep the realms up to date with
	// the desired state
	actionRunner := NewClusterAndKeycloakActionRunner(ctx, r.client, r.scheme, keycloakClient, authenticated)

	// Run all actions to keep the realms updated
	createClientAction.Run(actionRunner)

	// when
	uid, err := client.CreateClient(keycloakClient, realm.Spec.Realm.Realm)

	// then
	// correct path expected on httptest server
	assert.NoError(t, err)
	assert.Equal(t, "4710701c-3f81-4b96-8ba0-f6b73651fbca", uid)
	assert.Equal(t, method, "POST")
	assert.Equal(t, keycloakClient.Spec.Client.ID, ID)
}
*/ 
