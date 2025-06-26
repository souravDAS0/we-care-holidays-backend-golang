package services

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"path/filepath"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/google/uuid"
)

// S3FileService implements FileService using AWS S3
type S3FileService struct {
	client     *s3.Client
	bucketName string
	baseURL    string
}

// S3Config holds configuration for the S3 service
type S3Config struct {
	AccessKey  string
	SecretKey  string
	Region     string
	BucketName string
	BaseURL    string // Optional, can be empty if using the default S3 URL format
}

// NewS3FileService creates a new S3FileService
func NewS3FileService(cfg S3Config) (*S3FileService, error) {
	// Create AWS credentials provider
	credProvider := credentials.NewStaticCredentialsProvider(
		cfg.AccessKey,
		cfg.SecretKey,
		"", // Session token, typically empty for static credentials
	)

	// Load AWS configuration
	awsCfg, err := config.LoadDefaultConfig(
		context.Background(),
		config.WithRegion(cfg.Region),
		config.WithCredentialsProvider(credProvider),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to load AWS configuration: %w", err)
	}

	// Create S3 client
	client := s3.NewFromConfig(awsCfg)

	// Determine base URL
	baseURL := cfg.BaseURL
	if baseURL == "" {
		// Use default S3 URL format if not provided
		baseURL = fmt.Sprintf("https://%s.s3.%s.amazonaws.com", cfg.BucketName, cfg.Region)
	}

	return &S3FileService{
		client:     client,
		bucketName: cfg.BucketName,
		baseURL:    baseURL,
	}, nil
}

// UploadFile uploads a file to S3 and returns the URL
func (s *S3FileService) UploadFile(ctx context.Context, file io.Reader, originalFileName string, contentType string) (string, error) {
	// Read file content
	buffer, err := io.ReadAll(file)
	if err != nil {
		return "", fmt.Errorf("failed to read file: %w", err)
	}

	// Generate a unique filename to avoid collisions
	fileExt := filepath.Ext(originalFileName)
	fileName := fmt.Sprintf("%s%s", uuid.New().String(), fileExt)

	// Set the folder path in S3 - organize by year/month
	now := time.Now()
	s3Path := fmt.Sprintf("uploads/%d/%02d/%s", now.Year(), now.Month(), fileName)

	// Upload the file to S3
	_, err = s.client.PutObject(ctx, &s3.PutObjectInput{
		Bucket:      aws.String(s.bucketName),
		Key:         aws.String(s3Path),
		Body:        bytes.NewReader(buffer),
		ContentType: aws.String(contentType),
	
	})

	if err != nil {
		return "", fmt.Errorf("failed to upload file to S3: %w", err)
	}

	// Return the URL to the uploaded file
	return s.GetFileURL(s3Path), nil
}

// DeleteFile deletes a file from S3
func (s *S3FileService) DeleteFile(ctx context.Context, fileURL string) error {
	// Extract the key from the URL
	key, err := s.extractKeyFromURL(fileURL)
	if err != nil {
		return err
	}

	// Delete the file from S3
	_, err = s.client.DeleteObject(ctx, &s3.DeleteObjectInput{
		Bucket: aws.String(s.bucketName),
		Key:    aws.String(key),
	})

	if err != nil {
		return fmt.Errorf("failed to delete file from S3: %w", err)
	}

	return nil
}

// // GetFileURL returns the full URL for a file in S3
// func (s *S3FileService) GetFileURL(fileName string) string {
// 	return fmt.Sprintf("%s/%s", s.baseURL, fileName)
// }

// extractKeyFromURL extracts the S3 key from a file URL
func (s *S3FileService) extractKeyFromURL(fileURL string) (string, error) {
	if !strings.HasPrefix(fileURL, s.baseURL) {
		return "", fmt.Errorf("invalid file URL: %s", fileURL)
	}

	// Remove base URL to get the key
	key := strings.TrimPrefix(fileURL, s.baseURL+"/")
	return key, nil
}


func (s *S3FileService) GetSignedURL(ctx context.Context, objectKey string, expiration time.Duration) (string, error) {
    presignClient := s3.NewPresignClient(s.client)
    
    presignedURL, err := presignClient.PresignGetObject(ctx, &s3.GetObjectInput{
        Bucket: aws.String(s.bucketName),
        Key:    aws.String(objectKey),
    }, func(opts *s3.PresignOptions) {
        opts.Expires = expiration
    })
    
    if err != nil {
        return "", fmt.Errorf("failed to generate signed URL: %w", err)
    }
    
    return presignedURL.URL, nil
}

func (s *S3FileService) GetFileURL(fileName string) string {
    // Generate a signed URL valid for 1 hour
    signedURL, err := s.GetSignedURL(context.Background(), fileName, time.Hour)
    if err != nil {
        // Fall back to direct URL if signing fails, but this won't work with private buckets
        return fmt.Sprintf("%s/%s", s.baseURL, fileName)
    }
    return signedURL
}