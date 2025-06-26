package mongodb

import (
	"context"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// SetupLocationIndexes creates indexes for locations
func SetupLocationIndexes(coll *mongo.Collection) error {
	models := []mongo.IndexModel{
		{
			Keys: bson.D{
				{Key: "name", Value: 1},
				{Key: "type", Value: 1},
			},
			Options: options.Index().
				SetName("idx_name_type_unique").
				SetUnique(true).
				SetPartialFilterExpression(bson.D{{Key: "deletedAt", Value: nil}}),
		},
		// Index on type for filtering by location type
		{
			Keys:    bson.D{{Key: "type", Value: 1}},
			Options: options.Index().SetName("idx_type"),
		},
		
		{ // Index on country (non-unique)
			Keys: bson.D{{Key: "country", Value: 1}},
			Options: options.Index().SetName("idx_country"),
		},
		{ // Index on state (non-unique)
			Keys: bson.D{{Key: "state", Value: 1}},
			Options: options.Index().SetName("idx_state"),
		},
		{ // Index on tags (non-unique)
			Keys: bson.D{{Key: "tags", Value: 1}},
			Options: options.Index().SetName("idx_tags"),
		},
		{ // Index on aliases (non-unique)
			Keys: bson.D{{Key: "aliases", Value: 1}},
			Options: options.Index().SetName("idx_aliases"),
		},
		{ // Index on createdAt (for sorting by creation time, non-unique)
			Keys: bson.D{{Key: "createdAt", Value: -1}},
			Options: options.Index().SetName("idx_created_desc"),
		},
		{ // Index on deletedAt (non-unique)
			Keys: bson.D{{Key: "deletedAt", Value: 1}},
			Options: options.Index().SetName("idx_deletedAt"),
		},
	}
	_, err := coll.Indexes().CreateMany(context.Background(), models)
	return err
}
