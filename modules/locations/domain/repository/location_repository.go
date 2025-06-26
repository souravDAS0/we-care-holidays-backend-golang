package repository

import (
	"context"

	"bitbucket.org/abhishek_fordel/we-care-holidays-backend-golang/internal/models"
	"bitbucket.org/abhishek_fordel/we-care-holidays-backend-golang/modules/locations/domain/entity"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// LocationRepository defines persistence operations
type LocationRepository interface {
	// FindAll retrieves locations matching the filter, with pagination.
	// By default it should exclude soft-deleted locations unless filter includes them.
	FindAll(ctx context.Context, filter map[string]interface{}, page, limit int) ([]*entity.Location, int64, error)

	// FindByID finds a single location by its ObjectID.
	// Should return nil if not found or soft-deleted (unless specifically included in filter).
	FindByID(ctx context.Context, id primitive.ObjectID) (*entity.Location, error)

	// Create inserts a new location and sets its ID on the entity.
	Create(ctx context.Context, loc *entity.Location) error

	// Update applies changes to an existing location document.
	Update(ctx context.Context, loc *entity.Location) error

	// SoftDelete marks a location as deleted by setting deletedAt timestamp.
	SoftDelete(ctx context.Context, id primitive.ObjectID) (bool, error)

	// BulkSoftDelete marks multiple locations as deleted and returns detailed results.
	BulkSoftDelete(ctx context.Context, ids []string) (*models.BulkDeleteResponse, error)

	// Restore un-does a soft-delete by clearing the deletedAt timestamp.
	Restore(ctx context.Context, id primitive.ObjectID) (bool, error)

	// BulkRestore restores multiple soft-deleted locations.
	BulkRestore(ctx context.Context, ids []string) (*models.BulkRestoreResponse, error)

	// HardDelete permanently removes a location from the database.
	HardDelete(ctx context.Context, id primitive.ObjectID) (bool, error)
}
