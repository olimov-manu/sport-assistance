package middlewares

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

func (m *Middleware) RequirePermissions(requiredPermissions ...string) gin.HandlerFunc {
	return func(c *gin.Context) {
		filteredRequired := make([]string, 0, len(requiredPermissions))
		for _, permission := range requiredPermissions {
			permission = strings.TrimSpace(permission)
			if permission != "" {
				filteredRequired = append(filteredRequired, permission)
			}
		}

		if len(filteredRequired) == 0 {
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

		missingPermissions := make([]string, 0)
		for _, permission := range filteredRequired {
			if _, hasPermission := permissionSet[permission]; !hasPermission {
				missingPermissions = append(missingPermissions, permission)
			}
		}

		if len(missingPermissions) > 0 {
			c.JSON(http.StatusForbidden, gin.H{
				"success": false,
				"error":   "insufficient permissions",
				"missing": missingPermissions,
			})
			c.Abort()
			return
		}

		c.Next()
	}
}
