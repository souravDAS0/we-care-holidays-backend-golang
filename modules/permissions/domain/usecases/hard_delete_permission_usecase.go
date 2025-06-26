package usecases

import (
	"context"

	"bitbucket.org/abhishek_fordel/we-care-holidays-backend-golang/modules/permissions/domain/repository"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// PermissionUseCase implements the Permission business logic
type HardDeletePermissionUseCase struct {
	repo repository.PermissionRepository
}

func NewHardDeletePermissionUseCase(repo repository.PermissionRepository) *HardDeletePermissionUseCase {
	return &HardDeletePermissionUseCase{
		repo: repo,
	}
}

// HardDeletePermission marks an Permission as deleted without removing it
func (uc *HardDeletePermissionUseCase) Execute(ctx context.Context, id primitive.ObjectID) (bool, error) {
	return uc.repo.HardDelete(ctx, id)
}
