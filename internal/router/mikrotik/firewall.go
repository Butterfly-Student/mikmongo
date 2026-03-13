package mikrotik

import (
	"github.com/gin-gonic/gin"
	mikrotikhdl "mikmongo/internal/handler/mikrotik"
)

func registerFirewallRoutes(parent *gin.RouterGroup, h *mikrotikhdl.Registry) {
	firewall := parent.Group("/firewall")
	{
		firewall.GET("/nat-rules", h.Firewall.ListNATRules)
		firewall.GET("/filter-rules", h.Firewall.ListFilterRules)
		firewall.GET("/address-lists", h.Firewall.ListAddressLists)
	}
}
