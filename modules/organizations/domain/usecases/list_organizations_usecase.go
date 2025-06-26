package usecases

import (
	"context"

	"bitbucket.org/abhishek_fordel/we-care-holidays-backend-golang/modules/organizations/domain/entity"
	"bitbucket.org/abhishek_fordel/we-care-holidays-backend-golang/modules/organizations/domain/repository"
)

// OrganizationUseCase implements the organization business logic
type ListOrganizationUseCase struct {
	repo repository.OrganizationRepository
}

func NewListOrganizationUseCase(repo repository.OrganizationRepository) *ListOrganizationUseCase {
	return &ListOrganizationUseCase{
		repo: repo,
	}
}

// ListOrganizations retrieves a list of organizations with pagination
func (uc *ListOrganizationUseCase) Execute(ctx context.Context, filter map[string]interface{}, page, limit int) ([]*entity.Organization, int64, error) {
	return uc.repo.FindAll(ctx, filter, page, limit)
}
