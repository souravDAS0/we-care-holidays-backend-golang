package server

import (
	"bitbucket.org/abhishek_fordel/we-care-holidays-backend-golang/container"
	userHandlers "bitbucket.org/abhishek_fordel/we-care-holidays-backend-golang/modules/users/presentation/http/handlers"
	userRoutes "bitbucket.org/abhishek_fordel/we-care-holidays-backend-golang/modules/users/presentation/http/routes"
	"github.com/gin-gonic/gin"
)

func registerUserRoutes(router *gin.RouterGroup, app *container.AppContainer) {
	userHandler := userHandlers.NewUserHandler(
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

	userRoutes.RegisterUserRoutes(router, userHandler, app)
}
