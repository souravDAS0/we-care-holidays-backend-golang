package dto

import (
	"errors"
	"strings"

	"bitbucket.org/abhishek_fordel/we-care-holidays-backend-golang/modules/users/domain/entity"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type CreateUserDto struct {
	FullName        string `json:"fullName" binding:"required" example:"John Doe"`
	Email           string `json:"email" binding:"required,email" example:"john@example.com"`
	Phone           string `json:"phone,omitempty" example:"+919876543210"`
	Status          string `json:"status,omitempty" example:"Active"`
	ProfilePhotoURL string `json:"profilePhotoUrl,omitempty" example:"https://example.com/photo.jpg"`
	RoleID          string `json:"roleId" binding:"required" example:"507f1f77bcf86cd799439011"`
	OrganizationID  string `json:"organizationId,omitempty" example:"507f1f77bcf86cd799439012"`
}

// Validation errors
var (
	ErrFullNameRequired      = errors.New("full name is required")
	ErrEmailRequired         = errors.New("email is required")
	ErrInvalidEmail          = errors.New("invalid email format")
	ErrInvalidPhone          = errors.New("invalid phone format")
	ErrRoleIDRequired        = errors.New("role ID is required")
	ErrInvalidRoleID         = errors.New("invalid role ID format")
	ErrInvalidOrganizationID = errors.New("invalid organization ID format")
	ErrInvalidStatus         = errors.New("invalid status")
)

func (dto *CreateUserDto) Validate() error {
	// Required fields validation
	if strings.TrimSpace(dto.FullName) == "" {
		return ErrFullNameRequired
	}

	// Email is required and must be valid
	if strings.TrimSpace(dto.Email) == "" {
		return ErrEmailRequired
	}
	if !isValidEmail(dto.Email) {
		return ErrInvalidEmail
	}

	// Phone validation - only if provided
	if dto.Phone != "" && !isValidPhone(dto.Phone) {
		return ErrInvalidPhone
	}

	// Role ID validation
	if strings.TrimSpace(dto.RoleID) == "" {
		return ErrRoleIDRequired
	}
	if _, err := primitive.ObjectIDFromHex(dto.RoleID); err != nil {
		return ErrInvalidRoleID
	}

	// Organization ID validation - only if provided
	if dto.OrganizationID != "" {
		if _, err := primitive.ObjectIDFromHex(dto.OrganizationID); err != nil {
			return ErrInvalidOrganizationID
		}
	}

	// Status validation - only if provided
	if dto.Status != "" {
		validStatuses := map[string]bool{
			"Invited":   true,
			"Active":    true,
			"Suspended": true,
			"Removed":   true,
		}
		if !validStatuses[dto.Status] {
			return ErrInvalidStatus
		}
	}

	return nil
}

// FIXED: ToEntity method now properly creates Email and Phone slices
func (dto *CreateUserDto) ToEntity() *entity.User {
	// Create emails slice - not just a string
	var emails []entity.Email
	if dto.Email != "" {
		emails = append(emails, entity.Email{
			Email:      strings.TrimSpace(dto.Email),
			IsVerified: false, // New emails start as unverified
		})
	}

	// Create phones slice - not just a string
	var phones []entity.Phone
	if dto.Phone != "" {
		phones = append(phones, entity.Phone{
			Number:     strings.TrimSpace(dto.Phone),
			IsVerified: false, // New phones start as unverified
		})
	}

	// Set default status
	status := entity.UserStatusInvited
	if dto.Status != "" {
		status = entity.UserStatus(dto.Status)
	}

	return &entity.User{
		ID:              primitive.NewObjectID(), // Generate new ID
		FullName:        strings.TrimSpace(dto.FullName),
		Emails:          emails, // FIXED: Now properly creates Email slice
		Phones:          phones, // FIXED: Now properly creates Phone slice
		Status:          status,
		ProfilePhotoURL: dto.ProfilePhotoURL,
		RoleID:          dto.RoleID,          // Keep as string as per entity
		OrganizationID:  dto.OrganizationID,  // Keep as string as per entity
		AuditTrail:      entity.AuditTrail{}, // Initialize empty audit trail
	}
}
