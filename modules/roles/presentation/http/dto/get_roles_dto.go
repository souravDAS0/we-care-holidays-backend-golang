// modules/roles/presentation/http/dto/get_roles_dto.go
package dto

import (
	"strconv"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// GetRolesDto defines the query parameters for listing roles
type GetRolesDto struct {
	// Pagination parameters
	Page  int `form:"page" json:"page"`
	Limit int `form:"limit" json:"limit"`

	// Filter parameters
	Name               string   `form:"name" json:"name"`
	HasPermissions     *bool    `form:"hasPermissions" json:"hasPermissions"`
	PermissionIDs      []string `form:"permissionIds" json:"permissionIds"`
	HasAllPermissions  *bool    `form:"hasAllPermissions" json:"hasAllPermissions"`  // If true, role must have ALL specified permissions
	SearchTerm         string   `form:"searchTerm" json:"searchTerm"`
	IncludeDeleted     bool     `form:"includeDeleted" json:"includeDeleted"`
	PermissionCount    *int     `form:"permissionCount" json:"permissionCount"`       // Exact permission count
	MinPermissionCount *int     `form:"minPermissionCount" json:"minPermissionCount"` // Minimum permission count
	MaxPermissionCount *int     `form:"maxPermissionCount" json:"maxPermissionCount"` // Maximum permission count
}

// NewGetRolesDto creates a new DTO from query parameters
func NewGetRolesDto(c *gin.Context) GetRolesDto {
	dto := GetRolesDto{}

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
	if hasPermissionsStr := c.Query("hasPermissions"); hasPermissionsStr != "" {
		if hasPermissions, err := strconv.ParseBool(hasPermissionsStr); err == nil {
			dto.HasPermissions = &hasPermissions
		}
	}

	if hasAllPermissionsStr := c.Query("hasAllPermissions"); hasAllPermissionsStr != "" {
		if hasAllPermissions, err := strconv.ParseBool(hasAllPermissionsStr); err == nil {
			dto.HasAllPermissions = &hasAllPermissions
		}
	}

	// Parse includeDeleted parameter
	includeDeleted := false
	if includeDeletedStr := c.Query("includeDeleted"); includeDeletedStr == "true" {
		includeDeleted = true
	}
	dto.IncludeDeleted = includeDeleted

	// Parse permission count parameters
	if permissionCountStr := c.Query("permissionCount"); permissionCountStr != "" {
		if permissionCount, err := strconv.Atoi(permissionCountStr); err == nil {
			dto.PermissionCount = &permissionCount
		}
	}

	if minPermissionCountStr := c.Query("minPermissionCount"); minPermissionCountStr != "" {
		if minPermissionCount, err := strconv.Atoi(minPermissionCountStr); err == nil {
			dto.MinPermissionCount = &minPermissionCount
		}
	}

	if maxPermissionCountStr := c.Query("maxPermissionCount"); maxPermissionCountStr != "" {
		if maxPermissionCount, err := strconv.Atoi(maxPermissionCountStr); err == nil {
			dto.MaxPermissionCount = &maxPermissionCount
		}
	}

	// Parse filter parameters
	dto.Name = c.Query("name")
	dto.SearchTerm = c.Query("searchTerm")

	// Parse permission IDs array
	if permissionIDsQuery := c.QueryArray("permissionIds"); len(permissionIDsQuery) > 0 {
		// Validate permission IDs format
		validPermissionIDs := make([]string, 0)
		for _, permID := range permissionIDsQuery {
			if _, err := primitive.ObjectIDFromHex(permID); err == nil {
				validPermissionIDs = append(validPermissionIDs, permID)
			}
		}
		dto.PermissionIDs = validPermissionIDs
	}

	return dto
}

// ToFilterMap converts the DTO to a map for filtering in the repository
func (dto *GetRolesDto) ToFilterMap() map[string]interface{} {
	filter := make(map[string]interface{})

	if dto.Name != "" {
		// Using a case-insensitive regex for partial name matching
		filter["name"] = map[string]interface{}{
			"$regex":   dto.Name,
			"$options": "i", // case-insensitive
		}
	}

	// Permissions existence filter
	if dto.HasPermissions != nil {
		if *dto.HasPermissions {
			filter["permissions"] = map[string]interface{}{
				"$exists": true,
				"$ne":     []interface{}{},
			}
		} else {
			filter["$or"] = []map[string]interface{}{
				{"permissions": map[string]interface{}{"$exists": false}},
				{"permissions": []interface{}{}},
			}
		}
	}

	// Permission count filters
	if dto.PermissionCount != nil {
		filter["permissions"] = map[string]interface{}{
			"$size": *dto.PermissionCount,
		}
	} else {
		// Min/Max permission count filters (only if exact count is not specified)
		if dto.MinPermissionCount != nil || dto.MaxPermissionCount != nil {
			// Use aggregation pipeline for array size comparison
			sizeFilter := make(map[string]interface{})
			if dto.MinPermissionCount != nil {
				sizeFilter["$gte"] = *dto.MinPermissionCount
			}
			if dto.MaxPermissionCount != nil {
				sizeFilter["$lte"] = *dto.MaxPermissionCount
			}
			
			// This will need to be handled in the repository layer using aggregation
			filter["_permissionCountRange"] = sizeFilter
		}
	}

	// Specific permission IDs filter
	if len(dto.PermissionIDs) > 0 {
		// Convert string IDs to ObjectIDs
		permissionObjectIDs := make([]primitive.ObjectID, 0, len(dto.PermissionIDs))
		for _, permID := range dto.PermissionIDs {
			if objID, err := primitive.ObjectIDFromHex(permID); err == nil {
				permissionObjectIDs = append(permissionObjectIDs, objID)
			}
		}

		if len(permissionObjectIDs) > 0 {
			if dto.HasAllPermissions != nil && *dto.HasAllPermissions {
				// Role must have ALL specified permissions
				filter["permissions"] = map[string]interface{}{
					"$all": permissionObjectIDs,
				}
			} else {
				// Role must have ANY of the specified permissions
				filter["permissions"] = map[string]interface{}{
					"$in": permissionObjectIDs,
				}
			}
		}
	}

	// Search term filter (searches in name and description)
	if dto.SearchTerm != "" {
		filter["$or"] = []map[string]interface{}{
			{
				"name": map[string]interface{}{
					"$regex":   dto.SearchTerm,
					"$options": "i",
				},
			},
			{
				"description": map[string]interface{}{
					"$regex":   dto.SearchTerm,
					"$options": "i",
				},
			},
		}
	}

	// Default behavior: only show non-deleted roles
	// When includeDeleted=true, skip adding the deletedAt filter so all roles will be returned
	if !dto.IncludeDeleted {
		filter["deletedAt"] = nil
	}

	return filter
}

// // RequiresAggregation returns true if the filter requires aggregation pipeline
// // This is needed for complex filters like permission count ranges
// func (dto *GetRolesDto) RequiresAggregation() bool {
// 	return (dto.MinPermissionCount != nil || dto.MaxPermissionCount != nil) && dto.PermissionCount == nil
// }