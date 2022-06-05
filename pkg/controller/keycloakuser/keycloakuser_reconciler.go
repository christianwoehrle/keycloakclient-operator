package keycloakuser

import (
	"fmt"

	"github.com/christianwoehrle/keycloakclient-operator/pkg/apis/keycloak/v1alpha1"
	"github.com/christianwoehrle/keycloakclient-operator/pkg/common"
)

func GetUserRealmRolesDesiredState(state *common.UserState, realmRoles []string, realmName string) []common.ClusterAction {
	var assignRoles []common.ClusterAction
	var removeRoles []common.ClusterAction

	for _, role := range realmRoles {
		// Is the role available for this user?
		roleRef := state.GetAvailableRealmRole(role)
		if roleRef == nil {
			continue
		}

		// Role requested but not assigned?
		if !containsRole(state.RealmRoles, role) {
			assignRoles = append(assignRoles, &common.AssignRealmRoleAction{
				UserID: state.User.ID,
				Ref:    roleRef,
				Realm:  realmName,
				Msg:    fmt.Sprintf("assign realm role %v to user %v", role, state.User.UserName),
			})
		}
	}

	for _, role := range state.RealmRoles {
		// Role assigned but not requested?
		if !containsRoleID(realmRoles, role.Name) {
			removeRoles = append(removeRoles, &common.RemoveRealmRoleAction{
				UserID: state.User.ID,
				Ref:    role,
				Realm:  realmName,
				Msg:    fmt.Sprintf("remove realm role %v from user %v", role.Name, state.User.UserName),
			})
		}
	}

	return append(assignRoles, removeRoles...)
}

func GetUserClientRolesDesiredState(state *common.UserState, clientRoles map[string][]string, realmName string) []common.ClusterAction {
	actions := []common.ClusterAction{}

	for _, client := range state.Clients {
		actions = append(actions, SyncRolesForClient(state, client.ClientID, clientRoles, realmName)...)
	}

	return actions
}

func SyncRolesForClient(state *common.UserState, clientID string, clientRoles map[string][]string, realmName string) []common.ClusterAction {
	var assignRoles []common.ClusterAction
	var removeRoles []common.ClusterAction

	for _, role := range clientRoles[clientID] {
		// Is the role available for this user?
		roleRef := state.GetAvailableClientRole(role, clientID)
		if roleRef == nil {
			continue
		}

		// Valid client?
		client := state.GetClientByID(clientID)
		if client == nil {
			continue
		}

		// Role requested but not assigned?
		if !containsRole(state.ClientRoles[clientID], role) {
			assignRoles = append(assignRoles, &common.AssignClientRoleAction{
				UserID:   state.User.ID,
				ClientID: client.ID,
				Ref:      roleRef,
				Realm:    realmName,
				Msg:      fmt.Sprintf("assign role %v of client %v to user %v", role, clientID, state.User.UserName),
			})
		}
	}

	for _, role := range state.ClientRoles[clientID] {
		// Valid client?
		client := state.GetClientByID(clientID)
		if client == nil {
			continue
		}

		// Role assigned but not requested?
		if !containsRoleID(clientRoles[clientID], role.Name) {
			removeRoles = append(removeRoles, &common.RemoveClientRoleAction{
				UserID:   state.User.ID,
				ClientID: client.ID,
				Ref:      role,
				Realm:    realmName,
				Msg:      fmt.Sprintf("remove role %v of client %v from user %v", role.Name, clientID, state.User.UserName),
			})
		}
	}

	return append(assignRoles, removeRoles...)
}

func containsRole(list []*v1alpha1.KeycloakUserRole, id string) bool {
	for _, item := range list {
		if item.ID == id {
			return true
		}
	}
	return false
}

func containsRoleID(list []string, id string) bool {
	for _, item := range list {
		if item == id {
			return true
		}
	}
	return false
}
