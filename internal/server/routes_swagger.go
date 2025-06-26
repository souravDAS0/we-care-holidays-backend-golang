// cmd/server/swagger_routes.go
package server

import (
	"bitbucket.org/abhishek_fordel/we-care-holidays-backend-golang/container"
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"

	// this blank import pulls in your generated swagger docs
	_ "bitbucket.org/abhishek_fordel/we-care-holidays-backend-golang/docs"
)

// registerSwaggerRoutes mounts the Swagger UI under /swagger/*any
// and redirects “/” → “/swagger/index.html”
func registerSwaggerRoutes(public *gin.RouterGroup, _ *container.AppContainer) {
    // serve swagger UI at /swagger/index.html, static under /swagger/*
    public.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

    // optional: make root redirect to swagger
    public.GET("/", func(c *gin.Context) {
        c.Redirect(302, public.BasePath()+"/swagger/index.html")
    })
}
