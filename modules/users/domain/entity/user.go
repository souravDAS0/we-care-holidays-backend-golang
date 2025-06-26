// modules/users/domain/entity/user.go
package entity

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"golang.org/x/crypto/bcrypt"
)

type UserStatus string

const (
	UserStatusInvited   UserStatus = "Invited"
	UserStatusActive    UserStatus = "Active"
	UserStatusSuspended UserStatus = "Suspended"
	UserStatusRemoved   UserStatus = "Removed"
)

// Phone represents a phone number for the user
type Phone struct {
	Number     string // Phone number (E.164)
	IsVerified bool   // Verification status
}

// Email represents an email address for the user
type Email struct {
	Email      string // Email address
	IsVerified bool   // Verification status
}

type AuditTrail struct {
	LastLoginAt     *time.Time `json:"lastLoginAt,omitempty" bson:"lastLoginAt,omitempty"`
	LastLoginIP     string     `json:"lastLoginIp" bson:"lastLoginIp"`
	LastLoginDevice string     `json:"lastLoginDevice" bson:"lastLoginDevice"`
}

type User struct {
	ID              primitive.ObjectID `json:"_id" bson:"_id"`
	FullName        string             `json:"fullName" bson:"fullName"`
	Emails          []Email            `json:"emails" bson:"emails"`
	Phones          []Phone            `json:"phones" bson:"phones"`
	Password        string             `json:"password" bson:"password"`
	Status          UserStatus         `json:"status" bson:"status"`
	ProfilePhotoURL string             `json:"profilePhotoUrl" bson:"profilePhotoUrl"`
	RoleID          string             `json:"roleId" bson:"roleId"`
	Role            string             `json:"role" bson:"role"` // For convenience
	OrganizationID  string             `json:"organizationId" bson:"organizationId, omitempty"`
	AuditTrail      AuditTrail         `json:"auditTrail" bson:"auditTrail"`
	CreatedAt       time.Time          `json:"createdAt" bson:"createdAt"`
	UpdatedAt       time.Time          `json:"updatedAt" bson:"updatedAt"`
	DeletedAt       *time.Time         `json:"deletedAt,omitempty" bson:"deletedAt,omitempty"`
}

// GetPrimaryEmail returns the first email (primary email)
func (u *User) GetPrimaryEmail() string {
	if len(u.Emails) > 0 {
		return u.Emails[0].Email
	}
	return ""
}

// GetPrimaryPhone returns the first phone (primary phone)
func (u *User) GetPrimaryPhone() string {
	if len(u.Phones) > 0 {
		return u.Phones[0].Number
	}
	return ""
}

// IsDeleted checks if the user is soft deleted
func (u *User) IsDeleted() bool {
	return u.DeletedAt != nil
}

// UpdateAuditTrail updates the audit trail with new login information
func (u *User) UpdateAuditTrail(ip, device string) {
	now := time.Now()
	u.AuditTrail.LastLoginAt = &now
	u.AuditTrail.LastLoginIP = ip
	u.AuditTrail.LastLoginDevice = device

}

func (u *User) ComparePassword(plain string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(plain))
	return err == nil
}
