package container

import (
	"bitbucket.org/abhishek_fordel/we-care-holidays-backend-golang/modules/permissions/data/datasource"
	"bitbucket.org/abhishek_fordel/we-care-holidays-backend-golang/modules/permissions/data/mongodb/repository"
	"bitbucket.org/abhishek_fordel/we-care-holidays-backend-golang/modules/permissions/domain/usecases"
)

type PermissionContainer struct {
	Repository              *repository.PermissionRepositoryMongo
	GetPermissionUseCase    *usecases.GetPermissionUseCase
	CreatePermissionUseCase *usecases.CreatePermissionUseCase
	ListPermissionsUseCase  *usecases.ListPermissionsUseCase
	UpdatePermissionUseCase *usecases.UpdatePermissionUseCase

	HardDeletePermissionUseCase *usecases.HardDeletePermissionUseCase
}

func (c *AppContainer) InjectPermissionContainer() {
	// Datasource
	permissionDS := datasource.NewMongoPermissionDatasource(c.MongoDatabase)

	// Repository
	permissionRepo := repository.NewPermissionRepositoryMongo(permissionDS)

	// Use cases
	getPermissionUC := usecases.NewGetPermissionUseCase(permissionRepo)
	createPermissionUC := usecases.NewCreatePermissionUseCase(permissionRepo)
	listPermissionUC := usecases.NewListPermissionUseCase(permissionRepo)
	updatePermissionUC := usecases.NewUpdatePermissionUseCase(permissionRepo)
	hardDeletePermissionUC := usecases.NewHardDeletePermissionUseCase(permissionRepo)

	c.Permission = &PermissionContainer{
		Repository:                  permissionRepo,
		GetPermissionUseCase:        getPermissionUC,
		CreatePermissionUseCase:     createPermissionUC,
		ListPermissionsUseCase:      listPermissionUC,
		UpdatePermissionUseCase:     updatePermissionUC,
		HardDeletePermissionUseCase: hardDeletePermissionUC,
	}
}
