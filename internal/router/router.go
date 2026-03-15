package router

import (
	"github.com/gin-gonic/gin"
	"mikmongo/internal/handler"
	"mikmongo/internal/middleware"
)

// New creates a new Gin router with all routes configured
func New(handlers *handler.Registry, mw *middleware.Registry) *gin.Engine {
	r := gin.New()

	r.Use(gin.Recovery())
	r.Use(mw.CORS)
	r.Use(mw.Logger)
	r.Use(mw.RequestID)

	// Public routes (health, login, registration, webhooks, portal)
	registerPublicRoutes(r, handlers, mw)

	// Admin API (JWT required)
	v1 := r.Group("/api/v1")
	v1.Use(mw.Auth.Authenticate())
	{
		registerAdminRoutes(v1, handlers)
	}

	return r
}
