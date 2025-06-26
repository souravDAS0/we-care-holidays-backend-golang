package handlers

import (
	"net/http"

	"bitbucket.org/abhishek_fordel/we-care-holidays-backend-golang/internal/middleware"
	"bitbucket.org/abhishek_fordel/we-care-holidays-backend-golang/modules/organizations/presentation/http/dto"
	"github.com/gin-gonic/gin"
)

// ListOrganizations godoc
//
//	@Summary		List organizations
//	@Description	Get all organizations with pagination and filtering
//	@Tags			organizations
//	@Accept			json
//	@Produce		json
//	@Param			page			query		int																			false	"Page number"		default(1)
//	@Param			limit			query		int																			false	"Items per page"	default(10)	maximum(100)
//	@Param			name			query		string																		false	"Filter by name (partial match, case-insensitive)"
//	@Param			slug			query		string																		false	"Filter by slug (partial match, case-insensitive)"
//	@Param			type			query		string																		false	"Filter by organization type (exact match)"
//	@Param			status			query		string																		false	"Filter by status (exact match)"
//	@Param			email			query		string																		false	"Filter by email (partial match, case-insensitive)"
//	@Param			phone			query		string																		false	"Filter by phone number (partial match, case-insensitive)"
//	@Param			website			query		string																		false	"Filter by website (partial match, case-insensitive)"
//	@Param			city			query		string																		false	"Filter by city (partial match, case-insensitive)"
//	@Param			state			query		string																		false	"Filter by state (partial match, case-insensitive)"
//	@Param			country			query		string																		false	"Filter by country (partial match, case-insensitive)"
//	@Param			includeDeleted	query		bool																		false	"Include soft-deleted organizations"	default(false)
//	@Success		200				{object}	models.SwaggerStandardResponse{data=dto.PaginatedOrganizationsResponse}	"Successful response with paginated organizations"
//	@Failure		400				{object}	models.SwaggerErrorResponse
//	@Failure		500				{object}	models.SwaggerErrorResponse
//	@Router			/organizations [get]
func (h *OrganizationHandler) ListOrganizations(c *gin.Context) {
	// Parse query parameters using DTO
	queryDto := dto.NewGetOrganizationsDto(c)

	// Call use case with filter from DTO
	organizations, total, err := h.ListOrganizationUseCase.Execute(
		c.Request.Context(),
		queryDto.ToFilterMap(),
		queryDto.Page,
		queryDto.Limit,
	)
	if err != nil {
		middleware.HandleError(c, middleware.NewAppError(
			middleware.ErrorCodeInternalServer,
			"Failed to fetch organizations",
			err,
			http.StatusInternalServerError,
		))
		return
	}

	// Prepare response
	response := gin.H{
		"items":          organizations,
		"page":           queryDto.Page,
		"limit":          queryDto.Limit,
		"total":          total,
		"totalPages":     (total + int64(queryDto.Limit) - 1) / int64(queryDto.Limit),
		"includeDeleted": queryDto.IncludeDeleted,
	}

	c.JSON(http.StatusOK, response)
}
