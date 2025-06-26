package handlers

import (
	"bitbucket.org/abhishek_fordel/we-care-holidays-backend-golang/commons/services"
	"bitbucket.org/abhishek_fordel/we-care-holidays-backend-golang/modules/users/domain/usecases"
)

// UserHandler handles HTTP requests for Users
type UserHandler struct {
	GetUserUseCase             *usecases.GetUserUseCase
	CreateUserUseCase          *usecases.CreateUserUseCase
	ListUsersUseCase           *usecases.ListUsersUseCase
	UpdateUserUseCase          *usecases.UpdateUserUseCase
	SoftDeleteUserUseCase      *usecases.SoftDeleteUserUseCase
	RestoreUserUseCase         *usecases.RestoreUserUseCase
	BulkSoftDeleteUsersUseCase *usecases.BulkSoftDeleteUsersUseCase
	HardDeleteUserUseCase      *usecases.HardDeleteUserUseCase
	BulkRestoreUsersUseCase    *usecases.BulkRestoreUsersUseCase
	UpdateUserStatusUseCase    *usecases.UpdateUserStatusUseCase
	fileService                services.FileService
	FindUserByEmailUsecase     *usecases.FindUserByEmailUsecase
}

func NewUserHandler(GetUserUseCase *usecases.GetUserUseCase,
	CreateUserUseCase *usecases.CreateUserUseCase,
	ListUsersUseCase *usecases.ListUsersUseCase,
	UpdateUserUseCase *usecases.UpdateUserUseCase,
	SoftDeleteUserUseCase *usecases.SoftDeleteUserUseCase,
	RestoreUserUseCase *usecases.RestoreUserUseCase,
	BulkSoftDeleteUsersUseCase *usecases.BulkSoftDeleteUsersUseCase,
	HardDeleteUserUseCase *usecases.HardDeleteUserUseCase,
	BulkRestoreUsersUseCase *usecases.BulkRestoreUsersUseCase,
	UpdateUserStatusUseCase *usecases.UpdateUserStatusUseCase,
	fileService services.FileService,
	FindUserByEmailUsecase *usecases.FindUserByEmailUsecase,
) *UserHandler {
	return &UserHandler{
		GetUserUseCase:             GetUserUseCase,
		CreateUserUseCase:          CreateUserUseCase,
		ListUsersUseCase:           ListUsersUseCase,
		UpdateUserUseCase:          UpdateUserUseCase,
		SoftDeleteUserUseCase:      SoftDeleteUserUseCase,
		RestoreUserUseCase:         RestoreUserUseCase,
		BulkSoftDeleteUsersUseCase: BulkSoftDeleteUsersUseCase,
		HardDeleteUserUseCase:      HardDeleteUserUseCase,
		BulkRestoreUsersUseCase:    BulkRestoreUsersUseCase,
		UpdateUserStatusUseCase:    UpdateUserStatusUseCase,
		fileService:                fileService,
		FindUserByEmailUsecase:     FindUserByEmailUsecase,
	}
}
