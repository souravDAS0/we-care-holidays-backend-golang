package mongodb

import (
	"context"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func SetupPermissionIndexes(coll *mongo.Collection) error {
	models := []mongo.IndexModel{
		{
			Keys: bson.D{
				{Key: "resource", Value: 1},
				{Key: "action", Value: 1},
				{Key: "scope", Value: 1},
			},
			Options: options.Index().
				SetUnique(true).
				SetName("unique_resource_action_scope").
				SetPartialFilterExpression(bson.D{
					{Key: "deletedAt", Value: nil},
				}),
		},

		// Index on resource for resource-based queries
		{
			Keys: bson.D{
				{Key: "resource", Value: 1},
			},
			Options: options.Index().
				SetName("resource_lookup"),
		},

		// Index on action for action-based queries
		{
			Keys: bson.D{
				{Key: "action", Value: 1},
			},
			Options: options.Index().
				SetName("action_lookup"),
		},

		// Index on scope for scope-based queries
		{
			Keys: bson.D{
				{Key: "scope", Value: 1},
			},
			Options: options.Index().
				SetName("scope_lookup"),
		},

		// Index on enabled for filtering active/inactive permissions
		{
			Keys: bson.D{
				{Key: "enabled", Value: 1},
			},
			Options: options.Index().
				SetName("enabled_filter"),
		},

		// Index on priority for sorting and priority-based queries
		{
			Keys: bson.D{
				{Key: "priority", Value: -1}, // Descending for higher priority first
			},
			Options: options.Index().
				SetName("priority_sort"),
		},

		// Compound index for resource-action queries (common use case)
		{
			Keys: bson.D{
				{Key: "resource", Value: 1},
				{Key: "action", Value: 1},
			},
			Options: options.Index().
				SetName("resource_action_lookup"),
		},

		// Compound index for resource-scope queries
		{
			Keys: bson.D{
				{Key: "resource", Value: 1},
				{Key: "scope", Value: 1},
			},
			Options: options.Index().
				SetName("resource_scope_lookup"),
		},

		// Compound index for action-scope queries
		{
			Keys: bson.D{
				{Key: "action", Value: 1},
				{Key: "scope", Value: 1},
			},
			Options: options.Index().
				SetName("action_scope_lookup"),
		},

		// Text index for searching resource and notes
		{
			Keys: bson.D{
				{Key: "resource", Value: "text"},
				{Key: "notes", Value: "text"},
			},
			Options: options.Index().
				SetName("text_search").
				SetWeights(bson.D{
					{Key: "resource", Value: 10},
					{Key: "notes", Value: 5},
				}),
		},

		// Index on createdAt for sorting and date filtering
		{
			Keys: bson.D{
				{Key: "createdAt", Value: -1},
			},
			Options: options.Index().
				SetName("createdAt_desc"),
		},

		// Index on updatedAt for sorting and date filtering
		{
			Keys: bson.D{
				{Key: "updatedAt", Value: -1},
			},
			Options: options.Index().
				SetName("updatedAt_desc"),
		},

		// Index on deletedAt for soft delete filtering
		{
			Keys: bson.D{
				{Key: "deletedAt", Value: 1},
			},
			Options: options.Index().
				SetName("deletedAt_filter").
				SetSparse(true),
		},

		// Compound index for active permissions (enabled and not deleted)
		{
			Keys: bson.D{
				{Key: "enabled", Value: 1},
				{Key: "deletedAt", Value: 1},
			},
			Options: options.Index().
				SetName("active_permissions_filter"),
		},

		// Compound index for permission evaluation (resource, action, scope, enabled, priority)
		{
			Keys: bson.D{
				{Key: "resource", Value: 1},
				{Key: "action", Value: 1},
				{Key: "scope", Value: 1},
				{Key: "enabled", Value: 1},
				{Key: "priority", Value: -1},
			},
			Options: options.Index().
				SetName("permission_evaluation"),
		},



		// Compound index for efficient pagination with resource filter
		{
			Keys: bson.D{
				{Key: "resource", Value: 1},
				{Key: "createdAt", Value: -1},
				{Key: "_id", Value: 1},
			},
			Options: options.Index().
				SetName("resource_created_pagination"),
		},

		// Compound index for efficient pagination with action filter
		{
			Keys: bson.D{
				{Key: "action", Value: 1},
				{Key: "createdAt", Value: -1},
				{Key: "_id", Value: 1},
			},
			Options: options.Index().
				SetName("action_created_pagination"),
		},

		// Compound index for efficient pagination with scope filter
		{
			Keys: bson.D{
				{Key: "scope", Value: 1},
				{Key: "createdAt", Value: -1},
				{Key: "_id", Value: 1},
			},
			Options: options.Index().
				SetName("scope_created_pagination"),
		},

		// Compound index for priority range queries with enabled filter
		{
			Keys: bson.D{
				{Key: "enabled", Value: 1},
				{Key: "priority", Value: -1},
			},
			Options: options.Index().
				SetName("enabled_priority_filter"),
		},



		// Compound index for complex permission queries with all filters
		{
			Keys: bson.D{
				{Key: "enabled", Value: 1},
				{Key: "resource", Value: 1},
				{Key: "action", Value: 1},
				{Key: "scope", Value: 1},
				{Key: "priority", Value: -1},
				{Key: "createdAt", Value: -1},
			},
			Options: options.Index().
				SetName("complex_permission_query"),
		},
	}

	if _, err := coll.Indexes().CreateMany(context.Background(), models); err != nil {
		return err
	}

	return nil
}