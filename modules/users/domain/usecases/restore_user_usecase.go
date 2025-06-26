package usecases

import (
	"context"

	"bitbucket.org/abhishek_fordel/we-care-holidays-backend-golang/modules/users/domain/repository"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type RestoreUserUseCase struct {
	repo repository.UserRepository
}

func NewRestoreUserUseCase(repo repository.UserRepository) *RestoreUserUseCase {
	return &RestoreUserUseCase{
		repo: repo,
	}
}

// RestoreUser permanently removes an organization (admin/cleanup only)
func (uc *RestoreUserUseCase) Execute(ctx context.Context, id primitive.ObjectID) (bool, error) {
	return uc.repo.Restore(ctx, id)
}
