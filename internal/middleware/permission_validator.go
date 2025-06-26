package middleware

import (
	"context"
	"fmt"
	"strings"

	permissionRepo "bitbucket.org/abhishek_fordel/we-care-holidays-backend-golang/modules/permissions/data/mongodb/repository"
	permissionEntity "bitbucket.org/abhishek_fordel/we-care-holidays-backend-golang/modules/permissions/domain/entity"
)

type PermissionValidator struct {
	permissionRepo *permissionRepo.PermissionRepositoryMongo
}

func NewPermissionValidator(permissionRepo *permissionRepo.PermissionRepositoryMongo) *PermissionValidator {
	return &PermissionValidator{permissionRepo: permissionRepo}
}

// ValidatePermissionString validates permission in format "resource:action:scope"
func (pv *PermissionValidator) ValidatePermissionString(ctx context.Context, permStr string) (*permissionEntity.Permission, error) {
	parts := strings.Split(permStr, ":")
	if len(parts) != 3 {
		return nil, fmt.Errorf("invalid permission format. Expected 'resource:action:scope', got: %s", permStr)
	}

	resource, action, scope := parts[0], parts[1], parts[2]

	// Validate action type
	if !pv.isValidAction(action) {
		return nil, fmt.Errorf("invalid action: %s. Valid actions are: read, write, create, update, delete", action)
	}

	// Validate scope type
	if !pv.isValidScope(scope) {
		return nil, fmt.Errorf("invalid scope: %s. Valid scopes are: global, organization, self", scope)
	}

	// Check if permission exists in database
	filter := map[string]interface{}{
		"resource": resource,
		"action":   action,
		"scope":    scope,
	}

	permissions, _, err := pv.permissionRepo.List(ctx, filter, 1, 1)
	if err != nil {
		return nil, fmt.Errorf("error checking permission existence: %w", err)
	}

	if len(permissions) == 0 {
		return nil, fmt.Errorf("permission does not exist: %s", permStr)
	}

	return permissions[0], nil
}

// ValidatePermissionStrings validates multiple permission strings
func (pv *PermissionValidator) ValidatePermissionStrings(ctx context.Context, permStrs []string) ([]*permissionEntity.Permission, error) {
	var permissions []*permissionEntity.Permission

	for _, permStr := range permStrs {
		perm, err := pv.ValidatePermissionString(ctx, permStr)
		if err != nil {
			return nil, err
		}
		permissions = append(permissions, perm)
	}

	return permissions, nil
}

// ConvertPermissionStringsToIDs converts permission strings to permission IDs
func (pv *PermissionValidator) ConvertPermissionStringsToIDs(ctx context.Context, permStrs []string) ([]string, error) {
	var permissionIDs []string

	for _, permStr := range permStrs {
		perm, err := pv.ValidatePermissionString(ctx, permStr)
		if err != nil {
			return nil, err
		}
		permissionIDs = append(permissionIDs, perm.ID.Hex())
	}

	return permissionIDs, nil
}

// GetAvailablePermissions returns all available permissions for reference
func (pv *PermissionValidator) GetAvailablePermissions(ctx context.Context) ([]*permissionEntity.Permission, error) {
	permissions, _, err := pv.permissionRepo.List(ctx, make(map[string]interface{}), 1, 1000)
	if err != nil {
		return nil, err
	}
	return permissions, nil
}

// GetPermissionsByResource returns all permissions for a specific resource
func (pv *PermissionValidator) GetPermissionsByResource(ctx context.Context, resource string) ([]*permissionEntity.Permission, error) {
	filter := map[string]interface{}{
		"resource": resource,
	}

	permissions, _, err := pv.permissionRepo.List(ctx, filter, 1, 1000)
	if err != nil {
		return nil, err
	}
	return permissions, nil
}

// Validate action type
func (pv *PermissionValidator) isValidAction(action string) bool {
	validActions := map[string]bool{
		"read":   true,
		"write":  true,
		"create": true,
		"update": true,
		"delete": true,
	}
	return validActions[action]
}

// Validate scope type
func (pv *PermissionValidator) isValidScope(scope string) bool {
	validScopes := map[string]bool{
		"global":       true,
		"organization": true,
		"self":         true,
	}
	return validScopes[scope]
}

// GetValidActions returns all valid action types
func (pv *PermissionValidator) GetValidActions() []string {
	return []string{"read", "write", "create", "update", "delete"}
}

// GetValidScopes returns all valid scope types
func (pv *PermissionValidator) GetValidScopes() []string {
	return []string{"global", "organization", "self"}
}
