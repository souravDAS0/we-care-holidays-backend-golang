package usecases

import (
	"context"

	"bitbucket.org/abhishek_fordel/we-care-holidays-backend-golang/modules/roles/domain/entity"
	"bitbucket.org/abhishek_fordel/we-care-holidays-backend-golang/modules/roles/domain/repository"
	"go.mongodb.org/mongo-driver/bson/primitive"
)


type GetRoleUseCase struct {
	repo repository.RoleRepository
}

func NewGetRoleUseCase(repo repository.RoleRepository) *GetRoleUseCase {
	return &GetRoleUseCase{
		repo: repo,
	}
}

// Execute retrieves a role by its ID
func (uc *GetRoleUseCase) Execute(ctx context.Context, id primitive.ObjectID) (*entity.Role, error) {
	return uc.repo.GetByID(ctx,id)
}