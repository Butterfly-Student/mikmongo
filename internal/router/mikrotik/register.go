package mikrotik

import (
	"github.com/gin-gonic/gin"
	mikrotikhdl "mikmongo/internal/handler/mikrotik"
	"mikmongo/internal/middleware"
)

// Register registers all MikroTik API routes on the given v1 group.
func Register(v1 *gin.RouterGroup, handlers *mikrotikhdl.Registry, mw *middleware.Registry) {
	mikrotikRoutes := v1.Group("/mikrotik/:router_id")
	mikrotikRoutes.Use(mw.MikrotikRouter.ValidateRouterID())
	{
		registerHotspotRoutes(mikrotikRoutes, handlers)
		registerPPPRoutes(mikrotikRoutes, handlers)
		registerQueueRoutes(mikrotikRoutes, handlers)
		registerFirewallRoutes(mikrotikRoutes, handlers)
		registerIPPoolRoutes(mikrotikRoutes, handlers)
		registerIPAddressRoutes(mikrotikRoutes, handlers)
		registerMonitorRoutes(mikrotikRoutes, handlers)
		registerReportRoutes(mikrotikRoutes, handlers)
	}

	// Script routes (no router_id required)
	registerScriptRoutes(v1, handlers, mw)
}
