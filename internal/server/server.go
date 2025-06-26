package server

import (
	"fmt"

	"bitbucket.org/abhishek_fordel/we-care-holidays-backend-golang/internal/middleware"

	"bitbucket.org/abhishek_fordel/we-care-holidays-backend-golang/container"
	"github.com/gin-gonic/gin"
)

// StartServer initializes and starts the Gin HTTP server
func StartServer(app *container.AppContainer) error {
	// Gin engine
	r := gin.New()

	// Global middlewares
	r.Use(
		middleware.ErrorHandler(),
		// middleware.RequestLoggerMiddleware(),
	)

	// Register routes
	registerRoutes(r, app)

	// Start server
	return r.Run(fmt.Sprintf(":%s", app.Config.Port))
}