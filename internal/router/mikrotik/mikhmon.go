package mikrotik

import (
	"github.com/gin-gonic/gin"
	"mikmongo/internal/handler"
)

// RegisterMikhmonRoutes registers Mikhmon-specific routes scoped to a router group.
// routerGroup is expected to be /api/v1/routers/:router_id
func RegisterMikhmonRoutes(routerGroup *gin.RouterGroup, handlers *handler.Registry) {
	if handlers.Mikhmon == nil {
		return
	}

	mikhmonGroup := routerGroup.Group("/mikhmon")

	// Vouchers
	mikhmonGroup.POST("/vouchers/generate", handlers.Mikhmon.Voucher.GenerateBatch)
	mikhmonGroup.GET("/vouchers", handlers.Mikhmon.Voucher.GetVouchers)
	mikhmonGroup.DELETE("/vouchers", handlers.Mikhmon.Voucher.RemoveBatch)

	// Profiles
	mikhmonGroup.POST("/profiles", handlers.Mikhmon.Profile.Create)
	mikhmonGroup.PUT("/profiles/:id", handlers.Mikhmon.Profile.Update)
	mikhmonGroup.POST("/profiles/generate-script", handlers.Mikhmon.Profile.GenerateScript)

	// Reports
	mikhmonGroup.POST("/reports", handlers.Mikhmon.Report.Add)
	mikhmonGroup.GET("/reports", handlers.Mikhmon.Report.GetReports)
	mikhmonGroup.GET("/reports/summary", handlers.Mikhmon.Report.GetSummary)

	// Expire Monitor
	mikhmonGroup.POST("/expire/setup", handlers.Mikhmon.Expire.Setup)
	mikhmonGroup.POST("/expire/disable", handlers.Mikhmon.Expire.Disable)
	mikhmonGroup.GET("/expire/status", handlers.Mikhmon.Expire.GetStatus)
	mikhmonGroup.GET("/expire/generate-script", handlers.Mikhmon.Expire.GenerateScript)
}
