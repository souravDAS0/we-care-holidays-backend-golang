package handlers

import (
	"net/http"

	"bitbucket.org/abhishek_fordel/we-care-holidays-backend-golang/internal/middleware"
	"bitbucket.org/abhishek_fordel/we-care-holidays-backend-golang/modules/users/presentation/http/dto"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// CreateUser godoc
//
//	@Summary		Create a new user
//	@Description	Create a new user with the provided data
//	@Tags			users
//	@Accept			json
//	@Produce		json
//	@Param			user	body		dto.CreateUserDto	true	"User data"
//	@Success		201		{object}	models.SwaggerStandardResponse{data=entity.User}
//	@Failure		400		{object}	models.SwaggerErrorResponse
//	@Failure		409		{object}	models.SwaggerErrorResponse
//	@Failure		422		{object}	models.SwaggerErrorResponse
//	@Failure		500		{object}	models.SwaggerErrorResponse
//	@Router			/users [post]
func (h *UserHandler) CreateUser(c *gin.Context) {
	var createDto dto.CreateUserDto

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
	user := createDto.ToEntity()

	// Call use case to create
	if err := h.CreateUserUseCase.Execute(c.Request.Context(), user); err != nil {
		// Handle specific error cases
		if err.Error() == "user with this email already exists" {
			middleware.HandleError(c, middleware.NewAppError(
				middleware.ErrorCodeConflict,
				"User with this email already exists",
				nil,
				http.StatusConflict,
			))
			return
		} else if err.Error() == "user with this phone already exists" {
			middleware.HandleError(c, middleware.NewAppError(
				middleware.ErrorCodeConflict,
				"User with this phone already exists",
				nil,
				http.StatusConflict,
			))
			return
		}

		middleware.HandleError(c, middleware.NewAppError(
			middleware.ErrorCodeInternalServer,
			"Failed to create user",
			err,
			http.StatusInternalServerError,
		))
		return
	}

	c.JSON(http.StatusCreated, user)
}

// GetUser godoc
//
//	@Summary		Get user by ID
//	@Description	Get single user by its ID
//	@Tags			users
//	@Accept			json
//	@Produce		json
//	@Param			id	path		string	true	"User ID"	example("507f1f77bcf86cd799439011")
//	@Success		200	{object}	models.SwaggerStandardResponse{data=entity.User}
//	@Failure		400	{object}	models.SwaggerErrorResponse
//	@Failure		404	{object}	models.SwaggerErrorResponse
//	@Failure		500	{object}	models.SwaggerErrorResponse
//	@Router			/users/{id} [get]
func (h *UserHandler) GetUser(c *gin.Context) {
	id := c.Param("id")

	objectId, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		middleware.HandleError(c, middleware.NewAppError(
			middleware.ErrorCodeInvalidRequest,
			"Invalid user ID",
			err,
			http.StatusBadRequest,
		))
		return
	}

	user, err := h.GetUserUseCase.Execute(c.Request.Context(), objectId)
	if err != nil {
		middleware.HandleError(c, middleware.NewAppError(
			middleware.ErrorCodeInternalServer,
			"Failed to fetch user",
			err,
			http.StatusInternalServerError,
		))
		return
	}

	if user == nil {
		middleware.HandleError(c, middleware.NewAppError(
			middleware.ErrorCodeNotFound,
			"User not found",
			nil,
			http.StatusNotFound,
		))
		return
	}

	c.JSON(http.StatusOK, user)
}

// UpdateUser godoc
//
//	@Summary		Update a user
//	@Description	Update an existing user by ID with partial data
//	@Tags			users
//	@Accept			json
//	@Produce		json
//	@Param			id		path		string				true	"User ID"
//	@Param			user	body		dto.UpdateUserDto	true	"User data to update"
//	@Success		200		{object}	models.SwaggerStandardResponse{data=entity.User}
//	@Failure		400		{object}	models.SwaggerErrorResponse
//	@Failure		404		{object}	models.SwaggerErrorResponse
//	@Failure		500		{object}	models.SwaggerErrorResponse
//	@Router			/users/{id} [put]
func (h *UserHandler) UpdateUser(c *gin.Context) {
	id := c.Param("id")

	objectId, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		middleware.HandleError(c, middleware.NewAppError(
			middleware.ErrorCodeInvalidRequest,
			"Invalid user ID",
			err,
			http.StatusBadRequest,
		))
		return
	}

	// Fetch existing user to apply updates to
	existingUser, err := h.GetUserUseCase.Execute(c.Request.Context(), objectId)
	if err != nil {
		middleware.HandleError(c, middleware.NewAppError(
			middleware.ErrorCodeInternalServer,
			"Failed to fetch user",
			err,
			http.StatusInternalServerError,
		))
		return
	}

	if existingUser == nil {
		middleware.HandleError(c, middleware.NewAppError(
			middleware.ErrorCodeNotFound,
			"User not found",
			nil,
			http.StatusNotFound,
		))
		return
	}

	// Parse update DTO
	var updateDto dto.UpdateUserDto
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

	// Apply updates to the existing user
	updateDto.ApplyUpdates(existingUser)

	// Call use case to update
	if err := h.UpdateUserUseCase.Execute(c.Request.Context(), existingUser); err != nil {
		middleware.HandleError(c, middleware.NewAppError(
			middleware.ErrorCodeInternalServer,
			"Failed to update user",
			err,
			http.StatusInternalServerError,
		))
		return
	}

	c.JSON(http.StatusOK, existingUser)
}

// DeleteUser godoc
//
//	@Summary		Delete a user
//	@Description	Soft-delete a user by ID (marks as deleted but keeps in database)
//	@Tags			users
//	@Accept			json
//	@Produce		json
//	@Param			id	path		string	true	"User ID"	example("507f1f77bcf86cd799439011")
//	@Success		200	{object}	models.SwaggerStandardResponse{data=object}
//	@Failure		400	{object}	models.SwaggerErrorResponse
//	@Failure		404	{object}	models.SwaggerErrorResponse
//	@Failure		500	{object}	models.SwaggerErrorResponse
//	@Router			/users/{id} [delete]
func (h *UserHandler) DeleteUser(c *gin.Context) {
	id := c.Param("id")

	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		middleware.HandleError(c, middleware.NewAppError(
			middleware.ErrorCodeInvalidRequest,
			"Invalid user ID",
			err,
			http.StatusBadRequest,
		))
		return
	}

	// Get the user before deletion to return its details if needed
	user, err := h.GetUserUseCase.Execute(c.Request.Context(), objectID)
	if err != nil {
		middleware.HandleError(c, middleware.NewAppError(
			middleware.ErrorCodeInternalServer,
			"Failed to fetch user",
			err,
			http.StatusInternalServerError,
		))
		return
	}

	if user == nil {
		middleware.HandleError(c, middleware.NewAppError(
			middleware.ErrorCodeNotFound,
			"User not found",
			nil,
			http.StatusNotFound,
		))
		return
	}

	// Perform soft delete
	deleted, err := h.SoftDeleteUserUseCase.Execute(c.Request.Context(), objectID)
	if err != nil {
		middleware.HandleError(c, middleware.NewAppError(
			middleware.ErrorCodeInternalServer,
			"Failed to delete user",
			err,
			http.StatusInternalServerError,
		))
		return
	}

	if !deleted {
		middleware.HandleError(c, middleware.NewAppError(
			middleware.ErrorCodeNotFound,
			"User not found or already deleted",
			nil,
			http.StatusNotFound,
		))
		return
	}

	// Return success response with the deleted user ID
	c.JSON(http.StatusOK, gin.H{
		"message": "User deleted successfully",
		"id":      id,
	})
}

// RestoreUser godoc
//
//	@Summary		Restore a deleted user
//	@Description	Restore a soft-deleted user
//	@Tags			users
//	@Accept			json
//	@Produce		json
//	@Param			id	path		string	true	"User ID"	example("507f1f77bcf86cd799439011")
//	@Success		200	{object}	models.SwaggerStandardResponse{data=entity.User}
//	@Failure		400	{object}	models.SwaggerErrorResponse
//	@Failure		404	{object}	models.SwaggerErrorResponse
//	@Failure		500	{object}	models.SwaggerErrorResponse
//	@Router			/users/{id}/restore [post]
func (h *UserHandler) RestoreUser(c *gin.Context) {
	id := c.Param("id")

	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		middleware.HandleError(c, middleware.NewAppError(
			middleware.ErrorCodeInvalidRequest,
			"Invalid user ID",
			err,
			http.StatusBadRequest,
		))
		return
	}

	restored, err := h.RestoreUserUseCase.Execute(c.Request.Context(), objectID)
	if err != nil {
		middleware.HandleError(c, middleware.NewAppError(
			middleware.ErrorCodeInternalServer,
			"Failed to restore user",
			err,
			http.StatusInternalServerError,
		))
		return
	}

	if !restored {
		middleware.HandleError(c, middleware.NewAppError(
			middleware.ErrorCodeNotFound,
			"User not found or not deleted",
			nil,
			http.StatusNotFound,
		))
		return
	}

	// Get the restored user to return its details
	user, err := h.GetUserUseCase.Execute(c.Request.Context(), objectID)
	if err != nil {
		middleware.HandleError(c, middleware.NewAppError(
			middleware.ErrorCodeInternalServer,
			"User restored but failed to fetch updated data",
			err,
			http.StatusInternalServerError,
		))
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "User restored successfully",
		"user":    user,
	})
}

// HardDeleteUser godoc
//
//	@Summary		Permanently delete a user
//	@Description	Hard-delete a user by ID (permanently removes from database)
//	@Tags			users
//	@Accept			json
//	@Produce		json
//	@Param			id	path		string	true	"User ID"	example("507f1f77bcf86cd799439011")
//	@Success		200	{object}	models.SwaggerStandardResponse{data=object}
//	@Failure		400	{object}	models.SwaggerErrorResponse
//	@Failure		404	{object}	models.SwaggerErrorResponse
//	@Failure		500	{object}	models.SwaggerErrorResponse
//	@Router			/users/{id}/hard-delete [delete]
func (h *UserHandler) HardDeleteUser(c *gin.Context) {
	id := c.Param("id")

	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		middleware.HandleError(c, middleware.NewAppError(
			middleware.ErrorCodeInvalidRequest,
			"Invalid user ID",
			err,
			http.StatusBadRequest,
		))
		return
	}

	deleted, err := h.HardDeleteUserUseCase.Execute(c.Request.Context(), objectID)
	if err != nil {
		middleware.HandleError(c, middleware.NewAppError(
			middleware.ErrorCodeInternalServer,
			"Failed to permanently delete user",
			err,
			http.StatusInternalServerError,
		))
		return
	}

	if !deleted {
		middleware.HandleError(c, middleware.NewAppError(
			middleware.ErrorCodeNotFound,
			"User not found",
			nil,
			http.StatusNotFound,
		))
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "User permanently deleted",
		"id":      id,
	})
}
