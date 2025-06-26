package usecases

import (
	"context"
	"errors"
	"time"

	"bitbucket.org/abhishek_fordel/we-care-holidays-backend-golang/modules/permissions/domain/entity"
	"bitbucket.org/abhishek_fordel/we-care-holidays-backend-golang/modules/permissions/domain/repository"
)

// PermissionUseCase implements the Permission business logic
type CreatePermissionUseCase struct {
	repo repository.PermissionRepository
}

func NewCreatePermissionUseCase(repo repository.PermissionRepository) *CreatePermissionUseCase {
	return &CreatePermissionUseCase{
		repo: repo,
	}
}

// CreatePermission creates a new Permission
func (uc *CreatePermissionUseCase) Execute(ctx context.Context, permission *entity.Permission) error {

	exists, err := uc.repo.ExistsByResourceAction(ctx, permission.Resource, permission.Action)
	if err != nil {
		return err
	}
	if exists {
		return errors.New("permission with this resource, action, and scope already exists")
	}

	// Set default values
	permission.CreatedAt = time.Now()
	permission.UpdatedAt = time.Now()

	return uc.repo.Create(ctx, permission)
}
