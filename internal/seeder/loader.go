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
	Roles       []RoleSeed       `json:"roles"`
	Permissions []PermissionSeed `json:"permissions"`
	Users       []UserSeed       `json:"users"`
}

// RoleSeed represents a role record in seed.json
type RoleSeed struct {
	Name        string   `json:"name"`
	Description string   `json:"description"`
	Permissions []string `json:"permissions,omitempty"`
	Scope       string   `json:"scope"`
	// IsSystem    bool     `json:"isSystem"`
}

// PermissionSeed represents a permission record in seed.json
type PermissionSeed struct {
	Resource    string `json:"resource"`
	Action      string `json:"action"`
	Description string `json:"description"`
}

type UserSeed struct {
	FullName        string      `json:"fullName"` // FIXED: was "name" in JSON
	Phones          []PhoneSeed `json:"phones"`
	Emails          []EmailSeed `json:"emails"`
	ProfilePhotoURL string      `json:"profilePhotoUrl"`          // FIXED: JSON tag
	Role            string      `json:"role"`                     // FIXED: JSON tag
	Status          string      `json:"status"`                   // FIXED: JSON tag
	Password        string      `json:"password"`                 // Added for seed data
	OrganizationID  string      `json:"organizationId,omitempty"` // FIXED: JSON tag
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
