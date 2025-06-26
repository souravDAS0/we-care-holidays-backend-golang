package usecases

import (
	"context"

	"bitbucket.org/abhishek_fordel/we-care-holidays-backend-golang/modules/users/domain/repository"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// UserUseCase implements the organization business logic
type HardDeleteUserUseCase struct {
	repo repository.UserRepository
}

func NewHardDeleteUserUseCase(repo repository.UserRepository) *HardDeleteUserUseCase {
	return &HardDeleteUserUseCase{
		repo: repo,
	}
}

// HardDeleteUser permanently removes an organization (admin/cleanup only)
func (uc *HardDeleteUserUseCase) Execute(ctx context.Context, id primitive.ObjectID) (bool, error) {
	return uc.repo.HardDelete(ctx, id)
}
