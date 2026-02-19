package middlewares

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func (m *Middleware) RequirePermission(permission string) gin.HandlerFunc {
	return func(c *gin.Context) {
		if permission == "" {
			c.Next()
			return
		}

		rawPermissions, exists := c.Get("permissions")
		if !exists {
			c.JSON(http.StatusForbidden, gin.H{
				"success": false,
				"error":   "permissions are missing in context",
			})
			c.Abort()
			return
		}

		userPermissions, ok := rawPermissions.([]string)
		if !ok {
			c.JSON(http.StatusForbidden, gin.H{
				"success": false,
				"error":   "invalid permissions format in context",
			})
			c.Abort()
			return
		}

		permissionSet := make(map[string]struct{}, len(userPermissions))
		for _, permission := range userPermissions {
			permissionSet[permission] = struct{}{}
		}

		if _, hasPermission := permissionSet[permission]; !hasPermission {
			c.JSON(http.StatusForbidden, gin.H{
				"success": false,
				"error":   "insufficient permissions",
			})
			c.Abort()
			return
		}

		c.Next()
	}
}
