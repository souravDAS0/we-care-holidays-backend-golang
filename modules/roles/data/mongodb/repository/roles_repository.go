package repository

import (
	"context"

	"bitbucket.org/abhishek_fordel/we-care-holidays-backend-golang/internal/models"
	"bitbucket.org/abhishek_fordel/we-care-holidays-backend-golang/modules/roles/data/datasource"
	"bitbucket.org/abhishek_fordel/we-care-holidays-backend-golang/modules/roles/data/mongodb/model"
	"bitbucket.org/abhishek_fordel/we-care-holidays-backend-golang/modules/roles/domain/entity"
	"bitbucket.org/abhishek_fordel/we-care-holidays-backend-golang/modules/roles/domain/repository"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

var _ repository.RoleRepository = (*RoleRepositoryMongo)(nil)

type RoleRepositoryMongo struct {
	datasource *datasource.MongoRoleDatasource
}



func NewRoleRepositoryMongo(ds *datasource.MongoRoleDatasource) *RoleRepositoryMongo {
	return &RoleRepositoryMongo{
		datasource: ds,
	}
}

// List implements repository.RoleRepository.
func (r *RoleRepositoryMongo) List(ctx context.Context, filter map[string]interface{}, page int, limit int) ([]*entity.Role, int64, error) {
	roleModels, totalCount, err := r.datasource.FindByFilters(ctx, filter, page, limit)
	if err != nil {
		return nil, 0, err
	}

	roleEntities := make([]*entity.Role, len(roleModels))
	for i, model := range roleModels {
		entity := model.ToEntity()
		roleEntities[i] = &entity
	}
	return roleEntities, totalCount, nil
}

// Create implements repository.RoleRepository.
func (r *RoleRepositoryMongo) Create(ctx context.Context, role *entity.Role) error {
	roleModel  := model.FromEntity(role)

	err := r.datasource.Insert(ctx, roleModel)
	if err != nil {
		return err
	}

	role.ID = roleModel.ID
	role.CreatedAt = roleModel.CreatedAt
	role.UpdatedAt = roleModel.UpdatedAt
	return nil
}

// GetByID implements repository.RoleRepository.
func (r *RoleRepositoryMongo) GetByID(ctx context.Context, id primitive.ObjectID) (*entity.Role, error) {
	roleModel , err := r.datasource.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if roleModel == nil {
		return nil, nil // No document found
	}
	roleEntity := roleModel.ToEntity()
	return &roleEntity, nil
}

// GetByName implements repository.RoleRepository.
func (r *RoleRepositoryMongo) GetByName(ctx context.Context, name string) (*entity.Role, error) {
	roleModel, err := r.datasource.FindByName(ctx, name)
	if err != nil {
		return nil, err
	}
	if roleModel == nil {
		return nil, nil // No document found
	}
	roleEntity := roleModel.ToEntity()
	return &roleEntity, nil

}



// Update implements repository.RoleRepository.
func (r *RoleRepositoryMongo) Update(ctx context.Context, role *entity.Role) error {
	roleModel := model.FromEntity(role)
	return r.datasource.Update(ctx, roleModel)
}


func (r *RoleRepositoryMongo) HardDelete(ctx context.Context, id primitive.ObjectID) (bool, error) {
	return r.datasource.HardDelete(ctx, id)
}

func (r *RoleRepositoryMongo) Restore(ctx context.Context, id primitive.ObjectID) (bool, error) {
	return r.datasource.Restore(ctx, id)
}

func (r *RoleRepositoryMongo) SoftDelete(ctx context.Context, id primitive.ObjectID) (bool, error) {
	return r.datasource.SoftDelete(ctx, id)
}


func (r *RoleRepositoryMongo) BulkSoftDelete(ctx context.Context, ids []string) (*models.BulkDeleteResponse, error) {
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

func (r *RoleRepositoryMongo) BulkRestore(ctx context.Context, ids []string) (*models.BulkRestoreResponse, error) {
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
func (r *RoleRepositoryMongo) ExistsByID(ctx context.Context, id primitive.ObjectID) (bool, error) {
	return r.datasource.ExistsByID(ctx, id)

}

// ExistsByName implements repository.RoleRepository.
func (r *RoleRepositoryMongo) ExistsByName(ctx context.Context, name string) (bool, error) {
	return r.datasource.ExistsByName(ctx, name)

}
