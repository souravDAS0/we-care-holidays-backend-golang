package model

import (
	"time"

	"bitbucket.org/abhishek_fordel/we-care-holidays-backend-golang/modules/locations/domain/entity"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// LocationModel is the Mongo schema
type LocationModel struct {
	ID          primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	Name        string             `bson:"name" json:"name"`
	Type        string             `bson:"type" json:"type"`
	Country     string             `bson:"country" json:"country"`
	State       string             `bson:"state" json:"state"`
	District    string             `bson:"district" json:"district"`
	Pincode     string             `bson:"pincode" json:"pincode"`
	Coordinates struct {
		Lat float64 `bson:"lat" json:"lat"`
		Lng float64 `bson:"lng" json:"lng"`
	} `bson:"coordinates" json:"coordinates"`
	GeoJSON     interface{}        `bson:"geojson,omitempty" json:"geojson,omitempty"`
	Tags        []string           `bson:"tags" json:"tags"`
	Description string             `bson:"description" json:"description"`
	Aliases     []string           `bson:"aliases" json:"aliases"`
	MediaURLs   struct {
		Photos []string `bson:"photos" json:"photos"`
		Videos []string `bson:"videos" json:"videos"`
	} `bson:"mediaUrls" json:"mediaUrls"`
	CreatedBy   primitive.ObjectID `bson:"createdBy" json:"createdBy"`
	CreatedAt   time.Time          `bson:"createdAt" json:"createdAt"`
	UpdatedAt   time.Time          `bson:"updatedAt" json:"updatedAt"`
	DeletedAt   *time.Time         `bson:"deletedAt,omitempty" json:"deletedAt,omitempty"`
}

func (m *LocationModel) CollectionName() string {
	return "locations"
}

// FromEntity maps domain→model
func FromEntity(e *entity.Location) *LocationModel {
	return &LocationModel{
		ID:        e.ID,
		Name:      e.Name,
		Type:      e.Type,
		Country:   e.Country,
		State:     e.State,
		District:  e.District,
		Pincode:   e.Pincode,
		Coordinates: struct {
			Lat float64 `bson:"lat" json:"lat"`
			Lng float64 `bson:"lng" json:"lng"`
		}{
			Lat: e.Coordinates.Lat,
			Lng: e.Coordinates.Lng,
		},
		GeoJSON:    e.GeoJSON,
		Tags:       e.Tags,
		Description: e.Description,
		Aliases:    e.Aliases,
		MediaURLs: struct {
			Photos []string `bson:"photos" json:"photos"`
			Videos []string `bson:"videos" json:"videos"`
		}{
			Photos: e.MediaURLs.Photos,
			Videos: e.MediaURLs.Videos,
		},
		CreatedBy: e.CreatedBy,
		CreatedAt: e.CreatedAt,
		UpdatedAt: e.UpdatedAt,
		DeletedAt: e.DeletedAt,
	}
}

// ToEntity maps model→domain
func (m *LocationModel) ToEntity() entity.Location {
	return entity.Location{
		ID:          m.ID,
		Name:        m.Name,
		Type:        m.Type,
		Country:     m.Country,
		State:       m.State,
		District:    m.District,
		Pincode:     m.Pincode,
		Coordinates: entity.Coordinates{Lat: m.Coordinates.Lat, Lng: m.Coordinates.Lng},
		GeoJSON:     m.GeoJSON,
		Tags:        m.Tags,
		Description: m.Description,
		Aliases:     m.Aliases,
		MediaURLs:   entity.MediaURLs{Photos: m.MediaURLs.Photos, Videos: m.MediaURLs.Videos},
		CreatedBy:   m.CreatedBy,
		CreatedAt:   m.CreatedAt,
		UpdatedAt:   m.UpdatedAt,
		DeletedAt:   m.DeletedAt,
	}
}
