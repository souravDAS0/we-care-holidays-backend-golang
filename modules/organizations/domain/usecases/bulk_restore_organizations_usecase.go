package usecases

import (
	"context"

	"bitbucket.org/abhishek_fordel/we-care-holidays-backend-golang/internal/models"
	"bitbucket.org/abhishek_fordel/we-care-holidays-backend-golang/modules/organizations/domain/repository"
)

// BulkRestoreOrganizationssUseCase implements the bulk restore business logic
type BulkRestoreOrganizationsUseCase struct {
	repo repository.OrganizationRepository
}

func NewBulkRestoreOrganizationsUseCase(repo repository.OrganizationRepository) *BulkRestoreOrganizationsUseCase {
	return &BulkRestoreOrganizationsUseCase{
		repo: repo,
	}
}

// Execute restores multiple permissions from soft-deleted state
func (uc *BulkRestoreOrganizationsUseCase) Execute(ctx context.Context, ids []string) (*models.BulkRestoreResponse, error) {
	return uc.repo.BulkRestore(ctx, ids)
}
