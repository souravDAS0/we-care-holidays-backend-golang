package seeder

import (
	"context"
	"fmt"
	"strconv"

	"bitbucket.org/abhishek_fordel/we-care-holidays-backend-golang/container"
	"bitbucket.org/abhishek_fordel/we-care-holidays-backend-golang/internal/logger"
	"go.uber.org/zap"
)

// RunSeeder triggers the seeding process.
//
// Parameters:
// - appContainer: global DI container.
// - force: if true, cleans existing data before seeding.
//
// Returns:
// - error: if seeding fails.
func RunSeeder(appContainer *container.AppContainer, force bool) error {
	logger.Log.Info("🚀 Seeder started", zap.String("force", strconv.FormatBool(force)))

	// Validate container
	if err := validateContainer(appContainer); err != nil {
		logger.Log.Error("❌ Container validation failed", zap.Error(err))
		return err
	}

	// If force mode, purge collections
	if force {
		if err := purgeSeedCollections(appContainer); err != nil {
			logger.Log.Error("❌ Failed to purge collections", zap.Error(err))
			return err
		}
		logger.Log.Info("✅ Purged existing seed collections")
	}

	// Load seed data from JSON
	data, err := LoadSeedData()
	if err != nil {
		logger.Log.Error("❌ Failed to load seed data", zap.Error(err))
		return err
	}

	logger.Log.Info("📋 Loaded seed data",
		zap.Int("permissions", len(data.Permissions)),
		zap.Int("organizations", len(data.Organizations)),
		zap.Int("roles", len(data.Roles)),
		zap.Int("users", len(data.Users)))

	// Execute seeding in proper order
	if err := ExecuteSeeder(appContainer, data); err != nil {
		logger.Log.Error("❌ Seeder execution failed", zap.Error(err))
		return err
	}

	logger.Log.Info("✅ Seeder completed successfully")
	return nil
}

// validateContainer ensures all required components are initialized
func validateContainer(appContainer *container.AppContainer) error {
	if appContainer == nil {
		return fmt.Errorf("app container is nil")
	}
	if appContainer.MongoDatabase == nil {
		return fmt.Errorf("mongo database is nil")
	}
	if appContainer.Permission == nil {
		return fmt.Errorf("permission container is nil")
	}
	if appContainer.Organization == nil {
		return fmt.Errorf("organization container is nil")
	}
	if appContainer.Role == nil {
		return fmt.Errorf("role container is nil")
	}
	if appContainer.User == nil {
		return fmt.Errorf("user container is nil")
	}
	return nil
}

// purgeSeedCollections drops existing seed collections
func purgeSeedCollections(appContainer *container.AppContainer) error {
	ctx := context.Background()

	// Order matters: delete dependents first
	collections := []string{
		"users",         // Users depend on roles and organizations
		"roles",         // Roles depend on permissions
		"organizations", // Organizations are independent (except for potential references)
		"permissions",   // Permissions are base entities
	}

	for _, coll := range collections {
		if err := appContainer.MongoDatabase.Collection(coll).Drop(ctx); err != nil {
			logger.Log.Warn("⚠️ Failed to purge collection (might not exist)", zap.String("collection", coll), zap.Error(err))
			// Don't return error here as collection might not exist
		} else {
			logger.Log.Info("💥 Collection purged", zap.String("collection", coll))
		}
	}

	return nil
}
