package middleware

import (
	casbincore "github.com/casbin/casbin/v3"
	"github.com/gin-contrib/requestid"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"mikmongo/pkg/jwt"
	"mikmongo/pkg/redis"
)

// Registry holds all middleware instances
type Registry struct {
	Auth           *AuthMiddleware
	Logger         gin.HandlerFunc
	RateLimit      *RateLimitMiddleware
	RequestID      gin.HandlerFunc
	CORS           gin.HandlerFunc
	PortalAuth     *PortalAuthMiddleware
	AgentPortalAuth *AgentPortalAuthMiddleware
	MikrotikRouter  *MikrotikRouterMiddleware
	RBAC            gin.HandlerFunc
}

// NewRegistry creates a new middleware registry
func NewRegistry(logger *zap.Logger, jwtService *jwt.Service, redisClient *redis.Client, enforcer *casbincore.Enforcer, allowedOrigins []string) *Registry {
	loggerMiddleware := NewLoggerMiddleware(logger)

	return &Registry{
		Auth:            NewAuthMiddleware(jwtService, redisClient),
		Logger:          loggerMiddleware.GinLogger(),
		RateLimit:       NewRateLimitMiddleware(redisClient),
		RequestID:       requestid.New(),
		CORS:            NewCORSMiddleware(allowedOrigins),
		PortalAuth:      NewPortalAuthMiddleware(jwtService),
		AgentPortalAuth: NewAgentPortalAuthMiddleware(jwtService),
		MikrotikRouter:  NewMikrotikRouterMiddleware(),
		RBAC:            CasbinMiddleware(enforcer),
	}
}
