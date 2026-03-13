package mikrotik

import (
	"github.com/gin-gonic/gin"
	mikrotikhdl "mikmongo/internal/handler/mikrotik"
)

func registerPPPRoutes(parent *gin.RouterGroup, h *mikrotikhdl.Registry) {
	ppp := parent.Group("/ppp")
	{
		// Secrets
		ppp.GET("/secrets", h.PPP.ListSecrets)
		ppp.POST("/secrets", h.PPP.CreateSecret)
		ppp.GET("/secrets/:id", h.PPP.GetSecret)
		ppp.GET("/secrets/by-name", h.PPP.GetSecretByName)
		ppp.PUT("/secrets/:id", h.PPP.UpdateSecret)
		ppp.DELETE("/secrets/:id", h.PPP.DeleteSecret)
		ppp.POST("/secrets/:id/enable", h.PPP.EnableSecret)
		ppp.POST("/secrets/:id/disable", h.PPP.DisableSecret)
		ppp.POST("/secrets/batch-delete", h.PPP.BatchDeleteSecrets)
		ppp.POST("/secrets/batch-enable", h.PPP.BatchEnableSecrets)
		ppp.POST("/secrets/batch-disable", h.PPP.BatchDisableSecrets)

		// Profiles
		ppp.GET("/profiles", h.PPP.ListProfiles)
		ppp.POST("/profiles", h.PPP.CreateProfile)
		ppp.GET("/profiles/:id", h.PPP.GetProfile)
		ppp.GET("/profiles/by-name", h.PPP.GetProfileByName)
		ppp.PUT("/profiles/:id", h.PPP.UpdateProfile)
		ppp.DELETE("/profiles/:id", h.PPP.DeleteProfile)
		ppp.POST("/profiles/:id/enable", h.PPP.EnableProfile)
		ppp.POST("/profiles/:id/disable", h.PPP.DisableProfile)
		ppp.POST("/profiles/batch-delete", h.PPP.BatchDeleteProfiles)
		ppp.POST("/profiles/batch-enable", h.PPP.BatchEnableProfiles)
		ppp.POST("/profiles/batch-disable", h.PPP.BatchDisableProfiles)

		// Active sessions
		ppp.GET("/active", h.PPP.ListActive)
		ppp.GET("/active/:id", h.PPP.GetActive)
		ppp.DELETE("/active/:id", h.PPP.DisconnectActive)
		ppp.POST("/active/batch-disconnect", h.PPP.BatchDisconnectActive)

		// WebSocket streams
		ppp.GET("/active/stream", h.WebSocket.StreamPPPActive)
		ppp.GET("/inactive/stream", h.WebSocket.StreamPPPInactive)
	}
}
