package usecases

import (
	"context"

	"bitbucket.org/abhishek_fordel/we-care-holidays-backend-golang/modules/permissions/domain/entity"
	"bitbucket.org/abhishek_fordel/we-care-holidays-backend-golang/modules/permissions/domain/repository"
)

// PermissionUseCase implements the Permission business logic
type ListPermissionsUseCase struct {
	repo repository.PermissionRepository
}

func NewListPermissionUseCase(repo repository.PermissionRepository) *ListPermissionsUseCase {
	return &ListPermissionsUseCase{
		repo: repo,
	}
}

// ListPermissions retrieves a list of Permissions with pagination
func (uc *ListPermissionsUseCase) Execute(ctx context.Context, filter map[string]interface{}, page, limit int) ([]*entity.Permission, int64, error) {
	return uc.repo.List(ctx, filter, page, limit)
}
