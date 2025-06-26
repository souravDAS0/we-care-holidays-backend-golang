package handlers

import (
	"bitbucket.org/abhishek_fordel/we-care-holidays-backend-golang/commons/services"
	"bitbucket.org/abhishek_fordel/we-care-holidays-backend-golang/modules/locations/domain/usecases"
)

// LocationHandler handles HTTP requests for locations
type LocationHandler struct {
	GetLocationUseCase             *usecases.GetLocationUseCase
	CreateLocationUseCase          *usecases.CreateLocationUseCase
	ListLocationsUseCase           *usecases.ListLocationsUseCase
	UpdateLocationUseCase          *usecases.UpdateLocationUseCase
	DeleteLocationUseCase          *usecases.DeleteLocationUseCase
	BulkSoftDeleteLocationsUseCase     *usecases.BulkSoftDeleteLocationsUseCase
	UploadLocationMediaUseCase     *usecases.UploadLocationMediaUseCase
	RestoreLocationUseCase       *usecases.RestoreLocationUseCase
	BulkRestoreLocationsUseCase  *usecases.BulkRestoreLocationsUseCase
	HardDeleteLocationUseCase    *usecases.HardDeleteLocationUseCase
	fileService                    services.FileService
}

// NewLocationHandler creates a new LocationHandler
func NewLocationHandler(
	GetLocationUseCase *usecases.GetLocationUseCase,
	fileService services.FileService,
	CreateLocationUseCase *usecases.CreateLocationUseCase,
	ListLocationsUseCase *usecases.ListLocationsUseCase,
	UpdateLocationUseCase *usecases.UpdateLocationUseCase,
	DeleteLocationUseCase *usecases.DeleteLocationUseCase,
	BulkSoftDeleteLocationsUseCase *usecases.BulkSoftDeleteLocationsUseCase,
	UploadLocationMediaUseCase *usecases.UploadLocationMediaUseCase,
	restoreUC *usecases.RestoreLocationUseCase,
	bulkRestoreUC *usecases.BulkRestoreLocationsUseCase,
	hardDeleteUC *usecases.HardDeleteLocationUseCase,
) *LocationHandler {
	return &LocationHandler{
		GetLocationUseCase:         GetLocationUseCase,
		fileService:                fileService,
		CreateLocationUseCase:      CreateLocationUseCase,
		ListLocationsUseCase:       ListLocationsUseCase,
		UpdateLocationUseCase:      UpdateLocationUseCase,
		DeleteLocationUseCase:      DeleteLocationUseCase,
		BulkSoftDeleteLocationsUseCase: BulkSoftDeleteLocationsUseCase,
		UploadLocationMediaUseCase: UploadLocationMediaUseCase,
		RestoreLocationUseCase:         restoreUC,
		BulkRestoreLocationsUseCase:    bulkRestoreUC,
		HardDeleteLocationUseCase:      hardDeleteUC,
	}
}

// Handler methods below (bind DTOs, validate, call use-cases, return JSON)
// e.g.:
// func (h *LocationHandler) ListLocations(c *gin.Context) { ... }
