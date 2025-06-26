package server

import (
	"bitbucket.org/abhishek_fordel/we-care-holidays-backend-golang/container"
	locHandlers "bitbucket.org/abhishek_fordel/we-care-holidays-backend-golang/modules/locations/presentation/http/handlers"
	locRoutes "bitbucket.org/abhishek_fordel/we-care-holidays-backend-golang/modules/locations/presentation/http/routes"
	"github.com/gin-gonic/gin"
)

func registerLocationRoutes(router *gin.RouterGroup, app *container.AppContainer) {
	// Create location handler using dependencies from the app container
	locHandler := locHandlers.NewLocationHandler(
		app.Location.GetLocationUseCase,
		app.FileService,
		app.Location.CreateLocationUseCase,
		app.Location.ListLocationsUseCase,
		app.Location.UpdateLocationUseCase,
		app.Location.DeleteLocationUseCase,
		app.Location.BulkSoftDeleteLocationsUseCase,
		app.Location.UploadLocationMediaUseCase,
		app.Location.RestoreLocationUseCase,
		app.Location.BulkRestoreLocationsUseCase,
		app.Location.HardDeleteLocationUseCase,
	)

	// Register location routes with the handler
	locRoutes.RegisterLocationRoutes(router, locHandler, app)
}
