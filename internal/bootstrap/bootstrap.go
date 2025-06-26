package bootstrap

import (
	"log"

	"bitbucket.org/abhishek_fordel/we-care-holidays-backend-golang/configs"
	"bitbucket.org/abhishek_fordel/we-care-holidays-backend-golang/container"
	"bitbucket.org/abhishek_fordel/we-care-holidays-backend-golang/internal/logger"
)

// Bootstrap sets up app config, DI container, etc.
func Bootstrap() *container.AppContainer {
	// Load app config
	cfg, err := configs.LoadConfig()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	// Initialize logger
	logger.InitLogger(cfg.Env)

	// Build global container
	appContainer := container.BuildAppContainer(cfg)

	// Inject module containers
	appContainer.InjectOrganizationContainer()
	appContainer.InjectPermissionContainer()
	appContainer.InjectRoleContainer()
	appContainer.InjectUserContainer()
	appContainer.InjectLocationContainer()

	appContainer.InjectRBACServices()

	return appContainer
}
