// modules/permissions/presentation/http/dto/create_permission_dto.go
package dto

import (
	"errors"
	"strings"

	"bitbucket.org/abhishek_fordel/we-care-holidays-backend-golang/modules/permissions/domain/entity"
)

// Validation errors
var (
	ErrResourceRequired = errors.New("resource is required")
	ErrActionRequired   = errors.New("action is required")
	ErrScopeRequired    = errors.New("scope is required")
	ErrInvalidAction    = errors.New("invalid action (must be read, write, manage, delete)")
	ErrInvalidScope     = errors.New("invalid scope (must be global, organization, self, custom)")
	ErrInvalidPriority  = errors.New("priority must be between 0 and 1000")
)

// CreatePermissionDto defines the required and optional fields for creating a permission
type CreatePermissionDto struct {
	Resource string `json:"resource" binding:"required" example:"organizations"`
	Action   string `json:"action" binding:"required" example:"read"`
}

// Validate performs validation on the CreatePermissionDto
func (dto *CreatePermissionDto) Validate() error {
	// Required fields validation
	if strings.TrimSpace(dto.Resource) == "" {
		return ErrResourceRequired
	}

	if strings.TrimSpace(dto.Action) == "" {
		return ErrActionRequired
	}

	// Action validation
	validActions := map[string]bool{
		"read":   true,
		"write":  true,
		"update": true,
		"delete": true,
	}
	if !validActions[dto.Action] {
		return ErrInvalidAction
	}

	return nil
}

// ToEntity converts the DTO to an entity
func (dto *CreatePermissionDto) ToEntity() *entity.Permission {
	permission := &entity.Permission{
		Resource: strings.TrimSpace(dto.Resource),
		Action:   entity.PermissionAction(dto.Action),
	}

	return permission
}
