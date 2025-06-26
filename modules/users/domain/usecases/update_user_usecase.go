package usecases

import (
	"context"
	"errors"

	"bitbucket.org/abhishek_fordel/we-care-holidays-backend-golang/internal/utils"
	"bitbucket.org/abhishek_fordel/we-care-holidays-backend-golang/modules/users/domain/entity"
	"bitbucket.org/abhishek_fordel/we-care-holidays-backend-golang/modules/users/domain/repository"
)

type UpdateUserUseCase struct {
	repo repository.UserRepository
}

func NewUpdateUserUseCase(repo repository.UserRepository) *UpdateUserUseCase {
	return &UpdateUserUseCase{
		repo: repo,
	}
}

// Execute updates an existing user
func (uc *UpdateUserUseCase) Execute(ctx context.Context, user *entity.User) error {
	// Check if user exists
	existingUser, err := uc.repo.GetByID(ctx, user.ID)
	if err != nil {
		return err
	}
	if existingUser == nil {
		return errors.New("user not found")
	}

	// If email is being updated, check for duplicates
	if len(user.Emails) > 0 && user.GetPrimaryEmail() != existingUser.GetPrimaryEmail() {
		exists, err := uc.repo.ExistsByEmail(ctx, user.GetPrimaryEmail())
		if err != nil {
			return err
		}
		if exists {
			return errors.New("user with this email already exists")
		}
	}

	// If phone is being updated, check for duplicates
	if len(user.Phones) > 0 && user.GetPrimaryPhone() != existingUser.GetPrimaryPhone() {
		exists, err := uc.repo.ExistsByPhone(ctx, user.GetPrimaryPhone())
		if err != nil {
			return err
		}
		if exists {
			return errors.New("user with this phone already exists")
		}
	}

	if user.Password != "" {
		hashedPassword, err := utils.HashPassword(user.Password)
		if err != nil {
			return err
		}
		user.Password = hashedPassword
	}

	return uc.repo.Update(ctx, user)
}
