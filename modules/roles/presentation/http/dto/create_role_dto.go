// modules/roles/presentation/http/dto/create_role_dto.go
package dto

import (
	"errors"
	"strings"

	"bitbucket.org/abhishek_fordel/we-care-holidays-backend-golang/modules/roles/domain/entity"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Validation errors
var (
	ErrNameRequired        = errors.New("name is required")
	ErrDescriptionRequired = errors.New("description is required")
	ErrInvalidPermissionID = errors.New("invalid permission ID format")
	ErrDuplicatePermission = errors.New("duplicate permission ID found")
)

// CreateRoleDto defines the required and optional fields for creating a role
type CreateRoleDto struct {
	// Required fields
	Name        string `json:"name" binding:"required" example:"Admin"`
	Description string `json:"description" binding:"required" example:"Administrator role with full access"`
	// Optional fields
	Scope       string   `json:"scope,omitempty" example:"organizations"`
	Permissions []string `json:"permissions,omitempty" example:"507f1f77bcf86cd799439011, 507f1f77bcf86cd799439012"`
}

// Validate performs validation on the CreateRoleDto
func (dto *CreateRoleDto) Validate() error {
	// Required fields validation
	if strings.TrimSpace(dto.Name) == "" {
		return ErrNameRequired
	}

	if strings.TrimSpace(dto.Description) == "" {
		return ErrDescriptionRequired
	}

	// Permission IDs validation
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

// ToEntity converts the DTO to an entity
func (dto *CreateRoleDto) ToEntity() *entity.Role {
	role := &entity.Role{
		Name:        strings.TrimSpace(dto.Name),
		Description: strings.TrimSpace(dto.Description),
		Permissions: dto.Permissions,
		IsSystem:    false,
	}

	if dto.Scope != "" {
		role.Scope = entity.RoleScope(dto.Scope)
	} else {
		role.Scope = "self"
	}

	return role
}
