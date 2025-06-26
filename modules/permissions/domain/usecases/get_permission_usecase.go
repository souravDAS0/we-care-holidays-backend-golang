package usecases

import (
	"context"

	"bitbucket.org/abhishek_fordel/we-care-holidays-backend-golang/modules/permissions/domain/entity"
	"bitbucket.org/abhishek_fordel/we-care-holidays-backend-golang/modules/permissions/domain/repository"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// PermissionUseCase implements the Permission business logic
type GetPermissionUseCase struct {
	repo repository.PermissionRepository
}

func NewGetPermissionUseCase(repo repository.PermissionRepository) *GetPermissionUseCase {
	return &GetPermissionUseCase{
		repo: repo,
	}
}

// GetPermission retrieves an Permission by its ID
func (uc *GetPermissionUseCase) Execute(ctx context.Context, id primitive.ObjectID) (*entity.Permission, error) {
	return uc.repo.GetByID(ctx, id)
}
