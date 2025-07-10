// loader.go
//
// Seed data loader. Reads and parses seed.json file.

package seeder

import (
	"encoding/json"
	"os"
)

// SeedData represents the entire seed.json structure
type SeedData struct {
	Permissions   []PermissionSeed   `json:"permissions"`
	Organizations []OrganizationSeed `json:"organizations"`
	Roles         []RoleSeed         `json:"roles"`
	Users         []UserSeed         `json:"users"`
}

// OrganizationSeed represents an organization record in seed.json
type OrganizationSeed struct {
	Name    string      `json:"name"`
	Slug    string      `json:"slug"`
	Type    string      `json:"type"`
	Email   string      `json:"email"`
	Phone   string      `json:"phone"`
	Website string      `json:"website"`
	TaxIDs  []string    `json:"taxIds"`
	Logo    string      `json:"logo"`
	Address AddressSeed `json:"address"`
	Status  string      `json:"status"`
}

// AddressSeed represents address data in seed.json
type AddressSeed struct {
	Street  string `json:"street"`
	City    string `json:"city"`
	State   string `json:"state"`
	Country string `json:"country"`
	Pincode string `json:"pincode"`
}

// RoleSeed represents a role record in seed.json
type RoleSeed struct {
	Name        string   `json:"name"`
	Description string   `json:"description"`
	Permissions []string `json:"permissions,omitempty"`
	Scope       string   `json:"scope"`
}

// PermissionSeed represents a permission record in seed.json
type PermissionSeed struct {
	Resource    string `json:"resource"`
	Action      string `json:"action"`
	Description string `json:"description"`
}

type UserSeed struct {
	FullName         string      `json:"fullName"`
	Phones           []PhoneSeed `json:"phones"`
	Emails           []EmailSeed `json:"emails"`
	ProfilePhotoURL  string      `json:"profilePhotoUrl"`
	Role             string      `json:"role"`
	Status           string      `json:"status"`
	Password         string      `json:"password"`
	OrganizationSlug string      `json:"organizationSlug"` // Changed from ID to slug for easier reference
}

// PhoneSeed represents user phone
type PhoneSeed struct {
	Number     string `json:"number"`
	IsVerified bool   `json:"isVerified"`
}

// EmailSeed represents user email
type EmailSeed struct {
	Email      string `json:"email"`
	IsVerified bool   `json:"isVerified"`
}

// LoadSeedData loads the seed data from the JSON file.
//
// Returns:
// - SeedData: parsed seed data.
// - error: if parsing fails.
func LoadSeedData() (*SeedData, error) {
	file, err := os.ReadFile("internal/seeder/data/seed.json")
	if err != nil {
		return nil, err
	}

	var data SeedData
	if err := json.Unmarshal(file, &data); err != nil {
		return nil, err
	}

	return &data, nil
}
