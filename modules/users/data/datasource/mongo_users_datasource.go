package datasource

import (
	"context"
	"time"

	"bitbucket.org/abhishek_fordel/we-care-holidays-backend-golang/modules/users/data/mongodb/model"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type MongoUserDatasource struct {
	collection *mongo.Collection
}

func NewMongoUserDatasource(db *mongo.Database) *MongoUserDatasource {
	collection := db.Collection(model.UserModel{}.CollectionName()) // Use the appropriate collection name

	// if err := mongodb.SetupUserIndexes(collection); err != nil {
	// 	log.Printf("⚠️ Failed to setup user indexes: %v", err)
	// }

	return &MongoUserDatasource{
		collection: collection,
	}
}

// Insert inserts a new user document into the collection
func (ds *MongoUserDatasource) Insert(ctx context.Context, user *model.UserModel) error {
	now := time.Now()
	user.CreatedAt = now
	user.UpdatedAt = now

	result, err := ds.collection.InsertOne(ctx, user)
	if err != nil {
		return err
	}

	// Capture MongoDB-assigned ObjectID
	user.ID = result.InsertedID.(primitive.ObjectID)
	return nil
}

// FindByFilters retrieves user documents with optional filters, pagination, and sorting
func (ds *MongoUserDatasource) FindByFilters(ctx context.Context, filters map[string]interface{}, page int, limit int) ([]model.UserModel, int64, error) {
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

	var users []model.UserModel
	if err := cursor.All(ctx, &users); err != nil {
		return nil, 0, err
	}

	return users, totalCount, nil
}

// FindByID finds an user by its ID
func (ds *MongoUserDatasource) FindByID(ctx context.Context, id primitive.ObjectID) (*model.UserModel, error) {
	filter := bson.M{"_id": id}

	var user model.UserModel
	err := ds.collection.FindOne(ctx, filter).Decode(&user)
	if err != nil {
		return nil, err
	}

	return &user, nil
}

func (ds *MongoUserDatasource) FindByEmail(ctx context.Context, email string) (*model.UserModel, error) {
	// Use "emails.email" to search within the emails array
	filter := bson.M{
		"emails.email": email,
		"deletedAt":    nil,
	}

	var user model.UserModel
	err := ds.collection.FindOne(ctx, filter).Decode(&user)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, nil // Return nil if not found
		}
		return nil, err
	}

	return &user, nil
}

func (ds *MongoUserDatasource) FindByPhone(ctx context.Context, phone string) (*model.UserModel, error) {
	filter := bson.M{
		"phones.number": phone, // Changed from "phone" to "phones.number"
		"deletedAt":     nil,
	}

	var user model.UserModel
	err := ds.collection.FindOne(ctx, filter).Decode(&user)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, nil
		}
		return nil, err
	}

	return &user, nil
}

// Update updates an user document
func (ds *MongoUserDatasource) Update(ctx context.Context, user *model.UserModel) error {
	user.UpdatedAt = time.Now()

	filter := bson.M{"_id": user.ID}
	update := bson.M{"$set": user}

	_, err := ds.collection.UpdateOne(ctx, filter, update)
	return err
}

// SoftDelete marks an user as deleted by setting deletedAt timestamp
func (ds *MongoUserDatasource) SoftDelete(ctx context.Context, id primitive.ObjectID) (bool, error) {
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

// Restore restores a soft-deleted user by setting deletedAt to nil
func (ds *MongoUserDatasource) Restore(ctx context.Context, id primitive.ObjectID) (bool, error) {
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

// HardDelete permanently removes an user from the database
func (ds *MongoUserDatasource) HardDelete(ctx context.Context, id primitive.ObjectID) (bool, error) {
	filter := bson.M{"_id": id}

	result, err := ds.collection.DeleteOne(ctx, filter)
	if err != nil {
		return false, err
	}

	return result.DeletedCount > 0, nil
}

// BulkSoftDelete marks multiple users as deleted
func (ds *MongoUserDatasource) BulkSoftDelete(ctx context.Context, ids []primitive.ObjectID) ([]primitive.ObjectID, error) {
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
			var org model.UserModel
			if err := cursor.Decode(&org); err == nil {
				updatedIDs = append(updatedIDs, org.ID)
			}
		}
	}

	return updatedIDs, nil
}

// UpdateStatus updates the status of an user
func (ds *MongoUserDatasource) UpdateStatus(ctx context.Context, id primitive.ObjectID, status string) error {
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

func (ds *MongoUserDatasource) BulkRestore(ctx context.Context, ids []primitive.ObjectID) ([]primitive.ObjectID, error) {
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
			var perm model.UserModel
			if err := cursor.Decode(&perm); err == nil {
				updatedIDs = append(updatedIDs, perm.ID)
			}
		}
	}

	return updatedIDs, nil
}

func (ds *MongoUserDatasource) ExistsByID(ctx context.Context, id primitive.ObjectID) (bool, error) {
	filter := bson.M{"_id": id}

	count, err := ds.collection.CountDocuments(ctx, filter)
	if err != nil {
		return false, err
	}

	return count > 0, nil
}

func (ds *MongoUserDatasource) ExistsByEmail(ctx context.Context, email string) (bool, error) {
	filter := bson.M{
		"emails.email": email, // Changed from "email" to "emails.email"
		"deletedAt":    nil,
	}

	count, err := ds.collection.CountDocuments(ctx, filter)
	if err != nil {
		return false, err
	}

	return count > 0, nil
}

// ExistsByPhone checks if a user exists with the given phone number
func (ds *MongoUserDatasource) ExistsByPhone(ctx context.Context, phone string) (bool, error) {
	filter := bson.M{
		"phones.number": phone, // Changed from "phone" to "phones.number"
		"deletedAt":     nil,
	}

	count, err := ds.collection.CountDocuments(ctx, filter)
	if err != nil {
		return false, err
	}

	return count > 0, nil
}
