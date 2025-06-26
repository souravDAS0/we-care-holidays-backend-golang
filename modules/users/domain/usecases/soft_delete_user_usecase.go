package usecases

import (
	"context"

	"bitbucket.org/abhishek_fordel/we-care-holidays-backend-golang/modules/users/domain/repository"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type SoftDeleteUserUseCase struct {
	repo repository.UserRepository
}

func NewSoftDeleteUserUseCase(repo repository.UserRepository) *SoftDeleteUserUseCase {
	return &SoftDeleteUserUseCase{
		repo: repo,
	}
}

// SoftDeleteUser permanently removes an organization (admin/cleanup only)
func (uc *SoftDeleteUserUseCase) Execute(ctx context.Context, id primitive.ObjectID) (bool, error) {
	return uc.repo.SoftDelete(ctx, id)
}
