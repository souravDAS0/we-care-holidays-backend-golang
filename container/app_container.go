package container

import (
	"context"
	"log"

	"bitbucket.org/abhishek_fordel/we-care-holidays-backend-golang/commons/services"
	"bitbucket.org/abhishek_fordel/we-care-holidays-backend-golang/configs"
	"bitbucket.org/abhishek_fordel/we-care-holidays-backend-golang/internal/middleware"
	"github.com/redis/go-redis/v9"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// AppContainer holds global shared dependencies
type AppContainer struct {
	Config              *configs.Config
	MongoClient         *mongo.Client
	MongoDatabase       *mongo.Database
	RedisClient         *redis.Client
	FileService         services.FileService
	RBACService         middleware.RBACService
	PermissionValidator *middleware.PermissionValidator

	// Module containers
	Permission   *PermissionContainer
	Role         *RoleContainer
	User         *UserContainer
	Organization *OrganizationContainer
	Location     *LocationContainer
}

func BuildAppContainer(cfg *configs.Config) *AppContainer {
	mongoClient, mongoDatabase := initMongo(cfg)
	redisClient := initRedis(cfg)
	fileService := initFileService(cfg)

	return &AppContainer{
		Config:        cfg,
		MongoClient:   mongoClient,
		MongoDatabase: mongoDatabase,
		RedisClient:   redisClient,
		FileService:   fileService,
	}
}

func initMongo(cfg *configs.Config) (*mongo.Client, *mongo.Database) {
	clientOptions := options.Client().ApplyURI(cfg.MongoURI)
	client, err := mongo.Connect(context.Background(), clientOptions)
	if err != nil {
		log.Fatalf("failed to connect to MongoDB: %v", err)
	}

	if err := client.Ping(context.Background(), nil); err != nil {
		log.Fatalf("failed to ping MongoDB: %v", err)
	}

	log.Println("MongoDB connected successfully")
	return client, client.Database(cfg.DBName)
}

// initRedis initializes Redis connection
func initRedis(cfg *configs.Config) *redis.Client {
	opt, err := redis.ParseURL(cfg.RedisURI)
	if err != nil {
		log.Fatalf("failed to parse Redis URI: %v", err)
	}

	client := redis.NewClient(opt)

	if err := client.Ping(context.Background()).Err(); err != nil {
		log.Fatalf("failed to connect to Redis: %v", err)
	}

	log.Println("Redis connected successfully")
	return client
}

// initFileService initializes file service (S3)
func initFileService(cfg *configs.Config) services.FileService {
	if cfg.S3AccessKey != "" && cfg.S3Bucket != "" {
		s3Config := services.S3Config{
			AccessKey:  cfg.S3AccessKey,
			SecretKey:  cfg.S3SecretKey,
			Region:     cfg.S3Region,
			BucketName: cfg.S3Bucket,
			BaseURL:    cfg.S3BaseURL,
		}

		fileService, err := services.NewS3FileService(s3Config)
		if err != nil {
			log.Printf("Warning: Failed to initialize S3 file service: %v", err)
			log.Println("File upload functionality will be disabled")
			return nil
		}

		log.Println("S3 file service initialized successfully")
		return fileService
	}

	// Return nil when S3 configuration is not provided
	log.Println("S3 configuration not provided, file service will be disabled")
	return nil
}

func (ac *AppContainer) InjectRBACServices() {
	// Add debugging logs
	log.Printf("User container: %v", ac.User)
	log.Printf("Role container: %v", ac.Role)
	log.Printf("Permission container: %v", ac.Permission)

	if ac.User != nil {
		log.Printf("User repository: %v", ac.User.Repository)
	}
	if ac.Role != nil {
		log.Printf("Role repository: %v", ac.Role.Repository)
	}
	if ac.Permission != nil {
		log.Printf("Permission repository: %v", ac.Permission.Repository)
	}

	// Ensure all required containers are available
	if ac.User == nil || ac.Role == nil || ac.Permission == nil {
		log.Fatal("Required module containers not available for RBAC initialization")
	}

	// Create RBAC service with individual repositories
	ac.RBACService = middleware.NewRBACService(
		ac.User.Repository,
		ac.Role.Repository,
		ac.Permission.Repository,
	)

	log.Printf("RBAC Service created: %v", ac.RBACService)
}
