package handlers

import (
	"net/http"

	"bitbucket.org/abhishek_fordel/we-care-holidays-backend-golang/internal/middleware"
	"bitbucket.org/abhishek_fordel/we-care-holidays-backend-golang/modules/organizations/presentation/http/dto"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// GetOrganization godoc
//
//	@Summary		Get organization by ID
//	@Description	Get single organization by its ID
//	@Tags			organizations
//	@Accept			json
//	@Produce		json
//	@Param			id	path		string	true	"Organization ID"	example("6824886e6b180b753cea43e9")
//	@Success		200	{object}	models.SwaggerStandardResponse{data=entity.Organization}
//	@Failure		400	{object}	models.SwaggerErrorResponse
//	@Failure		404	{object}	models.SwaggerErrorResponse
//	@Failure		500	{object}	models.SwaggerErrorResponse
//	@Router			/organizations/{id} [get]
func (h *OrganizationHandler) GetOrganization(c *gin.Context) {
	id := c.Param("id")

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

	if organization == nil {
		middleware.HandleError(c, middleware.NewAppError(
			middleware.ErrorCodeNotFound,
			"Organization not found",
			nil,
			http.StatusNotFound,
		))
		return
	}

	c.JSON(http.StatusOK, organization)
}

// CreateOrganization godoc
//
//	@Summary		Create a new organization
//	@Description	Create a new organization with the provided data
//	@Tags			organizations
//	@Accept			json
//	@Produce		json
//	@Param			organization	body		dto.CreateOrganizationDto	true	"Organization data"
//	@Success		201				{object}	models.SwaggerStandardResponse{data=entity.Organization}
//	@Failure		400				{object}	models.SwaggerErrorResponse
//	@Failure		422				{object}	models.SwaggerErrorResponse
//	@Failure		500				{object}	models.SwaggerErrorResponse
//	@Router			/organizations [post]
func (h *OrganizationHandler) CreateOrganization(c *gin.Context) {
	var createDto dto.CreateOrganizationDto

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
	organization := createDto.ToEntity()

	// Call use case to create
	if err := h.CreateOrganizationUseCase.Execute(c.Request.Context(), organization); err != nil {
		// Handle specific error cases
		if err.Error() == "organization with this slug already exists" {
			middleware.HandleError(c, middleware.NewAppError(
				middleware.ErrorCodeConflict,
				"Organization with this slug already exists",
				nil,
				http.StatusConflict,
			))
			return
		}

		middleware.HandleError(c, middleware.NewAppError(
			middleware.ErrorCodeInternalServer,
			"Failed to create organization",
			err,
			http.StatusInternalServerError,
		))
		return
	}

	c.JSON(http.StatusCreated, organization)
}

// UpdateOrganization godoc
//
//	@Summary		Update an organization
//	@Description	Update an existing organization by ID with partial data
//	@Tags			organizations
//	@Accept			json
//	@Produce		json
//	@Param			id				path		string						true	"Organization ID"
//	@Param			organization	body		dto.UpdateOrganizationDto	true	"Organization data to update"
//	@Success		200				{object}	models.SwaggerStandardResponse{data=entity.Organization}
//	@Failure		400				{object}	models.SwaggerErrorResponse
//	@Failure		404				{object}	models.SwaggerErrorResponse
//	@Failure		500				{object}	models.SwaggerErrorResponse
//	@Router			/organizations/{id} [put]
func (h *OrganizationHandler) UpdateOrganization(c *gin.Context) {
	id := c.Param("id")

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

	// Fetch existing organization to apply updates to
	existingOrg, err := h.GetOrganizationUseCase.Execute(c.Request.Context(), objectId)
	if err != nil {
		middleware.HandleError(c, middleware.NewAppError(
			middleware.ErrorCodeInternalServer,
			"Failed to fetch organization",
			err,
			http.StatusInternalServerError,
		))
		return
	}

	if existingOrg == nil {
		middleware.HandleError(c, middleware.NewAppError(
			middleware.ErrorCodeNotFound,
			"Organization not found",
			nil,
			http.StatusNotFound,
		))
		return
	}

	// Parse update DTO
	var updateDto dto.UpdateOrganizationDto
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

	// Apply updates to the existing organization
	updateDto.ApplyUpdates(existingOrg)

	// Call use case to update
	if err := h.UpdateOrganizationUseCase.Execute(c.Request.Context(), existingOrg); err != nil {
		middleware.HandleError(c, middleware.NewAppError(
			middleware.ErrorCodeInternalServer,
			"Failed to update organization",
			err,
			http.StatusInternalServerError,
		))
		return
	}

	c.JSON(http.StatusOK, existingOrg)
}

// DeleteOrganization godoc
//
//	@Summary		Delete an organization
//	@Description	Soft-delete an organization by ID (marks as deleted but keeps in database)
//	@Tags			organizations
//	@Accept			json
//	@Produce		json
//	@Param			id	path		string	true	"Organization ID"	example("6824886e6b180b753cea43e9")
//	@Success		200	{object}	models.SwaggerStandardResponse{data=object}
//	@Failure		400	{object}	models.SwaggerErrorResponse
//	@Failure		404	{object}	models.SwaggerErrorResponse
//	@Failure		500	{object}	models.SwaggerErrorResponse
//	@Router			/organizations/{id} [delete]
func (h *OrganizationHandler) DeleteOrganization(c *gin.Context) {
	id := c.Param("id")

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

	// Get the organization before deletion to return its details if needed
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

	// Perform soft delete instead of hard delete
	deleted, err := h.SoftDeleteOrganizationUseCase.Execute(c.Request.Context(), objectID)
	if err != nil {
		middleware.HandleError(c, middleware.NewAppError(
			middleware.ErrorCodeInternalServer,
			"Failed to delete organization",
			err,
			http.StatusInternalServerError,
		))
		return
	}

	if !deleted {
		middleware.HandleError(c, middleware.NewAppError(
			middleware.ErrorCodeNotFound,
			"Organization not found or already deleted",
			nil,
			http.StatusNotFound,
		))
		return
	}

	// Return success response with the deleted organization ID
	c.JSON(http.StatusOK, gin.H{
		"message": "Organization deleted successfully",
		"id":      id,
	})
}

// RestoreOrganization godoc
//
//	@Summary		Restore a deleted organization
//	@Description	Restore a soft-deleted organization
//	@Tags			organizations
//	@Accept			json
//	@Produce		json
//	@Param			id	path		string	true	"Organization ID"	example("6824886e6b180b753cea43e9")
//	@Success		200	{object}	models.SwaggerStandardResponse{data=entity.Organization}
//	@Failure		400	{object}	models.SwaggerErrorResponse
//	@Failure		404	{object}	models.SwaggerErrorResponse
//	@Failure		500	{object}	models.SwaggerErrorResponse
//	@Router			/organizations/{id}/restore [post]
func (h *OrganizationHandler) RestoreOrganization(c *gin.Context) {
	id := c.Param("id")

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

	restored, err := h.RestoreOrganizationUseCase.Execute(c.Request.Context(), objectID)
	if err != nil {
		middleware.HandleError(c, middleware.NewAppError(
			middleware.ErrorCodeInternalServer,
			"Failed to restore organization",
			err,
			http.StatusInternalServerError,
		))
		return
	}

	if !restored {
		middleware.HandleError(c, middleware.NewAppError(
			middleware.ErrorCodeNotFound,
			"Organization not found or not deleted",
			nil,
			http.StatusNotFound,
		))
		return
	}

	// Get the restored organization to return its details
	organization, err := h.GetOrganizationUseCase.Execute(c.Request.Context(), objectID)
	if err != nil {
		middleware.HandleError(c, middleware.NewAppError(
			middleware.ErrorCodeInternalServer,
			"Organization restored but failed to fetch updated data",
			err,
			http.StatusInternalServerError,
		))
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message":      "Organization restored successfully",
		"organization": organization,
	})
}

// HardDeleteOrganization godoc
//
//	@Summary		Permanently delete an organization
//	@Description	Hard-delete an organization by ID (permanently removes from database)
//	@Tags			organizations
//	@Accept			json
//	@Produce		json
//	@Param			id	path		string	true	"Organization ID"	example("6824886e6b180b753cea43e9")
//	@Success		200	{object}	models.SwaggerStandardResponse{data=object}
//	@Failure		400	{object}	models.SwaggerErrorResponse
//	@Failure		404	{object}	models.SwaggerErrorResponse
//	@Failure		500	{object}	models.SwaggerErrorResponse
//	@Router			/organizations/{id}/hard-delete [delete]
func (h *OrganizationHandler) HardDeleteOrganization(c *gin.Context) {
	id := c.Param("id")

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

	deleted, err := h.HardDeleteOrganizationUseCase.Execute(c.Request.Context(), objectID)
	if err != nil {
		middleware.HandleError(c, middleware.NewAppError(
			middleware.ErrorCodeInternalServer,
			"Failed to permanently delete organization",
			err,
			http.StatusInternalServerError,
		))
		return
	}

	if !deleted {
		middleware.HandleError(c, middleware.NewAppError(
			middleware.ErrorCodeNotFound,
			"Organization not found",
			nil,
			http.StatusNotFound,
		))
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Organization permanently deleted",
		"id":      id,
	})
}
