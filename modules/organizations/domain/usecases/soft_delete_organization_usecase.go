package usecases

import (
	"context"

	"bitbucket.org/abhishek_fordel/we-care-holidays-backend-golang/modules/organizations/domain/repository"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// OrganizationUseCase implements the organization business logic
type SoftDeleteOrganizationUseCase struct {
	repo repository.OrganizationRepository
}

func NewSoftDeleteOrganizationUseCase(repo repository.OrganizationRepository) *SoftDeleteOrganizationUseCase {
	return &SoftDeleteOrganizationUseCase{
		repo: repo,
	}
}

// SoftDeleteOrganization marks an organization as deleted without removing it
func (uc *SoftDeleteOrganizationUseCase) Execute(ctx context.Context, id primitive.ObjectID) (bool, error) {
	return uc.repo.SoftDelete(ctx, id)
}
