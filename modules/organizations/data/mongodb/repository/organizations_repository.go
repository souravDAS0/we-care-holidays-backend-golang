// internal/modules/organizations/data/mongodb/repository/organizations_repository.go
package repository

import (
	"context"

	"bitbucket.org/abhishek_fordel/we-care-holidays-backend-golang/internal/models"
	"bitbucket.org/abhishek_fordel/we-care-holidays-backend-golang/modules/organizations/data/datasource"
	"bitbucket.org/abhishek_fordel/we-care-holidays-backend-golang/modules/organizations/data/mongodb/model"
	"bitbucket.org/abhishek_fordel/we-care-holidays-backend-golang/modules/organizations/domain/entity"
	"bitbucket.org/abhishek_fordel/we-care-holidays-backend-golang/modules/organizations/domain/repository"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Ensure interface compliance
var _ repository.OrganizationRepository = (*OrganizationRepositoryMongo)(nil)

// OrganizationRepositoryMongo implements the domain OrganizationRepository interface
type OrganizationRepositoryMongo struct {
	datasource *datasource.MongoOrganizationDatasource
}

// NewOrganizationRepositoryMongo creates a new instance of OrganizationRepositoryMongo
func NewOrganizationRepositoryMongo(ds *datasource.MongoOrganizationDatasource) *OrganizationRepositoryMongo {
	return &OrganizationRepositoryMongo{
		datasource: ds,
	}
}

// FindAll retrieves organizations with filtering and pagination
func (r *OrganizationRepositoryMongo) FindAll(ctx context.Context, filter map[string]interface{}, page, limit int) ([]*entity.Organization, int64, error) {
	organizationModels, totalCount, err := r.datasource.FindByFilters(ctx, filter, page, limit)
	if err != nil {
		return nil, 0, err
	}

	organizationEntities := make([]*entity.Organization, len(organizationModels))
	for i, model := range organizationModels {
		entity := model.ToEntity()
		organizationEntities[i] = &entity
	}

	return organizationEntities, totalCount, nil
}

// FindByID finds an organization by its ID
func (r *OrganizationRepositoryMongo) FindByID(ctx context.Context, id primitive.ObjectID) (*entity.Organization, error) {
	organizationModel, err := r.datasource.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}

	if organizationModel == nil {
		return nil, nil
	}

	organizationEntity := organizationModel.ToEntity()
	return &organizationEntity, nil
}

// FindBySlug finds an organization by its slug
func (r *OrganizationRepositoryMongo) FindBySlug(ctx context.Context, slug string) (*entity.Organization, error) {
	organizationModel, err := r.datasource.FindBySlug(ctx, slug)
	if err != nil {
		return nil, err
	}

	if organizationModel == nil {
		return nil, nil
	}

	organizationEntity := organizationModel.ToEntity()
	return &organizationEntity, nil
}

// Create inserts a new organization
func (r *OrganizationRepositoryMongo) Create(ctx context.Context, organization *entity.Organization) error {
	organizationModel := model.FromEntity(organization)

	err := r.datasource.Insert(ctx, organizationModel)
	if err != nil {
		return err
	}

	// Set back the generated ID in entity for downstream use
	organization.ID = organizationModel.ID
	organization.CreatedAt = organizationModel.CreatedAt
	organization.UpdatedAt = organizationModel.UpdatedAt
	return nil
}

// Update updates an existing organization
func (r *OrganizationRepositoryMongo) Update(ctx context.Context, organization *entity.Organization) error {
	organizationModel := model.FromEntity(organization)
	return r.datasource.Update(ctx, organizationModel)
}

// SoftDelete marks an organization as deleted without removing it from the database
func (r *OrganizationRepositoryMongo) SoftDelete(ctx context.Context, id primitive.ObjectID) (bool, error) {
	return r.datasource.SoftDelete(ctx, id)
}

// Restore restores a soft-deleted organization
func (r *OrganizationRepositoryMongo) Restore(ctx context.Context, id primitive.ObjectID) (bool, error) {
	return r.datasource.Restore(ctx, id)
}

// HardDelete permanently removes an organization
func (r *OrganizationRepositoryMongo) HardDelete(ctx context.Context, id primitive.ObjectID) (bool, error) {
	return r.datasource.HardDelete(ctx, id)
}

// BulkSoftDelete marks multiple organizations as deleted
func (r *OrganizationRepositoryMongo) BulkSoftDelete(ctx context.Context, ids []string) (*models.BulkDeleteResponse, error) {
	result := &models.BulkDeleteResponse{
		RequestedIDs: ids,
		DeletedIDs:   []string{},
		InvalidIDs:   []string{},
		NotFoundIDs:  []string{},
		DeletedCount: 0,
	}

	// Convert string IDs to ObjectIDs and separate invalid ones
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

	// Perform bulk soft delete
	deletedObjectIDs, err := r.datasource.BulkSoftDelete(ctx, validObjectIDs)
	if err != nil {
		return result, err
	}

	// Convert deleted ObjectIDs back to strings
	for _, objID := range deletedObjectIDs {
		result.DeletedIDs = append(result.DeletedIDs, objID.Hex())
	}
	result.DeletedCount = len(result.DeletedIDs)

	// Find IDs that were requested but not deleted (not found)
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

// UpdateStatus updates the status of an organization
func (r *OrganizationRepositoryMongo) UpdateStatus(ctx context.Context, id primitive.ObjectID, status string) error {
	return r.datasource.UpdateStatus(ctx, id, status)
}



func (r *OrganizationRepositoryMongo) BulkRestore(ctx context.Context, ids []string) (*models.BulkRestoreResponse, error) {
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

// ExistsByID implements repository.RoleRepository.
func (r *OrganizationRepositoryMongo) ExistsByID(ctx context.Context, id primitive.ObjectID) (bool, error) {
	return r.datasource.ExistsByID(ctx, id)

}