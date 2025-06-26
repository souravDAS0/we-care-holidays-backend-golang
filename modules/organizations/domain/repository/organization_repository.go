package repository

import (
	"context"

	"bitbucket.org/abhishek_fordel/we-care-holidays-backend-golang/internal/models"
	"bitbucket.org/abhishek_fordel/we-care-holidays-backend-golang/modules/organizations/domain/entity"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// OrganizationRepository defines the interface for organization persistence
type OrganizationRepository interface {
	// FindAll retrieves organizations with filtering and pagination
	// The filter should exclude soft-deleted organizations by default
	FindAll(ctx context.Context, filter map[string]interface{}, page, limit int) ([]*entity.Organization, int64, error)
	
	// FindByID finds an organization by its ID
	// This should not return soft-deleted organizations unless explicitly asked
	FindByID(ctx context.Context, id primitive.ObjectID) (*entity.Organization, error)
	
	// FindBySlug finds an organization by its slug
	// This should not return soft-deleted organizations
	FindBySlug(ctx context.Context, slug string) (*entity.Organization, error)
	
	// Create inserts a new organization
	Create(ctx context.Context, organization *entity.Organization) error
	
	// Update updates an existing organization
	Update(ctx context.Context, organization *entity.Organization) error
	
	// SoftDelete marks an organization as deleted without removing it from the database
	SoftDelete(ctx context.Context, id primitive.ObjectID) (bool, error)
	
	// Restore restores a soft-deleted organization
	Restore(ctx context.Context, id primitive.ObjectID) (bool, error)
	
	// HardDelete permanently removes an organization (primarily for admin/cleanup)
	HardDelete(ctx context.Context, id primitive.ObjectID) (bool, error)
	
	// BulkSoftDelete marks multiple organizations as deleted
	BulkSoftDelete(ctx context.Context, ids []string) (*models.BulkDeleteResponse, error)
	
	// UpdateStatus updates the status of an organization
	UpdateStatus(ctx context.Context, id primitive.ObjectID, status string) error

	BulkRestore(ctx context.Context, ids []string) (*models.BulkRestoreResponse, error) 
	
	ExistsByID(ctx context.Context, id primitive.ObjectID) (bool, error)
	
}