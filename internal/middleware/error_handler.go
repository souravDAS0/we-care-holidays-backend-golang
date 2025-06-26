package middleware

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

// ErrorCode represents standard application error codes
type ErrorCode string

// Standard error codes
const (
	ErrorCodeInvalidRequest     ErrorCode = "INVALID_REQUEST"
	ErrorCodeUnauthorized       ErrorCode = "UNAUTHORIZED"
	ErrorCodeForbidden          ErrorCode = "FORBIDDEN"
	ErrorCodeNotFound           ErrorCode = "NOT_FOUND"
	ErrorCodeConflict           ErrorCode = "CONFLICT"
	ErrorCodeInternalServer     ErrorCode = "INTERNAL_SERVER_ERROR"
	ErrorCodeServiceUnavailable ErrorCode = "SERVICE_UNAVAILABLE"
	ErrorCodeValidationFailed   ErrorCode = "VALIDATION_FAILED"
)

// AppError represents a standard application error
type AppError struct {
	Code    ErrorCode
	Message string
	Err     error
	Status  int
}

// Error implements the error interface
func (e *AppError) Error() string {
	if e.Err != nil {
		return fmt.Sprintf("%s: %v", e.Message, e.Err)
	}
	return e.Message
}

// ErrorWithContext extends AppError with a request context
type ErrorWithContext struct {
	*AppError
	Context gin.H
}

// NewAppError creates a new application error
func NewAppError(code ErrorCode, message string, err error, status int) *AppError {
	return &AppError{
		Code:    code,
		Message: message,
		Err:     err,
		Status:  status,
	}
}

// CommonErrors maps common error types to standardized AppErrors
var CommonErrors = map[string]*AppError{
	"mongo: no documents in result": NewAppError(
		ErrorCodeNotFound,
		"Resource not found",
		nil,
		http.StatusNotFound,
	),
	"context canceled": NewAppError(
		ErrorCodeServiceUnavailable,
		"Request canceled by client",
		nil,
		http.StatusServiceUnavailable,
	),
}

// ErrorHandler middleware for handling errors in a standardized way
func ErrorHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()

		// If there are errors
		if len(c.Errors) > 0 {
			err := c.Errors.Last().Err
			var appErr *AppError

			// Check if it's already an AppError
			if e, ok := err.(*AppError); ok {
				appErr = e
			} else if e, ok := err.(*ErrorWithContext); ok {
				appErr = e.AppError
			} else {
				// Check if it matches any common errors
				errorMsg := err.Error()
				for key, commonErr := range CommonErrors {
					if strings.Contains(errorMsg, key) {
						appErr = commonErr
						break
					}
				}

				// If not found in common errors, use a generic one
				if appErr == nil {
					appErr = NewAppError(
						ErrorCodeInternalServer,
						"An unexpected error occurred",
						err,
						http.StatusInternalServerError,
					)
				}
			}

			// If there was a custom status in the context, use it
			if statusCode, exists := c.Get("ErrorStatusCode"); exists {
				if code, ok := statusCode.(int); ok {
					appErr.Status = code
				}
			}

			// Create standard response
			response := StandardResponse{
				Success:      false,
				StatusCode:   appErr.Status,
				ErrorMessage: appErr.Message,
			}

			// Abort with the error
			c.AbortWithStatusJSON(appErr.Status, response)
		}
	}
}

// HandleError helper function to handle errors in handlers
func HandleError(c *gin.Context, err error) {
	if err != nil {
		c.Error(err)
	}
}