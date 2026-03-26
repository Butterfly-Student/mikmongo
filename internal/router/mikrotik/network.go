package mikrotik

import (
	"github.com/gin-gonic/gin"
	"mikmongo/internal/handler"
)

// RegisterNetworkRoutes registers Queue, Firewall, and IP routes scoped to /api/v1/routers/:router_id
func RegisterNetworkRoutes(routerGroup *gin.RouterGroup, handlers *handler.Registry) {
	if handlers.Mikrotik == nil {
		return
	}

	// Queue
	queue := routerGroup.Group("/queue")
	{
		queue.GET("/simple", handlers.Mikrotik.Queue.GetSimpleQueues)
	}

	// Firewall
	firewall := routerGroup.Group("/firewall")
	{
		firewall.GET("/filter", handlers.Mikrotik.Firewall.GetFilterRules)
		firewall.GET("/nat", handlers.Mikrotik.Firewall.GetNATRules)
		firewall.GET("/address-list", handlers.Mikrotik.Firewall.GetAddressLists)
	}

	// IP
	ip := routerGroup.Group("/ip")
	{
		ip.GET("/pools", handlers.Mikrotik.IP.GetPools)
		ip.POST("/pools", handlers.Mikrotik.IP.AddPool)
		ip.DELETE("/pools/:id", handlers.Mikrotik.IP.RemovePool)
		ip.GET("/addresses", handlers.Mikrotik.IP.GetAddresses)
	}
}
