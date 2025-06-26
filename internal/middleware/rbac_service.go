package middleware

import (
	"context"
	"fmt"

	// permissionEntity "bitbucket.org/abhishek_fordel/we-care-holidays-backend-golang/modules/permissions/domain/entity"
	roleEntity "bitbucket.org/abhishek_fordel/we-care-holidays-backend-golang/modules/roles/domain/entity"

	permissionRepo "bitbucket.org/abhishek_fordel/we-care-holidays-backend-golang/modules/permissions/data/mongodb/repository"
	roleRepo "bitbucket.org/abhishek_fordel/we-care-holidays-backend-golang/modules/roles/data/mongodb/repository"
	userRepo "bitbucket.org/abhishek_fordel/we-care-holidays-backend-golang/modules/users/data/mongodb/repository"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type RBACService interface {
	GetUserPermissions(ctx context.Context, userID primitive.ObjectID) ([]Permission, roleEntity.RoleScope, error)
	ValidatePermission(ctx context.Context, authCtx *AuthContext, resource, action string) bool
	GetScopeFilter(ctx context.Context, authCtx *AuthContext, resource string) map[string]interface{}
	CanCreateRole(ctx context.Context, authCtx *AuthContext, targetScope roleEntity.RoleScope) bool
	ValidateRolePermissions(ctx context.Context, authCtx *AuthContext, permissionIDs []string) error
}

type rbacService struct {
	userRepo       *userRepo.UserRepositoryMongo
	roleRepo       *roleRepo.RoleRepositoryMongo
	permissionRepo *permissionRepo.PermissionRepositoryMongo
}

func NewRBACService(
	userRepo *userRepo.UserRepositoryMongo,
	roleRepo *roleRepo.RoleRepositoryMongo,
	permissionRepo *permissionRepo.PermissionRepositoryMongo,
) RBACService {
	return &rbacService{
		userRepo:       userRepo,
		roleRepo:       roleRepo,
		permissionRepo: permissionRepo,
	}
}

func (r *rbacService) GetUserPermissions(ctx context.Context, userID primitive.ObjectID) ([]Permission, roleEntity.RoleScope, error) {
	if r == nil {
		return nil, "", fmt.Errorf("rbac service is nil")
	}
	if r.userRepo == nil {
		return nil, "", fmt.Errorf("user repository is nil")
	}
	if r.roleRepo == nil {
		return nil, "", fmt.Errorf("role repository is nil")
	}
	if r.permissionRepo == nil {
		return nil, "", fmt.Errorf("permission repository is nil")
	}

	// Get user to find their role
	user, err := r.userRepo.GetByID(ctx, userID)
	if err != nil {
		return nil, "", err
	}
	if user == nil {
		return nil, "", fmt.Errorf("user not found")
	}

	// Get role
	roleID, err := primitive.ObjectIDFromHex(user.RoleID)
	if err != nil {
		return nil, "", err
	}

	role, err := r.roleRepo.GetByID(ctx, roleID)
	if err != nil {
		return nil, "", err
	}
	if role == nil {
		return nil, "", fmt.Errorf("role not found")
	}

	var permissions []Permission
	for _, permIDStr := range role.Permissions {
		permID, err := primitive.ObjectIDFromHex(permIDStr)
		if err != nil {
			continue
		}

		perm, err := r.permissionRepo.GetByID(ctx, permID)
		if err != nil {
			continue
		}
		if perm != nil {
			permissions = append(permissions, Permission{
				Resource: perm.Resource,
				Action:   string(perm.Action),
			})
		}
	}

	return permissions, role.Scope, nil // Return role scope
}

func (r *rbacService) ValidatePermission(ctx context.Context, authCtx *AuthContext, resource, action string) bool {
	for _, perm := range authCtx.Permissions {
		if perm.Resource == resource && perm.Action == action {
			return true
		}
	}
	return false
}

func (r *rbacService) GetScopeFilter(ctx context.Context, authCtx *AuthContext, resource string) map[string]interface{} {
	filter := make(map[string]interface{})

	// Check if user has read permission for this resource
	if !r.ValidatePermission(ctx, authCtx, resource, "read") {
		// No permission - return filter that matches nothing
		filter["_id"] = primitive.NewObjectID()
		return filter
	}

	// Apply scope-based filtering based on role scope
	switch authCtx.RoleScope {
	case roleEntity.RoleScopeGlobal:
		// No additional filter - user can see all
	case roleEntity.RoleScopeOrganization:
		if authCtx.OrganizationID != nil {
			filter["organizationId"] = *authCtx.OrganizationID
		} else {
			// If user has no org, they see nothing
			filter["_id"] = primitive.NewObjectID()
		}
	case roleEntity.RoleScopeSelf:
		// User can only see their own data
		filter["userId"] = authCtx.UserID
	default:
		// Unknown scope - return filter that matches nothing
		filter["_id"] = primitive.NewObjectID()
	}

	return filter
}

// Fixed: Dynamic role creation validation
func (r *rbacService) CanCreateRole(ctx context.Context, authCtx *AuthContext, targetScope roleEntity.RoleScope) bool {
	// Check if user has roles:create permission
	if !r.ValidatePermission(ctx, authCtx, "roles", "create") {
		return false
	}

	return r.isScopeAllowed(targetScope, authCtx.RoleScope)

}

// Check if target scope is allowed given the maximum scope
func (r *rbacService) isScopeAllowed(targetScope, userScope roleEntity.RoleScope) bool {
	scopeHierarchy := map[roleEntity.RoleScope]int{
		roleEntity.RoleScopeSelf:         1,
		roleEntity.RoleScopeOrganization: 2,
		roleEntity.RoleScopeGlobal:       3,
	}

	targetLevel := scopeHierarchy[targetScope]
	userLevel := scopeHierarchy[userScope]

	return targetLevel <= userLevel
}

// Fixed: Enhanced permission validation for dynamic roles
func (r *rbacService) ValidateRolePermissions(ctx context.Context, authCtx *AuthContext, permissionIDs []string) error {

	for _, permIDStr := range permissionIDs {
		permID, err := primitive.ObjectIDFromHex(permIDStr)
		if err != nil {
			return fmt.Errorf("invalid permission ID: %s", permIDStr)
		}

		perm, err := r.permissionRepo.GetByID(ctx, permID)
		if err != nil {
			return err
		}

		if perm == nil {
			return fmt.Errorf("permission not found: %s", permIDStr)
		}

		// Simple check - user must have the permission they're trying to assign
		if !r.ValidatePermission(ctx, authCtx, perm.Resource, string(perm.Action)) {
			return fmt.Errorf("cannot assign permission you don't possess: %s", perm.String())
		}
	}

	return nil
}
