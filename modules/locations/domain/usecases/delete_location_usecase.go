package usecases

import (
	"context"

	repo "bitbucket.org/abhishek_fordel/we-care-holidays-backend-golang/modules/locations/domain/repository"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// DeleteLocationUseCase soft-deletes a single location
type DeleteLocationUseCase struct{ Repo repo.LocationRepository }

func NewDeleteLocationUseCase(r repo.LocationRepository) *DeleteLocationUseCase {
	return &DeleteLocationUseCase{Repo: r}
}

func (uc *DeleteLocationUseCase) Execute(ctx context.Context, id primitive.ObjectID) (bool, error) {
	return uc.Repo.SoftDelete(ctx, id)
}
