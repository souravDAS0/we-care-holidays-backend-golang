package services

import (
	"context"
	"io"
)

// FileService defines the interface for file operations
type FileService interface {
	// UploadFile uploads a file and returns the URL
	UploadFile(ctx context.Context, file io.Reader, fileName string, contentType string) (string, error)
	
	// DeleteFile deletes a file by URL
	DeleteFile(ctx context.Context, fileURL string) error
	
	// GetFileURL returns the full URL for a file path
	GetFileURL(fileName string) string
}
