package usecases

import (
	"context"

	"bitbucket.org/abhishek_fordel/we-care-holidays-backend-golang/internal/models"
	"bitbucket.org/abhishek_fordel/we-care-holidays-backend-golang/modules/users/domain/repository"
)

// BulkRestoreUsersUseCase implements the bulk restore business logic
type BulkRestoreUsersUseCase struct {
	repo repository.UserRepository
}

func NewBulkRestoreUsersUseCase(repo repository.UserRepository) *BulkRestoreUsersUseCase {
	return &BulkRestoreUsersUseCase{
		repo: repo,
	}
}

// Execute restores multiple permissions from soft-deleted state
func (uc *BulkRestoreUsersUseCase) Execute(ctx context.Context, ids []string) (*models.BulkRestoreResponse, error) {
	return uc.repo.BulkRestore(ctx, ids)
}
