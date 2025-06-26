package dto

import "errors"

var (
	ErrNameRequired      = errors.New("organization name is required")
	ErrEmailRequired     = errors.New("email is required")
	ErrTypeRequired      = errors.New("organization type is required")
	ErrInvalidEmail      = errors.New("invalid email address")
	ErrInvalidType       = errors.New("invalid organization type (must be SUPPLIER, TRAVEL_AGENT, PLATFORM)")
	ErrInvalidStatus     = errors.New("invalid organization status: must be one of Pending, Approved, Suspended, or Archived")
	ErrAddressIncomplete = errors.New("address is incomplete (city and country are required)")
)

// AddressDto represents a physical address for DTOs
type AddressDto struct {
	Street  *string `json:"street,omitempty" example:"123 Main Street, Building A"`
	City    *string `json:"city,omitempty" example:"Siliguri"`
	State   *string `json:"state,omitempty" example:"West Bengal"`
	Country *string `json:"country,omitempty" example:"India"`
	Pincode *string `json:"pincode,omitempty" example:"734001"`
}


