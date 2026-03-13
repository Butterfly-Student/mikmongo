package mikrotik

import (
	"github.com/gin-gonic/gin"
	mikrotikhdl "mikmongo/internal/handler/mikrotik"
)

func registerHotspotRoutes(parent *gin.RouterGroup, h *mikrotikhdl.Registry) {
	hotspot := parent.Group("/hotspot")
	{
		// Users
		hotspot.GET("/users", h.Hotspot.ListUsers)
		hotspot.POST("/users", h.Hotspot.CreateUser)
		hotspot.GET("/users/:id", h.Hotspot.GetUser)
		hotspot.PUT("/users/:id", h.Hotspot.UpdateUser)
		hotspot.DELETE("/users/:id", h.Hotspot.DeleteUser)
		hotspot.GET("/users/by-name", h.Hotspot.GetUserByName)
		hotspot.GET("/users/by-comment", h.Hotspot.GetUsersByComment)
		hotspot.POST("/users/:id/enable", h.Hotspot.EnableUser)
		hotspot.POST("/users/:id/disable", h.Hotspot.DisableUser)
		hotspot.POST("/users/:id/reset-counters", h.Hotspot.ResetUserCounters)
		hotspot.POST("/users/batch-delete", h.Hotspot.BatchDeleteUsers)
		hotspot.POST("/users/batch-enable", h.Hotspot.BatchEnableUsers)
		hotspot.POST("/users/batch-disable", h.Hotspot.BatchDisableUsers)
		hotspot.POST("/users/batch-reset-counters", h.Hotspot.BatchResetUserCounters)
		hotspot.GET("/users/count", h.Hotspot.GetUsersCount)

		// Profiles
		hotspot.GET("/profiles", h.Hotspot.ListProfiles)
		hotspot.POST("/profiles", h.Hotspot.CreateProfile)
		hotspot.GET("/profiles/:id", h.Hotspot.GetProfile)
		hotspot.GET("/profiles/by-name", h.Hotspot.GetProfileByName)
		hotspot.PUT("/profiles/:id", h.Hotspot.UpdateProfile)
		hotspot.DELETE("/profiles/:id", h.Hotspot.DeleteProfile)
		hotspot.POST("/profiles/:id/enable", h.Hotspot.EnableProfile)
		hotspot.POST("/profiles/:id/disable", h.Hotspot.DisableProfile)
		hotspot.POST("/profiles/batch-delete", h.Hotspot.BatchDeleteProfiles)
		hotspot.POST("/profiles/batch-enable", h.Hotspot.BatchEnableProfiles)
		hotspot.POST("/profiles/batch-disable", h.Hotspot.BatchDisableProfiles)

		// Active sessions
		hotspot.GET("/active", h.Hotspot.ListActive)
		hotspot.GET("/active/count", h.Hotspot.GetActiveCount)
		hotspot.DELETE("/active/:id", h.Hotspot.DeleteActive)
		hotspot.POST("/active/batch-delete", h.Hotspot.BatchDeleteActive)

		// Hosts
		hotspot.GET("/hosts", h.Hotspot.ListHosts)
		hotspot.DELETE("/hosts/:id", h.Hotspot.DeleteHost)

		// Servers
		hotspot.GET("/servers", h.Hotspot.ListServers)

		// IP Bindings
		hotspot.GET("/ip-bindings", h.Hotspot.ListIPBindings)
		hotspot.POST("/ip-bindings", h.Hotspot.CreateIPBinding)
		hotspot.DELETE("/ip-bindings/:id", h.Hotspot.DeleteIPBinding)
		hotspot.POST("/ip-bindings/:id/enable", h.Hotspot.EnableIPBinding)
		hotspot.POST("/ip-bindings/:id/disable", h.Hotspot.DisableIPBinding)

		// WebSocket streams
		hotspot.GET("/active/stream", h.WebSocket.StreamHotspotActive)
		hotspot.GET("/inactive/stream", h.WebSocket.StreamHotspotInactive)
	}
}
