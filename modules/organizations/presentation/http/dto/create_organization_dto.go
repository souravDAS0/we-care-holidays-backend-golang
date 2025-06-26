package dto

import (
	"net/mail"
	"strings"

	"bitbucket.org/abhishek_fordel/we-care-holidays-backend-golang/modules/organizations/domain/entity"
)

// Validation errors

// CreateOrganizationDto defines the required and optional fields for creating an organization
type CreateOrganizationDto struct {
	// Required fields
	Name  string `json:"name" binding:"required" example:"WeCare Holidays"`
	Email string `json:"email" binding:"required" example:"contact@wecareholidays.com"`
	Type  string `json:"type" binding:"required" example:"SUPPLIER"` // SUPPLIER, TRAVEL_AGENT, PLATFORM

	// Optional fields
	Slug          string     `json:"slug,omitempty" example:"wecare-holidays"`
	Phone         string     `json:"phone,omitempty" example:"+91 1234567890"`
	Website       string     `json:"website,omitempty" example:"https://www.wecareholidays.com"`
	TaxIDs        []string   `json:"taxIds,omitempty" example:"['GST123456', 'PAN1234567']"`
	Logo          string     `json:"logo,omitempty" example:"https://storage.example.com/logos/wecare.png"`
	Address       AddressDto `json:"address,omitempty"`
	Status        string     `json:"status,omitempty" example:"Approved"` // If not provided, defaults to "Pending"
}

// Validate performs validation on the CreateOrganizationDto
func (dto *CreateOrganizationDto) Validate() error {
	// Required fields validation
	if strings.TrimSpace(dto.Name) == "" {
		return ErrNameRequired
	}

	if strings.TrimSpace(dto.Email) == "" {
		return ErrEmailRequired
	}

	// Email format validation
	if _, err := mail.ParseAddress(dto.Email); err != nil {
		return ErrInvalidEmail
	}

	// Type validation
	validTypes := map[string]bool{
		"SUPPLIER":    true,
		"TRAVEL_AGENT":    true,
		"PLATFORM": true,
	}
	if !validTypes[dto.Type] {
		return ErrInvalidType
	}

	// Status validation (if provided)
	if dto.Status != "" {
		validStatuses := map[string]bool{
			"Pending":   true,
			"Approved":  true,
			"Suspended": true,
			"Archived":  true,
		}
		if !validStatuses[dto.Status] {
			return ErrInvalidStatus
		}
	}

	// Address validation - if any address field is provided, required fields must be present
	if dto.Address.Street != nil || dto.Address.City != nil || dto.Address.State != nil ||
		dto.Address.Country != nil || dto.Address.Pincode != nil {
		if dto.Address.City == nil || dto.Address.Country == nil {
			return ErrAddressIncomplete
		}
	}

	return nil
}

// ToEntity converts the DTO to an entity
func (dto *CreateOrganizationDto) ToEntity() *entity.Organization {
	org := &entity.Organization{
		Name:  dto.Name,
		Email: dto.Email,
		Type:  dto.Type,
		Slug:  dto.Slug,
	}


	org.Phone = dto.Phone
	org.Website = dto.Website
	org.TaxIDs = dto.TaxIDs
	org.Logo = dto.Logo

	// Convert Address
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

	// Set status if provided, otherwise it will be set to default in the use case
	if dto.Status != "" {
		org.Status = dto.Status
	}

	return org
}