package usecases

import (
	"context"

	"bitbucket.org/abhishek_fordel/we-care-holidays-backend-golang/modules/roles/domain/entity"
	"bitbucket.org/abhishek_fordel/we-care-holidays-backend-golang/modules/roles/domain/repository"
)


type ListRolesUseCase struct {
	repo repository.RoleRepository
}

func NewListRolesUseCase(repo repository.RoleRepository) *ListRolesUseCase {
	return &ListRolesUseCase{
		repo: repo,
	}
}	



func (uc *ListRolesUseCase) Execute(ctx context.Context, filter map[string]interface{}, page, limit int) ([]*entity.Role, int64, error) {
	return uc.repo.List(ctx, filter, page, limit)
}