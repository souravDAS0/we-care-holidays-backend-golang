package middleware

import (
	"bytes"
	"encoding/json"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

// StandardResponse defines the standard response structure
type StandardResponse struct {
	Success      bool        `json:"success"`
	StatusCode   int         `json:"statusCode"`
	ErrorMessage string      `json:"errorMessage,omitempty"`
	Data         interface{} `json:"data,omitempty"`
}

// ResponseInterceptor intercepts all responses and formats them
// to match the standard response structure
func ResponseInterceptor() gin.HandlerFunc {
	return func(c *gin.Context) {

		if strings.HasPrefix(c.Request.URL.Path, "/swagger") {
			c.Next()
			return
		}
		
		// Replace the ResponseWriter with our custom writer
		// that captures the response but doesn't actually write it
		originalWriter := c.Writer
		bodyCapture := &responseCapture{
			ResponseWriter: originalWriter,
			body:           new(bytes.Buffer),
		}
		c.Writer = bodyCapture

		// Process the request
		c.Next()

		// Get status code and headers
		statusCode := bodyCapture.status
		if statusCode == 0 {
			statusCode = http.StatusOK
		}

		// Check if this is a success response based on status code
		success := statusCode >= 200 && statusCode < 300

		// Get error message if any
		var errorMessage string
		if len(c.Errors) > 0 {
			errorMessage = c.Errors.String()
		}

		// Parse the response body as JSON if possible
		var responseData interface{}
		if bodyCapture.body.Len() > 0 {
			// Check if it's already our standard format
			var stdResponse StandardResponse
			if err := json.Unmarshal(bodyCapture.body.Bytes(), &stdResponse); err == nil {
				if stdResponse.StatusCode > 0 {
					// It's already our format, use it directly
					originalWriter.WriteHeader(statusCode)
					originalWriter.Write(bodyCapture.body.Bytes())
					return
				}
			}

			// Otherwise parse as generic JSON
			if err := json.Unmarshal(bodyCapture.body.Bytes(), &responseData); err != nil {
				// If not valid JSON, use as string
				responseData = bodyCapture.body.String()
			} else {
				// Check for error message in the response data
				if m, ok := responseData.(map[string]interface{}); ok {
					if errMsg, exists := m["error"]; exists && errorMessage == "" {
						if errStr, ok := errMsg.(string); ok {
							errorMessage = errStr
							delete(m, "error")
							if len(m) == 0 {
								responseData = nil
							}
						}
					}
					
					// Check for errorMessage field in response data
					if errMsg, exists := m["errorMessage"]; exists && errorMessage == "" {
						if errStr, ok := errMsg.(string); ok {
							errorMessage = errStr
							delete(m, "errorMessage")
							if len(m) == 0 {
								responseData = nil
							}
						}
					}
				}
			}
		}

		// If there's an error message, ensure success is false
		if errorMessage != "" {
			success = false
			
			// If we have a success status code but an error message, 
			// update the status code to an error code if it's not already
			if statusCode >= 200 && statusCode < 300 {
				statusCode = http.StatusInternalServerError
			}
		}

		// Create the standardized response
		response := StandardResponse{
			Success:      success,
			StatusCode:   statusCode,
			ErrorMessage: errorMessage,
			Data:         responseData,
		}

		// Write to the original writer
		for k, values := range bodyCapture.Header() {
			for _, v := range values {
				originalWriter.Header().Add(k, v)
			}
		}
		originalWriter.Header().Set("Content-Type", "application/json")
		jsonResponse, _ := json.Marshal(response)
		originalWriter.WriteHeader(statusCode)
		originalWriter.Write(jsonResponse)
	}
}

// responseCapture is a ResponseWriter that captures the response
type responseCapture struct {
	gin.ResponseWriter
	body   *bytes.Buffer
	status int
}

// Write captures the response but doesn't write it
func (w *responseCapture) Write(b []byte) (int, error) {
	return w.body.Write(b)
}

// WriteString captures the response but doesn't write it
func (w *responseCapture) WriteString(s string) (int, error) {
	return w.body.WriteString(s)
}

// WriteHeader captures the status code but doesn't write it
func (w *responseCapture) WriteHeader(code int) {
	w.status = code
}

// Status returns the captured status code
func (w *responseCapture) Status() int {
	return w.status
}