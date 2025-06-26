package repository

import (
	"context"

	"bitbucket.org/abhishek_fordel/we-care-holidays-backend-golang/modules/permissions/data/datasource"
	"bitbucket.org/abhishek_fordel/we-care-holidays-backend-golang/modules/permissions/data/mongodb/model"
	"bitbucket.org/abhishek_fordel/we-care-holidays-backend-golang/modules/permissions/domain/entity"
	"bitbucket.org/abhishek_fordel/we-care-holidays-backend-golang/modules/permissions/domain/repository"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

var _ repository.PermissionRepository = (*PermissionRepositoryMongo)(nil)

type PermissionRepositoryMongo struct {
	datasource *datasource.MongoPermissionDatasource
}

func NewPermissionRepositoryMongo(ds *datasource.MongoPermissionDatasource) *PermissionRepositoryMongo {
	return &PermissionRepositoryMongo{
		datasource: ds,
	}
}

func (r *PermissionRepositoryMongo) List(ctx context.Context, filter map[string]interface{}, page, limit int) ([]*entity.Permission, int64, error) {
	permissionModels, totalCount, err := r.datasource.FindByFilters(ctx, filter, page, limit)
	if err != nil {
		return nil, 0, err
	}

	permissionEntities := make([]*entity.Permission, len(permissionModels))
	for i, model := range permissionModels {
		entity := model.ToEntity()
		permissionEntities[i] = &entity
	}
	return permissionEntities, totalCount, nil
}

func (r *PermissionRepositoryMongo) Create(ctx context.Context, permission *entity.Permission) error {
	permissionModel := model.FromEntity(permission)

	err := r.datasource.Insert(ctx, permissionModel)
	if err != nil {
		return err
	}

	permission.ID = permissionModel.ID
	permission.CreatedAt = permissionModel.CreatedAt
	permission.UpdatedAt = permissionModel.UpdatedAt
	return nil
}

func (r *PermissionRepositoryMongo) GetByID(ctx context.Context, id primitive.ObjectID) (*entity.Permission, error) {
	permissionModel, err := r.datasource.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}

	if permissionModel == nil {
		return nil, nil
	}

	permissionEntity := permissionModel.ToEntity()
	return &permissionEntity, nil
}

func (r *PermissionRepositoryMongo) Update(ctx context.Context, permission *entity.Permission) error {
	permissionModel := model.FromEntity(permission)

	return r.datasource.Update(ctx, permissionModel)

}

func (r *PermissionRepositoryMongo) HardDelete(ctx context.Context, id primitive.ObjectID) (bool, error) {
	return r.datasource.HardDelete(ctx, id)
}

func (r *PermissionRepositoryMongo) ExistsByResourceAction(ctx context.Context, resource string, action entity.PermissionAction) (bool, error) {
	return r.datasource.ExistsByResourceAction(ctx, resource, action)
}

func (r *PermissionRepositoryMongo) ExistsByResourceActionExcluding(ctx context.Context, resource string, action entity.PermissionAction, excludeID primitive.ObjectID) (bool, error) {
	return r.datasource.ExistsByResourceActionExcluding(ctx, resource, action, excludeID)
}

func (r *PermissionRepositoryMongo) ExistsByID(ctx context.Context, id primitive.ObjectID) (bool, error) {
	return r.datasource.ExistsByID(ctx, id)
}
