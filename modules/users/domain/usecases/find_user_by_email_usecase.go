package usecases

import (
	"context"

	"bitbucket.org/abhishek_fordel/we-care-holidays-backend-golang/modules/users/domain/entity"
	"bitbucket.org/abhishek_fordel/we-care-holidays-backend-golang/modules/users/domain/repository"
)

type FindUserByEmailUsecase struct {
	repo repository.UserRepository
}

func NewFindUserByEmailUsecase(repo repository.UserRepository) *FindUserByEmailUsecase {
	return &FindUserByEmailUsecase{
		repo: repo,
	}
}

// Execute retrieves a role by its ID
func (uc *FindUserByEmailUsecase) Execute(ctx context.Context, email string) (*entity.User, error) {
	return uc.repo.GetByEmail(ctx, email)
}
