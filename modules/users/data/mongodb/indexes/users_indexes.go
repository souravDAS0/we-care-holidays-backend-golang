package mongodb

import (
	"context"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// SetupUserIndexes creates the indexes for users collection
func SetupUserIndexes(coll *mongo.Collection) error {
	models := []mongo.IndexModel{
		// Unique index on primary email for active users
		{
			Keys: bson.D{{Key: "emails.0", Value: 1}}, // First email is primary
			Options: options.Index().
				SetName("idx_primary_email_unique").
				SetUnique(true).
				SetPartialFilterExpression(bson.D{{Key: "deletedAt", Value: nil}}),
		},

		// Index on all emails for lookups
		{
			Keys:    bson.D{{Key: "emails", Value: 1}},
			Options: options.Index().SetName("idx_emails"),
		},

		// Index on all phones for lookups
		{
			Keys:    bson.D{{Key: "phones", Value: 1}},
			Options: options.Index().SetName("idx_phones"),
		},

		// Index on roleId for role-based queries
		{
			Keys:    bson.D{{Key: "roleId", Value: 1}},
			Options: options.Index().SetName("idx_roleId"),
		},

		// Index on organizationId for organization-based queries
		{
			Keys:    bson.D{{Key: "organizationId", Value: 1}},
			Options: options.Index().SetName("idx_organizationId"),
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

		// Text index for fullName and email search
		{
			Keys: bson.D{
				{Key: "fullName", Value: "text"},
				{Key: "emails", Value: "text"},
			},
			Options: options.Index().
				SetName("idx_fullName_emails_text").
				SetWeights(bson.D{
					{Key: "fullName", Value: 10},
					{Key: "emails", Value: 8},
				}),
		},

		// Index on createdAt for sorting
		{
			Keys:    bson.D{{Key: "createdAt", Value: -1}},
			Options: options.Index().SetName("idx_created_desc"),
		},

		// Index on updatedAt for sorting
		{
			Keys:    bson.D{{Key: "updatedAt", Value: -1}},
			Options: options.Index().SetName("idx_updated_desc"),
		},

		// Compound index for role and organization queries
		{
			Keys: bson.D{
				{Key: "roleId", Value: 1},
				{Key: "organizationId", Value: 1},
			},
			Options: options.Index().SetName("idx_role_organization"),
		},

		// Compound index for status and organization queries
		{
			Keys: bson.D{
				{Key: "status", Value: 1},
				{Key: "organizationId", Value: 1},
			},
			Options: options.Index().SetName("idx_status_organization"),
		},

		// Compound index for active users (status + deletedAt)
		{
			Keys: bson.D{
				{Key: "status", Value: 1},
				{Key: "deletedAt", Value: 1},
			},
			Options: options.Index().SetName("idx_active_users"),
		},

		// Index on last login for audit trail queries
		{
			Keys:    bson.D{{Key: "auditTrail.lastLoginAt", Value: -1}},
			Options: options.Index().SetName("idx_last_login").SetSparse(true),
		},

		// Compound index for pagination with role filter
		{
			Keys: bson.D{
				{Key: "roleId", Value: 1},
				{Key: "createdAt", Value: -1},
				{Key: "_id", Value: 1},
			},
			Options: options.Index().SetName("idx_role_created_pagination"),
		},

		// Compound index for pagination with organization filter
		{
			Keys: bson.D{
				{Key: "organizationId", Value: 1},
				{Key: "createdAt", Value: -1},
				{Key: "_id", Value: 1},
			},
			Options: options.Index().SetName("idx_org_created_pagination"),
		},

		// Compound index for pagination with status filter
		{
			Keys: bson.D{
				{Key: "status", Value: 1},
				{Key: "createdAt", Value: -1},
				{Key: "_id", Value: 1},
			},
			Options: options.Index().SetName("idx_status_created_pagination"),
		},

		// Index for login method queries
		{
			Keys:    bson.D{{Key: "loginMethods.passwordLogin.enabled", Value: 1}},
			Options: options.Index().SetName("idx_password_login_enabled"),
		},

		{
			Keys:    bson.D{{Key: "loginMethods.emailOtpLogin.enabled", Value: 1}},
			Options: options.Index().SetName("idx_email_otp_enabled"),
		},

		{
			Keys:    bson.D{{Key: "loginMethods.phoneOtpLogin.enabled", Value: 1}},
			Options: options.Index().SetName("idx_phone_otp_enabled"),
		},
	}

	_, err := coll.Indexes().CreateMany(context.Background(), models)
	return err
}
