// modules/users/data/mongodb/model/user_model.go
package model

import (
	"time"

	"bitbucket.org/abhishek_fordel/we-care-holidays-backend-golang/modules/users/domain/entity"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// CollectionName returns the MongoDB collection name
func (UserModel) CollectionName() string {
	return "users"
}

// UserModel represents the MongoDB document structure for users
type UserModel struct {
	ID              primitive.ObjectID `bson:"_id,omitempty"`
	FullName        string             `bson:"fullName"`
	Emails          []EmailInfo        `bson:"emails"`
	Phones          []PhoneInfo        `bson:"phones"`
	Password        string             `bson:"password"`
	Status          string             `bson:"status"`
	ProfilePhotoURL string             `bson:"profilePhotoUrl"`
	Role            string             `bson:"role"`
	RoleID          primitive.ObjectID `bson:"roleId"`
	OrganizationID  primitive.ObjectID `bson:"organizationId,omitempty"`
	AuditTrail      AuditTrailModel    `bson:"auditTrail"`
	CreatedAt       time.Time          `bson:"createdAt"`
	UpdatedAt       time.Time          `bson:"updatedAt"`
	DeletedAt       *time.Time         `bson:"deletedAt,omitempty"`
}

// PhoneInfo stores a phone number and its verification/OTP info
type PhoneInfo struct {
	Number     string `bson:"number" json:"number" validate:"required,e164"`
	IsVerified bool   `bson:"isVerified" json:"isVerified"`
}

// EmailInfo stores an email address and its verification/OTP info
type EmailInfo struct {
	Email      string `bson:"email" json:"email" validate:"required,email"`
	IsVerified bool   `bson:"isVerified" json:"isVerified"`
}

type AuditTrailModel struct {
	LastLoginAt     *time.Time `bson:"lastLoginAt,omitempty"`
	LastLoginIP     string     `bson:"lastLoginIp"`
	LastLoginDevice string     `bson:"lastLoginDevice"`
}

// ToEntity converts UserModel to domain entity
func (m *UserModel) ToEntity() entity.User {
	phones := make([]entity.Phone, len(m.Phones))
	for i, p := range m.Phones {
		phones[i] = entity.Phone{
			Number:     p.Number,
			IsVerified: p.IsVerified,
		}
	}

	emails := make([]entity.Email, len(m.Emails))
	for i, e := range m.Emails {
		emails[i] = entity.Email{
			Email:      e.Email,
			IsVerified: e.IsVerified,
		}
	}

	return entity.User{
		ID:              m.ID,
		FullName:        m.FullName,
		Emails:          emails,
		Phones:          phones,
		Status:          entity.UserStatus(m.Status),
		Password:        m.Password,
		ProfilePhotoURL: m.ProfilePhotoURL,
		Role:            m.Role,
		RoleID:          m.RoleID.Hex(),
		OrganizationID:  m.OrganizationID.Hex(),
		AuditTrail:      m.toEntityAuditTrail(),
		CreatedAt:       m.CreatedAt,
		UpdatedAt:       m.UpdatedAt,
		DeletedAt:       m.DeletedAt,
	}
}

// FromEntity converts domain entity to UserModel
func FromEntity(user *entity.User) *UserModel {
	roleID, _ := primitive.ObjectIDFromHex(user.RoleID)
	organizationId, _ := primitive.ObjectIDFromHex(user.OrganizationID)

	phones := make([]PhoneInfo, len(user.Phones))
	for i, p := range user.Phones {
		phones[i] = PhoneInfo{
			Number:     p.Number,
			IsVerified: p.IsVerified,
		}
	}

	emails := make([]EmailInfo, len(user.Emails))
	for i, e := range user.Emails {
		emails[i] = EmailInfo{
			Email:      e.Email,
			IsVerified: e.IsVerified,
		}
	}

	return &UserModel{
		ID:              user.ID,
		FullName:        user.FullName,
		Emails:          emails,
		Phones:          phones,
		Password:        user.Password,
		Status:          string(user.Status),
		ProfilePhotoURL: user.ProfilePhotoURL,
		Role:            user.Role,
		RoleID:          roleID,
		OrganizationID:  organizationId,
		AuditTrail:      fromEntityAuditTrail(user.AuditTrail),
		CreatedAt:       user.CreatedAt,
		UpdatedAt:       user.UpdatedAt,
		DeletedAt:       user.DeletedAt,
	}
}

func (m *UserModel) toEntityAuditTrail() entity.AuditTrail {

	return entity.AuditTrail{
		LastLoginAt:     m.AuditTrail.LastLoginAt,
		LastLoginIP:     m.AuditTrail.LastLoginIP,
		LastLoginDevice: m.AuditTrail.LastLoginDevice,
	}
}

func fromEntityAuditTrail(auditTrail entity.AuditTrail) AuditTrailModel {

	return AuditTrailModel{
		LastLoginAt:     auditTrail.LastLoginAt,
		LastLoginIP:     auditTrail.LastLoginIP,
		LastLoginDevice: auditTrail.LastLoginDevice,
	}
}
