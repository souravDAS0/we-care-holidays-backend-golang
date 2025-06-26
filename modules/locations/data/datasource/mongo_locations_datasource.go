package datasource

import (
	"context"
	"log"
	"time"

	indexes "bitbucket.org/abhishek_fordel/we-care-holidays-backend-golang/modules/locations/data/mongodb/indexes"
	"bitbucket.org/abhishek_fordel/we-care-holidays-backend-golang/modules/locations/data/mongodb/model"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// MongoLocationDatasource handles raw MongoDB operations for locations.
type MongoLocationDatasource struct {
	collection *mongo.Collection
}

// NewMongoLocationDatasource creates a new datasource, sets up indexes.
func NewMongoLocationDatasource(db *mongo.Database) *MongoLocationDatasource {
	coll := db.Collection((&model.LocationModel{}).CollectionName())
	if err := indexes.SetupLocationIndexes(coll); err != nil {
		log.Printf("Failed to setup location indexes: %v", err)
	}
	return &MongoLocationDatasource{collection: coll}
}

// Insert inserts a new location model, setting CreatedAt/UpdatedAt.
func (ds *MongoLocationDatasource) Insert(ctx context.Context, lm *model.LocationModel) error {
	now := time.Now()
	lm.CreatedAt, lm.UpdatedAt = now, now

	res, err := ds.collection.InsertOne(ctx, lm)
	if err != nil {
		return err
	}
	lm.ID = res.InsertedID.(primitive.ObjectID)
	return nil
}

// FindByFilters retrieves location documents matching filters (excluding soft-deleted),
// paginated and sorted by createdAt desc.
func (ds *MongoLocationDatasource) FindByFilters(ctx context.Context, filters map[string]interface{}, page, limit int) ([]model.LocationModel, int64, error) {
	// always exclude soft-deleted
	filters["deletedAt"] = nil

	total, err := ds.collection.CountDocuments(ctx, filters)
	if err != nil {
		return nil, 0, err
	}

	opts := options.Find().
		SetSkip(int64((page - 1) * limit)).
		SetLimit(int64(limit)).
		SetSort(bson.D{{Key: "createdAt", Value: -1}})

	cur, err := ds.collection.Find(ctx, filters, opts)
	if err != nil {
		return nil, 0, err
	}
	defer cur.Close(ctx)

	var out []model.LocationModel
	if err := cur.All(ctx, &out); err != nil {
		return nil, 0, err
	}
	return out, total, nil
}

// FindByID finds a location by its ObjectID.
// Returns (nil, nil) if not found.
func (ds *MongoLocationDatasource) FindByID(ctx context.Context, id primitive.ObjectID) (*model.LocationModel, error) {
	var lm model.LocationModel
	err := ds.collection.FindOne(ctx, bson.M{"_id": id}).Decode(&lm)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, nil
		}
		return nil, err
	}
	return &lm, nil
}

// Update replaces an existing location document (sets UpdatedAt).
func (ds *MongoLocationDatasource) Update(ctx context.Context, lm *model.LocationModel) error {
	lm.UpdatedAt = time.Now()
	_, err := ds.collection.UpdateOne(
		ctx,
		bson.M{"_id": lm.ID},
		bson.M{"$set": lm},
	)
	return err
}

// SoftDelete marks a document as deleted by setting deletedAt and updatedAt.
func (ds *MongoLocationDatasource) SoftDelete(ctx context.Context, id primitive.ObjectID) (bool, error) {
	res, err := ds.collection.UpdateOne(
		ctx,
		bson.M{"_id": id, "deletedAt": nil},
		bson.M{"$set": bson.M{"deletedAt": time.Now(), "updatedAt": time.Now()}},
	)
	if err != nil {
		return false, err
	}
	return res.ModifiedCount > 0, nil
}

// Restore clears the deletedAt timestamp, restoring a soft-deleted document.
func (ds *MongoLocationDatasource) Restore(ctx context.Context, id primitive.ObjectID) (bool, error) {
	res, err := ds.collection.UpdateOne(
		ctx,
		bson.M{"_id": id, "deletedAt": bson.M{"$ne": nil}},
		bson.M{"$set": bson.M{"deletedAt": nil, "updatedAt": time.Now()}},
	)
	if err != nil {
		return false, err
	}
	return res.ModifiedCount > 0, nil
}

// HardDelete permanently removes a document from the collection.
func (ds *MongoLocationDatasource) HardDelete(ctx context.Context, id primitive.ObjectID) (bool, error) {
	res, err := ds.collection.DeleteOne(ctx, bson.M{"_id": id})
	if err != nil {
		return false, err
	}
	return res.DeletedCount > 0, nil
}

// BulkSoftDelete marks multiple documents as deleted, then returns the IDs that were actually updated.
func (ds *MongoLocationDatasource) BulkSoftDelete(ctx context.Context, ids []primitive.ObjectID) ([]primitive.ObjectID, error) {
	filter := bson.M{"_id": bson.M{"$in": ids}, "deletedAt": nil}
	update := bson.M{"$set": bson.M{"deletedAt": time.Now(), "updatedAt": time.Now()}}

	if _, err := ds.collection.UpdateMany(ctx, filter, update); err != nil {
		return nil, err
	}

	// now fetch those with deletedAt != nil
	cur, err := ds.collection.Find(ctx, bson.M{
		"_id":        bson.M{"$in": ids},
		"deletedAt": bson.M{"$ne": nil},
	})
	if err != nil {
		return nil, err
	}
	defer cur.Close(ctx)

	var out []primitive.ObjectID
	for cur.Next(ctx) {
		var lm model.LocationModel
		if err := cur.Decode(&lm); err == nil {
			out = append(out, lm.ID)
		}
	}
	return out, nil
}



func (ds *MongoLocationDatasource) BulkRestore(ctx context.Context, ids []primitive.ObjectID) ([]primitive.ObjectID, error) {
	filter := bson.M{
		"_id":       bson.M{"$in": ids},
		"deletedAt": bson.M{"$ne": nil}, // Only restore soft-deleted items
	}
	update := bson.M{
		"$set": bson.M{
			"deletedAt": nil,
			"updatedAt": time.Now(),
		},
	}

	result, err := ds.collection.UpdateMany(ctx, filter, update)
	if err != nil {
		return nil, err
	}

	// Return IDs that were actually updated
	var updatedIDs []primitive.ObjectID
	if result.ModifiedCount > 0 {
		// Find the documents that were actually restored
		cursor, err := ds.collection.Find(ctx, bson.M{
			"_id":       bson.M{"$in": ids},
			"deletedAt": nil,
		})
		if err != nil {
			return nil, err
		}
		defer cursor.Close(ctx)

		for cursor.Next(ctx) {
			var perm model.LocationModel
			if err := cursor.Decode(&perm); err == nil {
				updatedIDs = append(updatedIDs, perm.ID)
			}
		}
	}

	return updatedIDs, nil
}