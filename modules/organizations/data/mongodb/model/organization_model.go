// internal/modules/organizations/data/mongodb/model/organization_model.go
package model

import (
	"time"

	"bitbucket.org/abhishek_fordel/we-care-holidays-backend-golang/modules/organizations/domain/entity"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// AddressModel represents embedded address in organization model
type AddressModel struct {
	Street  string `bson:"street" json:"street"`
	City    string `bson:"city" json:"city"`
	State   string `bson:"state" json:"state"`
	Country string `bson:"country" json:"country"`
	Pincode string `bson:"pincode" json:"pincode"`
}

// OrganizationModel represents the MongoDB organization schema
type OrganizationModel struct {
	ID            primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	Name          string             `bson:"name" json:"name"`
	Slug          string             `bson:"slug" json:"slug"`
	Type          string             `bson:"type" json:"type"`
	Email         string             `bson:"email" json:"email"`
	Phone         string             `bson:"phone" json:"phone"`
	Website       string             `bson:"website" json:"website"`
	TaxIDs        []string           `bson:"taxIds" json:"taxIds"`
	Logo          string             `bson:"logo" json:"logo"`
	Address       AddressModel       `bson:"address" json:"address"`
	Status        string             `bson:"status" json:"status"`
	CreatedAt     time.Time          `bson:"createdAt" json:"createdAt"`
	UpdatedAt     time.Time          `bson:"updatedAt" json:"updatedAt"`
	DeletedAt     *time.Time         `bson:"deletedAt,omitempty" json:"deletedAt,omitempty"`
}

// CollectionName returns the MongoDB collection name
func (OrganizationModel) CollectionName() string {
	return "organizations"
}

// FromEntity maps entity.Organization to OrganizationModel
func FromEntity(entity *entity.Organization) *OrganizationModel {
	return &OrganizationModel{
		ID:            entity.ID,
		Name:          entity.Name,
		Slug:          entity.Slug,
		Type:          entity.Type,
		Email:         entity.Email,
		Phone:         entity.Phone,
		Website:       entity.Website,
		TaxIDs:        entity.TaxIDs,
		Logo:          entity.Logo,
		Address: AddressModel{
			Street:  entity.Address.Street,
			City:    entity.Address.City,
			State:   entity.Address.State,
			Country: entity.Address.Country,
			Pincode: entity.Address.Pincode,
		},
		Status:    entity.Status,
		CreatedAt: entity.CreatedAt,
		UpdatedAt: entity.UpdatedAt,
		DeletedAt: entity.DeletedAt,
	}
}

// ToEntity maps OrganizationModel to entity.Organization
func (m *OrganizationModel) ToEntity() entity.Organization {
	return entity.Organization{
		ID:            m.ID,
		Name:          m.Name,
		Slug:          m.Slug,
		Type:          m.Type,
		Email:         m.Email,
		Phone:         m.Phone,
		Website:       m.Website,
		TaxIDs:        m.TaxIDs,
		Logo:          m.Logo,
		Address: entity.Address{
			Street:  m.Address.Street,
			City:    m.Address.City,
			State:   m.Address.State,
			Country: m.Address.Country,
			Pincode: m.Address.Pincode,
		},
		Status:    m.Status,
		CreatedAt: m.CreatedAt,
		UpdatedAt: m.UpdatedAt,
		DeletedAt: m.DeletedAt,
	}
}
