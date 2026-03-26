package mikrotik

import (
	"github.com/gin-gonic/gin"
	"mikmongo/internal/handler"
)

// RegisterHotspotRoutes registers Hotspot routes scoped to /api/v1/routers/:router_id
func RegisterHotspotRoutes(routerGroup *gin.RouterGroup, handlers *handler.Registry) {
	if handlers.Mikrotik == nil {
		return
	}

	hotspot := routerGroup.Group("/hotspot")
	{
		hotspot.GET("/profiles", handlers.Mikrotik.Hotspot.GetProfiles)
		hotspot.POST("/profiles", handlers.Mikrotik.Hotspot.AddProfile)
		hotspot.GET("/profiles/:name", handlers.Mikrotik.Hotspot.GetProfileByName)
		hotspot.DELETE("/profiles/:id", handlers.Mikrotik.Hotspot.RemoveProfile)

		hotspot.GET("/users", handlers.Mikrotik.Hotspot.GetUsers)
		hotspot.POST("/users", handlers.Mikrotik.Hotspot.AddUser)
		hotspot.GET("/users/:name", handlers.Mikrotik.Hotspot.GetUserByName)
		hotspot.DELETE("/users/:id", handlers.Mikrotik.Hotspot.RemoveUser)

		hotspot.GET("/active", handlers.Mikrotik.Hotspot.GetActive)
		hotspot.GET("/hosts", handlers.Mikrotik.Hotspot.GetHosts)
		hotspot.GET("/servers", handlers.Mikrotik.Hotspot.GetServers)

		// WebSocket streaming
		wsGroup := hotspot.Group("/ws")
		{
			wsGroup.GET("/active", handlers.Mikrotik.HotspotWS.ListenActive)
			wsGroup.GET("/inactive", handlers.Mikrotik.HotspotWS.ListenInactive)
		}
	}
}
