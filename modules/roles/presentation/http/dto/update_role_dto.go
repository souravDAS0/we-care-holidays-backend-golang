// modules/roles/presentation/http/dto/update_role_dto.go
package dto

import (
	"strings"

	"bitbucket.org/abhishek_fordel/we-care-holidays-backend-golang/modules/roles/domain/entity"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// UpdateRoleDto defines the fields for updating a role
// All fields are optional to support partial updates
type UpdateRoleDto struct {
	// All fields are pointers to distinguish between nil (not provided) and empty values
	Name        *string  `json:"name,omitempty"`
	Description *string  `json:"description,omitempty"`
	Permissions []string `json:"permissions,omitempty"`
}

// Validate performs validation on the UpdateRoleDto
func (dto *UpdateRoleDto) Validate() error {
	// Name validation if provided
	if dto.Name != nil && strings.TrimSpace(*dto.Name) == "" {
		return ErrNameRequired
	}

	// Description validation if provided
	if dto.Description != nil && strings.TrimSpace(*dto.Description) == "" {
		return ErrDescriptionRequired
	}

	// Permission IDs validation if provided
	if dto.Permissions != nil {
		seenPermissions := make(map[string]bool)
		for _, permID := range dto.Permissions {
			// Check for valid ObjectID format
			if _, err := primitive.ObjectIDFromHex(permID); err != nil {
				return ErrInvalidPermissionID
			}

			// Check for duplicates
			if seenPermissions[permID] {
				return ErrDuplicatePermission
			}
			seenPermissions[permID] = true
		}
	}

	return nil
}

// ApplyUpdates applies the update DTO to an existing role entity
// Only updates fields that are provided in the DTO
func (dto *UpdateRoleDto) ApplyUpdates(role *entity.Role) {
	// Update only the fields that are provided in the DTO
	if dto.Name != nil {
		role.Name = strings.TrimSpace(*dto.Name)
	}
	if dto.Description != nil {
		role.Description = strings.TrimSpace(*dto.Description)
	}

	// // Update permissions if provided
	// if dto.Permissions != nil {
	// 	role.Permissions = make([]primitive.ObjectID, len(dto.Permissions))
	// 	for i, permID := range dto.Permissions {
	// 		// We already validated the format in Validate(), so this should not error
	// 		objID, _ := primitive.ObjectIDFromHex(permID)
	// 		role.Permissions[i] = objID
	// 	}
	// }
}

// ToUpdateEntity creates a partial entity for update operations
// This is useful when implementing patch-style updates where we only want to send
// the fields that need to be updated to the repository
func (dto *UpdateRoleDto) ToUpdateEntity(id primitive.ObjectID) *entity.Role {
	role := &entity.Role{
		ID: id,
	}
	dto.ApplyUpdates(role)
	return role
}