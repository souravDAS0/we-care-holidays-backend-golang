package usecases

import (
	"context"
	"errors"
	"time"

	"bitbucket.org/abhishek_fordel/we-care-holidays-backend-golang/internal/utils"
	"bitbucket.org/abhishek_fordel/we-care-holidays-backend-golang/modules/organizations/domain/entity"
	"bitbucket.org/abhishek_fordel/we-care-holidays-backend-golang/modules/organizations/domain/repository"
)

// OrganizationUseCase implements the organization business logic
type CreateOrganizationUseCase struct {
	repo repository.OrganizationRepository
}

func NewCreateOrganizationUseCase(repo repository.OrganizationRepository) *CreateOrganizationUseCase {
	return &CreateOrganizationUseCase{
		repo: repo,
	}
}

// CreateOrganization creates a new organization
func (uc *CreateOrganizationUseCase) Execute(ctx context.Context, org *entity.Organization) error {
	// Handle slug generation if not provided
	if org.Slug == "" {
		// Generate slug from name
		org.Slug = utils.GenerateSlug(org.Name)
	}

	// Check if slug exists (to prevent duplicates)
	existingOrg, err := uc.repo.FindBySlug(ctx, org.Slug)
	if err != nil {
		return err
	}
	if existingOrg != nil {
		return errors.New("organization with this slug already exists")
	}

	// Set default values
	org.CreatedAt = time.Now()
	org.UpdatedAt = time.Now()
	if org.Status == "" {
		org.Status = "Pending" // Default status
	}

	// Ensure DeletedAt is nil for new organizations
	org.DeletedAt = nil

	return uc.repo.Create(ctx, org)
}
