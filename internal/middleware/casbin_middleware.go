package middleware

import (
	"net/http"

	casbincore "github.com/casbin/casbin/v3"
	"github.com/gin-gonic/gin"
)

// CasbinMiddleware returns a Gin handler that enforces RBAC using the provided Casbin enforcer.
// It reads the "role" context key set by AuthMiddleware.
func CasbinMiddleware(enforcer *casbincore.Enforcer) gin.HandlerFunc {
	return func(c *gin.Context) {
		roleVal, exists := c.Get("role")
		if !exists {
			c.JSON(http.StatusForbidden, gin.H{
				"error":   "forbidden",
				"message": "insufficient permissions",
			})
			c.Abort()
			return
		}
		role, _ := roleVal.(string)
		path := c.Request.URL.Path
		method := c.Request.Method

		allowed, err := enforcer.Enforce(role, path, method)
		if err != nil || !allowed {
			c.JSON(http.StatusForbidden, gin.H{
				"error":   "forbidden",
				"message": "insufficient permissions",
			})
			c.Abort()
			return
		}
		c.Next()
	}
}
