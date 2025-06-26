// modules/permissions/domain/repository/permission_repository.go
package repository

import (
	"context"

	"bitbucket.org/abhishek_fordel/we-care-holidays-backend-golang/modules/permissions/domain/entity"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type PermissionRepository interface {
	Create(ctx context.Context, permission *entity.Permission) error
	GetByID(ctx context.Context, id primitive.ObjectID) (*entity.Permission, error)
	Update(ctx context.Context, permission *entity.Permission) error

	HardDelete(ctx context.Context, id primitive.ObjectID) (bool, error)

	List(ctx context.Context, filter map[string]interface{}, page, limit int) ([]*entity.Permission, int64, error)

	ExistsByResourceAction(ctx context.Context, resource string, action entity.PermissionAction) (bool, error)
	ExistsByResourceActionExcluding(ctx context.Context, resource string, action entity.PermissionAction, excludeID primitive.ObjectID) (bool, error)
	ExistsByID(ctx context.Context, id primitive.ObjectID) (bool, error)
}
