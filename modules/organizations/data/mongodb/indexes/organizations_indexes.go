// internal/modules/organizations/data/mongodb/indexes/organizations_indexes.go
package mongodb

import (
	"context"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// SetupOrganizationIndexes creates the indexes for organizations collection
func SetupOrganizationIndexes(coll *mongo.Collection) error {
	models := []mongo.IndexModel{
		// Unique index on slug for active organizations
		{
			Keys: bson.D{{Key: "slug", Value: 1}},
			Options: options.Index().
				SetName("idx_slug_unique").
				SetUnique(true).
				SetPartialFilterExpression(bson.D{{Key: "deletedAt", Value: nil}}),
		},
		// Index on email for lookups
		{
			Keys:    bson.D{{Key: "email", Value: 1}},
			Options: options.Index().SetName("idx_email"),
		},
		// Index on type for filtering by organization type
		{
			Keys:    bson.D{{Key: "type", Value: 1}},
			Options: options.Index().SetName("idx_type"),
		},

		// Index on status for filtering by status
		{
			Keys:    bson.D{{Key: "status", Value: 1}},
			Options: options.Index().SetName("idx_status"),
		},

		// Index on deletedAt for soft delete queries (sparse)
		{
			Keys:    bson.D{{Key: "deletedAt", Value: 1}},
			Options: options.Index().SetName("idx_deletedAt").SetSparse(true),
		},
		// Text index for name and email search
		{
			Keys: bson.D{
				{Key: "name", Value: "text"},
				{Key: "email", Value: "text"},
			},
			Options: options.Index().SetName("idx_name_email_text"),
		},
		// Index on address fields for location-based queries
		{
			Keys: bson.D{
				{Key: "address.country", Value: 1},
				{Key: "address.state", Value: 1},
				{Key: "address.city", Value: 1},
			},
			Options: options.Index().SetName("idx_address_location"),
		},

		// Default sort index for created date
		{
			Keys:    bson.D{{Key: "createdAt", Value: -1}},
			Options: options.Index().SetName("idx_created_desc"),
		},

		// Index on updatedAt for sorting
		{
			Keys:    bson.D{{Key: "updatedAt", Value: -1}},
			Options: options.Index().SetName("idx_updated_desc"),
		},

		// Compound index for common queries (type + status)
		{
			Keys: bson.D{
				{Key: "type", Value: 1},
				{Key: "status", Value: 1},
			},
			Options: options.Index().SetName("idx_type_status"),
		},
		// Compound index for active organizations (status + deletedAt)
		{
			Keys: bson.D{
				{Key: "status", Value: 1},
				{Key: "deletedAt", Value: 1},
			},
			Options: options.Index().SetName("idx_active_organizations"),
		},

		// Index on phone for lookups
		{
			Keys:    bson.D{{Key: "phone", Value: 1}},
			Options: options.Index().SetName("idx_phone"),
		},

		// Index on website for lookups
		{
			Keys:    bson.D{{Key: "website", Value: 1}},
			Options: options.Index().SetName("idx_website"),
		},

		// Compound index for pagination with filters
		{
			Keys: bson.D{
				{Key: "type", Value: 1},
				{Key: "createdAt", Value: -1},
				{Key: "_id", Value: 1},
			},
			Options: options.Index().SetName("idx_type_created_pagination"),
		},

		// Compound index for status-based pagination
		{
			Keys: bson.D{
				{Key: "status", Value: 1},
				{Key: "createdAt", Value: -1},
				{Key: "_id", Value: 1},
			},
			Options: options.Index().SetName("idx_status_created_pagination"),
		},
	}

	_, err := coll.Indexes().CreateMany(context.Background(), models)
	return err
}
