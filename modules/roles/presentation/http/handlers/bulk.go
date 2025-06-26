package handlers

import (
	"net/http"

	"bitbucket.org/abhishek_fordel/we-care-holidays-backend-golang/internal/middleware"
	"bitbucket.org/abhishek_fordel/we-care-holidays-backend-golang/internal/models"
	"github.com/gin-gonic/gin"
)

// BulkSoftDeleteRoles godoc
//
//	@Summary		Delete multiple roles
//	@Description	Soft-delete multiple roles by their IDs
//	@Tags			roles
//	@Accept			json
//	@Produce		json
//	@Param			request	body		models.BulkDeleteDto	true	"IDs to delete"
//	@Success		200		{object}	models.SwaggerStandardResponse{data=models.BulkDeleteResponse}
//	@Failure		400		{object}	models.SwaggerErrorResponse
//	@Failure		500		{object}	models.SwaggerErrorResponse
//	@Router			/roles/bulk-delete [delete]
func (h *RoleHandler) BulkSoftDeleteRoles(c *gin.Context) {
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
	result, err := h.BulkSoftDeleteRolesUseCase.Execute(c.Request.Context(), bulkDeleteDto.IDs)
	if err != nil {
		middleware.HandleError(c, middleware.NewAppError(
			middleware.ErrorCodeInternalServer,
			"Failed to delete roles",
			err,
			http.StatusInternalServerError,
		))
		return
	}

	success := result.DeletedCount > 0
	message := "Roles deletion processed"

	if !success {
		message = "No roles were deleted"

		// Add more specific details to the message
		if len(result.InvalidIDs) > 0 && len(result.NotFoundIDs) == 0 {
			message = "No roles were deleted: all provided IDs were invalid"
		} else if len(result.InvalidIDs) == 0 && len(result.NotFoundIDs) > 0 {
			message = "No roles were deleted: none of the provided IDs were found"
		} else if len(result.InvalidIDs) > 0 && len(result.NotFoundIDs) > 0 {
			message = "No roles were deleted: some IDs were invalid and others were not found"
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

// BulkRestoreRoles godoc
//
//	@Summary		Restore multiple deleted roles
//	@Description	Restore multiple soft-deleted roles by their IDs
//	@Tags			roles
//	@Accept			json
//	@Produce		json
//	@Param			request	body		models.BulkDeleteDto	true	"IDs to restore"
//	@Success		200		{object}	models.SwaggerStandardResponse{data=models.BulkRestoreResponse}
//	@Failure		400		{object}	models.SwaggerErrorResponse
//	@Failure		500		{object}	models.SwaggerErrorResponse
//	@Router			/roles/bulk-restore [post]
func (h *RoleHandler) BulkRestoreRoles(c *gin.Context) {
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

	result, err := h.BulkRestoreRolesUseCase.Execute(c.Request.Context(), bulkRestoreDto.IDs)
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
