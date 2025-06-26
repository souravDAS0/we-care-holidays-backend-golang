package entity

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Coordinates represents latitude/longitude
type Coordinates struct {
	Lat float64 `json:"lat" bson:"lat"`
	Lng float64 `json:"lng" bson:"lng"`
}

// MediaURLs holds arrays of media links
type MediaURLs struct {
	Photos []string `json:"photos" bson:"photos"`
	Videos []string `json:"videos" bson:"videos"`
}

// Location is the domain entity for a geo-tagged place
type Location struct {
	ID          primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	Name        string             `json:"name" bson:"name"`
	Type        string             `json:"type" bson:"type"` // city | scenic_spot | ...
	Country     string             `json:"country" bson:"country"`
	State       string             `json:"state" bson:"state"`
	District    string             `json:"district" bson:"district"`
	Pincode     string             `json:"pincode" bson:"pincode"`
	Coordinates Coordinates        `json:"coordinates" bson:"coordinates"`
	GeoJSON     interface{}        `json:"geojson,omitempty" bson:"geojson,omitempty"`
	Tags        []string           `json:"tags" bson:"tags"`
	Description string             `json:"description" bson:"description"`
	Aliases     []string           `json:"aliases" bson:"aliases"`
	MediaURLs   MediaURLs          `json:"mediaUrls" bson:"mediaUrls"`
	CreatedBy   primitive.ObjectID `json:"createdBy" bson:"createdBy"`
	CreatedAt   time.Time          `json:"createdAt" bson:"createdAt"`
	UpdatedAt   time.Time          `json:"updatedAt" bson:"updatedAt"`
	DeletedAt   *time.Time         `json:"deletedAt,omitempty" bson:"deletedAt,omitempty"`
}

// IsDeleted returns true if soft-deleted
func (l *Location) IsDeleted() bool {
	return l.DeletedAt != nil
}