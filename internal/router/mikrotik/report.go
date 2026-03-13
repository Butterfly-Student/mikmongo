package mikrotik

import (
	"github.com/gin-gonic/gin"
	mikrotikhdl "mikmongo/internal/handler/mikrotik"
)

func registerReportRoutes(parent *gin.RouterGroup, h *mikrotikhdl.Registry) {
	reports := parent.Group("/reports")
	{
		reports.GET("/sales", h.Report.ListSalesReports)
		reports.GET("/sales/by-day", h.Report.GetSalesReportsByDay)
		reports.POST("/sales", h.Report.CreateSalesReport)
		reports.GET("/sales/summary", h.Report.GetReportSummary)
	}
}
