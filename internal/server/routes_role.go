package server

import (
	"bitbucket.org/abhishek_fordel/we-care-holidays-backend-golang/container"
	"bitbucket.org/abhishek_fordel/we-care-holidays-backend-golang/modules/roles/presentation/http/handlers"
	"bitbucket.org/abhishek_fordel/we-care-holidays-backend-golang/modules/roles/presentation/http/routes"
	"github.com/gin-gonic/gin"
)

func registerRoleRoutes(router *gin.RouterGroup, app *container.AppContainer) {
	roleHandler := handlers.NewRoleHandler(
		app.Role.GetRoleUseCase,
		app.Role.CreateRoleUseCase,
		app.Role.ListRolesUseCase,
		app.Role.UpdateRoleUseCase,
		app.Role.SoftDeleteRoleUseCase,
		app.Role.RestoreRoleUseCase,
		app.Role.BulkSoftDeleteRolesUseCase,
		app.Role.HardDeleteRoleUseCase,
		app.Role.BulkRestoreRolesUseCase,
		app.Permission.ListPermissionsUseCase,
		// app.RBACService,
		// app.PermissionValidator,
	)

	routes.RegisterRoleRoutes(router, roleHandler, app)
}
