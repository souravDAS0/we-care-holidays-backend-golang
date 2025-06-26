// modules/permissions/presentation/http/dto/update_permission_dto.go
package dto

import (
	"strings"

	"bitbucket.org/abhishek_fordel/we-care-holidays-backend-golang/modules/permissions/domain/entity"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// UpdatePermissionDto defines the fields for updating a permission
// All fields are optional to support partial updates
type UpdatePermissionDto struct {
	Resource *string `json:"resource,omitempty"`
	Action   *string `json:"action,omitempty"`
}

// Validate performs validation on the UpdatePermissionDto
func (dto *UpdatePermissionDto) Validate() error {
	// Resource validation if provided
	if dto.Resource != nil && strings.TrimSpace(*dto.Resource) == "" {
		return ErrResourceRequired
	}

	// Action validation if provided
	if dto.Action != nil {
		if strings.TrimSpace(*dto.Action) == "" {
			return ErrActionRequired
		}
		validActions := map[string]bool{
			"read":   true,
			"write":  true,
			"manage": true,
			"delete": true,
		}
		if !validActions[*dto.Action] {
			return ErrInvalidAction
		}
	}

	return nil
}

// ApplyUpdates applies the update DTO to an existing permission entity
// Only updates fields that are provided in the DTO
func (dto *UpdatePermissionDto) ApplyUpdates(permission *entity.Permission) {
	// Update only the fields that are provided in the DTO
	if dto.Resource != nil {
		permission.Resource = strings.TrimSpace(*dto.Resource)
	}
	if dto.Action != nil {
		permission.Action = entity.PermissionAction(*dto.Action)
	}

}

// ToUpdateEntity creates a partial entity for update operations
// This is useful when implementing patch-style updates where we only want to send
// the fields that need to be updated to the repository
func (dto *UpdatePermissionDto) ToUpdateEntity(id primitive.ObjectID) *entity.Permission {
	permission := &entity.Permission{
		ID: id,
	}
	dto.ApplyUpdates(permission)
	return permission
}
