package usecases

import (
	"context"

	repo "bitbucket.org/abhishek_fordel/we-care-holidays-backend-golang/modules/locations/domain/repository"
	en "bitbucket.org/abhishek_fordel/we-care-holidays-backend-golang/modules/locations/domain/entity"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// GetLocationUseCase fetches a single location
type GetLocationUseCase struct { Repo repo.LocationRepository }

func NewGetLocationUseCase(r repo.LocationRepository) *GetLocationUseCase {
	return &GetLocationUseCase{Repo: r}
}

func (uc *GetLocationUseCase) Execute(ctx context.Context, id primitive.ObjectID) (*en.Location, error) {
	return uc.Repo.FindByID(ctx, id)
}
