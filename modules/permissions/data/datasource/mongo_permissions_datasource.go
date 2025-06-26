package datasource

import (
	"context"
	"log"
	"time"

	mongodb "bitbucket.org/abhishek_fordel/we-care-holidays-backend-golang/modules/permissions/data/mongodb/indexes"
	"bitbucket.org/abhishek_fordel/we-care-holidays-backend-golang/modules/permissions/data/mongodb/model"
	"bitbucket.org/abhishek_fordel/we-care-holidays-backend-golang/modules/permissions/domain/entity"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// MongoPermissionDatasource handles raw MongoDB operations for Permissions
type MongoPermissionDatasource struct {
	collection *mongo.Collection
}

// NewMongoPermissionDatasource creates a new instance of the Permission datasource
func NewMongoPermissionDatasource(db *mongo.Database) *MongoPermissionDatasource {
	collection := db.Collection(model.PermissionModel{}.CollectionName())

	if err := mongodb.SetupPermissionIndexes(collection); err != nil {
		log.Printf("⚠️ Failed to setup Permission indexes: %v", err)
	}
	return &MongoPermissionDatasource{
		collection: collection,
	}
}

// Insert inserts a new permission document into the collection
func (ds *MongoPermissionDatasource) Insert(ctx context.Context, permission *model.PermissionModel) error {
	now := time.Now()
	permission.CreatedAt = now
	permission.UpdatedAt = now

	result, err := ds.collection.InsertOne(ctx, permission)
	if err != nil {
		return err
	}

	// Capture MongoDB-assigned ObjectID
	permission.ID = result.InsertedID.(primitive.ObjectID)
	return nil
}

// FindByFilters retrieves permission documents with optional filters, pagination, and sorting
func (ds *MongoPermissionDatasource) FindByFilters(ctx context.Context, filters map[string]interface{}, page int, limit int) ([]model.PermissionModel, int64, error) {
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

	var permissions []model.PermissionModel
	if err := cursor.All(ctx, &permissions); err != nil {
		return nil, 0, err
	}

	return permissions, totalCount, nil
}

// FindByID retrieves a permission document by its ID
func (ds *MongoPermissionDatasource) FindByID(ctx context.Context, id primitive.ObjectID) (*model.PermissionModel, error) {
	var permission model.PermissionModel
	filter := bson.M{"_id": id}

	err := ds.collection.FindOne(ctx, filter).Decode(&permission)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, nil // No document found
		}
		return nil, err // Other errors
	}

	return &permission, nil
}

// Update updates an permission document
func (ds *MongoPermissionDatasource) Update(ctx context.Context, permission *model.PermissionModel) error {
	permission.UpdatedAt = time.Now()

	filter := bson.M{"_id": permission.ID}
	update := bson.M{"$set": permission}

	_, err := ds.collection.UpdateOne(ctx, filter, update)
	return err
}

// HardDelete permanently removes an permission from the database
func (ds *MongoPermissionDatasource) HardDelete(ctx context.Context, id primitive.ObjectID) (bool, error) {
	filter := bson.M{"_id": id}

	result, err := ds.collection.DeleteOne(ctx, filter)
	if err != nil {
		return false, err
	}

	return result.DeletedCount > 0, nil
}

func (ds *MongoPermissionDatasource) ExistsByResourceAction(ctx context.Context, resource string, action entity.PermissionAction) (bool, error) {
	filter := bson.M{
		"resource":  resource,
		"action":    string(action),
		"deletedAt": nil,
	}

	count, err := ds.collection.CountDocuments(ctx, filter)
	if err != nil {
		return false, err
	}

	return count > 0, nil
}

func (ds *MongoPermissionDatasource) ExistsByResourceActionExcluding(ctx context.Context, resource string, action entity.PermissionAction, excludeID primitive.ObjectID) (bool, error) {
	filter := bson.M{
		"resource":  resource,
		"action":    string(action),
		"deletedAt": nil,
		"_id":       bson.M{"$ne": excludeID},
	}

	count, err := ds.collection.CountDocuments(ctx, filter)
	if err != nil {
		return false, err
	}

	return count > 0, nil
}

func (ds *MongoPermissionDatasource) ExistsByID(ctx context.Context, id primitive.ObjectID) (bool, error) {
	filter := bson.M{"_id": id}

	count, err := ds.collection.CountDocuments(ctx, filter)
	if err != nil {
		return false, err
	}

	return count > 0, nil
}
