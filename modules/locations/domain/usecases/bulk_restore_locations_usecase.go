package usecases

import (
	"context"

	"bitbucket.org/abhishek_fordel/we-care-holidays-backend-golang/internal/models"
	repo "bitbucket.org/abhishek_fordel/we-care-holidays-backend-golang/modules/locations/domain/repository"
)

// BulkRestoreLocationsUseCase restores multiple soft-deleted locations.
type BulkRestoreLocationsUseCase struct {
    repo repo.LocationRepository
}

// NewBulkRestoreLocationsUseCase creates a BulkRestoreLocationsUseCase.
func NewBulkRestoreLocationsUseCase(r repo.LocationRepository) *BulkRestoreLocationsUseCase {
    return &BulkRestoreLocationsUseCase{repo: r}
}

// Execute attempts to restore each ID in the list,
// returning a BulkDeleteResponse that reports which succeeded, invalid, or not found.
func (uc *BulkRestoreLocationsUseCase) Execute(
    ctx context.Context,
    ids []string,
) (*models.BulkRestoreResponse, error) {
    return uc.repo.BulkRestore(ctx, ids)
}
