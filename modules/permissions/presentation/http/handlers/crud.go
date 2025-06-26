package handlers

import (
	"net/http"

	"bitbucket.org/abhishek_fordel/we-care-holidays-backend-golang/internal/middleware"
	"bitbucket.org/abhishek_fordel/we-care-holidays-backend-golang/modules/permissions/presentation/http/dto"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// GetPermission godoc
//
//	@Summary		Get a permission by ID
//	@Description	Retrieves a single permission by its unique identifier
//	@Tags			permissions
//	@Accept			json
//	@Produce		json
//	@Param			id	path		string	true	"Permission ID"	format(objectid)	example("6824886e6b180b753cea43e9")
//	@Success		200	{object}	models.SwaggerStandardResponse{data=entity.Permission}
//	@Failure		400	{object}	models.SwaggerErrorResponse
//	@Failure		404	{object}	models.SwaggerErrorResponse
//	@Failure		500	{object}	models.SwaggerErrorResponse
//	@Router			/permissions/{id} [get]
func (h *PermissionHandler) GetPermission(c *gin.Context) {
	id := c.Param("id")

	objectId, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		middleware.HandleError(c, middleware.NewAppError(
			middleware.ErrorCodeInvalidRequest,
			"Invalid Permission ID",
			err,
			http.StatusBadRequest,
		))
		return
	}

	Permission, err := h.GetPermissionUseCase.Execute(c.Request.Context(), objectId)
	if err != nil {
		middleware.HandleError(c, middleware.NewAppError(
			middleware.ErrorCodeInternalServer,
			"Failed to fetch permission",
			err,
			http.StatusInternalServerError,
		))
		return
	}

	if Permission == nil {
		middleware.HandleError(c, middleware.NewAppError(
			middleware.ErrorCodeNotFound,
			"permission not found",
			nil,
			http.StatusNotFound,
		))
		return
	}

	c.JSON(http.StatusOK, Permission)
}

// CreatePermission godoc
//
//	@Summary		Create a new permission
//	@Description	Creates a new permission with the provided details
//	@Tags			permissions
//	@Accept			json
//	@Produce		json
//	@Param			permission	body	dto.CreatePermissionDto	true	"Permission creation data"
//	@Success		201			{object}	models.SwaggerStandardResponse{data=entity.Permission}
//	@Failure		400			{object}	models.SwaggerErrorResponse
//	@Failure		422			{object}	models.SwaggerErrorResponse
//	@Failure		500			{object}	models.SwaggerErrorResponse
//	@Router		/permissions [post]
func (h *PermissionHandler) CreatePermission(c *gin.Context) {
	var createDto dto.CreatePermissionDto

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
	permission := createDto.ToEntity()

	// Call use case to create
	if err := h.CreatePermissionUseCase.Execute(c.Request.Context(), permission); err != nil {

		middleware.HandleError(c, middleware.NewAppError(
			middleware.ErrorCodeInternalServer,
			"Failed to create permission",
			err,
			http.StatusInternalServerError,
		))
		return
	}

	c.JSON(http.StatusCreated, permission)
}

// UpdatePermission godoc
//
//	@Summary		Update an existing permission
//	@Description	Updates an existing permission with the provided details
//	@Tags			permissions
//	@Accept			json
//	@Produce		json
//	@Param			id			path		string					true	"Permission ID"	format(objectid)
//	@Param			permission	body		dto.UpdatePermissionDto	true	"Permission update data"
//	@Success		200			{object}	models.SwaggerStandardResponse{data=entity.Permission}
//	@Failure		400			{object}	models.SwaggerErrorResponse
//	@Failure		404			{object}	models.SwaggerErrorResponse
//	@Failure		500			{object}	models.SwaggerErrorResponse
//	@Router			/permissions/{id} [put]
func (h *PermissionHandler) UpdatePermission(c *gin.Context) {
	id := c.Param("id")

	objectId, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		middleware.HandleError(c, middleware.NewAppError(
			middleware.ErrorCodeInvalidRequest,
			"Invalid permission ID",
			err,
			http.StatusBadRequest,
		))
		return
	}

	// Fetch existing permission to apply updates to
	existingPermission, err := h.GetPermissionUseCase.Execute(c.Request.Context(), objectId)
	if err != nil {
		middleware.HandleError(c, middleware.NewAppError(
			middleware.ErrorCodeInternalServer,
			"Failed to fetch permission",
			err,
			http.StatusInternalServerError,
		))
		return
	}

	if existingPermission == nil {
		middleware.HandleError(c, middleware.NewAppError(
			middleware.ErrorCodeNotFound,
			"Permission not found",
			nil,
			http.StatusNotFound,
		))
		return
	}

	// Parse update DTO
	var updateDto dto.UpdatePermissionDto
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

	// Apply updates to the existing permission
	updateDto.ApplyUpdates(existingPermission)

	// Call use case to update
	if err := h.UpdatePermissionUseCase.Execute(c.Request.Context(), existingPermission); err != nil {
		middleware.HandleError(c, middleware.NewAppError(
			middleware.ErrorCodeInternalServer,
			"Failed to update permission",
			err,
			http.StatusInternalServerError,
		))
		return
	}

	c.JSON(http.StatusOK, existingPermission)
}

// HardDeletePermission godoc
//
//	@Summary		Permanently delete a permission
//	@Description	Permanently deletes a permission from the database (cannot be undone)
//	@Tags			permissions
//	@Accept			json
//	@Produce		json
//	@Param			id	path		string					true	"Permission ID"	format(objectid)
//	@Success		200	{object}	models.SwaggerStandardResponse{data=object}
//	@Failure		400	{object}	models.SwaggerErrorResponse
//	@Failure		404	{object}	models.SwaggerErrorResponse
//	@Failure		500	{object}	models.SwaggerErrorResponse
//	@Router			/permissions/{id}/hard-delete [delete]
func (h *PermissionHandler) HardDeletePermission(c *gin.Context) {
	id := c.Param("id")

	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		middleware.HandleError(c, middleware.NewAppError(
			middleware.ErrorCodeInvalidRequest,
			"Invalid permission ID",
			err,
			http.StatusBadRequest,
		))
		return
	}

	deleted, err := h.HardDeletePermissionUseCase.Execute(c.Request.Context(), objectID)
	if err != nil {
		middleware.HandleError(c, middleware.NewAppError(
			middleware.ErrorCodeInternalServer,
			"Failed to permanently delete permission",
			err,
			http.StatusInternalServerError,
		))
		return
	}

	if !deleted {
		middleware.HandleError(c, middleware.NewAppError(
			middleware.ErrorCodeNotFound,
			"Permission not found",
			nil,
			http.StatusNotFound,
		))
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Permission permanently deleted",
		"id":      id,
	})
}
