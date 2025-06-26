// modules/users/presentation/http/dto/get_users_dto.go
package dto

import (
	"strconv"

	"github.com/gin-gonic/gin"
)

// GetUsersDto defines the query parameters for listing users
type GetUsersDto struct {
	// Pagination parameters
	Page  int `form:"page" json:"page"`
	Limit int `form:"limit" json:"limit"`

	// Filter parameters
	FullName        string `form:"fullName" json:"fullName"`
	Email           string `form:"email" json:"email"`
	Phone           string `form:"phone" json:"phone"`
	Status          string `form:"status" json:"status"`
	RoleID          string `form:"roleId" json:"roleId"`
	OrganizationID  string `form:"organizationId" json:"organizationId"`
	Search          string `form:"search" json:"search"`
	IncludeDeleted  bool   `form:"includeDeleted" json:"includeDeleted"`
}

// NewGetUsersDto creates a new DTO from query parameters
func NewGetUsersDto(c *gin.Context) GetUsersDto {
	dto := GetUsersDto{}

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
	dto.FullName = c.Query("fullName")
	dto.Email = c.Query("email")
	dto.Phone = c.Query("phone")
	dto.Status = c.Query("status")
	dto.RoleID = c.Query("roleId")
	dto.OrganizationID = c.Query("organizationId")
	dto.Search = c.Query("search")

	return dto
}

// ToFilterMap converts the DTO to a map for filtering in the repository
func (dto *GetUsersDto) ToFilterMap() map[string]interface{} {
	filter := make(map[string]interface{})

	if dto.FullName != "" {
		// Using a case-insensitive regex for partial name matching
		filter["fullName"] = map[string]interface{}{
			"$regex":   dto.FullName,
			"$options": "i", // case-insensitive
		}
	}

	if dto.Email != "" {
		// Using a case-insensitive regex for partial email matching
		filter["emails"] = map[string]interface{}{
			"$regex":   dto.Email,
			"$options": "i", // case-insensitive
		}
	}

	if dto.Phone != "" {
		filter["phones"] = map[string]interface{}{
			"$regex":   dto.Phone,
			"$options": "i",
		}
	}

	if dto.Status != "" {
		filter["status"] = dto.Status
	}

	if dto.RoleID != "" {
		filter["roleId"] = dto.RoleID
	}

	if dto.OrganizationID != "" {
		filter["organizationId"] = dto.OrganizationID
	}

	// Search across multiple fields
	if dto.Search != "" {
		filter["$or"] = []map[string]interface{}{
			{"fullName": map[string]interface{}{"$regex": dto.Search, "$options": "i"}},
			{"emails": map[string]interface{}{"$regex": dto.Search, "$options": "i"}},
			{"phones": map[string]interface{}{"$regex": dto.Search, "$options": "i"}},
		}
	}

	// Default behavior: only show non-deleted users
	// When includeDeleted=true, skip adding the deletedAt filter so all users will be returned
	if !dto.IncludeDeleted {
		filter["deletedAt"] = nil
	}

	return filter
}