package mikrotik

import (
	"github.com/gin-gonic/gin"
	"mikmongo/internal/handler"
)

// RegisterMonitorRoutes registers Monitor REST and WebSocket routes scoped to /api/v1/routers/:router_id
func RegisterMonitorRoutes(routerGroup *gin.RouterGroup, handlers *handler.Registry) {
	if handlers.Mikrotik == nil {
		return
	}

	monitor := routerGroup.Group("/monitor")
	{
		// REST endpoints
		monitor.GET("/system-resource", handlers.Mikrotik.Monitor.GetSystemResource)
		monitor.GET("/interfaces", handlers.Mikrotik.Monitor.GetInterfaces)

		// WebSocket streaming endpoints
		wsGroup := monitor.Group("/ws")
		{
			wsGroup.GET("/system-resource", handlers.Mikrotik.MonitorWS.SystemResource)
			wsGroup.GET("/traffic/:name", handlers.Mikrotik.MonitorWS.Traffic)
			wsGroup.GET("/logs", handlers.Mikrotik.MonitorWS.Logs)
			wsGroup.GET("/ping", handlers.Mikrotik.MonitorWS.Ping)
		}
	}
}
