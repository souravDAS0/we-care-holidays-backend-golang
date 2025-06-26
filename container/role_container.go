package container

import (
	"bitbucket.org/abhishek_fordel/we-care-holidays-backend-golang/modules/roles/data/datasource"
	"bitbucket.org/abhishek_fordel/we-care-holidays-backend-golang/modules/roles/data/mongodb/repository"
	"bitbucket.org/abhishek_fordel/we-care-holidays-backend-golang/modules/roles/domain/usecases"
)

type RoleContainer struct {
	Repository                 *repository.RoleRepositoryMongo
	GetRoleUseCase             *usecases.GetRoleUseCase
	CreateRoleUseCase          *usecases.CreateRoleUseCase
	ListRolesUseCase           *usecases.ListRolesUseCase
	UpdateRoleUseCase          *usecases.UpdateRoleUseCase
	SoftDeleteRoleUseCase      *usecases.SoftDeleteRoleUseCase
	RestoreRoleUseCase         *usecases.RestoreRoleUseCase
	BulkSoftDeleteRolesUseCase *usecases.BulkSoftDeleteRolesUseCase
	HardDeleteRoleUseCase      *usecases.HardDeleteRoleUseCase
	BulkRestoreRolesUseCase    *usecases.BulkRestoreRolesUseCase
}

func (c *AppContainer) InjectRoleContainer() {
	if c.Permission == nil {
		panic("Permission container must be injected before Role container")
	}

	// Datasource
	roleDS := datasource.NewMongoRoleDatasource(c.MongoDatabase)
	// Repository
	roleRepo := repository.NewRoleRepositoryMongo(roleDS)

	permissionRepo := c.Permission.Repository

	// Use cases
	getRoleUC := usecases.NewGetRoleUseCase(roleRepo)
	createRoleUC := usecases.NewCreateRoleUseCase(roleRepo, permissionRepo)
	listRolesUC := usecases.NewListRolesUseCase(roleRepo)
	updateRoleUC := usecases.NewUpdateRoleUseCase(roleRepo, permissionRepo)
	softDeleteRoleUC := usecases.NewSoftDeleteRoleUseCase(roleRepo)
	restoreRoleUC := usecases.NewRestoreRoleUseCase(roleRepo)
	bulkSoftDeleteRolesUC := usecases.NewBulkSoftDeleteRolesUseCase(roleRepo)
	hardDeleteRoleUC := usecases.NewHardDeleteRoleUseCase(roleRepo)
	bulkRestoreRolesUC := usecases.NewBulkRestoreRolesUseCase(roleRepo)

	// Assign to container
	c.Role = &RoleContainer{
		Repository:                 roleRepo,
		GetRoleUseCase:             getRoleUC,
		CreateRoleUseCase:          createRoleUC,
		ListRolesUseCase:           listRolesUC,
		UpdateRoleUseCase:          updateRoleUC,
		SoftDeleteRoleUseCase:      softDeleteRoleUC,
		RestoreRoleUseCase:         restoreRoleUC,
		BulkSoftDeleteRolesUseCase: bulkSoftDeleteRolesUC,
		HardDeleteRoleUseCase:      hardDeleteRoleUC,
		BulkRestoreRolesUseCase:    bulkRestoreRolesUC,
	}
}
