package container

import (
	"bitbucket.org/abhishek_fordel/we-care-holidays-backend-golang/modules/organizations/data/datasource"
	"bitbucket.org/abhishek_fordel/we-care-holidays-backend-golang/modules/organizations/data/mongodb/repository"
	"bitbucket.org/abhishek_fordel/we-care-holidays-backend-golang/modules/organizations/domain/usecases"
)


type OrganizationContainer struct {
	Repository *repository.OrganizationRepositoryMongo
	GetOrganizationUseCase             *usecases.GetOrganizationUseCase
	CreateOrganizationUseCase          *usecases.CreateOrganizationUseCase
	ListOrganizationUseCase            *usecases.ListOrganizationUseCase
	UpdateOrganizationUseCase          *usecases.UpdateOrganizationUseCase
	UpdateOrganizationStatusUseCase    *usecases.UpdateOrganizationStatusUseCase
	SoftDeleteOrganizationUseCase      *usecases.SoftDeleteOrganizationUseCase
	RestoreOrganizationUseCase         *usecases.RestoreOrganizationUseCase
	BulkSoftDeleteOrganizationsUseCase *usecases.BulkSoftDeleteOrganizationsUseCase
	HardDeleteOrganizationUseCase      *usecases.HardDeleteOrganizationUseCase
	BulkRestoreOrganizationsUseCase 	*usecases.BulkRestoreOrganizationsUseCase
}


func (c *AppContainer) InjectOrganizationContainer() {
	// Datasource
	organizationDS := datasource.NewMongoOrganizationDatasource(c.MongoDatabase)

	// Repository
	organizationRepo := repository.NewOrganizationRepositoryMongo(organizationDS)
	// Use cases
	getOrganizationUC := usecases.NewGetOrganizationUseCase(organizationRepo)
	createOrganizationUC := usecases.NewCreateOrganizationUseCase(organizationRepo)
	listOrganizationUC := usecases.NewListOrganizationUseCase(organizationRepo)
	updateOrganizationUC := usecases.NewUpdateOrganizationUseCase(organizationRepo)
	updateOrganizationStatusUC := usecases.NewUpdateOrganizationStatusUseCase(organizationRepo)
	softDeleteOrganizationUC := usecases.NewSoftDeleteOrganizationUseCase(organizationRepo)
	restoreOrganizationUC := usecases.NewRestoreOrganizationUseCase(organizationRepo)
	bulkSoftDeleteOrganizationsUC := usecases.NewBulkSoftDeleteOrganizationsUseCase(organizationRepo)
	hardDeleteOrganizationUC := usecases.NewHardDeleteOrganizationUseCase(organizationRepo)
	bulkRestoreOrganizationsUC 	:= usecases.NewBulkRestoreOrganizationsUseCase(organizationRepo)

	// Assign to container
	c.Organization = &OrganizationContainer{
		Repository: organizationRepo,
		GetOrganizationUseCase:             getOrganizationUC,
		CreateOrganizationUseCase:          createOrganizationUC,
		ListOrganizationUseCase:            listOrganizationUC,
		UpdateOrganizationUseCase:          updateOrganizationUC,
		UpdateOrganizationStatusUseCase:    updateOrganizationStatusUC,
		SoftDeleteOrganizationUseCase:      softDeleteOrganizationUC,
		RestoreOrganizationUseCase:         restoreOrganizationUC,
		BulkSoftDeleteOrganizationsUseCase: bulkSoftDeleteOrganizationsUC,
		HardDeleteOrganizationUseCase:      hardDeleteOrganizationUC,
		BulkRestoreOrganizationsUseCase: 	bulkRestoreOrganizationsUC,
	}
}