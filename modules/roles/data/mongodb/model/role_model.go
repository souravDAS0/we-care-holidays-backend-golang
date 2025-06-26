package model

import (
	"time"

	"bitbucket.org/abhishek_fordel/we-care-holidays-backend-golang/modules/roles/domain/entity"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// RoleModel represents the MongoDB role schema
type RoleModel struct {
	ID          primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	Name        string             `bson:"name" json:"name"`
	Description string             `bson:"description" json:"description"`
	Permissions []string           `bson:"permissions" json:"permissions"`
	Scope       string             `json:"scope" bson:"scope"`
	IsSystem    bool               `json:"isSystem" bson:"isSystem"`
	CreatedBy   string             `json:"createdBy" bson:"createdBy"`
	CreatedAt   time.Time          `bson:"createdAt" json:"createdAt"`
	UpdatedAt   time.Time          `bson:"updatedAt" json:"updatedAt"`
	DeletedAt   *time.Time         `bson:"deletedAt,omitempty" json:"deletedAt,omitempty"`
}

// CollectionName returns the MongoDB collection name
func (RoleModel) CollectionName() string {
	return "roles"
}

// FromEntity maps entity.Role to RoleModel
func FromEntity(entity *entity.Role) *RoleModel {
	model := &RoleModel{
		ID:          entity.ID,
		Name:        entity.Name,
		Description: entity.Description,
		Permissions: entity.Permissions,
		Scope:       string(entity.Scope),
		IsSystem:    entity.IsSystem,
		CreatedBy:   entity.CreatedBy,
		CreatedAt:   entity.CreatedAt,
		UpdatedAt:   entity.UpdatedAt,
		DeletedAt:   entity.DeletedAt,
	}

	return model
}

// ToEntity maps RoleModel to entity.Role
func (m *RoleModel) ToEntity() entity.Role {
	role := entity.Role{
		ID:          m.ID,
		Name:        m.Name,
		Description: m.Description,
		Permissions: m.Permissions,
		Scope:       entity.RoleScope(m.Scope),
		IsSystem:    m.IsSystem,
		CreatedBy:   m.CreatedBy,
		CreatedAt:   m.CreatedAt,
		UpdatedAt:   m.UpdatedAt,
		DeletedAt:   m.DeletedAt,
	}

	return role
}
