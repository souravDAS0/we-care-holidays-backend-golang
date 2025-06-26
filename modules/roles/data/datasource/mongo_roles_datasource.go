package datasource

import (
	"context"
	"time"

	"bitbucket.org/abhishek_fordel/we-care-holidays-backend-golang/modules/roles/data/mongodb/model"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)


type MongoRoleDatasource struct {
	collection *mongo.Collection
}


func NewMongoRoleDatasource(db *mongo.Database) *MongoRoleDatasource {
	collection := db.Collection(model.RoleModel{}.CollectionName()) // Use the appropriate collection name

	return &MongoRoleDatasource{
		collection: collection,
	}
}


// Insert inserts a new role document into the collection
func (ds *MongoRoleDatasource) Insert(ctx context.Context, role *model.RoleModel) error {
	now := time.Now()
	role.CreatedAt = now
	role.UpdatedAt = now

	result, err := ds.collection.InsertOne(ctx, role)
	if err != nil {
		return err
	}

	// Capture MongoDB-assigned ObjectID
	role.ID = result.InsertedID.(primitive.ObjectID)
	return nil
}

// FindByFilters retrieves role documents with optional filters, pagination, and sorting
func (ds *MongoRoleDatasource) FindByFilters(ctx context.Context, filters map[string]interface{}, page int, limit int) ([]model.RoleModel, int64, error) {
	// Step 1: Count total for pagination
	totalCount, err := ds.collection.CountDocuments(ctx, filters)
	if err != nil {
		return nil, 0, err
	}

	// Step 2: Prepare options
	opts := options.Find().
		SetSkip(int64((page - 1) * limit)).
		SetLimit(int64(limit)).
		SetSort(bson.D{{Key: "createdAt", Value: -1}}) // Default sort by createdAt descending

	// Step 3: Execute query
	cursor, err := ds.collection.Find(ctx, filters, opts)
	if err != nil {
		return nil, 0, err
	}
	defer cursor.Close(ctx)

	var roles []model.RoleModel
	if err := cursor.All(ctx, &roles); err != nil {
		return nil, 0, err
	}

	return roles, totalCount, nil
}

// FindByID retrieves a role document by its ID
func (ds *MongoRoleDatasource) FindByID(ctx context.Context, id primitive.ObjectID) (*model.RoleModel, error) {
	var role model.RoleModel
	filter := bson.M{"_id": id}

	err := ds.collection.FindOne(ctx, filter).Decode(&role)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, nil // No document found
		}
		return nil, err // Other errors
	}

	return &role, nil
}

// FindByName retrieves a role document by name (only non-deleted roles)
func (ds *MongoRoleDatasource) FindByName(ctx context.Context, name string) (*model.RoleModel, error) {
	var role model.RoleModel
	filter := bson.M{
		"name":      name,
		"deletedAt": nil, // Only find non-deleted roles
	}

	err := ds.collection.FindOne(ctx, filter).Decode(&role)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, nil // No document found
		}
		return nil, err // Other errors
	}

	return &role, nil
}



// Update updates an role document
func (ds *MongoRoleDatasource) Update(ctx context.Context, role *model.RoleModel) error {
	role.UpdatedAt = time.Now()

	filter := bson.M{"_id": role.ID}
	update := bson.M{"$set": role}

	_, err := ds.collection.UpdateOne(ctx, filter, update)
	return err
}

// SoftDelete marks an role as deleted by setting deletedAt timestamp
func (ds *MongoRoleDatasource) SoftDelete(ctx context.Context, id primitive.ObjectID) (bool, error) {
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

// Restore restores a soft-deleted role by setting deletedAt to nil
func (ds *MongoRoleDatasource) Restore(ctx context.Context, id primitive.ObjectID) (bool, error) {
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

// HardDelete permanently removes an role from the database
func (ds *MongoRoleDatasource) HardDelete(ctx context.Context, id primitive.ObjectID) (bool, error) {
	filter := bson.M{"_id": id}

	result, err := ds.collection.DeleteOne(ctx, filter)
	if err != nil {
		return false, err
	}

	return result.DeletedCount > 0, nil
}

// BulkSoftDelete marks multiple organizations as deleted
func (ds *MongoRoleDatasource) BulkSoftDelete(ctx context.Context, ids []primitive.ObjectID) ([]primitive.ObjectID, error) {
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
			var org model.RoleModel
			if err := cursor.Decode(&org); err == nil {
				updatedIDs = append(updatedIDs, org.ID)
			}
		}
	}

	return updatedIDs, nil
}


func (ds *MongoRoleDatasource) BulkRestore(ctx context.Context, ids []primitive.ObjectID) ([]primitive.ObjectID, error) {
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
			var perm model.RoleModel
			if err := cursor.Decode(&perm); err == nil {
				updatedIDs = append(updatedIDs, perm.ID)
			}
		}
	}

	return updatedIDs, nil
}


func (ds *MongoRoleDatasource) ExistsByID(ctx context.Context, id primitive.ObjectID) (bool, error) {
	filter := bson.M{"_id": id}
	
	count, err := ds.collection.CountDocuments(ctx, filter)
	if err != nil {
		return false, err
	}
	
	return count > 0, nil
}

func (ds *MongoRoleDatasource) ExistsByName(ctx context.Context, name string) (bool, error) {
	filter := bson.M{
		"name":      name,
		"deletedAt": nil, 
	}
	
	count, err := ds.collection.CountDocuments(ctx, filter)
	if err != nil {
		return false, err
	}
	
	return count > 0, nil
}
