package configs

import (
	"log"
	"os"
	"path/filepath"
	"strconv"

	"github.com/joho/godotenv"
)

type Config struct {
	Port         string
	Env          string
	MongoURI     string
	DBName       string
	RedisURI     string
	JWTSecret    string
	JWTExpiresIn int

	// S3 Configuration
	S3AccessKey string
	S3SecretKey string
	S3Region    string
	S3Bucket    string
	S3BaseURL   string

	// File Upload Limits
	MaxFileSize int64
}

var AppConfig *Config

func LoadConfig() (*Config, error) {
	// Print current working directory for debugging
	currentDir, _ := os.Getwd()
	log.Printf("Current working directory: %s", currentDir)

	// Try multiple possible locations for .env file
	envPaths := []string{
		".env",                              // Current directory
		"../.env",                           // Parent directory
		"../../.env",                        // Go up two directories
		filepath.Join("configs", ".env"),    // configs directory
		filepath.Join("cmd", "api", ".env"), // cmd/api directory
	}

	envLoaded := false
	for _, path := range envPaths {
		if _, err := os.Stat(path); err == nil {
			err := godotenv.Load(path)
			if err == nil {
				log.Printf("Loaded .env from %s", path)
				envLoaded = true
				break
			}
		}
	}

	if !envLoaded {
		// Try one more time with godotenv's default behavior
		err := godotenv.Load()
		if err != nil {
			log.Println("No .env file found, using environment variables")
		} else {
			log.Println("Loaded .env file")
			envLoaded = true
		}
	}

	// Parse JWT expiration time
	jwtExpires, err := strconv.Atoi(GetEnv("JWT_EXPIRES_IN", "24"))
	if err != nil {
		jwtExpires = 24 // Default to 24 hours
	}

	// Parse max file size (default 5MB)
	maxFileSize, err := strconv.ParseInt(GetEnv("MAX_FILE_SIZE", "5"), 10, 64)
	if err != nil {
		maxFileSize = 5 // Default to 5MB
	}
	maxFileSize = maxFileSize * 1024 * 1024 // Convert to bytes

	AppConfig = &Config{
		Port:         GetEnv("PORT", "8080"),
		Env:          GetEnv("ENV", "development"),
		MongoURI:     GetEnv("MONGODB_URI", "mongodb://mongo:27017"),
		DBName:       GetEnv("DB_NAME", "wecare_holidays"),
		RedisURI:     GetEnv("REDIS_URI", "redis://redis:6379"),
		JWTSecret:    GetEnv("JWT_SECRET", ""),
		JWTExpiresIn: jwtExpires,

		// S3 Configuration
		S3AccessKey: GetEnv("AWS_ACCESS_KEY", ""),
		S3SecretKey: GetEnv("AWS_SECRET_KEY", ""),
		S3Region:    GetEnv("AWS_REGION", "ap-south-1"),
		S3Bucket:    GetEnv("S3_BUCKET", ""),
		S3BaseURL:   GetEnv("S3_BASE_URL", ""),

		// File Upload Limits
		MaxFileSize: maxFileSize,
	}
	return AppConfig, nil
}

func GetEnv(key, fallback string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return fallback
}
