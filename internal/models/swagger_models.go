package models

// swagger_models.go
// This file contains models used only for Swagger documentation purposes.
// These types are not used in the actual code, but provide documentation
// for the standardized response format used by the API.

// StandardResponse models for Swagger documentation

// SwaggerStandardResponse represents the standard API response format
// @Description Standard API response format
type SwaggerStandardResponse struct {
	// Indicates whether the request was successful
	Success bool `json:"success" example:"true"`
	// HTTP status code
	StatusCode int `json:"statusCode" example:"200"`
	// Error message (only present if success is false)
	ErrorMessage string `json:"errorMessage,omitempty" example:""`
	// Response payload
	Data interface{} `json:"data,omitempty"`
}

// SwaggerErrorResponse documents the error response format
// @Description Error response model
type SwaggerErrorResponse struct {
	// Indicates request failure
	Success bool `json:"success" example:"false"`
	// HTTP status code
	StatusCode int `json:"statusCode" example:"400"`
	// Human-readable error message
	ErrorMessage string `json:"errorMessage" example:"Invalid request parameters"`
	// Optional error details
	Data interface{} `json:"data,omitempty"`
}

// SwaggerPaginatedResponse documents the paginated response format
// @Description Success response with pagination
type SwaggerPaginatedResponse struct {
	// Indicates successful request
	Success bool `json:"success" example:"true"`
	// HTTP status code
	StatusCode int `json:"statusCode" example:"200"`
	// Response data with pagination information
	Data struct {
		// Array of items
		Items []interface{} `json:"items"`
		// Current page number
		Page int `json:"page" example:"1"`
		// Number of items per page
		Limit int `json:"limit" example:"10"`
		// Total number of items
		Total int64 `json:"total" example:"42"`
		// Total number of pages
		TotalPages int `json:"total_pages" example:"5"`
	} `json:"data"`
}
