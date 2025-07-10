// internal/middleware/scoped_rbac_middleware.go
package middleware

import (
	"errors"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"

	"bitbucket.org/abhishek_fordel/we-care-holidays-backend-golang/internal/constants"
	"bitbucket.org/abhishek_fordel/we-care-holidays-backend-golang/internal/logger"
)

// ScopedRBACContext holds the scoped access information
type ScopedRBACContext struct {
	UserID         string
	RoleID         string
	RoleName       string
	RoleScope      string
	OrganizationID string
	Permissions    []string
	IsGlobalAdmin  bool
}

// ScopedRBACMiddleware creates a middleware for scoped role-based access control
func ScopedRBACMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Get authenticated user context (assumes auth middleware ran before this)
		userContext, exists := c.Get("user")
		if !exists {
			logger.Log.Error("User context not found in scoped RBAC middleware")
			c.JSON(constants.HTTPUnauthorized, gin.H{"error": "Authentication required"})
			c.Abort()
			return
		}

		// Type assertion to get user data (adjust based on your auth context structure)
		userData, ok := userContext.(map[string]interface{})
		if !ok {
			logger.Log.Error("Invalid user context type in scoped RBAC middleware")
			c.JSON(constants.HTTPUnauthorized, gin.H{"error": "Invalid authentication context"})
			c.Abort()
			return
		}

		// Extract user information
		userID, _ := userData["user_id"].(string)
		roleID, _ := userData["role_id"].(string)
		roleName, _ := userData["role"].(string)
		roleScope, _ := userData["role_scope"].(string)
		organizationID, _ := userData["organization_id"].(string)
		permissions, _ := userData["permissions"].([]string)

		// Determine if user is global admin
		isGlobalAdmin := roleScope == "global" || roleName == "PLATFORM_ADMIN"

		// Create scoped RBAC context
		scopedContext := &ScopedRBACContext{
			UserID:         userID,
			RoleID:         roleID,
			RoleName:       roleName,
			RoleScope:      roleScope,
			OrganizationID: organizationID,
			Permissions:    permissions,
			IsGlobalAdmin:  isGlobalAdmin,
		}

		// Add to context for use in handlers
		c.Set("scoped_rbac", scopedContext)

		logger.Log.Debug("Scoped RBAC context set",
			zap.String("user_id", userID),
			zap.String("role", roleName),
			zap.String("scope", roleScope),
			zap.String("org_id", organizationID),
			zap.Bool("is_global_admin", isGlobalAdmin))

		c.Next()
	}
}

// RequireScopedPermission creates a middleware that checks for specific permissions with organization scoping
func RequireScopedPermission(resource, action string) gin.HandlerFunc {
	return func(c *gin.Context) {
		scopedRBAC, exists := c.Get("scoped_rbac")
		if !exists {
			logger.Log.Error("Scoped RBAC context not found")
			c.JSON(constants.HTTPForbidden, gin.H{"error": "Access control context missing"})
			c.Abort()
			return
		}

		rbacContext, ok := scopedRBAC.(*ScopedRBACContext)
		if !ok {
			logger.Log.Error("Invalid scoped RBAC context type")
			c.JSON(constants.HTTPForbidden, gin.H{"error": "Invalid access control context"})
			c.Abort()
			return
		}

		// Check if user has the required permission
		requiredPermission := resource + ":" + action
		if !hasPermission(rbacContext.Permissions, requiredPermission) {
			logger.Log.Warn("Permission denied",
				zap.String("user_id", rbacContext.UserID),
				zap.String("required_permission", requiredPermission),
				zap.Strings("user_permissions", rbacContext.Permissions))
			c.JSON(constants.HTTPForbidden, gin.H{"error": "Insufficient permissions"})
			c.Abort()
			return
		}

		// For organization-scoped roles, add organization filtering
		if rbacContext.RoleScope == "organization" && !rbacContext.IsGlobalAdmin {
			// Add organization filter to query parameters or context
			c.Set("organization_filter", rbacContext.OrganizationID)

			// Log scoped access
			logger.Log.Debug("Organization-scoped access granted",
				zap.String("user_id", rbacContext.UserID),
				zap.String("permission", requiredPermission),
				zap.String("organization_id", rbacContext.OrganizationID))
		}

		c.Next()
	}
}

// RequireOrganizationAccess ensures user can only access their own organization's data
func RequireOrganizationAccess() gin.HandlerFunc {
	return func(c *gin.Context) {
		scopedRBAC, exists := c.Get("scoped_rbac")
		if !exists {
			c.JSON(constants.HTTPForbidden, gin.H{"error": "Access control context missing"})
			c.Abort()
			return
		}

		rbacContext := scopedRBAC.(*ScopedRBACContext)

		// Global admins can access any organization
		if rbacContext.IsGlobalAdmin {
			c.Next()
			return
		}

		// Get organization ID from URL parameter (adjust based on your routing)
		requestedOrgID := c.Param("id")
		if requestedOrgID == "" {
			// Fallback to query parameter or organizationId param
			requestedOrgID = c.Param("organizationId")
			if requestedOrgID == "" {
				requestedOrgID = c.Query("organizationId")
			}
		}

		// If no organization specified in request, allow (will be filtered by middleware)
		if requestedOrgID == "" {
			c.Next()
			return
		}

		// Check if user is trying to access their own organization
		if requestedOrgID != rbacContext.OrganizationID {
			logger.Log.Warn("Organization access denied",
				zap.String("user_id", rbacContext.UserID),
				zap.String("user_org_id", rbacContext.OrganizationID),
				zap.String("requested_org_id", requestedOrgID))
			c.JSON(constants.HTTPForbidden, gin.H{"error": "Access denied to this organization"})
			c.Abort()
			return
		}

		c.Next()
	}
}

// GetScopedRBACContext retrieves the scoped RBAC context from gin context
func GetScopedRBACContext(c *gin.Context) (*ScopedRBACContext, error) {
	scopedRBAC, exists := c.Get("scoped_rbac")
	if !exists {
		return nil, errors.New("scoped RBAC context not found")
	}

	rbacContext, ok := scopedRBAC.(*ScopedRBACContext)
	if !ok {
		return nil, errors.New("invalid scoped RBAC context type")
	}

	return rbacContext, nil
}

// hasPermission checks if user has a specific permission
func hasPermission(userPermissions []string, requiredPermission string) bool {
	for _, perm := range userPermissions {
		if perm == requiredPermission || perm == "*" {
			return true
		}
	}
	return false
}

// ApplyOrganizationFilter applies organization filtering to database queries
func ApplyOrganizationFilter(c *gin.Context, filter map[string]interface{}) map[string]interface{} {
	orgFilter, exists := c.Get("organization_filter")
	if !exists {
		return filter
	}

	organizationID, ok := orgFilter.(string)
	if !ok || organizationID == "" {
		return filter
	}

	// Add organization filter to query
	if filter == nil {
		filter = make(map[string]interface{})
	}
	filter["organizationId"] = organizationID

	return filter
}
