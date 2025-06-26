package handlers

import (
	"net/http"

	"bitbucket.org/abhishek_fordel/we-care-holidays-backend-golang/internal/middleware"
	"bitbucket.org/abhishek_fordel/we-care-holidays-backend-golang/modules/roles/domain/entity"
	"bitbucket.org/abhishek_fordel/we-care-holidays-backend-golang/modules/roles/presentation/http/dto"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// GetRole godoc
//
//	@Summary		Get role by ID
//	@Description	Get single role by its ID
//	@Tags			roles
//	@Accept			json
//	@Produce		json
//	@Param			id	path		string	true	"role ID"	example("6824886e6b180b753cea43e9")
//	@Success		200	{object}	models.SwaggerStandardResponse{data=entity.Role}
//	@Failure		400	{object}	models.SwaggerErrorResponse
//	@Failure		404	{object}	models.SwaggerErrorResponse
//	@Failure		500	{object}	models.SwaggerErrorResponse
//	@Router			/roles/{id} [get]
func (h *RoleHandler) GetRole(c *gin.Context) {
	id := c.Param("id")

	objectId, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		middleware.HandleError(c, middleware.NewAppError(
			middleware.ErrorCodeInvalidRequest,
			"Invalid Role ID",
			err,
			http.StatusBadRequest,
		))
		return
	}

	role, err := h.GetRoleUseCase.Execute(c.Request.Context(), objectId)
	if err != nil {
		middleware.HandleError(c, middleware.NewAppError(
			middleware.ErrorCodeInternalServer,
			"Failed to fetch role",
			err,
			http.StatusInternalServerError,
		))
		return
	}

	if role == nil {
		middleware.HandleError(c, middleware.NewAppError(
			middleware.ErrorCodeNotFound,
			"Role not found",
			nil,
			http.StatusNotFound,
		))
		return
	}

	c.JSON(http.StatusOK, role)
}

// CreateRole godoc
//
//	@Summary		Create a new role
//	@Description	Create a new role with the provided data
//	@Tags			roles
//	@Accept			json
//	@Produce		json
//	@Param			role	body		dto.CreateRoleDto	true	"Role data"
//	@Success		201			{object}	models.SwaggerStandardResponse{data=entity.Role}
//	@Failure		400			{object}	models.SwaggerErrorResponse
//	@Failure		422			{object}	models.SwaggerErrorResponse
//	@Failure		500			{object}	models.SwaggerErrorResponse
//	@Router			/roles [post]
func (h *RoleHandler) CreateRole(c *gin.Context) {
	ctx := c.Request.Context()
	authCtx := middleware.GetAuthContext(ctx)

	var createDto dto.CreateRoleDto

	if err := c.ShouldBindJSON(&createDto); err != nil {
		middleware.HandleError(c, middleware.NewAppError(
			middleware.ErrorCodeInvalidRequest,
			"Invalid request body",
			err,
			http.StatusBadRequest,
		))
		return
	}

	// Validate permission strings and convert to IDs
	var permissionIDs []string
	for _, permStr := range createDto.Permissions {
		perm, err := h.permissionValidator.ValidatePermissionString(ctx, permStr)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error":      "Invalid permission: " + permStr,
				"details":    err.Error(),
				"suggestion": "Format should be 'resource:action:scope' (e.g., 'users:read:organization')",
			})
			return
		}

		// Find the permission ID by resource, action, and scope
		existingPerms, _, err := h.ListPermissionsUseCase.Execute(ctx, map[string]interface{}{
			"resource": perm.Resource,
			"action":   string(perm.Action),
		}, 1, 1)
		if err != nil || len(existingPerms) == 0 {
			c.JSON(http.StatusBadRequest, gin.H{
				"error":      "Permission not found: " + permStr,
				"suggestion": "Check available permissions at GET /api/v1/permissions",
			})
			return
		}

		permissionIDs = append(permissionIDs, existingPerms[0].ID.Hex())
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

	// Validate that user can assign the requested permissions
	if err := h.rbacService.ValidateRolePermissions(ctx, authCtx, permissionIDs); err != nil {
		c.JSON(http.StatusForbidden, gin.H{
			"error":      err.Error(),
			"suggestion": "You can only assign permissions that you possess and within your scope level",
		})
		return
	}

	// Check if role name already exists
	existing, _, err := h.ListRolesUseCase.Execute(ctx, map[string]interface{}{
		"name": createDto.Name,
	}, 1, 1)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to check existing roles"})
		return
	}
	if len(existing) > 0 {
		c.JSON(http.StatusConflict, gin.H{
			"error":      "Role name already exists",
			"suggestion": "Choose a different role name",
		})
		return
	}

	// Create the role
	role := &entity.Role{
		ID:          primitive.NewObjectID(),
		Name:        createDto.Name,
		Description: createDto.Description,
		Scope:       entity.RoleScope(createDto.Scope),
		Permissions: permissionIDs,
		IsSystem:    false,
		CreatedBy:   authCtx.UserID.Hex(),
	}

	// Call use case to create
	if err := h.CreateRoleUseCase.Execute(ctx, role); err != nil {

		middleware.HandleError(c, middleware.NewAppError(
			middleware.ErrorCodeInternalServer,
			"Failed to create role",
			err,
			http.StatusInternalServerError,
		))
		return
	}

	c.JSON(http.StatusCreated, role)
}

// UpdateRole godoc
//
//	@Summary		Update an existing role
//	@Description	Update a role by ID with provided data
//	@Tags			roles
//	@Accept			json
//	@Produce		json
//	@Param			id			path		string					true	"Role ID"
//	@Param			role	body		dto.UpdateRoleDto	true	"Role data to update"
//	@Success		200			{object}	models.SwaggerStandardResponse{data=entity.Role}
//	@Failure		400			{object}	models.SwaggerErrorResponse
//	@Failure		404			{object}	models.SwaggerErrorResponse
//	@Failure		500			{object}	models.SwaggerErrorResponse
//	@Router			/roles/{id} [put]
func (h *RoleHandler) UpdateRole(c *gin.Context) {
	id := c.Param("id")

	objectId, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		middleware.HandleError(c, middleware.NewAppError(
			middleware.ErrorCodeInvalidRequest,
			"Invalid role ID",
			err,
			http.StatusBadRequest,
		))
		return
	}

	// Fetch existing role to apply updates to
	existingRole, err := h.GetRoleUseCase.Execute(c.Request.Context(), objectId)
	if err != nil {
		middleware.HandleError(c, middleware.NewAppError(
			middleware.ErrorCodeInternalServer,
			"Failed to fetch role",
			err,
			http.StatusInternalServerError,
		))
		return
	}

	if existingRole == nil {
		middleware.HandleError(c, middleware.NewAppError(
			middleware.ErrorCodeNotFound,
			"Role not found",
			nil,
			http.StatusNotFound,
		))
		return
	}

	// Parse update DTO
	var updateDto dto.UpdateRoleDto
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

	// Apply updates to the existing role
	updateDto.ApplyUpdates(existingRole)

	// Call use case to update
	if err := h.UpdateRoleUseCase.Execute(c.Request.Context(), existingRole); err != nil {
		middleware.HandleError(c, middleware.NewAppError(
			middleware.ErrorCodeInternalServer,
			"Failed to update role",
			err,
			http.StatusInternalServerError,
		))
		return
	}

	c.JSON(http.StatusOK, existingRole)
}

// SoftDeleteRole godoc
//
//	@Summary		Delete a role
//	@Description	Soft-delete a role by ID (marks as deleted but keeps in database)
//	@Tags			roles
//	@Accept			json
//	@Produce		json
//	@Param			id	path		string	true	"Role ID"	example("6824886e6b180b753cea43e9")
//	@Success		200	{object}	models.SwaggerStandardResponse{data=object}
//	@Failure		400	{object}	models.SwaggerErrorResponse
//	@Failure		404	{object}	models.SwaggerErrorResponse
//	@Failure		500	{object}	models.SwaggerErrorResponse
//	@Router			/roles/{id} [delete]
func (h *RoleHandler) SoftDeleteRole(c *gin.Context) {
	id := c.Param("id")

	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		middleware.HandleError(c, middleware.NewAppError(
			middleware.ErrorCodeInvalidRequest,
			"Invalid role ID",
			err,
			http.StatusBadRequest,
		))
		return
	}

	// Get the role before deletion to return its details if needed
	role, err := h.GetRoleUseCase.Execute(c.Request.Context(), objectID)
	if err != nil {
		middleware.HandleError(c, middleware.NewAppError(
			middleware.ErrorCodeInternalServer,
			"Failed to fetch role",
			err,
			http.StatusInternalServerError,
		))
		return
	}

	if role == nil {
		middleware.HandleError(c, middleware.NewAppError(
			middleware.ErrorCodeNotFound,
			"Role not found",
			nil,
			http.StatusNotFound,
		))
		return
	}

	// Perform soft delete instead of hard delete
	deleted, err := h.SoftDeleteRoleUseCase.Execute(c.Request.Context(), objectID)
	if err != nil {
		middleware.HandleError(c, middleware.NewAppError(
			middleware.ErrorCodeInternalServer,
			"Failed to delete role",
			err,
			http.StatusInternalServerError,
		))
		return
	}

	if !deleted {
		middleware.HandleError(c, middleware.NewAppError(
			middleware.ErrorCodeNotFound,
			"Role not found or already deleted",
			nil,
			http.StatusNotFound,
		))
		return
	}

	// Return success response with the deleted role ID
	c.JSON(http.StatusOK, gin.H{
		"message": "Role deleted successfully",
		"id":      id,
	})
}

// RestoreRole godoc
//
//	@Summary		Restore a deleted role
//	@Description	Restore a soft-deleted role by ID (clears the deletedAt timestamp)
//	@Tags			roles
//	@Accept			json
//	@Produce		json
//	@Param			id	path		string	true	"Role ID"	example("6824886e6b180b753cea43e9")
//	@Success		200	{object}	models.SwaggerStandardResponse{data=entity.Role}
//	@Failure		400	{object}	models.SwaggerErrorResponse
//	@Failure		404	{object}	models.SwaggerErrorResponse
//	@Failure		500	{object}	models.SwaggerErrorResponse
//	@Router			/roles/{id}/restore [post]
func (h *RoleHandler) RestoreRole(c *gin.Context) {
	id := c.Param("id")

	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		middleware.HandleError(c, middleware.NewAppError(
			middleware.ErrorCodeInvalidRequest,
			"Invalid role ID",
			err,
			http.StatusBadRequest,
		))
		return
	}

	restored, err := h.RestoreRoleUseCase.Execute(c.Request.Context(), objectID)
	if err != nil {
		middleware.HandleError(c, middleware.NewAppError(
			middleware.ErrorCodeInternalServer,
			"Failed to restore role",
			err,
			http.StatusInternalServerError,
		))
		return
	}

	if !restored {
		middleware.HandleError(c, middleware.NewAppError(
			middleware.ErrorCodeNotFound,
			"Role not found or not deleted",
			nil,
			http.StatusNotFound,
		))
		return
	}

	// Get the restored role to return its details
	role, err := h.GetRoleUseCase.Execute(c.Request.Context(), objectID)
	if err != nil {
		middleware.HandleError(c, middleware.NewAppError(
			middleware.ErrorCodeInternalServer,
			"Role restored but failed to fetch updated data",
			err,
			http.StatusInternalServerError,
		))
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Role restored successfully",
		"role":    role,
	})
}

// HardDeleteRole godoc
//
//	@Summary		Permanently delete a role
//	@Description	Permanently delete a role by ID (removes it from the database)
//	@Tags			roles
//	@Accept			json
//	@Produce		json
//	@Param			id	path		string	true	"Role ID"	example("6824886e6b180b753cea43e9")
//	@Success		200	{object}	models.SwaggerStandardResponse{data=object}
//	@Failure		400	{object}	models.SwaggerErrorResponse
//	@Failure		404	{object}	models.SwaggerErrorResponse
//	@Failure		500	{object}	models.SwaggerErrorResponse
//	@Router			/roles/{id}/hard-delete [delete]
func (h *RoleHandler) HardDeleteRole(c *gin.Context) {
	id := c.Param("id")

	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		middleware.HandleError(c, middleware.NewAppError(
			middleware.ErrorCodeInvalidRequest,
			"Invalid role ID",
			err,
			http.StatusBadRequest,
		))
		return
	}

	deleted, err := h.HardDeleteRoleUseCase.Execute(c.Request.Context(), objectID)
	if err != nil {
		middleware.HandleError(c, middleware.NewAppError(
			middleware.ErrorCodeInternalServer,
			"Failed to permanently delete role",
			err,
			http.StatusInternalServerError,
		))
		return
	}

	if !deleted {
		middleware.HandleError(c, middleware.NewAppError(
			middleware.ErrorCodeNotFound,
			"Role not found",
			nil,
			http.StatusNotFound,
		))
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Role permanently deleted",
		"id":      id,
	})
}
