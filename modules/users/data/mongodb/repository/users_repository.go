package repository

import (
	"context"
	"fmt"

	"bitbucket.org/abhishek_fordel/we-care-holidays-backend-golang/internal/models"
	"bitbucket.org/abhishek_fordel/we-care-holidays-backend-golang/modules/users/data/datasource"
	"bitbucket.org/abhishek_fordel/we-care-holidays-backend-golang/modules/users/data/mongodb/model"
	"bitbucket.org/abhishek_fordel/we-care-holidays-backend-golang/modules/users/domain/entity"
	"bitbucket.org/abhishek_fordel/we-care-holidays-backend-golang/modules/users/domain/repository"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Ensure interface compliance
var _ repository.UserRepository = (*UserRepositoryMongo)(nil)

type UserRepositoryMongo struct {
	datasource *datasource.MongoUserDatasource
}

func NewUserRepositoryMongo(ds *datasource.MongoUserDatasource) *UserRepositoryMongo {
	return &UserRepositoryMongo{
		datasource: ds,
	}
}

// List implements repository.UserRepository.
func (u *UserRepositoryMongo) List(ctx context.Context, filter map[string]interface{}, page int, limit int) ([]*entity.User, int64, error) {
	userModels, totalCount, err := u.datasource.FindByFilters(ctx, filter, page, limit)
	if err != nil {
		return nil, 0, err
	}
	userEntities := make([]*entity.User, len(userModels))
	for i, model := range userModels {
		entity := model.ToEntity()
		userEntities[i] = &entity
	}
	return userEntities, totalCount, nil
}

// Create implements repository.UserRepository.
func (u *UserRepositoryMongo) Create(ctx context.Context, user *entity.User) error {

	userModel := model.FromEntity(user)

	err := u.datasource.Insert(ctx, userModel)
	if err != nil {
		return err
	}

	// Set back the generated ID in entity for downstream use
	user.ID = userModel.ID
	user.CreatedAt = userModel.CreatedAt
	user.UpdatedAt = userModel.UpdatedAt
	return nil
}

// GetByID implements repository.UserRepository.
func (u *UserRepositoryMongo) GetByID(ctx context.Context, id primitive.ObjectID) (*entity.User, error) {
	// Add defensive nil check
	if u == nil {
		return nil, fmt.Errorf("user repository is nil")
	}
	if u.datasource == nil {
		return nil, fmt.Errorf("user datasource is nil - check container initialization")
	}

	userModel, err := u.datasource.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}

	if userModel == nil {
		return nil, nil
	}

	userEntity := userModel.ToEntity()
	return &userEntity, nil
}

// GetByEmail implements repository.UserRepository.
func (u *UserRepositoryMongo) GetByEmail(ctx context.Context, email string) (*entity.User, error) {
	userModel, err := u.datasource.FindByEmail(ctx, email)
	if err != nil {
		return nil, err
	}

	if userModel == nil {
		return nil, nil
	}

	userEntity := userModel.ToEntity()
	return &userEntity, nil
}

// GetByPhone implements repository.UserRepository.
func (u *UserRepositoryMongo) GetByPhone(ctx context.Context, phone string) (*entity.User, error) {
	userModel, err := u.datasource.FindByPhone(ctx, phone)
	if err != nil {
		return nil, err
	}

	if userModel == nil {
		return nil, nil
	}

	userEntity := userModel.ToEntity()
	return &userEntity, nil
}

// HardDelete implements repository.UserRepository.
func (u *UserRepositoryMongo) HardDelete(ctx context.Context, id primitive.ObjectID) (bool, error) {
	return u.datasource.HardDelete(ctx, id)
}

// Restore implements repository.UserRepository.
func (u *UserRepositoryMongo) Restore(ctx context.Context, id primitive.ObjectID) (bool, error) {
	return u.datasource.Restore(ctx, id)
}

// SoftDelete implements repository.UserRepository.
func (u *UserRepositoryMongo) SoftDelete(ctx context.Context, id primitive.ObjectID) (bool, error) {
	return u.datasource.SoftDelete(ctx, id)
}

// Update implements repository.UserRepository.
func (u *UserRepositoryMongo) Update(ctx context.Context, user *entity.User) error {
	userModel := model.FromEntity(user)
	return u.datasource.Update(ctx, userModel)
}

// BulkSoftDelete implements repository.UserRepository.
func (u *UserRepositoryMongo) BulkSoftDelete(ctx context.Context, ids []string) (*models.BulkDeleteResponse, error) {
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
	deletedObjectIDs, err := u.datasource.BulkSoftDelete(ctx, validObjectIDs)
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

func (u *UserRepositoryMongo) BulkRestore(ctx context.Context, ids []string) (*models.BulkRestoreResponse, error) {
	result := &models.BulkRestoreResponse{
		RequestedIDs:  ids,
		RestoredIDs:   []string{},
		InvalidIDs:    []string{},
		NotFoundIDs:   []string{},
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

	restoredObjectIDs, err := u.datasource.BulkRestore(ctx, validObjectIDs)

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

// ExistsByEmail implements repository.UserRepository.
func (u *UserRepositoryMongo) ExistsByEmail(ctx context.Context, email string) (bool, error) {
	return u.datasource.ExistsByEmail(ctx, email)
}

// ExistsByID implements repository.UserRepository.
func (u *UserRepositoryMongo) ExistsByID(ctx context.Context, id primitive.ObjectID) (bool, error) {
	return u.datasource.ExistsByID(ctx, id)
}

// ExistsByPhone implements repository.UserRepository.
func (u *UserRepositoryMongo) ExistsByPhone(ctx context.Context, phone string) (bool, error) {
	return u.datasource.ExistsByPhone(ctx, phone)
}
