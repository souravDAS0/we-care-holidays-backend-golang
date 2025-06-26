package routes

import (
	"bitbucket.org/abhishek_fordel/we-care-holidays-backend-golang/internal/constants"
	"bitbucket.org/abhishek_fordel/we-care-holidays-backend-golang/modules/users/presentation/http/handlers"
	"github.com/gin-gonic/gin"
)

func RegisterUserRoutes(router *gin.RouterGroup, handler *handlers.UserHandler) {
	userGroup := router.Group(constants.UserBasePath)
	{
		userGroup.GET(constants.ListUsersPath, handler.ListUsers)
		userGroup.POST(constants.CreateUserPath, handler.CreateUser)

		userGroup.DELETE(constants.BulkDeleteUsersPath, handler.BulkDeleteUsers)
		userGroup.POST(constants.BulkRestoreUsersPath, handler.BulkRestoreUsers)

		userGroup.GET(constants.GetUserPath, handler.GetUser)
		userGroup.PUT(constants.UpdateUserPath, handler.UpdateUser)
		userGroup.DELETE(constants.DeleteUserPath, handler.DeleteUser)

		userGroup.POST(constants.RestoreUserPath, handler.RestoreUser)

		userGroup.POST(constants.UploadUserAvatarPath, handler.UploadUserProfilePhoto)

		userGroup.DELETE(constants.HardDeleteUserPath, handler.HardDeleteUser)

	}

}
