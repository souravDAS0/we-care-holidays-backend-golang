package server

import (
	"bitbucket.org/abhishek_fordel/we-care-holidays-backend-golang/container"
	"github.com/gin-gonic/gin"
)

// registerRoutes sets up all routes (public and private)
func registerRoutes(r *gin.Engine, app *container.AppContainer) {

	registerPublicRoutes(r, app)
	registerPrivateRoutes(r, app)

}
