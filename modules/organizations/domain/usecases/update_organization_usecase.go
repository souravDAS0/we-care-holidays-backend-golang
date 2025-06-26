package usecases

import (
	"context"
	"errors"
	"time"

	"bitbucket.org/abhishek_fordel/we-care-holidays-backend-golang/modules/organizations/domain/entity"
	"bitbucket.org/abhishek_fordel/we-care-holidays-backend-golang/modules/organizations/domain/repository"
)

// OrganizationUseCase implements the organization business logic
type UpdateOrganizationUseCase struct {
	repo repository.OrganizationRepository
}

func NewUpdateOrganizationUseCase(repo repository.OrganizationRepository) *UpdateOrganizationUseCase {
	return &UpdateOrganizationUseCase{
		repo: repo,
	}
}

// UpdateOrganization updates an existing organization
func (uc *UpdateOrganizationUseCase) Execute(ctx context.Context, org *entity.Organization) error {
	if org.ID.IsZero() {
		return errors.New("organization ID is required")
	}

	// If the slug is being updated, check for duplicates
	if org.Slug != "" {
		existingOrg, err := uc.repo.FindBySlug(ctx, org.Slug)
		if err != nil {
			return err
		}

		// Check if the slug is already used by another organization
		if existingOrg != nil && existingOrg.ID != org.ID {
			return errors.New("organization with this slug already exists")
		}
	}

	// Set updated timestamp
	org.UpdatedAt = time.Now()

	return uc.repo.Update(ctx, org)
}
