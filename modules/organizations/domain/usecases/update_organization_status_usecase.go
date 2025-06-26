package usecases

import (
	"context"
	"errors"

	"bitbucket.org/abhishek_fordel/we-care-holidays-backend-golang/modules/organizations/domain/repository"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// OrganizationStatusUseCase implements the organization business logic
type UpdateOrganizationStatusUseCase struct {
	repo repository.OrganizationRepository
}

func NewUpdateOrganizationStatusUseCase(repo repository.OrganizationRepository) *UpdateOrganizationStatusUseCase {
	return &UpdateOrganizationStatusUseCase{
		repo: repo,
	}
}

// UpdateOrganizationStatus updates the status of an organization
func (uc *UpdateOrganizationStatusUseCase) Execute(ctx context.Context, id primitive.ObjectID, status string) error {
	// Validate status (should be one of: "Pending", "Approved", "Suspended", "Archived")
	validStatuses := map[string]bool{
		"Pending":   true,
		"Approved":  true,
		"Suspended": true,
		"Archived":  true,
	}

	if !validStatuses[status] {
		return errors.New("invalid organization status")
	}

	return uc.repo.UpdateStatus(ctx, id, status)
}
