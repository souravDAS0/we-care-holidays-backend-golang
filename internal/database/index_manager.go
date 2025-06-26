// internal/database/index_manager.go
package database

import (
	"context"
	"fmt"
	"log"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"

	locMongodb "bitbucket.org/abhishek_fordel/we-care-holidays-backend-golang/modules/locations/data/mongodb/indexes"
	orgMongodb "bitbucket.org/abhishek_fordel/we-care-holidays-backend-golang/modules/organizations/data/mongodb/indexes"
	permMongodb "bitbucket.org/abhishek_fordel/we-care-holidays-backend-golang/modules/permissions/data/mongodb/indexes"
	usersMongodb "bitbucket.org/abhishek_fordel/we-care-holidays-backend-golang/modules/users/data/mongodb/indexes"
)

type IndexManager struct {
	db *mongo.Database
}

func NewIndexManager(db *mongo.Database) *IndexManager {
	return &IndexManager{db: db}
}

// DropConflictingIndexes drops indexes that might conflict with new ones
func (im *IndexManager) DropConflictingIndexes() error {
	collections := map[string][]string{
		"organizations": {
			"idx_deleted_at", // Conflicting name
			"idx_slug",       // May conflict with new unique constraint
		},
		"permissions": {
			"unique_resource_action_scope", // Has invalid partial filter
		},
		"locations": {
			"idx_name_type", // Has invalid partial filter
		},
		"users": {
			"idx_primary_email", // May conflict with new constraint
		},
	}

	for collectionName, indexNames := range collections {
		coll := im.db.Collection(collectionName)
		
		for _, indexName := range indexNames {
			if err := im.dropIndexIfExists(coll, indexName); err != nil {
				log.Printf("Warning: Failed to drop index %s in collection %s: %v", indexName, collectionName, err)
				// Continue with other indexes even if one fails
			}
		}
	}

	return nil
}

// dropIndexIfExists drops an index if it exists, ignoring "index not found" errors
func (im *IndexManager) dropIndexIfExists(coll *mongo.Collection, indexName string) error {
	ctx := context.Background()
	
	// Check if index exists first
	cursor, err := coll.Indexes().List(ctx)
	if err != nil {
		return fmt.Errorf("failed to list indexes: %w", err)
	}
	defer cursor.Close(ctx)

	var indexExists bool
	for cursor.Next(ctx) {
		var index bson.M
		if err := cursor.Decode(&index); err != nil {
			continue
		}
		
		if name, ok := index["name"].(string); ok && name == indexName {
			indexExists = true
			break
		}
	}

	if !indexExists {
		log.Printf("Index %s does not exist, skipping drop", indexName)
		return nil
	}

	// Drop the index
	_, err = coll.Indexes().DropOne(ctx, indexName)
	if err != nil {
		// Check if it's a "index not found" error and ignore it
		
		return fmt.Errorf("failed to drop index %s: %w", indexName, err)
	}

	log.Printf("Successfully dropped index: %s", indexName)
	return nil
}

// SetupAllIndexes sets up indexes for all collections
func (im *IndexManager) SetupAllIndexes() error {
	log.Println("Setting up database indexes...")

	// Drop conflicting indexes first
	if err := im.DropConflictingIndexes(); err != nil {
		log.Printf("Warning: Some conflicting indexes could not be dropped: %v", err)
	}

	// Setup indexes for each collection
	collections := map[string]func(*mongo.Collection) error{
		"organizations": orgMongodb.SetupOrganizationIndexes,
		"users":         usersMongodb.SetupUserIndexes,
		"permissions":   permMongodb.SetupPermissionIndexes,
		"locations":    locMongodb.SetupLocationIndexes,
	}

	for collectionName, setupFunc := range collections {
		coll := im.db.Collection(collectionName)
		
		log.Printf("Setting up %s indexes...", collectionName)
		if err := setupFunc(coll); err != nil {
			log.Printf("⚠️ Failed to setup %s indexes: %v", collectionName, err)
			// Continue with other collections even if one fails
		} else {
			log.Printf("✅ Successfully set up %s indexes", collectionName)
		}
	}

	log.Println("Index setup completed")
	return nil
}

// ListAllIndexes lists all indexes in the database for debugging
func (im *IndexManager) ListAllIndexes() error {
	collections := []string{"organizations", "users", "permissions", "locations", "roles"}
	
	for _, collectionName := range collections {
		fmt.Printf("\n=== Indexes for %s ===\n", collectionName)
		coll := im.db.Collection(collectionName)
		
		cursor, err := coll.Indexes().List(context.Background())
		if err != nil {
			fmt.Printf("Error listing indexes for %s: %v\n", collectionName, err)
			continue
		}
		defer cursor.Close(context.Background())

		for cursor.Next(context.Background()) {
			var index bson.M
			if err := cursor.Decode(&index); err != nil {
				continue
			}
			
			name := index["name"].(string)
			keys := index["key"].(bson.M)
			
			fmt.Printf("Name: %s, Keys: %v", name, keys)
			
			if unique, ok := index["unique"].(bool); ok && unique {
				fmt.Printf(", Unique: true")
			}
			
			if sparse, ok := index["sparse"].(bool); ok && sparse {
				fmt.Printf(", Sparse: true")
			}
			
			if partialFilter, ok := index["partialFilterExpression"].(bson.M); ok {
				fmt.Printf(", PartialFilter: %v", partialFilter)
			}
			
			fmt.Println()
		}
	}
	
	return nil
}

// Usage in your main application setup
func SetupDatabaseIndexes(db *mongo.Database) error {
	indexManager := NewIndexManager(db)
	return indexManager.SetupAllIndexes()
}

// Debug function to list all indexes
func ListDatabaseIndexes(db *mongo.Database) error {
	indexManager := NewIndexManager(db)
	return indexManager.ListAllIndexes()
}