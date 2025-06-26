package handlers

import (
	"net/http"

	"bitbucket.org/abhishek_fordel/we-care-holidays-backend-golang/internal/middleware"
	"bitbucket.org/abhishek_fordel/we-care-holidays-backend-golang/modules/roles/presentation/http/dto"
	"github.com/gin-gonic/gin"
)

// ListRoles godoc
//
//	@Summary		List roles
//	@Description	Get all roles with pagination and filtering
//	@Tags			roles
//	@Accept			json
//	@Produce		json
//	@Param			page				query		int		false	"Page number"													default(1)
//	@Param			limit				query		int		false	"Items per page"											default(20)	maximum(100)
//	@Param			name				query		string	false	"Filter by name (partial match, case-insensitive)"
//	@Param			hasPermissions		query		bool	false	"Filter roles that have permissions (true) or no permissions (false)"
//	@Param			hasAllPermissions	query		bool	false	"If true, role must have ALL specified permissions; if false, role must have ANY"
//	@Param			searchTerm			query		string	false	"Search in role name and description (case-insensitive)"
//	@Param			includeDeleted		query		bool	false	"Include soft-deleted roles"								default(false)
//	@Param			permissionCount		query		int		false	"Filter by exact permission count"
//	@Param			minPermissionCount	query		int		false	"Filter by minimum permission count"
//	@Param			maxPermissionCount	query		int		false	"Filter by maximum permission count"
//	@Success		200					{object}	models.SwaggerStandardResponse{data=dto.PaginatedRolesResponse}	"Successful response with paginated roles"
//	@Failure		400					{object}	models.SwaggerErrorResponse
//	@Failure		500					{object}	models.SwaggerErrorResponse
//	@Router			/roles [get]
func (h *RoleHandler) ListRoles(c *gin.Context) {
	// Parse query parameters using DTO
	queryDto := dto.NewGetRolesDto(c)

	// Call use case with filter from DTO
	roles, total, err := h.ListRolesUseCase.Execute(
		c.Request.Context(),
		queryDto.ToFilterMap(),
		queryDto.Page,
		queryDto.Limit,
	)
	if err != nil {
		middleware.HandleError(c, middleware.NewAppError(
			middleware.ErrorCodeInternalServer,
			"Failed to fetch roles",
			err,
			http.StatusInternalServerError,
		))
		return
	}

	// Prepare response
	response := gin.H{
		"items":          roles,
		"page":           queryDto.Page,
		"limit":          queryDto.Limit,
		"total":          total,
		"totalPages":     (total + int64(queryDto.Limit) - 1) / int64(queryDto.Limit),
		"includeDeleted": queryDto.IncludeDeleted,
	}

	c.JSON(http.StatusOK, response)
}
