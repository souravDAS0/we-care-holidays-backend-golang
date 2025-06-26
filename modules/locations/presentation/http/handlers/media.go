package handlers

import (
	"net/http"
	"path/filepath"
	"strings"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"

	"bitbucket.org/abhishek_fordel/we-care-holidays-backend-golang/internal/middleware"
)

// UploadLocationMedia handles POST /locations/:id/media
func (h *LocationHandler) UploadLocationMedia(c *gin.Context) {
	idStr := c.Param("id")
	oid, err := primitive.ObjectIDFromHex(idStr)
	if err != nil {
		middleware.HandleError(c, middleware.NewAppError(
			middleware.ErrorCodeInvalidRequest,
			"Invalid location ID",
			err,
			http.StatusBadRequest,
		))
		return
	}

	// Parse the multipart form
	if err := c.Request.ParseMultipartForm(32 << 20); err != nil {
		middleware.HandleError(c, middleware.NewAppError(
			middleware.ErrorCodeInvalidRequest,
			"Failed to parse multipart form",
			err,
			http.StatusBadRequest,
		))
		return
	}

	form := c.Request.MultipartForm
	files := form.File["file"]
	if len(files) == 0 {
		middleware.HandleError(c, middleware.NewAppError(
			middleware.ErrorCodeInvalidRequest,
			"No media files provided",
			nil,
			http.StatusBadRequest,
		))
		return
	}

	var photoURLs []string
	var videoURLs []string

	allowed := map[string]string{
		".jpg": "image/jpeg", ".jpeg": "image/jpeg", ".png": "image/png",
		".svg": "image/svg+xml", ".webp": "image/webp", ".gif": "image/gif",
		".mp4": "video/mp4", ".mov": "video/quicktime",
	}

	for _, fileHeader := range files {
		if fileHeader.Size > 10*1024*1024 {
			continue // Skip large files (optional: collect skipped errors)
		}

		ext := strings.ToLower(filepath.Ext(fileHeader.Filename))
		contentType, ok := allowed[ext]
		if !ok {
			continue // Skip unsupported types
		}

		file, err := fileHeader.Open()
		if err != nil {
			continue // Could log or collect failed uploads
		}
		defer file.Close()

		url, err := h.fileService.UploadFile(c.Request.Context(), file, fileHeader.Filename, contentType)
		if err != nil {
			continue
		}

		if strings.HasPrefix(contentType, "image/") {
			photoURLs = append(photoURLs, url)
		} else if strings.HasPrefix(contentType, "video/") {
			videoURLs = append(videoURLs, url)
		}
	}

	// If no files uploaded successfully, return error
	if len(photoURLs) == 0 && len(videoURLs) == 0 {
		middleware.HandleError(c, middleware.NewAppError(
			middleware.ErrorCodeInvalidRequest,
			"None of the media files were valid or successfully uploaded",
			nil,
			http.StatusBadRequest,
		))
		return
	}

	// Persist media URLs
	if err := h.UploadLocationMediaUseCase.Execute(c.Request.Context(), oid.Hex(), photoURLs, videoURLs); err != nil {
		middleware.HandleError(c, middleware.NewAppError(
			middleware.ErrorCodeInternalServer,
			"Failed to save media URLs",
			err,
			http.StatusInternalServerError,
		))
		return
	}

	// Return updated location
	location, err := h.GetLocationUseCase.Execute(c.Request.Context(), oid)
	if err != nil || location == nil {
		middleware.HandleError(c, middleware.NewAppError(
			middleware.ErrorCodeInternalServer,
			"Uploaded but failed to fetch updated location",
			err,
			http.StatusInternalServerError,
		))
		return
	}

	c.JSON(http.StatusOK, location)
}
