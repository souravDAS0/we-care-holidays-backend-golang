package usecases

import (
	"context"

	"bitbucket.org/abhishek_fordel/we-care-holidays-backend-golang/modules/roles/domain/repository"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type RestoreRoleUseCase struct {
	repo repository.RoleRepository
}

func NewRestoreRoleUseCase(repo repository.RoleRepository) *RestoreRoleUseCase {
	return &RestoreRoleUseCase{
		repo: repo,
	}
}

// RestoreRole restores a soft-deleted Role
func (uc *RestoreRoleUseCase) Execute(ctx context.Context, id primitive.ObjectID) (bool, error) {
	return uc.repo.Restore(ctx, id)
}
