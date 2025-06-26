package handlers

import (
	"net/http"

	"bitbucket.org/abhishek_fordel/we-care-holidays-backend-golang/internal/middleware"
	"bitbucket.org/abhishek_fordel/we-care-holidays-backend-golang/modules/locations/presentation/http/dto"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// GetLocation godoc
//
//	@Summary		Get location by ID
//	@Description	Get single location by its ID
//	@Tags			locations
//	@Accept			json
//	@Produce		json
//	@Param			id	path		string	true	"Location ID"	example("6824886e6b180b753cea43e9")
//	@Success		200	{object}	models.SwaggerStandardResponse{data=entity.Location}
//	@Failure		400	{object}	models.SwaggerErrorResponse
//	@Failure		404	{object}	models.SwaggerErrorResponse
//	@Failure		500	{object}	models.SwaggerErrorResponse
//	@Router			/locations/{id} [get]
func (h *LocationHandler) GetLocation(c *gin.Context) {
	id := c.Param("id")

	objectId, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		middleware.HandleError(c, middleware.NewAppError(
			middleware.ErrorCodeInvalidRequest,
			"Invalid location ID",
			err,
			http.StatusBadRequest,
		))
		return
	}

	location, err := h.GetLocationUseCase.Execute(c.Request.Context(), objectId)
	if err != nil {
		middleware.HandleError(c, middleware.NewAppError(
			middleware.ErrorCodeInternalServer,
			"Failed to fetch location",
			err,
			http.StatusInternalServerError,
		))
		return
	}

	if location == nil {
		middleware.HandleError(c, middleware.NewAppError(
			middleware.ErrorCodeNotFound,
			"Location not found",
			nil,
			http.StatusNotFound,
		))
		return
	}

	c.JSON(http.StatusOK, location)
}

// CreateLocation godoc
//
//	@Summary		Create a new location
//	@Description	Create a new location with the provided data
//	@Tags			locations
//	@Accept			json
//	@Produce		json
//	@Param			location	body		dto.CreateLocationDto	true	"Location data"
//	@Success		201			{object}	models.SwaggerStandardResponse{data=entity.Location}
//	@Failure		400			{object}	models.SwaggerErrorResponse
//	@Failure		422			{object}	models.SwaggerErrorResponse
//	@Failure		500			{object}	models.SwaggerErrorResponse
//	@Router			/locations [post]
func (h *LocationHandler) CreateLocation(c *gin.Context) {
	var createDto dto.CreateLocationDto

	if err := c.ShouldBindJSON(&createDto); err != nil {
		middleware.HandleError(c, middleware.NewAppError(
			middleware.ErrorCodeInvalidRequest,
			"Invalid request body",
			err,
			http.StatusBadRequest,
		))
		return
	}

	// Validate the DTO
	if err := createDto.Validate(); err != nil {
		middleware.HandleError(c, middleware.NewAppError(
			middleware.ErrorCodeValidationFailed,
			err.Error(),
			nil,
			http.StatusBadRequest,
		))
		return
	}

	// Convert to entity
	location := createDto.ToEntity()

	// Call use case to create
	if err := h.CreateLocationUseCase.Execute(c.Request.Context(), location); err != nil {
		middleware.HandleError(c, middleware.NewAppError(
			middleware.ErrorCodeInternalServer,
			"Failed to create location",
			err,
			http.StatusInternalServerError,
		))
		return
	}

	c.JSON(http.StatusCreated, location)
}

// UpdateLocation godoc
//
//	@Summary		Update an existing location
//	@Description	Update a location by ID with provided data
//	@Tags			locations
//	@Accept			json
//	@Produce		json
//	@Param			id			path		string					true	"Location ID"
//	@Param			location	body		dto.UpdateLocationDto	true	"Location data to update"
//	@Success		200			{object}	models.SwaggerStandardResponse{data=entity.Location}
//	@Failure		400			{object}	models.SwaggerErrorResponse
//	@Failure		404			{object}	models.SwaggerErrorResponse
//	@Failure		500			{object}	models.SwaggerErrorResponse
//	@Router			/locations/{id} [put]
func (h *LocationHandler) UpdateLocation(c *gin.Context) {
	id := c.Param("id")

	objectId, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		middleware.HandleError(c, middleware.NewAppError(
			middleware.ErrorCodeInvalidRequest,
			"Invalid location ID",
			err,
			http.StatusBadRequest,
		))
		return
	}

	// Fetch existing location to apply updates
	existingLoc, err := h.GetLocationUseCase.Execute(c.Request.Context(), objectId)
	if err != nil {
		middleware.HandleError(c, middleware.NewAppError(
			middleware.ErrorCodeInternalServer,
			"Failed to fetch location",
			err,
			http.StatusInternalServerError,
		))
		return
	}

	if existingLoc == nil {
		middleware.HandleError(c, middleware.NewAppError(
			middleware.ErrorCodeNotFound,
			"Location not found",
			nil,
			http.StatusNotFound,
		))
		return
	}

	// Parse update DTO
	var updateDto dto.UpdateLocationDto
	if err := c.ShouldBindJSON(&updateDto); err != nil {
		middleware.HandleError(c, middleware.NewAppError(
			middleware.ErrorCodeInvalidRequest,
			"Invalid request body",
			err,
			http.StatusBadRequest,
		))
		return
	}

	// Validate the update DTO
	if err := updateDto.Validate(); err != nil {
		middleware.HandleError(c, middleware.NewAppError(
			middleware.ErrorCodeValidationFailed,
			err.Error(),
			nil,
			http.StatusBadRequest,
		))
		return
	}

	// Apply updates to the existing location
	updateDto.ApplyUpdates(existingLoc)

	// Call use case to update
	if err := h.UpdateLocationUseCase.Execute(c.Request.Context(), existingLoc); err != nil {
		middleware.HandleError(c, middleware.NewAppError(
			middleware.ErrorCodeInternalServer,
			"Failed to update location",
			err,
			http.StatusInternalServerError,
		))
		return
	}

	c.JSON(http.StatusOK, existingLoc)
}

// DeleteLocation godoc
//
//	@Summary		Delete a location
//	@Description	Soft-delete a location by ID (marks as deleted but keeps in database)
//	@Tags			locations
//	@Accept			json
//	@Produce		json
//	@Param			id	path		string	true	"Location ID"	example("6824886e6b180b753cea43e9")
//	@Success		200	{object}	models.SwaggerStandardResponse{data=object}
//	@Failure		400	{object}	models.SwaggerErrorResponse
//	@Failure		404	{object}	models.SwaggerErrorResponse
//	@Failure		500	{object}	models.SwaggerErrorResponse
//	@Router			/locations/{id} [delete]
func (h *LocationHandler) DeleteLocation(c *gin.Context) {
	id := c.Param("id")

	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		middleware.HandleError(c, middleware.NewAppError(
			middleware.ErrorCodeInvalidRequest,
			"Invalid location ID",
			err,
			http.StatusBadRequest,
		))
		return
	}

	// Get the location before deletion to return its details if needed
	location, err := h.GetLocationUseCase.Execute(c.Request.Context(), objectID)
	if err != nil {
		middleware.HandleError(c, middleware.NewAppError(
			middleware.ErrorCodeInternalServer,
			"Failed to fetch location",
			err,
			http.StatusInternalServerError,
		))
		return
	}

	if location == nil {
		middleware.HandleError(c, middleware.NewAppError(
			middleware.ErrorCodeNotFound,
			"Location not found",
			nil,
			http.StatusNotFound,
		))
		return
	}

	// Perform soft delete
	deleted, err := h.DeleteLocationUseCase.Execute(c.Request.Context(), objectID)
	if err != nil {
		middleware.HandleError(c, middleware.NewAppError(
			middleware.ErrorCodeInternalServer,
			"Failed to delete location",
			err,
			http.StatusInternalServerError,
		))
		return
	}

	if !deleted {
		middleware.HandleError(c, middleware.NewAppError(
			middleware.ErrorCodeNotFound,
			"Location not found or already deleted",
			nil,
			http.StatusNotFound,
		))
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Location deleted successfully",
		"id":      id,
	})
}

// RestoreLocation godoc
//
//	@Summary		Restore a deleted location
//	@Description	Restore a soft-deleted location by ID (clears the deletedAt timestamp)
//	@Tags			locations
//	@Accept			json
//	@Produce		json
//	@Param			id	path		string	true	"Location ID"	example("6824886e6b180b753cea43e9")
//	@Success		200	{object}	models.SwaggerStandardResponse{data=entity.Location}
//	@Failure		400	{object}	models.SwaggerErrorResponse
//	@Failure		404	{object}	models.SwaggerErrorResponse
//	@Failure		500	{object}	models.SwaggerErrorResponse
//	@Router			/locations/{id}/restore [post]
func (h *LocationHandler) RestoreLocation(c *gin.Context) {
	id := c.Param("id")

	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		middleware.HandleError(c, middleware.NewAppError(
			middleware.ErrorCodeInvalidRequest,
			"Invalid location ID",
			err,
			http.StatusBadRequest,
		))
		return
	}

	// Call use case to restore the location
	restored, err := h.RestoreLocationUseCase.Execute(c.Request.Context(), objectID)
	if err != nil {
		middleware.HandleError(c, middleware.NewAppError(
			middleware.ErrorCodeInternalServer,
			"Failed to restore location",
			err,
			http.StatusInternalServerError,
		))
		return
	}

	if !restored {
		middleware.HandleError(c, middleware.NewAppError(
			middleware.ErrorCodeNotFound,
			"Location not found or not deleted",
			nil,
			http.StatusNotFound,
		))
		return
	}

	// Get the restored location to return its details
	location, err := h.GetLocationUseCase.Execute(c.Request.Context(), objectID)
	if err != nil {
		middleware.HandleError(c, middleware.NewAppError(
			middleware.ErrorCodeInternalServer,
			"Location restored but failed to fetch updated data",
			err,
			http.StatusInternalServerError,
		))
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message":  "Location restored successfully",
		"location": location,
	})
}

// HardDeleteLocation godoc
//
//	@Summary		Permanently delete a location
//	@Description	Permanently delete a location by ID (removes it from the database)
//	@Tags			locations
//	@Accept			json
//	@Produce		json
//	@Param			id	path		string	true	"Location ID"	example("6824886e6b180b753cea43e9")
//	@Success		200	{object}	models.SwaggerStandardResponse{data=object}
//	@Failure		400	{object}	models.SwaggerErrorResponse
//	@Failure		404	{object}	models.SwaggerErrorResponse
//	@Failure		500	{object}	models.SwaggerErrorResponse
//	@Router			/locations/{id}/hard-delete [delete]
func (h *LocationHandler) HardDeleteLocation(c *gin.Context) {
	id := c.Param("id")

	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		middleware.HandleError(c, middleware.NewAppError(
			middleware.ErrorCodeInvalidRequest,
			"Invalid location ID",
			err,
			http.StatusBadRequest,
		))
		return
	}

	// Perform the hard delete
	deleted, err := h.HardDeleteLocationUseCase.Execute(c.Request.Context(), objectID)
	if err != nil {
		middleware.HandleError(c, middleware.NewAppError(
			middleware.ErrorCodeInternalServer,
			"Failed to permanently delete location",
			err,
			http.StatusInternalServerError,
		))
		return
	}

	if !deleted {
		middleware.HandleError(c, middleware.NewAppError(
			middleware.ErrorCodeNotFound,
			"Location not found",
			nil,
			http.StatusNotFound,
		))
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Location permanently deleted",
		"id":      id,
	})
}
