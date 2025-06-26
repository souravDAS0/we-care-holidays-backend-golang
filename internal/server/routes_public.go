package server

import (
	"net/http"
	"time"

	"bitbucket.org/abhishek_fordel/we-care-holidays-backend-golang/container"
	"bitbucket.org/abhishek_fordel/we-care-holidays-backend-golang/internal/constants"
	"bitbucket.org/abhishek_fordel/we-care-holidays-backend-golang/modules/users/presentation/http/handlers"

	"github.com/gin-gonic/gin"
)

func registerPublicRoutes(r *gin.Engine, app *container.AppContainer) {
	public := r.Group(constants.AppBasePath)

	r.GET("/", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "WeCare Holidays API is running",
			"swagger": "Available at /api/v1/swagger/index.html",
		})
	})

	// Health Check Route
	public.GET(constants.HealthCheckRoute, func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"status":    "UP",
			"timestamp": time.Now().Format(time.RFC3339),
			"services": gin.H{
				"api":     "available",
				"version": "1.0.0",
			},
		})
	})

	userHandler := handlers.NewUserHandler(
		app.User.GetUserUseCase,
		app.User.CreateUserUseCase,
		app.User.ListUsersUseCase,
		app.User.UpdateUserUseCase,
		app.User.SoftDeleteUserUseCase,
		app.User.RestoreUserUseCase,
		app.User.BulkSoftDeleteUsersUseCase,
		app.User.HardDeleteUserUseCase,
		app.User.BulkRestoreUsersUseCase,
		app.User.UpdateUserStatusUseCase,
		app.FileService,
		app.User.FindUserByEmailUsecase,
	)

	public.POST("/users/login", userHandler.Login)

	registerSwaggerRoutes(public, app)
}
