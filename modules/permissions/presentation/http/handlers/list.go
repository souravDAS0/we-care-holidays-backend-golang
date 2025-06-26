package handlers

import (
	"net/http"

	"bitbucket.org/abhishek_fordel/we-care-holidays-backend-golang/internal/middleware"
	"bitbucket.org/abhishek_fordel/we-care-holidays-backend-golang/modules/permissions/presentation/http/dto"
	"github.com/gin-gonic/gin"
)

// ListPermissions godoc
//
//	@Summary		List permissions
//	@Description	Get all permissions with pagination and filtering
//	@Tags			permissions
//	@Accept			json
//	@Produce		json
//	@Param			page			query		int																			false	"Page number"		default(1)
//	@Param			limit			query		int																			false	"Items per page"	default(10)	maximum(100)
//	@Param			name			query		string																		false	"Filter by name (partial match, case-insensitive)"
//	@Param			description		query		string																		false	"Filter by description (partial match, case-insensitive)"
//	@Param			resource		query		string																		false	"Filter by resource (exact match)"
//	@Param			action			query		string																		false	"Filter by action (exact match)"
//	@Param			includeDeleted	query		bool																		false	"Include soft-deleted permissions"	default(false)
//	@Success		200				{object}	models.SwaggerStandardResponse{data=dto.PaginatedPermissionsResponse}	"Successful response with paginated permissions"
//	@Failure		400				{object}	models.SwaggerErrorResponse
//	@Failure		500				{object}	models.SwaggerErrorResponse
//	@Router			/permissions [get]
func (h *PermissionHandler) ListPermissions(c *gin.Context) {
	// Parse query parameters using DTO
	queryDto := dto.NewGetPermissionsDto(c)

	// Call use case with filter from DTO
	permissions, total, err := h.ListPermissionsUseCase.Execute(
		c.Request.Context(),
		queryDto.ToFilterMap(),
		queryDto.Page,
		queryDto.Limit,
	)
	if err != nil {
		middleware.HandleError(c, middleware.NewAppError(
			middleware.ErrorCodeInternalServer,
			"Failed to fetch permissions",
			err,
			http.StatusInternalServerError,
		))
		return
	}

	// Prepare response
	response := gin.H{
		"items":          permissions,
		"page":           queryDto.Page,
		"limit":          queryDto.Limit,
		"total":          total,
		"totalPages":     (total + int64(queryDto.Limit) - 1) / int64(queryDto.Limit),
		"includeDeleted": queryDto.IncludeDeleted,
	}

	c.JSON(http.StatusOK, response)
}
