package usecases

import (
	"context"

	"bitbucket.org/abhishek_fordel/we-care-holidays-backend-golang/internal/models"
	"bitbucket.org/abhishek_fordel/we-care-holidays-backend-golang/modules/roles/domain/repository"
)

// BulkRestoreRolesUseCase implements the bulk restore business logic
type BulkRestoreRolesUseCase struct {
	repo repository.RoleRepository
}

func NewBulkRestoreRolesUseCase(repo repository.RoleRepository) *BulkRestoreRolesUseCase {
	return &BulkRestoreRolesUseCase{
		repo: repo,
	}
}

// Execute restores multiple permissions from soft-deleted state
func (uc *BulkRestoreRolesUseCase) Execute(ctx context.Context, ids []string) (*models.BulkRestoreResponse, error) {
	return uc.repo.BulkRestore(ctx, ids)
}
