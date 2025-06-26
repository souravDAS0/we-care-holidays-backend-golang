package usecases

import (
	"context"

	repo "bitbucket.org/abhishek_fordel/we-care-holidays-backend-golang/modules/locations/domain/repository"
	en "bitbucket.org/abhishek_fordel/we-care-holidays-backend-golang/modules/locations/domain/entity"
)

// ListLocationsUseCase retrieves paginated, filtered list
type ListLocationsUseCase struct { Repo repo.LocationRepository }

func NewListLocationsUseCase(r repo.LocationRepository) *ListLocationsUseCase {
	return &ListLocationsUseCase{Repo: r}
}

func (uc *ListLocationsUseCase) Execute(ctx context.Context, filter map[string]interface{}, page, limit int) ([]*en.Location, int64, error) {
	return uc.Repo.FindAll(ctx, filter, page, limit)
}
