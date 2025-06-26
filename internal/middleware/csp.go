package middleware

import (
	"github.com/gin-gonic/gin"
)

// ContentSecurityPolicy adds CSP headers to allow Swagger UI to function properly
func ContentSecurityPolicy() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Header("Content-Security-Policy", "default-src 'self'; img-src 'self' data:; style-src 'self' 'unsafe-inline'; script-src 'self' 'unsafe-inline' 'unsafe-eval'; connect-src 'self'")
		c.Next()
	}
}