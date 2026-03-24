package middleware

import (
	"strings"

	"github.com/gin-gonic/gin"
	"mikmongo/pkg/jwt"
	"mikmongo/pkg/response"
)

// AgentPortalAuthMiddleware validates agent portal JWT tokens
type AgentPortalAuthMiddleware struct {
	jwtService *jwt.Service
}

// NewAgentPortalAuthMiddleware creates a new agent portal auth middleware
func NewAgentPortalAuthMiddleware(jwtService *jwt.Service) *AgentPortalAuthMiddleware {
	return &AgentPortalAuthMiddleware{jwtService: jwtService}
}

// AuthenticateAgentPortal validates agent portal JWT token
func (m *AgentPortalAuthMiddleware) AuthenticateAgentPortal() gin.HandlerFunc {
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

		if claims.Role != "agent_portal" {
			response.Forbidden(c, "not an agent portal token")
			c.Abort()
			return
		}

		c.Set("agent_id", claims.UserID)
		c.Next()
	}
}
