package dto

import (
	"strconv"

	"github.com/gin-gonic/gin"
)

// GetOrganizationsDto defines the query parameters for listing organizations
type GetOrganizationsDto struct {
	// Pagination parameters
	Page  int `form:"page" json:"page"`
	Limit int `form:"limit" json:"limit"`

	// Filter parameters
	Name          string `form:"name" json:"name"`
	Slug          string `form:"slug" json:"slug"`
	Type          string `form:"type" json:"type"`
	Status        string `form:"status" json:"status"`
	Email         string `form:"email" json:"email"`
	Phone         string `form:"phone" json:"phone"`
	Website       string `form:"website" json:"website"`
	City          string `form:"city" json:"city"`
	State         string `form:"state" json:"state"`
	Country       string `form:"country" json:"country"`
	IncludeDeleted bool   `form:"includeDeleted" json:"includeDeleted"`
}

// NewGetOrganizationsDto creates a new DTO from query parameters
func NewGetOrganizationsDto(c *gin.Context) GetOrganizationsDto {
	dto := GetOrganizationsDto{}

	// Parse pagination parameters with defaults
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
	dto.Name = c.Query("name")
	dto.Slug = c.Query("slug")
	dto.Type = c.Query("type")
	dto.Status = c.Query("status")
	dto.Email = c.Query("email")
	dto.Phone = c.Query("phone")
	dto.Website = c.Query("website")
	
	
	// Parse address-related filters
	dto.City = c.Query("city")
	dto.State = c.Query("state")
	dto.Country = c.Query("country")

	return dto
}

// ToFilterMap converts the DTO to a map for filtering in the repository
func (dto *GetOrganizationsDto) ToFilterMap() map[string]interface{} {
	filter := make(map[string]interface{})

	if dto.Name != "" {
		// Using a case-insensitive regex for partial name matching
		filter["name"] = map[string]interface{}{
			"$regex":   dto.Name,
			"$options": "i", // case-insensitive
		}
	}

	if dto.Slug != "" {
		filter["slug"] = map[string]interface{}{
			"$regex":   dto.Slug,
			"$options": "i",
		}
	}

	if dto.Type != "" {
		filter["type"] = dto.Type
	}

	if dto.Status != "" {
		filter["status"] = dto.Status
	}

	if dto.Email != "" {
		// Using a case-insensitive regex for partial email matching
		filter["email"] = map[string]interface{}{
			"$regex":   dto.Email,
			"$options": "i", // case-insensitive
		}
	}
	
	if dto.Phone != "" {
		filter["phone"] = map[string]interface{}{
			"$regex":   dto.Phone,
			"$options": "i",
		}
	}
	
	if dto.Website != "" {
		filter["website"] = map[string]interface{}{
			"$regex":   dto.Website,
			"$options": "i",
		}
	}

	
	// Address-related filters
	if dto.City != "" {
		filter["address.city"] = map[string]interface{}{
			"$regex":   dto.City,
			"$options": "i",
		}
	}
	
	if dto.State != "" {
		filter["address.state"] = map[string]interface{}{
			"$regex":   dto.State,
			"$options": "i",
		}
	}
	
	if dto.Country != "" {
		filter["address.country"] = map[string]interface{}{
			"$regex":   dto.Country,
			"$options": "i",
		}
	}

	// Default behavior: only show non-deleted organizations
	// When includeDeleted=true, skip adding the deletedAt filter so all organizations will be returned
	if !dto.IncludeDeleted {
		filter["deletedAt"] = nil
	}

	return filter
}