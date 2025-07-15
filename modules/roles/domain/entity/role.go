// modules/roles/domain/entity/role.go
package entity

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type RoleScope string

const (
	RoleScopeGlobal       RoleScope = "global"       // Can access all data across all organizations
	RoleScopeOrganization RoleScope = "organization" // Can access data within their organization only
	RoleScopeSelf         RoleScope = "self"         // Can access only their own data
)

type Role struct {
	ID             primitive.ObjectID  `json:"_id" bson:"_id"`
	Name           string              `json:"name" bson:"name"`
	Description    string              `json:"description" bson:"description"`
	Permissions    []string            `json:"permissions" bson:"permissions"` // Permission IDs
	Scope          RoleScope           `json:"scope" bson:"scope"`
	OrganizationID *primitive.ObjectID `json:"organizationId,omitempty" bson:"organizationId,omitempty"` // nil for system roles
	IsSystem       bool                `json:"isSystem" bson:"isSystem"`
	CreatedBy      string              `json:"createdBy" bson:"createdBy"` // User ID who created this role
	CreatedAt      time.Time           `json:"createdAt" bson:"createdAt"`
	UpdatedAt      time.Time           `json:"updatedAt" bson:"updatedAt"`
	DeletedAt      *time.Time          `json:"deletedAt,omitempty" bson:"deletedAt,omitempty"`
}

// IsDeleted checks if the role is soft deleted
func (r *Role) IsDeleted() bool {
	return r.DeletedAt != nil
}
