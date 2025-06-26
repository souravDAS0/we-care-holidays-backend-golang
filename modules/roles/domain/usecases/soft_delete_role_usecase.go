package usecases

import (
	"context"

	"bitbucket.org/abhishek_fordel/we-care-holidays-backend-golang/modules/roles/domain/repository"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// RoleUseCase implements the Role business logic
type SoftDeleteRoleUseCase struct {
	repo repository.RoleRepository
}

func NewSoftDeleteRoleUseCase(repo repository.RoleRepository) *SoftDeleteRoleUseCase {
	return &SoftDeleteRoleUseCase{
		repo: repo,
	}
}

// SoftDeleteRole marks an Role as deleted without removing it
func (uc *SoftDeleteRoleUseCase) Execute(ctx context.Context, id primitive.ObjectID) (bool, error) {
	return uc.repo.SoftDelete(ctx, id)
}
