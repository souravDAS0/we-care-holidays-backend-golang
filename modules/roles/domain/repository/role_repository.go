// modules/roles/domain/repository/role_repository.go
package repository

import (
	"context"

	"bitbucket.org/abhishek_fordel/we-care-holidays-backend-golang/internal/models"
	"bitbucket.org/abhishek_fordel/we-care-holidays-backend-golang/modules/roles/domain/entity"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type RoleRepository interface {
	// Basic CRUD operations
	Create(ctx context.Context, role *entity.Role) error
	GetByID(ctx context.Context, id primitive.ObjectID) (*entity.Role, error)
	GetByName(ctx context.Context, name string) (*entity.Role, error)
	Update(ctx context.Context, role *entity.Role) error
	
	// Listing and filtering
	List(ctx context.Context, filter map[string]interface{}, page, limit int) ([]*entity.Role, int64, error)
	
	// Bulk operations
	BulkSoftDelete(ctx context.Context, ids []string) (*models.BulkDeleteResponse, error) 
	BulkRestore(ctx context.Context, ids []string) (*models.BulkRestoreResponse, error) 
	
	// Soft delete operations
	SoftDelete(ctx context.Context, id primitive.ObjectID) (bool, error)
	Restore(ctx context.Context, id primitive.ObjectID) (bool, error)
	HardDelete(ctx context.Context, id primitive.ObjectID) (bool, error)
	
	// Permission operations
	// AddPermission(ctx context.Context, roleID, permissionID primitive.ObjectID) error
	// RemovePermission(ctx context.Context, roleID, permissionID primitive.ObjectID) error
	// SetPermissions(ctx context.Context, roleID primitive.ObjectID, permissionIDs []primitive.ObjectID) error
	// GetRolesByPermission(ctx context.Context, permissionID primitive.ObjectID) ([]*entity.Role, error)
	
	// Default role operations
	// GetDefaultRoles(ctx context.Context) ([]*entity.Role, error)
	// SetAsDefault(ctx context.Context, roleID primitive.ObjectID) error
	// UnsetAsDefault(ctx context.Context, roleID primitive.ObjectID) error
	
	// Assignment checks
	// IsRoleAssignedToUsers(ctx context.Context, roleID primitive.ObjectID) (bool, error)
	// GetAssignedUserCount(ctx context.Context, roleID primitive.ObjectID) (int64, error)
	
	// Existence checks
	ExistsByName(ctx context.Context, name string) (bool, error)
	ExistsByID(ctx context.Context, id primitive.ObjectID) (bool, error)
}
