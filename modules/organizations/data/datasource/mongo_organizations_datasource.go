// internal/modules/organizations/data/datasource/mongo_organizations_datasource.go
package datasource

import (
	"context"
	"log"
	"time"

	mongodb "bitbucket.org/abhishek_fordel/we-care-holidays-backend-golang/modules/organizations/data/mongodb/indexes"
	"bitbucket.org/abhishek_fordel/we-care-holidays-backend-golang/modules/organizations/data/mongodb/model"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// MongoOrganizationDatasource handles raw MongoDB operations for organizations
type MongoOrganizationDatasource struct {
	collection *mongo.Collection
}

// NewMongoOrganizationDatasource creates a new instance of the organization datasource
func NewMongoOrganizationDatasource(db *mongo.Database) *MongoOrganizationDatasource {
	collection := db.Collection(model.OrganizationModel{}.CollectionName())

	if err := mongodb.SetupOrganizationIndexes(collection); err != nil {
		log.Printf("⚠️ Failed to setup organization indexes: %v", err)
	}
	return &MongoOrganizationDatasource{
		collection: collection,
	}
}

// Insert inserts a new organization document into the collection
func (ds *MongoOrganizationDatasource) Insert(ctx context.Context, organization *model.OrganizationModel) error {
	now := time.Now()
	organization.CreatedAt = now
	organization.UpdatedAt = now

	result, err := ds.collection.InsertOne(ctx, organization)
	if err != nil {
		return err
	}

	// Capture MongoDB-assigned ObjectID
	organization.ID = result.InsertedID.(primitive.ObjectID)
	return nil
}

// FindByFilters retrieves organization documents with optional filters, pagination, and sorting
func (ds *MongoOrganizationDatasource) FindByFilters(ctx context.Context, filters map[string]interface{}, page int, limit int) ([]model.OrganizationModel, int64, error) {
	// Step 1: Count total for pagination
	totalCount, err := ds.collection.CountDocuments(ctx, filters)
	if err != nil {
		return nil, 0, err
	}

	// Step 2: Prepare options
	opts := options.Find().
		SetSkip(int64((page - 1) * limit)).
		SetLimit(int64(limit)).
		SetSort(bson.D{{Key: "createdAt", Value: -1}}) // Default sort by creation date descending

	// Step 3: Query
	cursor, err := ds.collection.Find(ctx, filters, opts)
	if err != nil {
		return nil, 0, err
	}
	defer cursor.Close(ctx)

	var organizations []model.OrganizationModel
	if err := cursor.All(ctx, &organizations); err != nil {
		return nil, 0, err
	}

	return organizations, totalCount, nil
}

// FindByID finds an organization by its ID
func (ds *MongoOrganizationDatasource) FindByID(ctx context.Context, id primitive.ObjectID) (*model.OrganizationModel, error) {
	filter := bson.M{"_id": id}

	var organization model.OrganizationModel
	err := ds.collection.FindOne(ctx, filter).Decode(&organization)
	if err != nil {
		return nil, err
	}

	return &organization, nil
}

// FindBySlug finds an organization by its slug
func (ds *MongoOrganizationDatasource) FindBySlug(ctx context.Context, slug string) (*model.OrganizationModel, error) {
	filter := bson.M{"slug": slug, "deletedAt": nil}

	var organization model.OrganizationModel
	err := ds.collection.FindOne(ctx, filter).Decode(&organization)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, nil // Return nil if not found
		}
		return nil, err
	}

	return &organization, nil
}

// Update updates an organization document
func (ds *MongoOrganizationDatasource) Update(ctx context.Context, organization *model.OrganizationModel) error {
	organization.UpdatedAt = time.Now()

	filter := bson.M{"_id": organization.ID}
	update := bson.M{"$set": organization}

	_, err := ds.collection.UpdateOne(ctx, filter, update)
	return err
}

// SoftDelete marks an organization as deleted by setting deletedAt timestamp
func (ds *MongoOrganizationDatasource) SoftDelete(ctx context.Context, id primitive.ObjectID) (bool, error) {
	filter := bson.M{"_id": id, "deletedAt": nil}
	update := bson.M{
		"$set": bson.M{
			"deletedAt": time.Now(),
			"updatedAt": time.Now(),
		},
	}

	result, err := ds.collection.UpdateOne(ctx, filter, update)
	if err != nil {
		return false, err
	}

	return result.ModifiedCount > 0, nil
}

// Restore restores a soft-deleted organization by setting deletedAt to nil
func (ds *MongoOrganizationDatasource) Restore(ctx context.Context, id primitive.ObjectID) (bool, error) {
	filter := bson.M{"_id": id, "deletedAt": bson.M{"$ne": nil}}
	update := bson.M{
		"$set": bson.M{
			"deletedAt": nil,
			"updatedAt": time.Now(),
		},
	}

	result, err := ds.collection.UpdateOne(ctx, filter, update)
	if err != nil {
		return false, err
	}

	return result.ModifiedCount > 0, nil
}

// HardDelete permanently removes an organization from the database
func (ds *MongoOrganizationDatasource) HardDelete(ctx context.Context, id primitive.ObjectID) (bool, error) {
	filter := bson.M{"_id": id}

	result, err := ds.collection.DeleteOne(ctx, filter)
	if err != nil {
		return false, err
	}

	return result.DeletedCount > 0, nil
}

// BulkSoftDelete marks multiple organizations as deleted
func (ds *MongoOrganizationDatasource) BulkSoftDelete(ctx context.Context, ids []primitive.ObjectID) ([]primitive.ObjectID, error) {
	filter := bson.M{
		"_id":       bson.M{"$in": ids},
		"deletedAt": nil,
	}
	update := bson.M{
		"$set": bson.M{
			"deletedAt": time.Now(),
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
		// Since we can't get the exact IDs that were updated from UpdateMany,
		// we'll assume all valid IDs were updated based on the filter
		cursor, err := ds.collection.Find(ctx, bson.M{
			"_id":       bson.M{"$in": ids},
			"deletedAt": bson.M{"$ne": nil},
		})
		if err != nil {
			return nil, err
		}
		defer cursor.Close(ctx)

		for cursor.Next(ctx) {
			var org model.OrganizationModel
			if err := cursor.Decode(&org); err == nil {
				updatedIDs = append(updatedIDs, org.ID)
			}
		}
	}

	return updatedIDs, nil
}

// UpdateStatus updates the status of an organization
func (ds *MongoOrganizationDatasource) UpdateStatus(ctx context.Context, id primitive.ObjectID, status string) error {
	filter := bson.M{"_id": id}
	update := bson.M{
		"$set": bson.M{
			"status":    status,
			"updatedAt": time.Now(),
		},
	}

	_, err := ds.collection.UpdateOne(ctx, filter, update)
	return err
}



func (ds *MongoOrganizationDatasource) BulkRestore(ctx context.Context, ids []primitive.ObjectID) ([]primitive.ObjectID, error) {
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
			var perm model.OrganizationModel
			if err := cursor.Decode(&perm); err == nil {
				updatedIDs = append(updatedIDs, perm.ID)
			}
		}
	}

	return updatedIDs, nil
}


func (ds *MongoOrganizationDatasource) ExistsByID(ctx context.Context, id primitive.ObjectID) (bool, error) {
	filter := bson.M{"_id": id}
	
	count, err := ds.collection.CountDocuments(ctx, filter)
	if err != nil {
		return false, err
	}
	
	return count > 0, nil
}