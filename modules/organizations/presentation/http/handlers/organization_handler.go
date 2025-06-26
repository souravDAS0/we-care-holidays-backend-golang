package handlers

import (
	"bitbucket.org/abhishek_fordel/we-care-holidays-backend-golang/commons/services"
	"bitbucket.org/abhishek_fordel/we-care-holidays-backend-golang/modules/organizations/domain/usecases"
)

// OrganizationHandler handles HTTP requests for organizations
type OrganizationHandler struct {
	GetOrganizationUseCase             *usecases.GetOrganizationUseCase
	CreateOrganizationUseCase          *usecases.CreateOrganizationUseCase
	ListOrganizationUseCase            *usecases.ListOrganizationUseCase
	UpdateOrganizationUseCase          *usecases.UpdateOrganizationUseCase
	UpdateOrganizationStatusUseCase    *usecases.UpdateOrganizationStatusUseCase
	SoftDeleteOrganizationUseCase      *usecases.SoftDeleteOrganizationUseCase
	RestoreOrganizationUseCase         *usecases.RestoreOrganizationUseCase
	BulkSoftDeleteOrganizationsUseCase *usecases.BulkSoftDeleteOrganizationsUseCase
	BulkRestoreOrganizationsUseCase 	*usecases.BulkRestoreOrganizationsUseCase
	HardDeleteOrganizationUseCase      *usecases.HardDeleteOrganizationUseCase
	fileService                        services.FileService
}

// NewOrganizationHandler creates a new organization handler
func NewOrganizationHandler(GetOrganizationUseCase *usecases.GetOrganizationUseCase,
	fileService services.FileService,
	CreateOrganizationUseCase *usecases.CreateOrganizationUseCase,
	ListOrganizationUseCase *usecases.ListOrganizationUseCase,
	UpdateOrganizationUseCase *usecases.UpdateOrganizationUseCase,
	UpdateOrganizationStatusUseCase *usecases.UpdateOrganizationStatusUseCase,
	SoftDeleteOrganizationUseCase *usecases.SoftDeleteOrganizationUseCase,
	RestoreOrganizationUseCase *usecases.RestoreOrganizationUseCase,
	BulkSoftDeleteOrganizationsUseCase *usecases.BulkSoftDeleteOrganizationsUseCase,
	HardDeleteOrganizationUseCase *usecases.HardDeleteOrganizationUseCase,
	BulkRestoreOrganizationsUseCase 	*usecases.BulkRestoreOrganizationsUseCase,
) *OrganizationHandler {
	return &OrganizationHandler{
		fileService:                        fileService,
		GetOrganizationUseCase:             GetOrganizationUseCase,
		CreateOrganizationUseCase:          CreateOrganizationUseCase,
		ListOrganizationUseCase:            ListOrganizationUseCase,
		UpdateOrganizationUseCase:          UpdateOrganizationUseCase,
		UpdateOrganizationStatusUseCase:    UpdateOrganizationStatusUseCase,
		SoftDeleteOrganizationUseCase:      SoftDeleteOrganizationUseCase,
		RestoreOrganizationUseCase:         RestoreOrganizationUseCase,
		BulkSoftDeleteOrganizationsUseCase: BulkSoftDeleteOrganizationsUseCase,
		HardDeleteOrganizationUseCase:      HardDeleteOrganizationUseCase,
		BulkRestoreOrganizationsUseCase: 	BulkRestoreOrganizationsUseCase,
	}
}
