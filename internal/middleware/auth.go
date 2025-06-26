package middleware

import (
	"net/http"

	"bitbucket.org/abhishek_fordel/we-care-holidays-backend-golang/configs"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// AuthMiddleware - Handles JWT validation and sets auth context
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

		// Get user permissions ONCE
		permissions, roleScope, err := rbacService.GetUserPermissions(ctx, userID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get user permissions"})
			c.Abort()
			return
		}

		// Create and set auth context
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
			}
		}

		ctx = SetAuthContext(ctx, authCtx)
		c.Request = c.Request.WithContext(ctx)
		c.Next()
	}
}
