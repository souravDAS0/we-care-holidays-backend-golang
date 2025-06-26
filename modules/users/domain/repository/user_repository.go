// modules/users/domain/repository/user_repository.go
package repository

import (
	"context"

	"bitbucket.org/abhishek_fordel/we-care-holidays-backend-golang/internal/models"
	"bitbucket.org/abhishek_fordel/we-care-holidays-backend-golang/modules/users/domain/entity"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type UserRepository interface {
	// Basic CRUD operations
	Create(ctx context.Context, user *entity.User) error
	GetByID(ctx context.Context, id primitive.ObjectID) (*entity.User, error)
	GetByEmail(ctx context.Context, email string) (*entity.User, error)
	GetByPhone(ctx context.Context, phone string) (*entity.User, error)
	Update(ctx context.Context, user *entity.User) error
	
	// Listing and filtering
	List(ctx context.Context, filter map[string]interface{}, page, limit int) ([]*entity.User, int64, error)
	
	// Bulk operations
	BulkSoftDelete(ctx context.Context, ids []string) (*models.BulkDeleteResponse, error)
	BulkRestore(ctx context.Context, ids []string) (*models.BulkRestoreResponse, error) 
	
	// Soft delete operations
	SoftDelete(ctx context.Context, id primitive.ObjectID) (bool, error)
	Restore(ctx context.Context, id primitive.ObjectID) (bool, error)
	HardDelete(ctx context.Context, id primitive.ObjectID) (bool, error)

	// Existence checks
	ExistsByEmail(ctx context.Context, email string) (bool, error)
	ExistsByPhone(ctx context.Context, phone string) (bool, error)
	ExistsByID(ctx context.Context, id primitive.ObjectID) (bool, error)
}
