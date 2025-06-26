package usecases

import (
	"context"

	"bitbucket.org/abhishek_fordel/we-care-holidays-backend-golang/modules/users/domain/entity"
	"bitbucket.org/abhishek_fordel/we-care-holidays-backend-golang/modules/users/domain/repository"
)

type ListUsersUseCase struct {
	repo repository.UserRepository
}

func NewListUsersUseCase(repo repository.UserRepository) *ListUsersUseCase {
	return &ListUsersUseCase{
		repo: repo,
	}
}

func (uc *ListUsersUseCase) Execute(ctx context.Context, filter map[string]interface{}, page, limit int) ([]*entity.User, int64, error) {
	return uc.repo.List(ctx, filter, page, limit)
}
