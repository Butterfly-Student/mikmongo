package middleware

import (
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

// NewCORSMiddleware creates a CORS middleware with appropriate settings.
// origins should come from config (e.g., ALLOWED_ORIGINS env var).
func NewCORSMiddleware(origins []string) gin.HandlerFunc {
	allowAll := false
	for _, o := range origins {
		if o == "*" {
			allowAll = true
			break
		}
	}

	cfg := cors.Config{
		AllowMethods:  []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowHeaders:  []string{"Origin", "Content-Type", "Accept", "Authorization", "X-Request-ID"},
		ExposeHeaders: []string{"Content-Length"},
		MaxAge:        12 * time.Hour,
	}

	if allowAll {
		cfg.AllowAllOrigins = true
	} else {
		cfg.AllowOrigins = origins
		cfg.AllowCredentials = true
	}

	return cors.New(cfg)
}
