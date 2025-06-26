package dto

import (
	"strconv"

	"github.com/gin-gonic/gin"
)

// ListLocationsDto defines the query parameters for listing locations
type ListLocationsDto struct {
	Page       int      `form:"page" json:"page" default:"1"`
	Limit      int      `form:"limit" json:"limit" default:"10"`
	
	Name       string  `form:"name" json:"name"`
	Type       string  `form:"type" json:"type"`
	Country    string  `form:"country" json:"country"`
	State      string  `form:"state" json:"state"`
	Tags       []string `form:"tags" json:"tags"`
	Aliases    []string `form:"aliases" json:"aliases"`
	IncludeDeleted bool   `form:"includeDeleted" json:"includeDeleted"`
}

func NewListLocationsDto(c *gin.Context) ListLocationsDto {

	dto := ListLocationsDto{}

	page, err := strconv.Atoi(c.DefaultQuery("page", "1"))
	if err != nil || page < 1 {
		page = 1
	}
	dto.Page = page
	limit, err := strconv.Atoi(c.DefaultQuery("limit", "10"))
	if err != nil || limit < 1 {
		limit = 10
	}
	// Cap the maximum limit to prevent performance issues
	if limit > 100 {
		limit = 100
	}
	dto.Limit = limit
	// Parse includeDeleted parameter
	includeDeleted := false
	if includeDeletedStr := c.Query("includeDeleted"); includeDeletedStr == "true" {
		includeDeleted = true
	}
	dto.IncludeDeleted = includeDeleted
	// Parse filter parameters
	if name := c.Query("name"); name != "" {
		dto.Name = name
	}
	if locationType := c.Query("type"); locationType != "" {
		dto.Type = locationType
	}
	if country := c.Query("country"); country != "" {
		dto.Country = country
	}
	if state := c.Query("state"); state != "" {
		dto.State = state
	}
	if tags := c.QueryArray("tags"); len(tags) > 0 {
		dto.Tags = tags
	}
	if aliases := c.QueryArray("aliases"); len(aliases) > 0 {
		dto.Aliases = aliases
	}
	return dto
}

// ToFilterMap converts the ListLocationsDto to a map for filtering in the repository
func (dto *ListLocationsDto) ToFilterMap() map[string]interface{} {
	filter := make(map[string]interface{})

	if dto.Name != "" {
		filter["name"] = map[string]interface{}{
			"$regex":   dto.Name,
			"$options": "i",
		}
	}
	if dto.Type != "" {
		filter["type"] = dto.Type
	}
	if dto.Country != "" {
		filter["country"] = dto.Country
	}
	if dto.State != "" {
		filter["state"] = dto.State
	}
	if len(dto.Tags) > 0 {
		filter["tags"] = map[string]interface{}{
			"$in": dto.Tags,
		}
	}
	if len(dto.Aliases) > 0 {
		filter["aliases"] = map[string]interface{}{
			"$in": dto.Aliases,
		}
	}
	if !dto.IncludeDeleted {
		filter["deletedAt"] = nil
	}

	return filter
}
