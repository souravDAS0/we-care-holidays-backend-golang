package usecases

import (
	"context"

	"bitbucket.org/abhishek_fordel/we-care-holidays-backend-golang/modules/organizations/domain/repository"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// OrganizationUseCase implements the organization business logic
type RestoreOrganizationUseCase struct {
	repo repository.OrganizationRepository
}

func NewRestoreOrganizationUseCase(repo repository.OrganizationRepository) *RestoreOrganizationUseCase {
	return &RestoreOrganizationUseCase{
		repo: repo,
	}
}

// RestoreOrganization restores a soft-deleted organization
func (uc *RestoreOrganizationUseCase) Execute(ctx context.Context, id primitive.ObjectID) (bool, error) {
	return uc.repo.Restore(ctx, id)
}
