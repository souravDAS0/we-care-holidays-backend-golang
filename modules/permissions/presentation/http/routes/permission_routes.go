package routes

import (
	"bitbucket.org/abhishek_fordel/we-care-holidays-backend-golang/internal/constants"
	"bitbucket.org/abhishek_fordel/we-care-holidays-backend-golang/modules/permissions/presentation/http/handlers"
	"github.com/gin-gonic/gin"
)

func RegisterPermissionRoutes(router *gin.RouterGroup, handler *handlers.PermissionHandler) {
	permissionGroup := router.Group(constants.PermissionBasePath)
	{
		permissionGroup.GET(constants.ListPermissionsPath, handler.ListPermissions)
		permissionGroup.POST(constants.CreatePermissionPath, handler.CreatePermission)

		permissionGroup.GET(constants.GetPermissionPath, handler.GetPermission)
		permissionGroup.PUT(constants.UpdatePermissionPath, handler.UpdatePermission)

		permissionGroup.DELETE(constants.HardDeletePermissionPath, handler.HardDeletePermission)

	}

}
