package usecases

import (
	"context"

	"bitbucket.org/abhishek_fordel/we-care-holidays-backend-golang/modules/roles/domain/repository"
	"go.mongodb.org/mongo-driver/bson/primitive"
)


type HardDeleteRoleUseCase struct {
	repo repository.RoleRepository
}

func NewHardDeleteRoleUseCase(repo repository.RoleRepository) *HardDeleteRoleUseCase {
	return &HardDeleteRoleUseCase{
		repo: repo,
	}
}

func (uc *HardDeleteRoleUseCase) Execute(ctx context.Context, id primitive.ObjectID) (bool, error) {
	return uc.repo.HardDelete(ctx,id)
}