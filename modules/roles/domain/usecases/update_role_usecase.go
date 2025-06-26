package usecases

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	permissionEntity "bitbucket.org/abhishek_fordel/we-care-holidays-backend-golang/modules/permissions/domain/entity"
	permissionRepo "bitbucket.org/abhishek_fordel/we-care-holidays-backend-golang/modules/permissions/domain/repository"
	"bitbucket.org/abhishek_fordel/we-care-holidays-backend-golang/modules/roles/domain/entity"
	"bitbucket.org/abhishek_fordel/we-care-holidays-backend-golang/modules/roles/domain/repository"
)

// RoleUseCase implements the Role business logic
type UpdateRoleUseCase struct {
	roleRepo       repository.RoleRepository
	permissionRepo permissionRepo.PermissionRepository
}

func NewUpdateRoleUseCase(roleRepo repository.RoleRepository, permissionRepo permissionRepo.PermissionRepository) *UpdateRoleUseCase {
	return &UpdateRoleUseCase{
		roleRepo:       roleRepo,
		permissionRepo: permissionRepo,
	}
}

// UpdateRole updates an existing Role
func (uc *UpdateRoleUseCase) Execute(ctx context.Context, role *entity.Role) error {
	if role.ID.IsZero() {
		return errors.New("role ID is required")
	}

	currentRole, err := uc.roleRepo.GetByID(ctx, role.ID)
	if err != nil {
		return err
	}
	if currentRole == nil {
		return errors.New("role not found")
	}

	// Check name uniqueness only if name is being changed
	if role.Name != currentRole.Name {
		exists, err := uc.roleRepo.ExistsByName(ctx, role.Name)
		if err != nil {
			return err
		}
		if exists {
			return errors.New("role with this name already exists")
		}
	}

	// Validate permission IDs if any are provided
	if len(role.Permissions) > 0 {
		if err := uc.validatePermissions(ctx, role.Permissions); err != nil {
			return err
		}
	}

	// Set updated timestamp
	role.UpdatedAt = time.Now()

	return uc.roleRepo.Update(ctx, role)
}

func (uc *UpdateRoleUseCase) validatePermissions(ctx context.Context, permissions []string) error {
	for _, perm := range permissions {
		resource := strings.Split(perm, ":")[0]
		action := strings.Split(perm, ":")[1]

		exists, err := uc.permissionRepo.ExistsByResourceAction(ctx, resource, permissionEntity.PermissionAction(action))
		if err != nil {
			return fmt.Errorf("failed to validate permission %s: %w", perm, err)
		}
		if !exists {
			return fmt.Errorf("permission  %s does not exist", perm)
		}
	}
	return nil
}
