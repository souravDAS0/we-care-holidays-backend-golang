package server

import (
	"bitbucket.org/abhishek_fordel/we-care-holidays-backend-golang/container"
	orgHandlers "bitbucket.org/abhishek_fordel/we-care-holidays-backend-golang/modules/organizations/presentation/http/handlers"
	orgRoutes "bitbucket.org/abhishek_fordel/we-care-holidays-backend-golang/modules/organizations/presentation/http/routes"
	"github.com/gin-gonic/gin"
)

func registerOrganizationRoutes( router *gin.RouterGroup, app *container.AppContainer){

	orgHandler := orgHandlers.NewOrganizationHandler(
		app.Organization.GetOrganizationUseCase,
		app.FileService,
		app.Organization.CreateOrganizationUseCase,
		app.Organization.ListOrganizationUseCase,
		app.Organization.UpdateOrganizationUseCase,
		app.Organization.UpdateOrganizationStatusUseCase,
		app.Organization.SoftDeleteOrganizationUseCase,
		app.Organization.RestoreOrganizationUseCase,
		app.Organization.BulkSoftDeleteOrganizationsUseCase,
		app.Organization.HardDeleteOrganizationUseCase,
		app.Organization.BulkRestoreOrganizationsUseCase,
	)

	orgRoutes.RegisterOrganizationRoutes(router, orgHandler)
}