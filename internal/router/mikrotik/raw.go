package mikrotik

import (
	"github.com/gin-gonic/gin"
	"mikmongo/internal/handler"
)

// RegisterRawRoutes registers raw RouterOS command routes scoped to /api/v1/routers/:router_id
func RegisterRawRoutes(routerGroup *gin.RouterGroup, handlers *handler.Registry) {
	if handlers.Mikrotik == nil {
		return
	}

	raw := routerGroup.Group("/raw")
	{
		raw.POST("/run", handlers.Mikrotik.Raw.Run)

		// WebSocket streaming endpoint
		raw.GET("/ws/listen", handlers.Mikrotik.Raw.ListenWS)
	}
}
