package usecases

import (
	"context"
	"errors"

	repo "bitbucket.org/abhishek_fordel/we-care-holidays-backend-golang/modules/locations/domain/repository"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// UploadLocationMediaUseCase updates media URLs for a location
type UploadLocationMediaUseCase struct {
	repo repo.LocationRepository
}

// Constructor
func NewUploadLocationMediaUseCase(r repo.LocationRepository) *UploadLocationMediaUseCase {
	return &UploadLocationMediaUseCase{repo: r}
}

// Execute handles persisting photo/video URLs to the location
func (uc *UploadLocationMediaUseCase) Execute(
	ctx context.Context,
	locationID string,
	photos []string,
	videos []string,
) error {
	// Convert to ObjectID
	objectID, err := primitive.ObjectIDFromHex(locationID)
	if err != nil {
		return errors.New("invalid location ID format")
	}

	// Fetch existing location
	location, err := uc.repo.FindByID(ctx, objectID)
	if err != nil {
		return err
	}
	if location == nil {
		return errors.New("location not found")
	}

	// Append media (don't overwrite)
	location.MediaURLs.Photos = append(location.MediaURLs.Photos, photos...)
	location.MediaURLs.Videos = append(location.MediaURLs.Videos, videos...)

	// Persist update
	return uc.repo.Update(ctx, location)
}
