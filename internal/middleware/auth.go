package middleware

import (
	"strings"

	"github.com/gin-gonic/gin"
	"mikmongo/pkg/jwt"
	"mikmongo/pkg/redis"
	"mikmongo/pkg/response"
)

// AuthMiddleware validates JWT tokens
type AuthMiddleware struct {
	jwtService  *jwt.Service
	redisClient *redis.Client
	internalKey string
}

// NewAuthMiddleware creates a new auth middleware
func NewAuthMiddleware(jwtService *jwt.Service, redisClient *redis.Client, internalKey string) *AuthMiddleware {
	return &AuthMiddleware{
		jwtService:  jwtService,
		redisClient: redisClient,
		internalKey: internalKey,
	}
}

// Authenticate validates the JWT token from Authorization header or query param.
// For WebSocket endpoints, browsers cannot set custom headers, so the token
// can be passed as ?token=<jwt> query parameter.
func (m *AuthMiddleware) Authenticate() gin.HandlerFunc {
	return func(c *gin.Context) {
		var token string

		authHeader := c.GetHeader("Authorization")
		if authHeader != "" {
			parts := strings.SplitN(authHeader, " ", 2)
			if len(parts) != 2 || strings.ToLower(parts[0]) != "bearer" {
				response.Unauthorized(c, "invalid authorization header format")
				c.Abort()
				return
			}
			token = parts[1]
		} else if qToken := c.Query("token"); qToken != "" {
			// Allow token via query param (WebSocket support)
			token = qToken
		} else {
			response.Unauthorized(c, "missing authorization header")
			c.Abort()
			return
		}

		// Internal key bypass — for WebSocket endpoints that don't carry user JWT
		if m.internalKey != "" && token == m.internalKey {
			c.Set("user_id", "internal")
			c.Set("email", "internal@system")
			c.Set("role", "superadmin")
			c.Next()
			return
		}

		claims, err := m.jwtService.Validate(token)
		if err != nil {
			response.Unauthorized(c, "invalid token")
			c.Abort()
			return
		}

		// Reject refresh tokens on regular endpoints
		if claims.TokenType == "refresh" {
			response.Unauthorized(c, "refresh token not allowed")
			c.Abort()
			return
		}

		// Check if token is blacklisted
		if claims.ID != "" {
			blacklisted, err := m.redisClient.IsBlacklisted(c.Request.Context(), claims.ID)
			if err == nil && blacklisted {
				response.Unauthorized(c, "token has been revoked")
				c.Abort()
				return
			}
		}

		// Check if password was changed after token was issued
		if claims.IssuedAt != nil {
			pwdChanged, err := m.redisClient.GetPasswordChangedAt(c.Request.Context(), claims.UserID)
			if err == nil && !pwdChanged.IsZero() && claims.IssuedAt.Time.Before(pwdChanged) {
				response.Unauthorized(c, "password has been changed, please login again")
				c.Abort()
				return
			}
		}

		c.Set("user_id", claims.UserID)
		c.Set("email", claims.Email)
		c.Set("role", claims.Role)
		c.Next()
	}
}

// Authorize checks if user has required role
func (m *AuthMiddleware) Authorize(roles ...string) gin.HandlerFunc {
	return func(c *gin.Context) {
		userRole, exists := c.Get("role")
		if !exists {
			response.Forbidden(c, "role not found")
			c.Abort()
			return
		}

		for _, role := range roles {
			if userRole == role {
				c.Next()
				return
			}
		}

		response.Forbidden(c, "insufficient permissions")
		c.Abort()
	}
}
