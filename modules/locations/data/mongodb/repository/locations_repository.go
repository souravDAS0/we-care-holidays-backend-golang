// internal/modules/locations/data/mongodb/repository/locations_repository.go
package repository

import (
	"context"

	"bitbucket.org/abhishek_fordel/we-care-holidays-backend-golang/internal/models"
	"bitbucket.org/abhishek_fordel/we-care-holidays-backend-golang/modules/locations/data/datasource"
	"bitbucket.org/abhishek_fordel/we-care-holidays-backend-golang/modules/locations/data/mongodb/model"
	"bitbucket.org/abhishek_fordel/we-care-holidays-backend-golang/modules/locations/domain/entity"
	"bitbucket.org/abhishek_fordel/we-care-holidays-backend-golang/modules/locations/domain/repository"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Ensure interface compliance
var _ repository.LocationRepository = (*LocationRepositoryMongo)(nil)

// LocationRepositoryMongo implements the LocationRepository interface using MongoDB.
type LocationRepositoryMongo struct {
	datasource *datasource.MongoLocationDatasource
}

// NewLocationRepositoryMongo creates a new LocationRepositoryMongo.
func NewLocationRepositoryMongo(ds *datasource.MongoLocationDatasource) *LocationRepositoryMongo {
	return &LocationRepositoryMongo{datasource: ds}
}

// FindAll retrieves locations matching the filter with pagination.
// Soft-deleted records are excluded unless filter overrides them.
func (r *LocationRepositoryMongo) FindAll(ctx context.Context, filter map[string]interface{}, page, limit int) ([]*entity.Location, int64, error) {
	models, total, err := r.datasource.FindByFilters(ctx, filter, page, limit)
	if err != nil {
		return nil, 0, err
	}

	locations := make([]*entity.Location, len(models))
	for i, m := range models {
		e := m.ToEntity()
		locations[i] = &e
	}
	return locations, total, nil
}

// FindByID fetches a single location by its ID.
// Returns (nil, nil) if not found or soft-deleted.
func (r *LocationRepositoryMongo) FindByID(ctx context.Context, id primitive.ObjectID) (*entity.Location, error) {
	m, err := r.datasource.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if m == nil {
		return nil, nil
	}
	e := m.ToEntity()
	return &e, nil
}

// Create inserts a new location and back-fills its ID and timestamps.
func (r *LocationRepositoryMongo) Create(ctx context.Context, loc *entity.Location) error {
	m := model.FromEntity(loc)
	if err := r.datasource.Insert(ctx, m); err != nil {
		return err
	}
	loc.ID = m.ID
	loc.CreatedAt = m.CreatedAt
	loc.UpdatedAt = m.UpdatedAt
	return nil
}

// Update replaces an existing location document.
func (r *LocationRepositoryMongo) Update(ctx context.Context, loc *entity.Location) error {
	m := model.FromEntity(loc)
	return r.datasource.Update(ctx, m)
}

// SoftDelete marks a location as deleted by setting deletedAt.
func (r *LocationRepositoryMongo) SoftDelete(ctx context.Context, id primitive.ObjectID) (bool, error) {
	return r.datasource.SoftDelete(ctx, id)
}

// Restore clears deletedAt, restoring a soft-deleted location.
func (r *LocationRepositoryMongo) Restore(ctx context.Context, id primitive.ObjectID) (bool, error) {
	return r.datasource.Restore(ctx, id)
}

// HardDelete permanently removes a location from the database.
func (r *LocationRepositoryMongo) HardDelete(ctx context.Context, id primitive.ObjectID) (bool, error) {
	return r.datasource.HardDelete(ctx, id)
}

// BulkSoftDelete marks multiple locations as deleted and returns detailed results.
func (r *LocationRepositoryMongo) BulkSoftDelete(ctx context.Context, ids []string) (*models.BulkDeleteResponse, error) {
	result := &models.BulkDeleteResponse{
		RequestedIDs: ids,
		DeletedIDs:   []string{},
		InvalidIDs:   []string{},
		NotFoundIDs:  []string{},
		DeletedCount: 0,
	}

	var validObjectIDs []primitive.ObjectID
	for _, idStr := range ids {
		objectID, err := primitive.ObjectIDFromHex(idStr)
		if err != nil {
			result.InvalidIDs = append(result.InvalidIDs, idStr)
			continue
		}
		validObjectIDs = append(validObjectIDs, objectID)
	}

	// If no valid IDs, return early
	if len(validObjectIDs) == 0 {
		return result, nil
	}

	deletedObjectIDs, err := r.datasource.BulkSoftDelete(ctx, validObjectIDs)

	if err != nil {
		return result, err
	}

	// Convert deleted ObjectIDs back to strings
	for _, objID := range deletedObjectIDs {
		result.DeletedIDs = append(result.DeletedIDs, objID.Hex())
	}
	result.DeletedCount = len(result.DeletedIDs)

	deletedIDsMap := make(map[string]bool)
	for _, deletedID := range result.DeletedIDs {
		deletedIDsMap[deletedID] = true
	}

	for _, objID := range validObjectIDs {
		idStr := objID.Hex()
		if !deletedIDsMap[idStr] {
			result.NotFoundIDs = append(result.NotFoundIDs, idStr)
		}
	}

	return result, nil
}

// BulkRestore restores multiple soft-deleted locations.
func (r *LocationRepositoryMongo) BulkRestore(ctx context.Context, ids []string) (*models.BulkRestoreResponse, error) {
	result := &models.BulkRestoreResponse{
		RequestedIDs: ids, 	
		RestoredIDs:  []string{},
		InvalidIDs:   []string{},
		NotFoundIDs:  []string{},
		RestoredCount: 0,
	}

	var validObjectIDs []primitive.ObjectID
	for _, idStr := range ids {
		objectID, err := primitive.ObjectIDFromHex(idStr)
		if err != nil {
			result.InvalidIDs = append(result.InvalidIDs, idStr)
			continue
		}
		validObjectIDs = append(validObjectIDs, objectID)
	}

	if len(validObjectIDs) == 0 {
		return result, nil
	}

	restoredObjectIDs, err := r.datasource.BulkRestore(ctx, validObjectIDs)

	if err != nil {
		return result, err
	}

	for _, objID := range restoredObjectIDs {
		result.RestoredIDs = append(result.RestoredIDs, objID.Hex())
	}
	result.RestoredCount = len(result.RestoredIDs)

	// Create a map for efficient lookup
	restoredIDsMap := make(map[string]bool)
	for _, restoredID := range result.RestoredIDs {
		restoredIDsMap[restoredID] = true
	}

	// Find IDs that were not restored (not found)
	for _, objID := range validObjectIDs {
		idStr := objID.Hex()
		if !restoredIDsMap[idStr] {
			result.NotFoundIDs = append(result.NotFoundIDs, idStr)
		}
	}

	return result, nil
}
