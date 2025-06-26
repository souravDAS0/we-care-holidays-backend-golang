package server

import (
	"bitbucket.org/abhishek_fordel/we-care-holidays-backend-golang/container"
	permissionHandlers "bitbucket.org/abhishek_fordel/we-care-holidays-backend-golang/modules/permissions/presentation/http/handlers"
	permissionRoutes "bitbucket.org/abhishek_fordel/we-care-holidays-backend-golang/modules/permissions/presentation/http/routes"
	"github.com/gin-gonic/gin"
)

func registerPermissionRoutes(router *gin.RouterGroup, app *container.AppContainer) {
	permissionHandler := permissionHandlers.NewPermissionHandler(
		app.Permission.GetPermissionUseCase,
		app.Permission.CreatePermissionUseCase,
		app.Permission.ListPermissionsUseCase,
		app.Permission.UpdatePermissionUseCase,
		app.Permission.HardDeletePermissionUseCase,
	)

	permissionRoutes.RegisterPermissionRoutes(router, permissionHandler)
}
