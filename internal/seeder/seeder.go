package seeder

import (
	"context"
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

	// Execute seeding
	if err := ExecuteSeeder(appContainer, data); err != nil {
		logger.Log.Error("❌ Seeder execution failed", zap.Error(err))
		return err
	}

	logger.Log.Info("✅ Seeder completed successfully")
	return nil
}

// purgeSeedCollections drops existing seed collections
func purgeSeedCollections(appContainer *container.AppContainer) error {
	ctx := context.Background()

	collections := []string{
		"permissions",
		"roles",
		"users",
	}

	for _, coll := range collections {
		if err := appContainer.MongoDatabase.Collection(coll).Drop(ctx); err != nil {
			logger.Log.Error("⚠️ Failed to purge collection", zap.String("collection", coll), zap.Error(err))
			return err
		}
		logger.Log.Info("💥 Collection purged", zap.String("collection", coll))
	}

	return nil
}
