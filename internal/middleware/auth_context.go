package middleware

import (
	"context"

	roleEntity "bitbucket.org/abhishek_fordel/we-care-holidays-backend-golang/modules/roles/domain/entity"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type ContextKey string

const (
	UserIDKey         ContextKey = "user_id"
	UserRoleKey       ContextKey = "user_role"
	OrganizationIDKey ContextKey = "organization_id"
	PermissionsKey    ContextKey = "permissions"
	TokenKey          ContextKey = "token"
)

type AuthContext struct {
	UserID         primitive.ObjectID   `json:"userId"`
	Role           string               `json:"role"`
	RoleScope      roleEntity.RoleScope `json:"roleScope"` // NEW: Role's scope
	Permissions    []Permission         `json:"permissions"`
	OrganizationID *primitive.ObjectID  `json:"organizationId,omitempty"`
	Token          string               `json:"token"`
}

type Permission struct {
	Resource string `json:"resource"`
	Action   string `json:"action"`
}

func GetAuthContext(ctx context.Context) *AuthContext {
	if authCtx, ok := ctx.Value("auth_context").(*AuthContext); ok {
		return authCtx
	}
	return nil
}

func SetAuthContext(ctx context.Context, authCtx *AuthContext) context.Context {
	return context.WithValue(ctx, "auth_context", authCtx)
}
