package container

import (
	"bitbucket.org/abhishek_fordel/we-care-holidays-backend-golang/modules/users/data/datasource"
	"bitbucket.org/abhishek_fordel/we-care-holidays-backend-golang/modules/users/data/mongodb/repository"
	"bitbucket.org/abhishek_fordel/we-care-holidays-backend-golang/modules/users/domain/usecases"
)

type UserContainer struct {
	Repository                 *repository.UserRepositoryMongo
	GetUserUseCase             *usecases.GetUserUseCase
	CreateUserUseCase          *usecases.CreateUserUseCase
	ListUsersUseCase           *usecases.ListUsersUseCase
	UpdateUserUseCase          *usecases.UpdateUserUseCase
	UpdateUserStatusUseCase    *usecases.UpdateUserStatusUseCase
	SoftDeleteUserUseCase      *usecases.SoftDeleteUserUseCase
	RestoreUserUseCase         *usecases.RestoreUserUseCase
	BulkSoftDeleteUsersUseCase *usecases.BulkSoftDeleteUsersUseCase
	BulkRestoreUsersUseCase    *usecases.BulkRestoreUsersUseCase
	HardDeleteUserUseCase      *usecases.HardDeleteUserUseCase
	FindUserByEmailUsecase     *usecases.FindUserByEmailUsecase
}

func (c *AppContainer) InjectUserContainer() {
	if c.Role == nil {
		panic("Role container must be injected before User container")
	}

	if c.Organization == nil {
		panic("Organization container must be injected before User container")
	}
	// Datasource
	userDS := datasource.NewMongoUserDatasource(c.MongoDatabase)
	// Repository
	userRepo := repository.NewUserRepositoryMongo(userDS)

	roleRepo := c.Role.Repository
	orgRepo := c.Organization.Repository
	// Use cases
	getUserUC := usecases.NewGetUserUseCase(userRepo)
	createUserUC := usecases.NewCreateUserUseCase(userRepo, roleRepo, orgRepo)
	listUserUC := usecases.NewListUsersUseCase(userRepo)
	updateUserUC := usecases.NewUpdateUserUseCase(userRepo)
	updateUserStatusUC := usecases.NewUpdateUserStatusUseCase(userRepo)
	softDeleteUserUC := usecases.NewSoftDeleteUserUseCase(userRepo)
	restoreUserUC := usecases.NewRestoreUserUseCase(userRepo)
	bulkSoftDeleteUsersUC := usecases.NewBulkSoftDeleteUsersUseCase(userRepo)
	hardDeleteUserUC := usecases.NewHardDeleteUserUseCase(userRepo)
	bulkRestoreUsersUC := usecases.NewBulkRestoreUsersUseCase(userRepo)
	findUserByEmailUC := usecases.NewFindUserByEmailUsecase(userRepo)

	// Assign to container
	c.User = &UserContainer{
		GetUserUseCase:             getUserUC,
		CreateUserUseCase:          createUserUC,
		ListUsersUseCase:           listUserUC,
		UpdateUserUseCase:          updateUserUC,
		UpdateUserStatusUseCase:    updateUserStatusUC,
		SoftDeleteUserUseCase:      softDeleteUserUC,
		RestoreUserUseCase:         restoreUserUC,
		BulkSoftDeleteUsersUseCase: bulkSoftDeleteUsersUC,
		HardDeleteUserUseCase:      hardDeleteUserUC,
		BulkRestoreUsersUseCase:    bulkRestoreUsersUC,
		FindUserByEmailUsecase:     findUserByEmailUC,
		Repository:                 userRepo,
	}
}
