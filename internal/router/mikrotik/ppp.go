package mikrotik

import (
	"github.com/gin-gonic/gin"
	"mikmongo/internal/handler"
)

// RegisterPPPRoutes registers PPP routes scoped to /api/v1/routers/:router_id
func RegisterPPPRoutes(routerGroup *gin.RouterGroup, handlers *handler.Registry) {
	if handlers.Mikrotik == nil {
		return
	}

	ppp := routerGroup.Group("/ppp")
	{
		ppp.GET("/profiles", handlers.Mikrotik.PPP.GetProfiles)
		ppp.POST("/profiles", handlers.Mikrotik.PPP.AddProfile)
		ppp.GET("/profiles/:name", handlers.Mikrotik.PPP.GetProfileByName)
		ppp.DELETE("/profiles/:id", handlers.Mikrotik.PPP.RemoveProfile)

		ppp.GET("/secrets", handlers.Mikrotik.PPP.GetSecrets)
		ppp.POST("/secrets", handlers.Mikrotik.PPP.AddSecret)
		ppp.GET("/secrets/:name", handlers.Mikrotik.PPP.GetSecretByName)
		ppp.DELETE("/secrets/:id", handlers.Mikrotik.PPP.RemoveSecret)

		ppp.GET("/active", handlers.Mikrotik.PPP.GetActive)

		// WebSocket streaming
		wsGroup := ppp.Group("/ws")
		{
			wsGroup.GET("/active", handlers.Mikrotik.PPPWS.ListenActive)
			wsGroup.GET("/inactive", handlers.Mikrotik.PPPWS.ListenInactive)
		}
	}
}
