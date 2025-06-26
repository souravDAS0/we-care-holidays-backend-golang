package routes

import (
	"github.com/gin-gonic/gin"

	"bitbucket.org/abhishek_fordel/we-care-holidays-backend-golang/internal/constants"
	"bitbucket.org/abhishek_fordel/we-care-holidays-backend-golang/modules/organizations/presentation/http/handlers"
)

// RegisterOrganizationRoutes registers all organization-related routes
func RegisterOrganizationRoutes(router *gin.RouterGroup, handler *handlers.OrganizationHandler) {
	orgGroup := router.Group(constants.OrganizationBasePath)
	{
		// List and create
		orgGroup.GET(constants.ListOrganizationsPath, handler.ListOrganizations)
		orgGroup.POST(constants.CreateOrganizationPath, handler.CreateOrganization)

		// Bulk operations
		orgGroup.DELETE(constants.BulkDeleteOrganizationsPath, handler.BulkDeleteOrganizations)
		orgGroup.POST(constants.BulkRestoreOrganizationsPath, handler.BulkRestoreOrganizations)

		// Single item operations
		orgGroup.GET(constants.GetOrganizationPath, handler.GetOrganization)
		orgGroup.PUT(constants.UpdateOrganizationPath, handler.UpdateOrganization)
		orgGroup.DELETE(constants.DeleteOrganizationPath, handler.DeleteOrganization)

		// Status update
		orgGroup.PUT(constants.UpdateStatusPath, handler.UpdateOrganizationStatus)

		// Logo upload
		orgGroup.POST(constants.UploadOrgLogoPath, handler.UploadOrganizationLogo)

		// Restore operation
		orgGroup.POST(constants.RestoreOrganizationPath, handler.RestoreOrganization)

		// Hard delete (admin/cleanup operations)
		// This is typically protected by admin-only middleware
		orgGroup.DELETE(constants.HardDeleteOrganizationPath, handler.HardDeleteOrganization)
	}
}
