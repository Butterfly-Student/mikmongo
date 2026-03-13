package middleware

import (
	"strings"

	"github.com/gin-gonic/gin"
	"mikmongo/pkg/jwt"
	"mikmongo/pkg/response"
)

// PortalAuthMiddleware validates portal JWT tokens for customer self-service
type PortalAuthMiddleware struct {
	jwtService *jwt.Service
}

// NewPortalAuthMiddleware creates a new portal auth middleware
func NewPortalAuthMiddleware(jwtService *jwt.Service) *PortalAuthMiddleware {
	return &PortalAuthMiddleware{jwtService: jwtService}
}

// AuthenticatePortal validates portal JWT token
func (m *PortalAuthMiddleware) AuthenticatePortal() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			response.Unauthorized(c, "missing authorization header")
			c.Abort()
			return
		}

		parts := strings.SplitN(authHeader, " ", 2)
		if len(parts) != 2 || strings.ToLower(parts[0]) != "bearer" {
			response.Unauthorized(c, "invalid authorization header format")
			c.Abort()
			return
		}

		claims, err := m.jwtService.Validate(parts[1])
		if err != nil {
			response.Unauthorized(c, "invalid token")
			c.Abort()
			return
		}

		// Portal tokens use role "portal"
		if claims.Role != "portal" {
			response.Forbidden(c, "not a portal token")
			c.Abort()
			return
		}

		c.Set("customer_id", claims.UserID)
		c.Next()
	}
}
