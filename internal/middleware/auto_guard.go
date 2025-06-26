package middleware

import (
	"log"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

type GuardOption func(*GuardConfig)

func WithOwnership(field, param string) GuardOption {
	return func(c *GuardConfig) {
		c.RequireOwnership = true
		c.OwnershipField = field
		c.OwnershipIDParam = param
	}
}

func WithCustomAction(action string) GuardOption {
	return func(c *GuardConfig) {
		c.RequiredAction = action
	}
}

func WithCustomResource(resource string) GuardOption {
	return func(c *GuardConfig) {
		c.RequiredResource = resource
	}
}

func AutoGuard(rbac RBACService, opts ...GuardOption) gin.HandlerFunc {
	return func(c *gin.Context) {
		path := c.FullPath()       // e.g., /locations/:id
		method := c.Request.Method // e.g., GET

		log.Printf("🛡️ AutoGuard - Path: %s, Method: %s", path, method)

		// Check if rbac service is nil
		if rbac == nil {
			log.Printf("❌ AutoGuard - RBAC service is nil!")
			c.JSON(http.StatusInternalServerError, gin.H{"error": "RBAC service not available"})
			c.Abort()
			return
		}

		resource := inferResourceFromPath(path)
		action := inferActionFromMethod(method, path)

		log.Printf("🛡️ AutoGuard - Inferred Resource: %s, Action: %s", resource, action)

		config := GuardConfig{
			RequireAuth:      true,
			RequiredResource: resource,
			RequiredAction:   action,
		}

		for _, opt := range opts {
			opt(&config)
		}

		log.Printf("🛡️ AutoGuard - Final Config: %+v", config)

		Middleware := MultiLayerGuard(rbac, config)
		Middleware(c)
	}
}

func inferResourceFromPath(path string) string {
	log.Printf("🔍 Inferring resource from path: %s", path)

	// Extract base resource from path: "/locations/:id" -> "locations"
	if path == "" || path == "/" {
		log.Printf("🔍 Empty or root path, returning 'unknown'")
		return "unknown"
	}

	parts := strings.Split(path, "/")
	log.Printf("🔍 Path parts: %v", parts)

	// Skip common prefixes like "api", "v1", "v2", etc.
	skipPrefixes := map[string]bool{
		"api": true,
		"v1":  true,
		"v2":  true,
		"v3":  true,
	}

	for _, part := range parts {
		if part != "" && !strings.HasPrefix(part, ":") && !skipPrefixes[part] {
			resource := part
			if !strings.HasSuffix(part, "s") {
				resource = part + "s" // make it plural
			}
			log.Printf("🔍 Found resource: %s", resource)
			return resource
		}
	}

	log.Printf("🔍 No valid resource found, returning 'unknown'")
	return "unknown"
}

func inferActionFromMethod(method string, path string) string {
	log.Printf("🔍 Inferring action from method: %s, path: %s", method, path)

	var action string
	switch method {
	case http.MethodGet:
		if strings.Contains(path, ":") {
			action = "read"
		} else {
			action = "list"
		}
	case http.MethodPost:
		if strings.Contains(path, "restore") {
			action = "restore"
		} else if strings.Contains(path, "media") {
			action = "upload"
		} else {
			action = "create"
		}
	case http.MethodPut:
		action = "update"
	case http.MethodDelete:
		if strings.Contains(path, "hard-delete") {
			action = "hard_delete"
		} else {
			action = "delete"
		}
	default:
		action = "unknown"
	}

	log.Printf("🔍 Inferred action: %s", action)
	return action
}
