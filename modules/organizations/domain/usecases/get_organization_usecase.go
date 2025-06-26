package usecases

import (
	"context"

	"bitbucket.org/abhishek_fordel/we-care-holidays-backend-golang/modules/organizations/domain/entity"
	"bitbucket.org/abhishek_fordel/we-care-holidays-backend-golang/modules/organizations/domain/repository"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// OrganizationUseCase implements the organization business logic
type GetOrganizationUseCase struct {
	repo repository.OrganizationRepository
}

func NewGetOrganizationUseCase(repo repository.OrganizationRepository) *GetOrganizationUseCase {
	return &GetOrganizationUseCase{
		repo: repo,
	}
}

// GetOrganization retrieves an organization by its ID
func (uc *GetOrganizationUseCase) Execute(ctx context.Context, id primitive.ObjectID) (*entity.Organization, error) {
	return uc.repo.FindByID(ctx, id)
}
