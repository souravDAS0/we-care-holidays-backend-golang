package dto

import (
	"errors"

	"bitbucket.org/abhishek_fordel/we-care-holidays-backend-golang/modules/locations/domain/entity"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// UpdateLocationDto defines the fields for updating a location
// All fields are pointers to distinguish between nil (not provided) and empty values
type UpdateLocationDto struct {
	Name        *string       `json:"name,omitempty"`
	Type        *LocationType `json:"type,omitempty"`
	Country     *string       `json:"country,omitempty"`
	State       *string       `json:"state,omitempty"`
	District    *string       `json:"district,omitempty"`
	Pincode     *string       `json:"pincode,omitempty"`
	Coordinates *CoordinatesDto `json:"coordinates,omitempty"`
	Geojson     interface{}   `json:"geojson,omitempty"`
	Tags        []string      `json:"tags,omitempty"`
	Description *string       `json:"description,omitempty"`
	Aliases     []string      `json:"aliases,omitempty"`
	MediaUrls   *MediaUrlsDto `json:"mediaUrls,omitempty"`
}

// Validate performs validation on the UpdateLocationDto
func (dto *UpdateLocationDto) Validate() error {
	if dto.Name != nil && *dto.Name == "" {
		return errors.New("location name cannot be empty")
	}
	if dto.Country != nil && *dto.Country == "" {
		return errors.New("location country cannot be empty")
	}
	if dto.State != nil && *dto.State == "" {
		return errors.New("location state cannot be empty")
	}
	if dto.District != nil && *dto.District == "" {
		return errors.New("location district cannot be empty")
	}
	if dto.Pincode != nil && *dto.Pincode == "" {
		return errors.New("location pincode cannot be empty")
	}
	return nil
}

// ApplyUpdates applies the update DTO to an existing location entity
func (dto *UpdateLocationDto) ApplyUpdates(loc *entity.Location) {
	if dto.Name != nil {
		loc.Name = *dto.Name
	}
	if dto.Type != nil {
		loc.Type = string(*dto.Type)
	}
	if dto.Country != nil {
		loc.Country = *dto.Country
	}
	if dto.State != nil {
		loc.State = *dto.State
	}
	if dto.District != nil {
		loc.District = *dto.District
	}
	if dto.Pincode != nil {
		loc.Pincode = *dto.Pincode
	}
	if dto.Coordinates != nil {
		loc.Coordinates = entity.Coordinates{
			Lat: dto.Coordinates.Lat,
			Lng: dto.Coordinates.Lng,
		}
	}
	if dto.Geojson != nil {
		loc.GeoJSON = dto.Geojson
	}
	if dto.Tags != nil {
		loc.Tags = dto.Tags
	}
	if dto.Description != nil {
		loc.Description = *dto.Description
	}
	if dto.Aliases != nil {
		loc.Aliases = dto.Aliases
	}
	if dto.MediaUrls != nil {
		loc.MediaURLs = entity.MediaURLs{
			Photos: dto.MediaUrls.Photos,
			Videos: dto.MediaUrls.Videos,
		}
	}
}

// ToEntity converts UpdateLocationDto to an entity
func (dto *UpdateLocationDto) ToEntity(id primitive.ObjectID) *entity.Location {
	loc := &entity.Location{
		ID: id,
	}
	dto.ApplyUpdates(loc)
	return loc
}
