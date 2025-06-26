package usecases

import (
    "context"

    repo "bitbucket.org/abhishek_fordel/we-care-holidays-backend-golang/modules/locations/domain/repository"
    "go.mongodb.org/mongo-driver/bson/primitive"
)

// HardDeleteLocationUseCase permanently removes a location.
type HardDeleteLocationUseCase struct {
    repo repo.LocationRepository
}

// NewHardDeleteLocationUseCase creates a HardDeleteLocationUseCase.
func NewHardDeleteLocationUseCase(r repo.LocationRepository) *HardDeleteLocationUseCase {
    return &HardDeleteLocationUseCase{repo: r}
}

// Execute deletes the location with the given ID.
// Returns true if a document was deleted.
func (uc *HardDeleteLocationUseCase) Execute(
    ctx context.Context,
    id primitive.ObjectID,
) (bool, error) {
    return uc.repo.HardDelete(ctx, id)
}
