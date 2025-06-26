package entity

import (
	"errors"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Validation errors
var (
	ErrNameRequired      = errors.New("organization name is required")
	ErrInvalidEmail      = errors.New("invalid email address")
	ErrInvalidType       = errors.New("invalid organization type")
	ErrInvalidStatus     = errors.New("invalid organization status")
	ErrAddressIncomplete = errors.New("address is incomplete")
)

// Organization represents a Platform, Supplier, or Travel agent
// @Description Organization entity representing clients, suppliers, or travel agents
type Organization struct {
	// Unique identifier for the organization
	ID            primitive.ObjectID `json:"id" bson:"_id,omitempty" example:"5f8d0c1b7ea3f0d0f3c8e1b9"`
	// Organization name
	Name          string    `json:"name" bson:"name" example:"WeCare Holidays"`
	// URL-friendly version of name
	Slug          string    `json:"slug" bson:"slug" example:"wecare-holidays"`
	// Organization type (SUPPLIER, TRAVEL_AGENT, PLATFORM)
	Type          string    `json:"type" bson:"type" example:"SUPPLIER"` 
	// Primary contact email
	Email         string    `json:"email" bson:"email" example:"contact@wecareholidays.com"`
	// Contact phone number
	Phone         string    `json:"phone" bson:"phone" example:"+91 1234567890"`
	// Organization website URL
	Website       string    `json:"website" bson:"website" example:"https://www.wecareholidays.com"`
	// Tax identification numbers
	TaxIDs        []string  `json:"taxIds" bson:"taxIds" example:"['GST123456', 'PAN1234567']"`
	// URL to organization logo
	Logo          string    `json:"logo" bson:"logo" example:"https://storage.example.com/logos/wecare.png"`
	// Physical address
	Address       Address   `json:"address" bson:"address"`
	// Current status (Pending, Approved, Suspended, Archived)
	Status        string    `json:"status" bson:"status" example:"Approved"` 
	// Creation timestamp
	CreatedAt     time.Time `json:"createdAt" bson:"createdAt"`
	// Last update timestamp
	UpdatedAt     time.Time `json:"updatedAt" bson:"updatedAt"`
	// Soft delete timestamp
	DeletedAt     *time.Time `json:"deletedAt,omitempty" bson:"deletedAt,omitempty"`
}

// Address represents a physical address
// @Description Physical address information
type Address struct {
	// Street address including building/apartment
	Street  string `json:"street" bson:"street" example:"123 Main Street, Building A"`
	// City name
	City    string `json:"city" bson:"city" example:"Siliguri"`
	// State or province
	State   string `json:"state" bson:"state" example:"West Bengal"`
	// Country name
	Country string `json:"country" bson:"country" example:"India"`
	// Postal/ZIP code
	Pincode string `json:"pincode" bson:"pincode" example:"734001"`
}

// IsDeleted returns true if the organization is soft-deleted
func (o *Organization) IsDeleted() bool {
	return o.DeletedAt != nil
}