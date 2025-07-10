package routes

import (
	"github.com/gin-gonic/gin"

	"bitbucket.org/abhishek_fordel/we-care-holidays-backend-golang/internal/constants"
	"bitbucket.org/abhishek_fordel/we-care-holidays-backend-golang/internal/middleware"
	"bitbucket.org/abhishek_fordel/we-care-holidays-backend-golang/modules/organizations/presentation/http/handlers"
)

// RegisterOrganizationRoutes registers all organization-related routes
func RegisterOrganizationRoutes(router *gin.RouterGroup, handler *handlers.OrganizationHandler) {
	orgGroup := router.Group(constants.OrganizationBasePath)
	orgGroup.Use(middleware.ScopedRBACMiddleware())
	{
		// List and create
		orgGroup.GET(constants.ListOrganizationsPath,
			middleware.RequireScopedPermission("organizations", "list"),
			handler.ListOrganizations)

		orgGroup.POST(constants.CreateOrganizationPath,
			middleware.RequireScopedPermission("organizations", "create"),
			handler.CreateOrganization)

		// Bulk operations
		orgGroup.DELETE(constants.BulkDeleteOrganizationsPath,
			middleware.RequireScopedPermission("organizations", "delete"),
			handler.BulkDeleteOrganizations)

		orgGroup.POST(constants.BulkRestoreOrganizationsPath,
			middleware.RequireScopedPermission("organizations", "create"),
			handler.BulkRestoreOrganizations)

		// Single item operations
		orgGroup.GET(constants.GetOrganizationPath,
			middleware.RequireScopedPermission("organizations", "read"),
			middleware.RequireOrganizationAccess(),
			handler.GetOrganization)

		orgGroup.PUT(constants.UpdateOrganizationPath,
			middleware.RequireScopedPermission("organizations", "update"),
			middleware.RequireOrganizationAccess(),
			handler.UpdateOrganization)

		orgGroup.DELETE(constants.DeleteOrganizationPath,
			middleware.RequireScopedPermission("organizations", "delete"),
			middleware.RequireOrganizationAccess(),
			handler.DeleteOrganization)

		// Status update
		orgGroup.PUT(constants.UpdateStatusPath,
			middleware.RequireScopedPermission("organizations", "update"),
			middleware.RequireOrganizationAccess(),
			handler.UpdateOrganizationStatus)

		// Logo upload
		orgGroup.POST(constants.UploadOrgLogoPath,
			middleware.RequireScopedPermission("organizations", "update"),
			middleware.RequireOrganizationAccess(),
			handler.UploadOrganizationLogo)

		// Restore operation
		orgGroup.POST(constants.RestoreOrganizationPath,
			middleware.RequireScopedPermission("organizations", "update"),
			middleware.RequireOrganizationAccess(),
			handler.RestoreOrganization)

		// Hard delete (admin/cleanup operations)
		// This is typically protected by admin-only middleware
		orgGroup.DELETE(constants.HardDeleteOrganizationPath,
			middleware.RequireScopedPermission("organizations", "delete"),
			handler.HardDeleteOrganization)
	}
}
