package dto

import (
	"strings"

	"bitbucket.org/abhishek_fordel/we-care-holidays-backend-golang/modules/users/domain/entity"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type UpdateUserDto struct {
	FullName        *string  `json:"fullName,omitempty"`
	Emails          []string `json:"emails,omitempty"` // Array of email strings
	Phones          []string `json:"phones,omitempty"` // Array of phone strings
	Status          *string  `json:"status,omitempty"`
	ProfilePhotoURL *string  `json:"profilePhotoUrl,omitempty"`
	RoleID          *string  `json:"roleId,omitempty"`
	OrganizationID  *string  `json:"organizationId,omitempty"`
}

func (dto *UpdateUserDto) Validate() error {
	// Validate emails if provided
	if dto.Emails != nil {
		for _, email := range dto.Emails {
			if email != "" && !isValidEmail(email) {
				return ErrInvalidEmail
			}
		}
	}

	// Validate phones if provided
	if dto.Phones != nil {
		for _, phone := range dto.Phones {
			if phone != "" && !isValidPhone(phone) {
				return ErrInvalidPhone
			}
		}
	}

	// Validate role ID if provided
	if dto.RoleID != nil && *dto.RoleID != "" {
		if _, err := primitive.ObjectIDFromHex(*dto.RoleID); err != nil {
			return ErrInvalidRoleID
		}
	}

	// Validate organization ID if provided
	if dto.OrganizationID != nil && *dto.OrganizationID != "" {
		if _, err := primitive.ObjectIDFromHex(*dto.OrganizationID); err != nil {
			return ErrInvalidOrganizationID
		}
	}

	// Validate status if provided
	if dto.Status != nil && *dto.Status != "" {
		validStatuses := map[string]bool{
			"Invited":   true,
			"Active":    true,
			"Suspended": true,
			"Removed":   true,
		}
		if !validStatuses[*dto.Status] {
			return ErrInvalidStatus
		}
	}

	return nil
}

// FIXED: ApplyUpdates method now properly handles Email and Phone entities
func (dto *UpdateUserDto) ApplyUpdates(user *entity.User) {
	if dto.FullName != nil {
		user.FullName = strings.TrimSpace(*dto.FullName)
	}

	// FIXED: Convert string slice to Email entity slice
	if dto.Emails != nil {
		var emails []entity.Email
		for _, emailStr := range dto.Emails {
			if emailStr != "" {
				// Preserve existing verification status if email exists
				isVerified := false
				for _, existingEmail := range user.Emails {
					if existingEmail.Email == emailStr {
						isVerified = existingEmail.IsVerified
						break
					}
				}
				emails = append(emails, entity.Email{
					Email:      strings.TrimSpace(emailStr),
					IsVerified: isVerified,
				})
			}
		}
		user.Emails = emails
	}

	// FIXED: Convert string slice to Phone entity slice
	if dto.Phones != nil {
		var phones []entity.Phone
		for _, phoneStr := range dto.Phones {
			if phoneStr != "" {
				// Preserve existing verification status if phone exists
				isVerified := false
				for _, existingPhone := range user.Phones {
					if existingPhone.Number == phoneStr {
						isVerified = existingPhone.IsVerified
						break
					}
				}
				phones = append(phones, entity.Phone{
					Number:     strings.TrimSpace(phoneStr),
					IsVerified: isVerified,
				})
			}
		}
		user.Phones = phones
	}

	if dto.Status != nil {
		user.Status = entity.UserStatus(*dto.Status)
	}

	if dto.ProfilePhotoURL != nil {
		user.ProfilePhotoURL = *dto.ProfilePhotoURL
	}

	if dto.RoleID != nil && *dto.RoleID != "" {
		user.RoleID = *dto.RoleID
	}

	if dto.OrganizationID != nil {
		user.OrganizationID = *dto.OrganizationID
	}
}
