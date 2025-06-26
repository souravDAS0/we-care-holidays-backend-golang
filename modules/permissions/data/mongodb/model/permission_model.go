package model

import (
	"time"

	"bitbucket.org/abhishek_fordel/we-care-holidays-backend-golang/modules/permissions/domain/entity"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// PermissionModel represents the MongoDB permission schema
type PermissionModel struct {
	ID          primitive.ObjectID `json:"_id" bson:"_id"`
	Resource    string             `json:"resource" bson:"resource"`       // e.g., "users", "roles"
	Action      string             `json:"action" bson:"action"`           // e.g., "read", "create"
	Description string             `json:"description" bson:"description"` // Human-readable description
	CreatedAt   time.Time          `json:"createdAt" bson:"createdAt"`
	UpdatedAt   time.Time          `json:"updatedAt" bson:"updatedAt"`
}

// CollectionName returns the MongoDB collection name
func (PermissionModel) CollectionName() string {
	return "permissions"
}

// FromEntity maps entity.Permission to PermissionModel
func FromEntity(entity *entity.Permission) *PermissionModel {
	model := &PermissionModel{
		ID:        entity.ID,
		Resource:  entity.Resource,
		Action:    string(entity.Action),
		Description: entity.Description,
		CreatedAt: entity.CreatedAt,
		UpdatedAt: entity.UpdatedAt,
	}

	return model
}

// ToEntity maps PermissionModel to entity.Permission
func (m *PermissionModel) ToEntity() entity.Permission {
	permission := entity.Permission{
		ID:        m.ID,
		Resource:  m.Resource,
		Action:    entity.PermissionAction(m.Action),
		Description: m.Description,
		CreatedAt: m.CreatedAt,
		UpdatedAt: m.UpdatedAt,
	}

	return permission
}
