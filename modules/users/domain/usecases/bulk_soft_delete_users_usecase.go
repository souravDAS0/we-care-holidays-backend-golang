package usecases

import (
	"context"

	"bitbucket.org/abhishek_fordel/we-care-holidays-backend-golang/internal/models"
	"bitbucket.org/abhishek_fordel/we-care-holidays-backend-golang/modules/users/domain/repository"
)

// BulkSoftDeleteUsersUseCase implements the bulk restore business logic
type BulkSoftDeleteUsersUseCase struct {
	repo repository.UserRepository
}

func NewBulkSoftDeleteUsersUseCase(repo repository.UserRepository) *BulkSoftDeleteUsersUseCase {
	return &BulkSoftDeleteUsersUseCase{
		repo: repo,
	}
}

// Execute restores multiple permissions from soft-deleted state
func (uc *BulkSoftDeleteUsersUseCase) Execute(ctx context.Context, ids []string) (*models.BulkDeleteResponse, error) {
	return uc.repo.BulkSoftDelete(ctx, ids)
}
