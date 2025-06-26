// modules/permissions/presentation/http/dto/get_permissions_dto.go
package dto

import (
	"strconv"

	"github.com/gin-gonic/gin"
)

// GetPermissionsDto defines the query parameters for listing permissions
type GetPermissionsDto struct {
	// Pagination parameters
	Page  int `form:"page" json:"page"`
	Limit int `form:"limit" json:"limit"`

	// Filter parameters
	Resource       string `form:"resource" json:"resource"`
	Action         string `form:"action" json:"action"`
	Scope          string `form:"scope" json:"scope"`
	Enabled        *bool  `form:"enabled" json:"enabled"`
	SearchTerm     string `form:"searchTerm" json:"searchTerm"`
	IncludeDeleted bool   `form:"includeDeleted" json:"includeDeleted"`
}

// NewGetPermissionsDto creates a new DTO from query parameters
func NewGetPermissionsDto(c *gin.Context) GetPermissionsDto {
	dto := GetPermissionsDto{}

	// Parse pagination parameters with defaults
	page, err := strconv.Atoi(c.DefaultQuery("page", "1"))
	if err != nil || page < 1 {
		page = 1
	}
	dto.Page = page

	limit, err := strconv.Atoi(c.DefaultQuery("limit", "20"))
	if err != nil || limit < 1 {
		limit = 20
	}
	// Cap the maximum limit to prevent performance issues
	if limit > 100 {
		limit = 100
	}
	dto.Limit = limit

	// Parse boolean parameters
	if enabledStr := c.Query("enabled"); enabledStr != "" {
		if enabled, err := strconv.ParseBool(enabledStr); err == nil {
			dto.Enabled = &enabled
		}
	}

	// Parse includeDeleted parameter
	includeDeleted := false
	if includeDeletedStr := c.Query("includeDeleted"); includeDeletedStr == "true" {
		includeDeleted = true
	}
	dto.IncludeDeleted = includeDeleted

	// Parse filter parameters
	dto.Resource = c.Query("resource")
	dto.Action = c.Query("action")
	dto.Scope = c.Query("scope")
	dto.SearchTerm = c.Query("searchTerm")

	return dto
}

// ToFilterMap converts the DTO to a map for filtering in the repository
func (dto *GetPermissionsDto) ToFilterMap() map[string]interface{} {
	filter := make(map[string]interface{})

	if dto.Resource != "" {
		// Using a case-insensitive regex for partial resource matching
		filter["resource"] = map[string]interface{}{
			"$regex":   dto.Resource,
			"$options": "i", // case-insensitive
		}
	}

	if dto.Action != "" {
		filter["action"] = dto.Action
	}

	if dto.Scope != "" {
		filter["scope"] = dto.Scope
	}

	if dto.Enabled != nil {
		filter["enabled"] = *dto.Enabled
	}

	// Search term filter (searches in resource and notes)
	if dto.SearchTerm != "" {
		filter["$or"] = []map[string]interface{}{
			{
				"resource": map[string]interface{}{
					"$regex":   dto.SearchTerm,
					"$options": "i",
				},
			},
			{
				"notes": map[string]interface{}{
					"$regex":   dto.SearchTerm,
					"$options": "i",
				},
			},
		}
	}

	// Default behavior: only show non-deleted permissions
	// When includeDeleted=true, skip adding the deletedAt filter so all permissions will be returned
	if !dto.IncludeDeleted {
		filter["deletedAt"] = nil
	}

	return filter
}
