package mikrotik

import (
	"github.com/gin-gonic/gin"
	mikrotikhdl "mikmongo/internal/handler/mikrotik"
)

func registerIPAddressRoutes(parent *gin.RouterGroup, h *mikrotikhdl.Registry) {
	ipAddresses := parent.Group("/ip-addresses")
	{
		ipAddresses.GET("", h.IPAddress.ListAddresses)
		ipAddresses.POST("", h.IPAddress.CreateAddress)
		ipAddresses.GET("/:id", h.IPAddress.GetAddress)
		ipAddresses.PUT("/:id", h.IPAddress.UpdateAddress)
		ipAddresses.DELETE("/:id", h.IPAddress.DeleteAddress)
		ipAddresses.POST("/:id/enable", h.IPAddress.EnableAddress)
		ipAddresses.POST("/:id/disable", h.IPAddress.DisableAddress)
		ipAddresses.GET("/interface/:interface", h.IPAddress.GetAddressesByInterface)
	}
}
