// modules/permissions/domain/entity/permission.go
package entity

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type PermissionAction string

const (
	PermissionActionRead   PermissionAction = "read"
	PermissionActionWrite  PermissionAction = "write"
	PermissionActionCreate PermissionAction = "create"
	PermissionActionUpdate PermissionAction = "update"
	PermissionActionDelete PermissionAction = "delete"
)

type Permission struct {
	ID          primitive.ObjectID `json:"_id" bson:"_id"`
	Resource    string             `json:"resource" bson:"resource"`       // e.g., "users", "roles"
	Action      PermissionAction   `json:"action" bson:"action"`           // e.g., "read", "create"
	Description string             `json:"description" bson:"description"` // Human-readable description
	CreatedAt   time.Time          `json:"createdAt" bson:"createdAt"`
	UpdatedAt   time.Time          `json:"updatedAt" bson:"updatedAt"`
}

func (p *Permission) String() string {
	return string(p.Resource) + ":" + string(p.Action)
}
