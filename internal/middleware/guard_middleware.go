package middleware

import (
	"log"
	"net/http"
	"strings"

	"bitbucket.org/abhishek_fordel/we-care-holidays-backend-golang/configs"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type GuardConfig struct {
	RequireAuth      bool
	RequiredResource string
	RequiredAction   string
	RequireOwnership bool
	OwnershipField   string // e.g., "userId", "organizationId"
	OwnershipIDParam string // e.g., "id", "userId"
}

func MultiLayerGuard(rbacService RBACService, config GuardConfig) gin.HandlerFunc {
	log.Printf("🛡️ MultiLayerGuard - Creating middleware with config: %+v", config)

	jwtValidator := NewJWTValidator(configs.GetEnv("JWT_SECRET", ""))

	if jwtValidator == nil {
		log.Printf("❌ MultiLayerGuard - JWT Validator is nil!")
	}

	if rbacService == nil {
		log.Printf("❌ MultiLayerGuard - RBAC Service is nil!")
	}

	return func(c *gin.Context) {
		log.Printf("🛡️ MultiLayerGuard - Processing request: %s %s", c.Request.Method, c.Request.URL.Path)

		ctx := c.Request.Context()

		// 1. Token Guard
		if config.RequireAuth {
			log.Printf("🔐 MultiLayerGuard - Authentication required")

			token := extractToken(c)
			if token == "" {
				log.Printf("❌ MultiLayerGuard - No token found")
				c.JSON(http.StatusUnauthorized, gin.H{"error": "Missing authentication token"})
				c.Abort()
				return
			}

			log.Printf("🔐 MultiLayerGuard - Token found, validating...")

			claims, err := jwtValidator.ValidateToken(token)
			if err != nil {
				log.Printf("❌ MultiLayerGuard - Token validation failed: %v", err)
				c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid or expired token"})
				c.Abort()
				return
			}

			log.Printf("✅ MultiLayerGuard - Token valid for user: %s", claims.UserID)

			userID, err := primitive.ObjectIDFromHex(claims.UserID)
			if err != nil {
				log.Printf("❌ MultiLayerGuard - Invalid user ID in token: %v", err)
				c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid user ID in token"})
				c.Abort()
				return
			}

			log.Printf("🔐 MultiLayerGuard - Getting user permissions for user: %s", userID.Hex())

			// Check if rbacService is nil before calling
			if rbacService == nil {
				log.Printf("❌ MultiLayerGuard - RBAC Service is nil when trying to get permissions!")
				c.JSON(http.StatusInternalServerError, gin.H{"error": "RBAC service not available"})
				c.Abort()
				return
			}

			permissions, roleScope, err := rbacService.GetUserPermissions(ctx, userID)
			if err != nil {
				log.Printf("❌ MultiLayerGuard - Failed to get user permissions: %v", err)
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get user permissions"})
				c.Abort()
				return
			}

			log.Printf("✅ MultiLayerGuard - Got permissions: %v, roleScope: %s", permissions, roleScope)

			authCtx := &AuthContext{
				UserID:      userID,
				Role:        claims.Role,
				RoleScope:   roleScope,
				Permissions: permissions,
				Token:       token,
			}

			if claims.OrganizationID != "" {
				orgID, err := primitive.ObjectIDFromHex(claims.OrganizationID)
				if err == nil {
					authCtx.OrganizationID = &orgID
					log.Printf("🏢 MultiLayerGuard - Organization ID set: %s", orgID.Hex())
				} else {
					log.Printf("⚠️ MultiLayerGuard - Invalid organization ID in token: %v", err)
				}
			}

			ctx = SetAuthContext(ctx, authCtx)
			c.Request = c.Request.WithContext(ctx)
			log.Printf("✅ MultiLayerGuard - Auth context set successfully")
		}

		authCtx := GetAuthContext(ctx)
		if authCtx == nil && config.RequireAuth {
			log.Printf("❌ MultiLayerGuard - Auth context is nil but authentication is required")
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Authentication required"})
			c.Abort()
			return
		}

		// 2. Permission Guard
		if config.RequiredResource != "" && config.RequiredAction != "" {
			log.Printf("🔒 MultiLayerGuard - Checking permission: %s:%s", config.RequiredResource, config.RequiredAction)

			if rbacService == nil {
				log.Printf("❌ MultiLayerGuard - RBAC Service is nil when validating permission!")
				c.JSON(http.StatusInternalServerError, gin.H{"error": "RBAC service not available"})
				c.Abort()
				return
			}

			if !rbacService.ValidatePermission(ctx, authCtx, config.RequiredResource, config.RequiredAction) {
				log.Printf("❌ MultiLayerGuard - Permission denied: %s:%s", config.RequiredResource, config.RequiredAction)
				c.JSON(http.StatusForbidden, gin.H{"error": "Insufficient permissions"})
				c.Abort()
				return
			}

			log.Printf("✅ MultiLayerGuard - Permission granted: %s:%s", config.RequiredResource, config.RequiredAction)
		}

		// 3. Ownership Guard
		if config.RequireOwnership && authCtx.Role != "PLATFORM_ADMIN" {
			log.Printf("👤 MultiLayerGuard - Checking ownership: field=%s, param=%s", config.OwnershipField, config.OwnershipIDParam)

			if !validateOwnership(c, authCtx, config) {
				log.Printf("❌ MultiLayerGuard - Ownership validation failed")
				c.JSON(http.StatusForbidden, gin.H{"error": "Access denied"})
				c.Abort()
				return
			}

			log.Printf("✅ MultiLayerGuard - Ownership validated successfully")
		}

		log.Printf("✅ MultiLayerGuard - All guards passed, proceeding to handler")
		c.Next()
	}
}

func extractToken(c *gin.Context) string {
	bearerToken := c.GetHeader("Authorization")
	log.Printf("🔐 extractToken - Authorization header: %s", bearerToken)

	parts := strings.Split(bearerToken, " ")
	if len(parts) == 2 && strings.ToLower(parts[0]) == "bearer" {
		log.Printf("🔐 extractToken - Token extracted successfully")
		return parts[1]
	}

	log.Printf("❌ extractToken - Invalid authorization header format")
	return ""
}

func validateOwnership(c *gin.Context, authCtx *AuthContext, config GuardConfig) bool {
	resourceID := c.Param(config.OwnershipIDParam)
	log.Printf("👤 validateOwnership - ResourceID: %s, Field: %s", resourceID, config.OwnershipField)

	if resourceID == "" {
		log.Printf("❌ validateOwnership - Resource ID is empty")
		return false
	}

	switch config.OwnershipField {
	case "userId":
		result := resourceID == authCtx.UserID.Hex()
		log.Printf("👤 validateOwnership - UserId check: %s == %s = %t", resourceID, authCtx.UserID.Hex(), result)
		return result
	case "organizationId":
		if authCtx.OrganizationID == nil {
			log.Printf("❌ validateOwnership - User has no organization ID")
			return false
		}
		result := resourceID == authCtx.OrganizationID.Hex()
		log.Printf("👤 validateOwnership - OrganizationId check: %s == %s = %t", resourceID, authCtx.OrganizationID.Hex(), result)
		return result
	}

	log.Printf("❌ validateOwnership - Unknown ownership field: %s", config.OwnershipField)
	return false
}
