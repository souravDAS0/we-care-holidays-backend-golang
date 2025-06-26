package dto

import (
	"net/mail"

	"bitbucket.org/abhishek_fordel/we-care-holidays-backend-golang/modules/organizations/domain/entity"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// UpdateOrganizationDto defines the fields for updating an organization
// All fields are optional to support partial updates
type UpdateOrganizationDto struct {
	// All fields are pointers to distinguish between nil (not provided) and empty values
	Name          *string    `json:"name,omitempty"`
	Slug          *string    `json:"slug,omitempty"`
	Type          *string    `json:"type,omitempty"`
	Email         *string    `json:"email,omitempty"`
	Phone         *string    `json:"phone,omitempty"`
	Website       *string    `json:"website,omitempty"`
	TaxIDs        []string   `json:"taxIds,omitempty"`
	Logo          *string    `json:"logo,omitempty"`
	Address       AddressDto `json:"address,omitempty"`
	Status        *string    `json:"status,omitempty"`
}

// Validate performs validation on the UpdateOrganizationDto
func (dto *UpdateOrganizationDto) Validate() error {
	// Email validation if provided
	if dto.Email != nil && *dto.Email != "" {
		if _, err := mail.ParseAddress(*dto.Email); err != nil {
			return ErrInvalidEmail
		}
	}

	// Type validation if provided
	if dto.Type != nil && *dto.Type != "" {
		validTypes := map[string]bool{
			"SUPPLIER":    true,
			"TRAVEL_AGENT":    true,
			"PLATFORM": true,
		}
		if !validTypes[*dto.Type] {
			return ErrInvalidType
		}
	}

	// Status validation if provided
	if dto.Status != nil && *dto.Status != "" {
		validStatuses := map[string]bool{
			"Pending":   true,
			"Approved":  true,
			"Suspended": true,
			"Archived":  true,
		}
		if !validStatuses[*dto.Status] {
			return ErrInvalidStatus
		}
	}

	// Address validation - if address is being updated and has fields, check for required fields
	hasAddressUpdate := dto.Address.Street != nil || dto.Address.City != nil ||
		dto.Address.State != nil || dto.Address.Country != nil || dto.Address.Pincode != nil

	if hasAddressUpdate {
		// For updates, we only validate if both city and country are provided together
		hasCity := dto.Address.City != nil
		hasCountry := dto.Address.Country != nil

		// If one is provided but not the other, require both
		if (hasCity && !hasCountry) || (!hasCity && hasCountry) {
			return ErrAddressIncomplete
		}
	}

	return nil
}

// ApplyUpdates applies the update DTO to an existing organization entity
// Only updates fields that are provided in the DTO
func (dto *UpdateOrganizationDto) ApplyUpdates(org *entity.Organization) {
	// Update only the fields that are provided in the DTO
	if dto.Name != nil {
		org.Name = *dto.Name
	}
	if dto.Slug != nil {
		org.Slug = *dto.Slug
	}
	if dto.Type != nil {
		org.Type = *dto.Type
	}
	if dto.Email != nil {
		org.Email = *dto.Email
	}
	if dto.Phone != nil {
		org.Phone = *dto.Phone
	}
	if dto.Website != nil {
		org.Website = *dto.Website
	}
	if dto.TaxIDs != nil {
		org.TaxIDs = dto.TaxIDs
	}
	if dto.Logo != nil {
		org.Logo = *dto.Logo
	}
	if dto.Status != nil {
		org.Status = *dto.Status
	}

	// Update address fields if provided
	if dto.Address.Street != nil {
		org.Address.Street = *dto.Address.Street
	}
	if dto.Address.City != nil {
		org.Address.City = *dto.Address.City
	}
	if dto.Address.State != nil {
		org.Address.State = *dto.Address.State
	}
	if dto.Address.Country != nil {
		org.Address.Country = *dto.Address.Country
	}
	if dto.Address.Pincode != nil {
		org.Address.Pincode = *dto.Address.Pincode
	}
}

// ToUpdateEntity creates a partial entity for update operations
// This is useful when implementing patch-style updates where we only want to send
// the fields that need to be updated to the repository
func (dto *UpdateOrganizationDto) ToUpdateEntity(id primitive.ObjectID) *entity.Organization {
	org := &entity.Organization{
		ID: id,
	}
	dto.ApplyUpdates(org)
	return org
}