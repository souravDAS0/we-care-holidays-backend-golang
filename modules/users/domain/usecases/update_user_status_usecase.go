package usecases

import (
	"context"
	"errors"

	"bitbucket.org/abhishek_fordel/we-care-holidays-backend-golang/modules/users/domain/entity"
	"bitbucket.org/abhishek_fordel/we-care-holidays-backend-golang/modules/users/domain/repository"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type UpdateUserStatusUseCase struct {
	repo repository.UserRepository
}

func NewUpdateUserStatusUseCase(repo repository.UserRepository) *UpdateUserStatusUseCase {
	return &UpdateUserStatusUseCase{
		repo: repo,
	}
}

// Execute updates the status of a user
func (uc *UpdateUserStatusUseCase) Execute(ctx context.Context, id primitive.ObjectID, status string) error {
	// Validate status
	validStatuses := map[string]bool{
		"Invited":   true,
		"Active":    true,
		"Suspended": true,
		"Removed":   true,
	}
	if !validStatuses[status] {
		return errors.New("invalid status")
	}

	// Check if user exists
	user, err := uc.repo.GetByID(ctx, id)
	if err != nil {
		return err
	}
	if user == nil {
		return errors.New("user not found")
	}

	// Update status
	user.Status = entity.UserStatus(status)
	return uc.repo.Update(ctx, user)
}
