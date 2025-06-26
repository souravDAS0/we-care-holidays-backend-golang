package usecases

import (
	"context"

	"bitbucket.org/abhishek_fordel/we-care-holidays-backend-golang/modules/organizations/domain/repository"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// OrganizationUseCase implements the organization business logic
type HardDeleteOrganizationUseCase struct {
	repo repository.OrganizationRepository
}

func NewHardDeleteOrganizationUseCase(repo repository.OrganizationRepository) *HardDeleteOrganizationUseCase {
	return &HardDeleteOrganizationUseCase{
		repo: repo,
	}
}

// HardDeleteOrganization permanently removes an organization (admin/cleanup only)
func (uc *HardDeleteOrganizationUseCase) Execute(ctx context.Context, id primitive.ObjectID) (bool, error) {
	return uc.repo.HardDelete(ctx, id)
}
