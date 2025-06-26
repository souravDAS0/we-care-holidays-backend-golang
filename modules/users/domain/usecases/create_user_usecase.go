package usecases

import (
	"context"
	"errors"
	"time"

	"bitbucket.org/abhishek_fordel/we-care-holidays-backend-golang/internal/utils"
	orgRepo "bitbucket.org/abhishek_fordel/we-care-holidays-backend-golang/modules/organizations/domain/repository"
	roleRepo "bitbucket.org/abhishek_fordel/we-care-holidays-backend-golang/modules/roles/domain/repository"
	"bitbucket.org/abhishek_fordel/we-care-holidays-backend-golang/modules/users/domain/entity"
	"bitbucket.org/abhishek_fordel/we-care-holidays-backend-golang/modules/users/domain/repository"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type CreateUserUseCase struct {
	userRepo repository.UserRepository
	roleRepo roleRepo.RoleRepository
	orgRepo  orgRepo.OrganizationRepository
}

func NewCreateUserUseCase(userRepo repository.UserRepository, roleRepo roleRepo.RoleRepository, orgRepo orgRepo.OrganizationRepository) *CreateUserUseCase {
	return &CreateUserUseCase{
		userRepo: userRepo,
		roleRepo: roleRepo,
		orgRepo:  orgRepo,
	}
}

// Execute creates a new user
func (uc *CreateUserUseCase) Execute(ctx context.Context, user *entity.User) error {
	// Set default status if not provided
	if user.Status == "" {
		user.Status = entity.UserStatusInvited
	}

	// Check if user with email already exists
	if len(user.Emails) > 0 {
		exists, err := uc.userRepo.ExistsByEmail(ctx, user.GetPrimaryEmail())
		if err != nil {
			return err
		}
		if exists {
			return errors.New("user with this email already exists")
		}
	}

	// Check if user with phone already exists
	if len(user.Phones) > 0 {
		exists, err := uc.userRepo.ExistsByPhone(ctx, user.GetPrimaryPhone())
		if err != nil {
			return err
		}
		if exists {
			return errors.New("user with this phone already exists")
		}
	}

	if user.RoleID == "" {
		return errors.New("role ID is required")
	}

	roleId, err := primitive.ObjectIDFromHex(user.RoleID)
	if err != nil {
		return err
	}

	// Validate role ID
	exists, err := uc.roleRepo.ExistsByID(ctx, roleId)
	if err != nil {
		return err
	}
	if !exists {
		return errors.New("role with this ID does not exist")
	}

	if user.OrganizationID != "" {

		organizationId, err := primitive.ObjectIDFromHex(user.OrganizationID)
		if err != nil {
			return err
		}
		orgExists, err := uc.orgRepo.ExistsByID(ctx, organizationId)
		if err != nil {
			return err
		}
		if !orgExists {
			return errors.New("organization with this ID does not exist")
		}
	}

	if user.Password != "" {
		hashedPassword, err := utils.HashPassword(user.Password)
		if err != nil {
			return err
		}
		user.Password = hashedPassword
	}

	user.CreatedAt = time.Now()
	user.UpdatedAt = time.Now()

	user.DeletedAt = nil

	return uc.userRepo.Create(ctx, user)
}
