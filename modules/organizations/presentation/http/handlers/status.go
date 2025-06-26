package handlers

import (
	"net/http"

	"bitbucket.org/abhishek_fordel/we-care-holidays-backend-golang/internal/middleware"
	"bitbucket.org/abhishek_fordel/we-care-holidays-backend-golang/modules/organizations/presentation/http/dto"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// UpdateOrganizationStatus godoc
//
//	@Summary		Update organization status
//	@Description	Update the status of an organization
//	@Tags			organizations
//	@Accept			json
//	@Produce		json
//	@Param			id	path		string	true	"Organization ID"	example("6824886e6b180b753cea43e9")
//	@Success		200	{object}	models.SwaggerStandardResponse{data=entity.Organization}
//	@Failure		400	{object}	models.SwaggerErrorResponse
//	@Failure		404	{object}	models.SwaggerErrorResponse
//	@Failure		500	{object}	models.SwaggerErrorResponse
//	@Router			/organizations/{id}/status [put]
func (h *OrganizationHandler) UpdateOrganizationStatus(c *gin.Context) {
	id := c.Param("id")

	var statusDto dto.OrgStatusUpdateDto

	if err := c.ShouldBindJSON(&statusDto); err != nil {
		middleware.HandleError(c, middleware.NewAppError(
			middleware.ErrorCodeInvalidRequest,
			"Invalid request body",
			err,
			http.StatusBadRequest,
		))
		return
	}

	// Validate the status
	if err := statusDto.Validate(); err != nil {
		middleware.HandleError(c, middleware.NewAppError(
			middleware.ErrorCodeValidationFailed,
			err.Error(),
			nil,
			http.StatusBadRequest,
		))
		return
	}

	objectId, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		middleware.HandleError(c, middleware.NewAppError(
			middleware.ErrorCodeInvalidRequest,
			"Invalid organization ID",
			err,
			http.StatusBadRequest,
		))
		return
	}

	// IMPORTANT: First check if the organization exists
	organization, err := h.GetOrganizationUseCase.Execute(c.Request.Context(), objectId)
	if err != nil {
		middleware.HandleError(c, middleware.NewAppError(
			middleware.ErrorCodeInternalServer,
			"Failed to fetch organization",
			err,
			http.StatusInternalServerError,
		))
		return
	}

	// Check if organization exists
	if organization == nil {
		middleware.HandleError(c, middleware.NewAppError(
			middleware.ErrorCodeNotFound,
			"Organization not found",
			nil,
			http.StatusNotFound,
		))
		return
	}

	// Now proceed with the update, knowing the organization exists
	if err := h.UpdateOrganizationStatusUseCase.Execute(c.Request.Context(), objectId, statusDto.Status); err != nil {
		middleware.HandleError(c, middleware.NewAppError(
			middleware.ErrorCodeInternalServer,
			"Failed to update organization status",
			err,
			http.StatusInternalServerError,
		))
		return
	}

	// Re-fetch the updated organization
	updatedOrganization, err := h.GetOrganizationUseCase.Execute(c.Request.Context(), objectId)
	if err != nil {
		middleware.HandleError(c, middleware.NewAppError(
			middleware.ErrorCodeInternalServer,
			"Organization status updated but failed to fetch updated data",
			err,
			http.StatusInternalServerError,
		))
		return
	}

	c.JSON(http.StatusOK, updatedOrganization)
}
