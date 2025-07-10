package middleware

import (
	"net/http"

	"bitbucket.org/abhishek_fordel/we-care-holidays-backend-golang/configs"
	"bitbucket.org/abhishek_fordel/we-care-holidays-backend-golang/internal/logger"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.uber.org/zap"
)

// AuthMiddleware - Enhanced version that works with ScopedRBACMiddleware
func AuthMiddleware(rbacService RBACService, cfg *configs.Config) gin.HandlerFunc {
	jwtValidator := NewJWTValidator(cfg.JWTSecret)

	return func(c *gin.Context) {
		ctx := c.Request.Context()

		// Extract and validate token
		token := extractToken(c)
		if token == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Missing authentication token"})
			c.Abort()
			return
		}

		claims, err := jwtValidator.ValidateToken(token)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid or expired token"})
			c.Abort()
			return
		}

		userID, err := primitive.ObjectIDFromHex(claims.UserID)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid user ID in token"})
			c.Abort()
			return
		}

		// Get user permissions and role information from RBAC service
		permissions, roleScope, err := rbacService.GetUserPermissions(ctx, userID)
		if err != nil {
			logger.Log.Error("Failed to get user permissions", 
				zap.String("user_id", claims.UserID), 
				zap.Error(err))
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get user permissions"})
			c.Abort()
			return
		}

		// Create structured auth context for your existing middleware
		authCtx := &AuthContext{
			UserID:      userID,
			Role:        claims.Role,
			RoleScope:   roleScope,
			Permissions: permissions,
			Token:       token,
		}

		var organizationID string
		if claims.OrganizationID != "" {
			orgID, err := primitive.ObjectIDFromHex(claims.OrganizationID)
			if err == nil {
				authCtx.OrganizationID = &orgID
				organizationID = claims.OrganizationID
			}
		}

		// Set structured auth context in request context (for your existing code)
		ctx = SetAuthContext(ctx, authCtx)
		c.Request = c.Request.WithContext(ctx)

		// Convert permissions to string slice for ScopedRBACMiddleware
		permissionStrings := make([]string, len(permissions))
		for i, perm := range permissions {
			permissionStrings[i] = perm.Resource + ":" + perm.Action
		}

		// Set user data in gin context for ScopedRBACMiddleware compatibility
		userData := map[string]interface{}{
			"user_id":         claims.UserID,
			"role":            claims.Role,
			"role_scope":      string(roleScope), // Convert to string
			"organization_id": organizationID,
			"permissions":     permissionStrings,
		}

		c.Set("user", userData)

		logger.Log.Debug("User authenticated successfully",
			zap.String("user_id", claims.UserID),
			zap.String("role", claims.Role),
			zap.String("role_scope", string(roleScope)),
			zap.String("organization_id", organizationID),
			zap.Int("permission_count", len(permissions)))

		c.Next()
	}
}