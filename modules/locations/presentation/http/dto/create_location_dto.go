package dto

import (
	"errors"

	"slices"

	"bitbucket.org/abhishek_fordel/we-care-holidays-backend-golang/modules/locations/domain/entity"
)

// LocationType represents the type of location
type LocationType string

const (
	City         LocationType = "city"
	ScenicSpot   LocationType = "scenic_spot"
	TransitHub   LocationType = "transit_hub"
	District     LocationType = "district"
	Village      LocationType = "village"
	Region       LocationType = "region"
)

// CreateLocationDto defines the required and optional fields for creating a location
type CreateLocationDto struct {
	// Required fields
	Name        string      `json:"name" binding:"required" example:"Siliguri"`
	Type        LocationType `json:"type" binding:"required" example:"city"`
	Country     string      `json:"country" binding:"required" example:"India"`
	State       string      `json:"state" binding:"required" example:"West Bengal"`
	District    string      `json:"district" binding:"required" example:"Darjeeling"`
	Pincode     string      `json:"pincode" binding:"required" example:"734001"`

	// Optional fields
	Coordinates CoordinatesDto `json:"coordinates,omitempty"`
	Geojson     interface{}    `json:"geojson,omitempty"`
	Tags        []string       `json:"tags,omitempty"`
	Description string         `json:"description,omitempty"`
	Aliases     []string       `json:"aliases,omitempty"`
	MediaURLs   MediaUrlsDto   `json:"mediaUrls,omitempty"`
}

// CoordinatesDto represents location coordinates
type CoordinatesDto struct {
	Lat float64 `json:"lat" example:"26.7270"`
	Lng float64 `json:"lng" example:"88.3950"`
}

func (c CoordinatesDto) IsValid() bool {
    return !(c.Lat == 0 && c.Lng == 0) && 
           c.Lat >= -90 && c.Lat <= 90 && 
           c.Lng >= -180 && c.Lng <= 180
}

// MediaUrlsDto represents media (photos/videos) URLs
type MediaUrlsDto struct {
	Photos []string `json:"photos,omitempty"`
	Videos []string `json:"videos,omitempty"`
}

// Validate validates the location creation DTO
func (dto *CreateLocationDto) Validate() error {
	if dto.Name == "" {
		return errors.New("location name is required")
	}
	if dto.Country == "" || dto.State == "" || dto.District == "" || dto.Pincode == "" {
		return errors.New("country, state, district, and pincode are required")
	}
	return nil
}

// ToEntity converts the CreateLocationDto to an entity
func (dto *CreateLocationDto) ToEntity() *entity.Location {
	location := &entity.Location{
		Name:        dto.Name,
		Type:        string(dto.Type),
		Country:     dto.Country,
		State:       dto.State,
		District:    dto.District,
		Pincode:     dto.Pincode,
		Description: dto.Description,
	}

	if dto.Geojson != nil {
		location.GeoJSON = dto.Geojson
	}

	if dto.Coordinates.IsValid() {
		location.Coordinates = entity.Coordinates{
			Lat: dto.Coordinates.Lat,
			Lng: dto.Coordinates.Lng,
		}
	} 

	if len(dto.MediaURLs.Photos) > 0 {
		if slices.Contains(dto.MediaURLs.Photos, "") {
				return nil
			}
		location.MediaURLs.Photos = dto.MediaURLs.Photos
	}

	if len(dto.MediaURLs.Videos) > 0 {
		if slices.Contains(dto.MediaURLs.Videos, "") {
				return nil
			}
		location.MediaURLs.Videos = dto.MediaURLs.Videos
	}

	if len(dto.Aliases) > 0 {
		if slices.Contains(dto.Aliases, "") {
				return nil
			}
		location.Aliases = dto.Aliases
	}
	if len(dto.Tags) > 0 {
		if slices.Contains(dto.Tags, "") {
				return nil
			}
		location.Tags = dto.Tags
	}

	return location
}

