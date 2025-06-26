package usecases

import (
	"context"
	"errors"
	"strings"
	"time"

	en "bitbucket.org/abhishek_fordel/we-care-holidays-backend-golang/modules/locations/domain/entity"
	repo "bitbucket.org/abhishek_fordel/we-care-holidays-backend-golang/modules/locations/domain/repository"
)

type CreateLocationUseCase struct {
	Repo repo.LocationRepository
}

func NewCreateLocationUseCase(r repo.LocationRepository) *CreateLocationUseCase {
	return &CreateLocationUseCase{Repo: r}
}

func (uc *CreateLocationUseCase) Execute(ctx context.Context, loc *en.Location) error {
	// Normalize name (optional: trim and lowercase)
	locationName := strings.TrimSpace(loc.Name)

	// Check for existing location with the same name (excluding soft-deleted)
	filter := map[string]interface{}{
		"name":       locationName,
		"deletedAt": nil, // exclude soft-deleted
	}

	existing, _, err := uc.Repo.FindAll(ctx, filter, 1, 1)
	if err != nil {
		return err
	}
	if len(existing) > 0 {
		return errors.New("location with the same name already exists")
	}

	// Set default timestamps
	loc.Name = locationName // normalize
	loc.CreatedAt = time.Now()
	loc.UpdatedAt = time.Now()
	loc.DeletedAt = nil

	return uc.Repo.Create(ctx, loc)
}
