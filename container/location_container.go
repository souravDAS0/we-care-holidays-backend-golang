package container

import (
	"bitbucket.org/abhishek_fordel/we-care-holidays-backend-golang/modules/locations/data/datasource"
	"bitbucket.org/abhishek_fordel/we-care-holidays-backend-golang/modules/locations/data/mongodb/repository"
	"bitbucket.org/abhishek_fordel/we-care-holidays-backend-golang/modules/locations/domain/usecases"
)

type LocationContainer struct {
	GetLocationUseCase          *usecases.GetLocationUseCase
	CreateLocationUseCase       *usecases.CreateLocationUseCase
	ListLocationsUseCase        *usecases.ListLocationsUseCase
	UpdateLocationUseCase       *usecases.UpdateLocationUseCase
	DeleteLocationUseCase       *usecases.DeleteLocationUseCase
	BulkSoftDeleteLocationsUseCase  *usecases.BulkSoftDeleteLocationsUseCase
	UploadLocationMediaUseCase  *usecases.UploadLocationMediaUseCase
	RestoreLocationUseCase      *usecases.RestoreLocationUseCase
	BulkRestoreLocationsUseCase *usecases.BulkRestoreLocationsUseCase
	HardDeleteLocationUseCase   *usecases.HardDeleteLocationUseCase
}

func (c *AppContainer) InjectLocationContainer() {
	// Datasource
	locationDS := datasource.NewMongoLocationDatasource(c.MongoDatabase)

	// Repository
	locationRepo := repository.NewLocationRepositoryMongo(locationDS)

	// Use cases
	getLocationUC := usecases.NewGetLocationUseCase(locationRepo)
	createLocationUC := usecases.NewCreateLocationUseCase(locationRepo)
	listLocationsUC := usecases.NewListLocationsUseCase(locationRepo)
	updateLocationUC := usecases.NewUpdateLocationUseCase(locationRepo)
	deleteLocationUC := usecases.NewDeleteLocationUseCase(locationRepo)
	bulkSoftDeleteLocationsUC := usecases.NewBulkSoftDeleteLocationsUseCase(locationRepo)
	uploadLocationMediaUC := usecases.NewUploadLocationMediaUseCase(locationRepo)
	restoreUC := usecases.NewRestoreLocationUseCase(locationRepo)
	bulkRestoreUC := usecases.NewBulkRestoreLocationsUseCase(locationRepo)
	hardDeleteUC := usecases.NewHardDeleteLocationUseCase(locationRepo)

	// Assign to container
	c.Location = &LocationContainer{
		GetLocationUseCase:         getLocationUC,
		CreateLocationUseCase:      createLocationUC,
		ListLocationsUseCase:       listLocationsUC,
		UpdateLocationUseCase:      updateLocationUC,
		DeleteLocationUseCase:      deleteLocationUC,
		BulkSoftDeleteLocationsUseCase: bulkSoftDeleteLocationsUC,
		UploadLocationMediaUseCase: uploadLocationMediaUC,
		RestoreLocationUseCase:      restoreUC,
		BulkRestoreLocationsUseCase: bulkRestoreUC,
		HardDeleteLocationUseCase:   hardDeleteUC,
	}
}
