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

// UploadOrganizationLogo godoc
//
//	@Summary		Upload organization logo
//	@Description	Upload a new logo for an organization and update its record
//	@Tags			organizations
//	@Accept			multipart/form-data
//	@Produce		json
//	@Param			id		path		string															true	"Organization ID"	example("6824886e6b180b753cea43e9")
//	@Param			file	formData	file															true	"Logo image file"
//	@Success		200		{object}	models.SwaggerStandardResponse{data=entity.Organization}	"Updated organization"
//	@Failure		400		{object}	models.SwaggerErrorResponse									"Invalid request"
//	@Failure		404		{object}	models.SwaggerErrorResponse									"Organization not found"
//	@Failure		413		{object}	models.SwaggerErrorResponse									"File too large"
//	@Failure		500		{object}	models.SwaggerErrorResponse									"Server error"
//	@Router			/organizations/{id}/logo [post]
func (h *OrganizationHandler) UploadOrganizationLogo(c *gin.Context) {
	id := c.Param("id")

	// Convert ID to ObjectID
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		middleware.HandleError(c, middleware.NewAppError(
			middleware.ErrorCodeInvalidRequest,
			"Invalid organization ID",
			err,
			http.StatusBadRequest,
		))
		return
	}

	// Get the organization to ensure it exists
	organization, err := h.GetOrganizationUseCase.Execute(c.Request.Context(), objectID)
	if err != nil {
		middleware.HandleError(c, middleware.NewAppError(
			middleware.ErrorCodeInternalServer,
			"Failed to fetch organization",
			err,
			http.StatusInternalServerError,
		))
		return
	}

	if organization == nil {
		middleware.HandleError(c, middleware.NewAppError(
			middleware.ErrorCodeNotFound,
			"Organization not found",
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
	logoURL, err := h.fileService.UploadFile(c.Request.Context(), file, header.Filename, contentType)
	if err != nil {
		middleware.HandleError(c, middleware.NewAppError(
			middleware.ErrorCodeInternalServer,
			"Failed to upload logo",
			err,
			http.StatusInternalServerError,
		))
		return
	}

	// Delete old logo if it exists
	if organization.Logo != "" {
		// Try to delete the old logo, but don't fail if it doesn't work
		_ = h.fileService.DeleteFile(c.Request.Context(), organization.Logo)
	}

	// Update organization with new logo URL
	organization.Logo = logoURL
	organization.UpdatedAt = time.Now()

	// Save updated organization
	if err := h.UpdateOrganizationUseCase.Execute(c.Request.Context(), organization); err != nil {
		middleware.HandleError(c, middleware.NewAppError(
			middleware.ErrorCodeInternalServer,
			"Failed to update organization with new logo",
			err,
			http.StatusInternalServerError,
		))
		return
	}

	// Return updated organization
	c.JSON(http.StatusOK, organization)
}
