package mikrotik

import (
	"github.com/gin-gonic/gin"
	mikrotikhdl "mikmongo/internal/handler/mikrotik"
	"mikmongo/internal/middleware"
)

func registerScriptRoutes(v1 *gin.RouterGroup, h *mikrotikhdl.Registry, mw *middleware.Registry) {
	scripts := v1.Group("/mikrotik/scripts")
	scripts.Use(mw.Auth.Authenticate())
	{
		scripts.POST("/on-login", h.Script.GenerateOnLogin)
		scripts.POST("/on-login/parse", h.Script.ParseOnLogin)
		scripts.POST("/expired-action", h.Script.GenerateExpiredAction)
		scripts.GET("/expire-monitor", h.Script.GenerateExpireMonitor)
		scripts.GET("/expire-modes", h.Script.GetExpireModes)
	}
}
