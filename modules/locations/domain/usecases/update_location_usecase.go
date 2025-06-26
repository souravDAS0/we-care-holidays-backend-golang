package usecases

import (
	"context"

	repo "bitbucket.org/abhishek_fordel/we-care-holidays-backend-golang/modules/locations/domain/repository"
	en "bitbucket.org/abhishek_fordel/we-care-holidays-backend-golang/modules/locations/domain/entity"
)

// UpdateLocationUseCase modifies an existing location
type UpdateLocationUseCase struct{ Repo repo.LocationRepository }

func NewUpdateLocationUseCase(r repo.LocationRepository) *UpdateLocationUseCase {
	return &UpdateLocationUseCase{Repo: r}
}

func (uc *UpdateLocationUseCase) Execute(ctx context.Context, loc *en.Location) error {
	return uc.Repo.Update(ctx, loc)
}