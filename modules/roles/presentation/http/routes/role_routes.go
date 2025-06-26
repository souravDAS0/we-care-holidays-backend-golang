package routes

import (
	"bitbucket.org/abhishek_fordel/we-care-holidays-backend-golang/container"
	"bitbucket.org/abhishek_fordel/we-care-holidays-backend-golang/internal/constants"
	"bitbucket.org/abhishek_fordel/we-care-holidays-backend-golang/internal/middleware"
	"bitbucket.org/abhishek_fordel/we-care-holidays-backend-golang/modules/roles/presentation/http/handlers"
	"github.com/gin-gonic/gin"
)

func RegisterRoleRoutes(router *gin.RouterGroup, handler *handlers.RoleHandler, app *container.AppContainer) {
	roleGroup := router.Group(constants.RoleBasePath)

	roleGroup.Use(middleware.AutoGuard(app.RBACService)) // middleware.WithOwnership("organizationId", "id"),

	{
		roleGroup.GET(constants.ListRolesPath, handler.ListRoles)
		roleGroup.POST(constants.CreateRolePath, handler.CreateRole)

		roleGroup.DELETE(constants.BulkDeleteRolesPath, handler.BulkSoftDeleteRoles)
		roleGroup.POST(constants.BulkRestoreRolesPath, handler.BulkRestoreRoles)

		roleGroup.GET(constants.GetRolePath, handler.GetRole)
		roleGroup.PUT(constants.UpdateRolePath, handler.UpdateRole)
		roleGroup.DELETE(constants.DeleteRolePath, handler.SoftDeleteRole)

		roleGroup.POST(constants.RestoreRolePath, handler.RestoreRole)
		roleGroup.DELETE(constants.HardDeleteRolePath, handler.HardDeleteRole)
	}
}
