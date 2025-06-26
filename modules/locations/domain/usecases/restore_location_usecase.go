package usecases

import (
	"context"

	repo "bitbucket.org/abhishek_fordel/we-care-holidays-backend-golang/modules/locations/domain/repository"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// RestoreLocationUseCase undoes a soft-delete by clearing deletedAt.
type RestoreLocationUseCase struct {
    repo repo.LocationRepository
}

// NewRestoreLocationUseCase creates a RestoreLocationUseCase.
func NewRestoreLocationUseCase(r repo.LocationRepository) *RestoreLocationUseCase {
    return &RestoreLocationUseCase{repo: r}
}

// Execute restores the location with the given ID.
// Returns true if it was restored, false if not found or already active.
func (uc *RestoreLocationUseCase) Execute(
    ctx context.Context,
    id primitive.ObjectID,
) (bool, error) {
    return uc.repo.Restore(ctx, id)
}
