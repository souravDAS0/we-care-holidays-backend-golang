package server

import (
	"bitbucket.org/abhishek_fordel/we-care-holidays-backend-golang/container"
	"bitbucket.org/abhishek_fordel/we-care-holidays-backend-golang/internal/constants"
	"bitbucket.org/abhishek_fordel/we-care-holidays-backend-golang/internal/middleware"
	"github.com/gin-gonic/gin"
)

func registerPrivateRoutes(r *gin.Engine, app *container.AppContainer) {
	private := r.Group(constants.AppBasePath)



	private.Use(
		middleware.ErrorHandler(),
		middleware.ResponseInterceptor(),
		middleware.SecureHeaders(),
		middleware.AuthMiddleware(app.RBACService, app.Config), // ← Auth happens ONCE here
	)



	registerUserRoutes(private, app)
	registerPermissionRoutes(private, app)
	registerRoleRoutes(private, app)
	registerOrganizationRoutes(private, app)
	registerLocationRoutes(private, app)
}
