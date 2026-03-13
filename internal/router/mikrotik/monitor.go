package mikrotik

import (
	"github.com/gin-gonic/gin"
	mikrotikhdl "mikmongo/internal/handler/mikrotik"
)

func registerMonitorRoutes(parent *gin.RouterGroup, h *mikrotikhdl.Registry) {
	monitor := parent.Group("/monitor")
	{
		monitor.GET("/system/resource", h.Monitor.GetSystemResource)
		monitor.GET("/system/health", h.Monitor.GetSystemHealth)
		monitor.GET("/system/identity", h.Monitor.GetSystemIdentity)
		monitor.GET("/system/clock", h.Monitor.GetSystemClock)
		monitor.GET("/system/routerboard", h.Monitor.GetRouterBoardInfo)
		monitor.GET("/interfaces", h.Monitor.ListInterfaces)
		monitor.GET("/logs", h.Monitor.GetLogs)
		monitor.GET("/logs/hotspot", h.Monitor.GetHotspotLogs)
		monitor.GET("/logs/ppp", h.Monitor.GetPPPLogs)
		monitor.POST("/logging/hotspot/enable", h.Monitor.EnableHotspotLogging)
		monitor.POST("/logging/ppp/enable", h.Monitor.EnablePPPLogging)
		monitor.POST("/ping", h.Monitor.Ping)

		// WebSocket streams
		monitor.GET("/system/resource/stream", h.WebSocket.StreamSystemResource)
		monitor.GET("/traffic/stream", h.WebSocket.StreamTraffic)
		monitor.GET("/logs/stream", h.WebSocket.StreamLogs)
		monitor.GET("/ping/stream", h.WebSocket.StreamPing)
	}
}
