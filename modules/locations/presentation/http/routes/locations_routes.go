package routes

import (
	"bitbucket.org/abhishek_fordel/we-care-holidays-backend-golang/container"
	"bitbucket.org/abhishek_fordel/we-care-holidays-backend-golang/internal/constants"
	"bitbucket.org/abhishek_fordel/we-care-holidays-backend-golang/internal/middleware"
	"bitbucket.org/abhishek_fordel/we-care-holidays-backend-golang/modules/locations/presentation/http/handlers"
	"github.com/gin-gonic/gin"
)

// RegisterLocationRoutes registers all location-related routes
func RegisterLocationRoutes(rg *gin.RouterGroup, h *handlers.LocationHandler, app *container.AppContainer) {
	n := rg.Group(constants.LocationBasePath)

	n.Use(middleware.AutoGuard(app.RBACService,
		middleware.WithOwnership("organizationId", "id"),
	))

	{
		// List and create
		n.GET(constants.ListLocationsPath, h.ListLocations)
		n.POST(constants.CreateLocationPath, h.CreateLocation)

		// Bulk actions
		n.DELETE(constants.BulkDeleteLocationsPath, h.BulkDeleteLocations)
		n.POST(constants.BulkRestoreLocationsPath, h.BulkRestoreLocations)

		// Single item operations
		n.GET(constants.GetLocationPath, h.GetLocation)
		n.PUT(constants.UpdateLocationPath, h.UpdateLocation)
		n.DELETE(constants.DeleteLocationPath, h.DeleteLocation)

		// Media upload
		n.POST(constants.UploadLocationMediaPath, h.UploadLocationMedia)
		n.POST(constants.RestoreLocationPath, h.RestoreLocation)
		n.DELETE(constants.HardDeleteLocationPath, h.HardDeleteLocation)
	}
}
