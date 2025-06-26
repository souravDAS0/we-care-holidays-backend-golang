package handlers

import (
	"net/http"
	"path/filepath"
	"strings"
	"time"

	"bitbucket.org/abhishek_fordel/we-care-holidays-backend-golang/internal/middleware"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// UploadUserProfilePhoto godoc
//
//	@Summary		Upload user profile photo
//	@Description	Upload a new profile photo for an user and update its record
//	@Tags			users
//	@Accept			multipart/form-data
//	@Produce		json
//	@Param			id		path		string															true	"User ID"	example("6824886e6b180b753cea43e9")
//	@Param			file	formData	file															true	"Logo image file"
//	@Success		200		{object}	models.SwaggerStandardResponse{data=entity.User}	"Updated user"
//	@Failure		400		{object}	models.SwaggerErrorResponse									"Invalid request"
//	@Failure		404		{object}	models.SwaggerErrorResponse									"User not found"
//	@Failure		413		{object}	models.SwaggerErrorResponse									"File too large"
//	@Failure		500		{object}	models.SwaggerErrorResponse									"Server error"
//	@Router			/users/{id}/profile-photo [post]
func (h *UserHandler) UploadUserProfilePhoto(c *gin.Context) {
	id := c.Param("id")

	// Convert ID to ObjectID
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

	// Get the user to ensure it exists
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

	// Get file from request
	file, header, err := c.Request.FormFile("file")
	if err != nil {
		middleware.HandleError(c, middleware.NewAppError(
			middleware.ErrorCodeInvalidRequest,
			"No file provided or invalid file",
			err,
			http.StatusBadRequest,
		))
		return
	}
	defer file.Close()

	// Check file size (limit to 5MB)
	if header.Size > 5*1024*1024 {
		middleware.HandleError(c, middleware.NewAppError(
			middleware.ErrorCodeInvalidRequest,
			"File too large, maximum size is 5MB",
			nil,
			http.StatusRequestEntityTooLarge,
		))
		return
	}

	// Check file type (allow only images)
	ext := strings.ToLower(filepath.Ext(header.Filename))
	validExtensions := map[string]bool{
		".jpg":  true,
		".jpeg": true,
		".png":  true,
		".gif":  true,
		".svg":  true,
		".webp": true,
	}

	if !validExtensions[ext] {
		middleware.HandleError(c, middleware.NewAppError(
			middleware.ErrorCodeInvalidRequest,
			"Invalid file type, only images are allowed",
			nil,
			http.StatusBadRequest,
		))
		return
	}

	// Determine content type
	contentType := ""
	switch ext {
	case ".jpg", ".jpeg":
		contentType = "image/jpeg"
	case ".png":
		contentType = "image/png"
	case ".gif":
		contentType = "image/gif"
	case ".svg":
		contentType = "image/svg+xml"
	case ".webp":
		contentType = "image/webp"
	}

	// Upload file to S3
	profilePhotoURL, err := h.fileService.UploadFile(c.Request.Context(), file, header.Filename, contentType)
	if err != nil {
		middleware.HandleError(c, middleware.NewAppError(
			middleware.ErrorCodeInternalServer,
			"Failed to upload profile photo",
			err,
			http.StatusInternalServerError,
		))
		return
	}

	// Delete old profile photo if it exists
	if user.ProfilePhotoURL != "" {
		// Try to delete the old profile photo, but don't fail if it doesn't work
		_ = h.fileService.DeleteFile(c.Request.Context(), user.ProfilePhotoURL)
	}

	// Update user with new profile photo URL
	user.ProfilePhotoURL = profilePhotoURL
	user.UpdatedAt = time.Now()

	// Save updated user
	if err := h.UpdateUserUseCase.Execute(c.Request.Context(), user); err != nil {
		middleware.HandleError(c, middleware.NewAppError(
			middleware.ErrorCodeInternalServer,
			"Failed to update user with new profile photo",
			err,
			http.StatusInternalServerError,
		))
		return
	}

	// Return updated user
	c.JSON(http.StatusOK, user)
}
