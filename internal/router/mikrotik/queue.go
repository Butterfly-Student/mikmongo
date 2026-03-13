package mikrotik

import (
	"github.com/gin-gonic/gin"
	mikrotikhdl "mikmongo/internal/handler/mikrotik"
)

func registerQueueRoutes(parent *gin.RouterGroup, h *mikrotikhdl.Registry) {
	queues := parent.Group("/queues")
	{
		queues.GET("/simple", h.Queue.ListSimpleQueues)
		queues.POST("/simple", h.Queue.CreateSimpleQueue)
		queues.GET("/simple/:id", h.Queue.GetSimpleQueue)
		queues.GET("/simple/by-name", h.Queue.GetSimpleQueueByName)
		queues.PUT("/simple/:id", h.Queue.UpdateSimpleQueue)
		queues.DELETE("/simple/:id", h.Queue.DeleteSimpleQueue)
		queues.POST("/simple/:id/enable", h.Queue.EnableSimpleQueue)
		queues.POST("/simple/:id/disable", h.Queue.DisableSimpleQueue)
		queues.POST("/simple/:id/reset-counters", h.Queue.ResetQueueCounters)
		queues.POST("/reset-all-counters", h.Queue.ResetAllQueueCounters)
		queues.GET("/names", h.Queue.GetAllQueues)
		queues.GET("/parent-names", h.Queue.GetAllParentQueues)

		// WebSocket streams
		queues.GET("/stats/stream", h.WebSocket.StreamQueueStats)
	}
}
