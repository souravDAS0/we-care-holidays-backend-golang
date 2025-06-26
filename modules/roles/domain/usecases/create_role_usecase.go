package usecases

import (
	"context"
	"errors"
	"fmt"
	"time"

	permissionRepo "bitbucket.org/abhishek_fordel/we-care-holidays-backend-golang/modules/permissions/domain/repository"
	"bitbucket.org/abhishek_fordel/we-care-holidays-backend-golang/modules/roles/domain/entity"
	"bitbucket.org/abhishek_fordel/we-care-holidays-backend-golang/modules/roles/domain/repository"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type CreateRoleUseCase struct {
	roleRepo       repository.RoleRepository
	permissionRepo permissionRepo.PermissionRepository
}

func NewCreateRoleUseCase(roleRepo repository.RoleRepository, permissionRepo permissionRepo.PermissionRepository) *CreateRoleUseCase {
	return &CreateRoleUseCase{
		roleRepo:       roleRepo,
		permissionRepo: permissionRepo,
	}
}

func (uc *CreateRoleUseCase) Execute(ctx context.Context, role *entity.Role) error {
	exists, err := uc.roleRepo.ExistsByName(ctx, role.Name)
	if err != nil {
		return err
	}

	if exists {
		return errors.New("role with this name already exists")
	}

	// Validate permission IDs if any are provided
	if len(role.Permissions) > 0 {
		if err := uc.validatePermissions(ctx, role.Permissions); err != nil {
			return err
		}
	}

	role.CreatedAt = time.Now()
	role.UpdatedAt = time.Now()

	role.DeletedAt = nil

	return uc.roleRepo.Create(ctx, role)
}

func (uc *CreateRoleUseCase) validatePermissions(ctx context.Context, permissionIDs []string) error {
	for _, permID := range permissionIDs {
		id, _ := primitive.ObjectIDFromHex(permID)
		exists, err := uc.permissionRepo.ExistsByID(ctx, id)
		if err != nil {
			return fmt.Errorf("failed to validate permission %s: %w", permID, err)
		}
		if !exists {
			return fmt.Errorf("permission with ID %s does not exist", permID)
		}
	}
	return nil
}
