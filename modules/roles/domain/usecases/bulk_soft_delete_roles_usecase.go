package usecases

import (
	"context"

	"bitbucket.org/abhishek_fordel/we-care-holidays-backend-golang/internal/models"
	"bitbucket.org/abhishek_fordel/we-care-holidays-backend-golang/modules/roles/domain/repository"
)

// RoleUseCase implements the Role business logic
type BulkSoftDeleteRolesUseCase struct {
	repo repository.RoleRepository
}

func NewBulkSoftDeleteRolesUseCase(repo repository.RoleRepository) *BulkSoftDeleteRolesUseCase {
	return &BulkSoftDeleteRolesUseCase{
		repo: repo,
	}
}

// BulkSoftDeleteRoles marks multiple Roles as deleted
func (uc *BulkSoftDeleteRolesUseCase) Execute(ctx context.Context, ids []string)  (*models.BulkDeleteResponse, error)  {
	return uc.repo.BulkSoftDelete(ctx, ids)
}
