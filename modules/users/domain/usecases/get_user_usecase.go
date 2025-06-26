package usecases

import (
	"context"

	"bitbucket.org/abhishek_fordel/we-care-holidays-backend-golang/modules/users/domain/entity"
	"bitbucket.org/abhishek_fordel/we-care-holidays-backend-golang/modules/users/domain/repository"
	"go.mongodb.org/mongo-driver/bson/primitive"
)


type GetUserUseCase struct {
	repo repository.UserRepository
}

func NewGetUserUseCase(repo repository.UserRepository) *GetUserUseCase {
	return &GetUserUseCase{
		repo: repo,
	}
}

// Execute retrieves a role by its ID
func (uc *GetUserUseCase) Execute(ctx context.Context, id primitive.ObjectID) (*entity.User, error) {
	return uc.repo.GetByID(ctx,id)
}