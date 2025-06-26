package handlers

import (
	"net/http"

	"bitbucket.org/abhishek_fordel/we-care-holidays-backend-golang/internal/middleware"
	"bitbucket.org/abhishek_fordel/we-care-holidays-backend-golang/modules/users/presentation/http/dto"

	"github.com/gin-gonic/gin"
)

// ListUsers godoc
//
//	@Summary		List users
//	@Description	Get all users with pagination and filtering
//	@Tags			users
//	@Accept			json
//	@Produce		json
//	@Param			page			query		int									false	"Page number"		default(1)
//	@Param			limit			query		int									false	"Items per page"	default(10)	maximum(100)
//	@Param			fullName		query		string								false	"Filter by full name (partial match, case-insensitive)"
//	@Param			email			query		string								false	"Filter by email (partial match, case-insensitive)"
//	@Param			phone			query		string								false	"Filter by phone number (partial match, case-insensitive)"
//	@Param			status			query		string								false	"Filter by status (exact match)"
//	@Param			roleId			query		string								false	"Filter by role ID (exact match)"
//	@Param			organizationId	query		string								false	"Filter by organization ID (exact match)"
//	@Param			search			query		string								false	"Search across name, email, and phone"
//	@Param			includeDeleted	query		bool								false	"Include soft-deleted users"	default(false)
//	@Success		200				{object}	models.SwaggerStandardResponse{data=dto.PaginatedUsersResponse}
//	@Failure		400				{object}	models.SwaggerErrorResponse
//	@Failure		500				{object}	models.SwaggerErrorResponse
//	@Router			/users [get]
func (h *UserHandler) ListUsers(c *gin.Context) {
	// Parse query parameters using DTO
	queryDto := dto.NewGetUsersDto(c)

	// Call use case with filter from DTO
	users, total, err := h.ListUsersUseCase.Execute(
		c.Request.Context(),
		queryDto.ToFilterMap(),
		queryDto.Page,
		queryDto.Limit,
	)
	if err != nil {
		middleware.HandleError(c, middleware.NewAppError(
			middleware.ErrorCodeInternalServer,
			"Failed to fetch users",
			err,
			http.StatusInternalServerError,
		))
		return
	}

	// Prepare response
	response := gin.H{
		"items":          users,
		"page":           queryDto.Page,
		"limit":          queryDto.Limit,
		"total":          total,
		"totalPages":     (total + int64(queryDto.Limit) - 1) / int64(queryDto.Limit),
		"includeDeleted": queryDto.IncludeDeleted,
	}

	c.JSON(http.StatusOK, response)
}
