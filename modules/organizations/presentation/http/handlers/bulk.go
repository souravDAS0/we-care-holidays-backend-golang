package handlers

import (
	"net/http"

	"bitbucket.org/abhishek_fordel/we-care-holidays-backend-golang/internal/middleware"
	"bitbucket.org/abhishek_fordel/we-care-holidays-backend-golang/internal/models"
	"github.com/gin-gonic/gin"
)

// BulkDeleteOrganizations godoc
//
//	@Summary		Delete multiple organizations
//	@Description	Soft-delete multiple organizations by their IDs
//	@Tags			organizations
//	@Accept			json
//	@Produce		json
//	@Param			request	body		models.BulkDeleteDto	true	"IDs to delete"
//	@Success		200		{object}	models.SwaggerStandardResponse{data=models.BulkDeleteResponse}
//	@Failure		400		{object}	models.SwaggerErrorResponse
//	@Failure		500		{object}	models.SwaggerErrorResponse
//	@Router			/organizations/bulk-delete [delete]
func (h *OrganizationHandler) BulkDeleteOrganizations(c *gin.Context) {
	var bulkDeleteDto models.BulkDeleteDto

	if err := c.ShouldBindJSON(&bulkDeleteDto); err != nil {
		middleware.HandleError(c, middleware.NewAppError(
			middleware.ErrorCodeInvalidRequest,
			"Invalid request body",
			err,
			http.StatusBadRequest,
		))
		return
	}

	// Validate the DTO
	if err := bulkDeleteDto.Validate(); err != nil {
		middleware.HandleError(c, middleware.NewAppError(
			middleware.ErrorCodeValidationFailed,
			err.Error(),
			nil,
			http.StatusBadRequest,
		))
		return
	}

	// Use soft delete instead of hard delete
	result, err := h.BulkSoftDeleteOrganizationsUseCase.Execute(c.Request.Context(), bulkDeleteDto.IDs)
	if err != nil {
		middleware.HandleError(c, middleware.NewAppError(
			middleware.ErrorCodeInternalServer,
			"Failed to delete organizations",
			err,
			http.StatusInternalServerError,
		))
		return
	}

	success := result.DeletedCount > 0
	message := "Organizations deletion processed"

	if !success {
		message = "No organizations were deleted"

		// Add more specific details to the message
		if len(result.InvalidIDs) > 0 && len(result.NotFoundIDs) == 0 {
			message = "No organizations were deleted: all provided IDs were invalid"
		} else if len(result.InvalidIDs) == 0 && len(result.NotFoundIDs) > 0 {
			message = "No organizations were deleted: none of the provided IDs were found"
		} else if len(result.InvalidIDs) > 0 && len(result.NotFoundIDs) > 0 {
			message = "No organizations were deleted: some IDs were invalid and others were not found"
		}

		middleware.HandleError(c, middleware.NewAppError(
			middleware.ErrorCodeNotFound,
			message,
			nil,
			http.StatusBadRequest,
		))
		return
	}

	response := models.BulkDeleteResponse{
		Message:      message,
		RequestedIDs: bulkDeleteDto.IDs,
		DeletedCount: result.DeletedCount,
		DeletedIDs:   result.DeletedIDs,
		InvalidIDs:   result.InvalidIDs,
		NotFoundIDs:  result.NotFoundIDs,
	}

	c.JSON(http.StatusOK, response)
}

// BulkRestoreOrganizations godoc
//
//	@Summary		Restore multiple deleted organizations
//	@Description	Restore multiple soft-deleted organizations by their IDs
//	@Tags			organizations
//	@Accept			json
//	@Produce		json
//	@Param			request	body		models.BulkDeleteDto	true	"IDs to restore"
//	@Success		200		{object}	models.SwaggerStandardResponse{data=models.BulkRestoreResponse}
//	@Failure		400		{object}	models.SwaggerErrorResponse
//	@Failure		500		{object}	models.SwaggerErrorResponse
//	@Router			/organizations/bulk-restore [post]
func (h *OrganizationHandler) BulkRestoreOrganizations(c *gin.Context) {
	// Reuse the BulkDeleteDto structure since it's just a list of IDs
	var bulkRestoreDto models.BulkDeleteDto

	if err := c.ShouldBindJSON(&bulkRestoreDto); err != nil {
		middleware.HandleError(c, middleware.NewAppError(
			middleware.ErrorCodeInvalidRequest,
			"Invalid request body",
			err,
			http.StatusBadRequest,
		))
		return
	}

	// Validate the DTO
	if err := bulkRestoreDto.Validate(); err != nil {
		middleware.HandleError(c, middleware.NewAppError(
			middleware.ErrorCodeValidationFailed,
			err.Error(),
			nil,
			http.StatusBadRequest,
		))
		return
	}

	result, err := h.BulkRestoreOrganizationsUseCase.Execute(c.Request.Context(), bulkRestoreDto.IDs)
	if err != nil {
		middleware.HandleError(c, middleware.NewAppError(
			middleware.ErrorCodeInternalServer,
			"Failed to restore roles",
			err,
			http.StatusInternalServerError,
		))
		return
	}

	success := result.RestoredCount > 0
	message := "Roles restoration processed"

	if !success {
		message = "No roles were restored"

		// Add more specific details to the message
		if len(result.InvalidIDs) > 0 && len(result.NotFoundIDs) == 0 {
			message = "No roles were restored: all provided IDs were invalid"
		} else if len(result.InvalidIDs) == 0 && len(result.NotFoundIDs) > 0 {
			message = "No roles were restored: none of the provided IDs were found or already active"
		} else if len(result.InvalidIDs) > 0 && len(result.NotFoundIDs) > 0 {
			message = "No roles were restored: some IDs were invalid and others were not found or already active"
		}

		middleware.HandleError(c, middleware.NewAppError(
			middleware.ErrorCodeNotFound,
			message,
			nil,
			http.StatusBadRequest,
		))
		return
	}

	response := models.BulkRestoreResponse{
		Message:       message,
		RequestedIDs:  bulkRestoreDto.IDs,
		RestoredCount: result.RestoredCount,
		RestoredIDs:   result.RestoredIDs,
		InvalidIDs:    result.InvalidIDs,
		NotFoundIDs:   result.NotFoundIDs,
	}

	c.JSON(http.StatusOK, response)
}
