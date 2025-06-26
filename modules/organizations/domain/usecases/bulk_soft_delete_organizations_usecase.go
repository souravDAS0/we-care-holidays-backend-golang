package usecases

import (
	"context"

	"bitbucket.org/abhishek_fordel/we-care-holidays-backend-golang/internal/models"
	"bitbucket.org/abhishek_fordel/we-care-holidays-backend-golang/modules/organizations/domain/repository"
)

// OrganizationUseCase implements the organization business logic
type BulkSoftDeleteOrganizationsUseCase struct {
	repo repository.OrganizationRepository
}

func NewBulkSoftDeleteOrganizationsUseCase(repo repository.OrganizationRepository) *BulkSoftDeleteOrganizationsUseCase {
	return &BulkSoftDeleteOrganizationsUseCase{
		repo: repo,
	}
}

// BulkSoftDeleteOrganizations marks multiple organizations as deleted
func (uc *BulkSoftDeleteOrganizationsUseCase) Execute(ctx context.Context, ids []string) (*models.BulkDeleteResponse, error) {
	return uc.repo.BulkSoftDelete(ctx, ids)
}
