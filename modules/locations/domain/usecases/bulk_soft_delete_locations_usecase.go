package usecases

import (
	"context"

	"bitbucket.org/abhishek_fordel/we-care-holidays-backend-golang/internal/models"
	repo "bitbucket.org/abhishek_fordel/we-care-holidays-backend-golang/modules/locations/domain/repository"
)

// BulkDeleteLocationsUseCase soft-deletes multiple
type BulkSoftDeleteLocationsUseCase struct{ Repo repo.LocationRepository }

func NewBulkSoftDeleteLocationsUseCase(r repo.LocationRepository) *BulkSoftDeleteLocationsUseCase {
	return &BulkSoftDeleteLocationsUseCase{Repo: r}
}

func (uc *BulkSoftDeleteLocationsUseCase) Execute(ctx context.Context, ids []string) (*models.BulkDeleteResponse, error) {
	return uc.Repo.BulkSoftDelete(ctx, ids)
}