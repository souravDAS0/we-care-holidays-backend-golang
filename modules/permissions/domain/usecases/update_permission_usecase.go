package usecases

import (
	"context"
	"errors"
	"time"

	"bitbucket.org/abhishek_fordel/we-care-holidays-backend-golang/modules/permissions/domain/entity"
	"bitbucket.org/abhishek_fordel/we-care-holidays-backend-golang/modules/permissions/domain/repository"
)

// PermissionUseCase implements the Permission business logic
type UpdatePermissionUseCase struct {
	repo repository.PermissionRepository
}

func NewUpdatePermissionUseCase(repo repository.PermissionRepository) *UpdatePermissionUseCase {
	return &UpdatePermissionUseCase{
		repo: repo,
	}
}

// UpdatePermission updates an existing Permission
func (uc *UpdatePermissionUseCase) Execute(ctx context.Context, permission *entity.Permission) error {
	if permission.ID.IsZero() {
		return errors.New("permission ID is required")
	}

	exists, err := uc.repo.ExistsByResourceActionExcluding(ctx,
		permission.Resource,
		permission.Action,
		permission.ID)
	if err != nil {
		return err
	}
	if exists {
		return errors.New("permission with this resource, action already exists")
	}

	// Set updated timestamp
	permission.UpdatedAt = time.Now()

	return uc.repo.Update(ctx, permission)
}
