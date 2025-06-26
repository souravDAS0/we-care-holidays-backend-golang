package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"bitbucket.org/abhishek_fordel/we-care-holidays-backend-golang/internal/middleware"
	"bitbucket.org/abhishek_fordel/we-care-holidays-backend-golang/modules/locations/presentation/http/dto"
)

// ListLocations godoc
//
//	@Summary		List locations
//	@Description	Get all locations with pagination and filtering
//	@Tags			locations
//	@Accept			json
//	@Produce		json
//	@Param			page			query		int																		false	"Page number"		default(1)
//	@Param			limit			query		int																		false	"Items per page"	default(10)	maximum(100)
//	@Param			name			query		string																	false	"Filter by name (partial match, case-insensitive)"
//	@Param			type			query		string																	false	"Filter by organization type (exact match)"
//	@Param			state			query		string																	false	"Filter by state (partial match, case-insensitive)"
//	@Param			country			query		string																	false	"Filter by country (partial match, case-insensitive)"
//	@Param			tags			query		[]string																false	"Filter by tags (exact match, any)"
//	@Param			aliases			query		[]string																false	"Filter by aliases (exact match, any)"
//	@Param			includeDeleted	query		bool																	false	"Include soft-deleted locations"	default(false)
//	@Success		200				{object}	models.SwaggerStandardResponse{data=dto.PaginatedLocationsResponse}	"Successful response with paginated locations"
//	@Failure		400				{object}	models.SwaggerErrorResponse
//	@Failure		500				{object}	models.SwaggerErrorResponse
//	@Router			/locations [get]
func (h *LocationHandler) ListLocations(c *gin.Context) {
	queryDto := dto.NewListLocationsDto(c)

	locations, total, err := h.ListLocationsUseCase.Execute(c.Request.Context(), queryDto.ToFilterMap(), queryDto.Page, queryDto.Limit)
	if err != nil {
		middleware.HandleError(c, middleware.NewAppError(middleware.ErrorCodeInternalServer, "Failed to list locations", err, http.StatusInternalServerError))
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"items":          locations,
		"page":           queryDto.Page,
		"limit":          queryDto.Limit,
		"total":          total,
		"totalPages":     (total + int64(queryDto.Limit) - 1) / int64(queryDto.Limit),
		"includeDeleted": queryDto.IncludeDeleted,
	})
}
