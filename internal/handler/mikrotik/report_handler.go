package mikrotik

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	mikrotiksvc "mikmongo/internal/service/mikrotik"
	"mikmongo/pkg/mikrotik/domain"
	"mikmongo/pkg/response"
)

// ReportHandler handles MikroTik Report HTTP requests
type ReportHandler struct {
	service *mikrotiksvc.ReportService
}

// NewReportHandler creates a new Report handler
func NewReportHandler(service *mikrotiksvc.ReportService) *ReportHandler {
	return &ReportHandler{service: service}
}

// getRouterIDReport extracts router ID from context
func getRouterIDReport(c *gin.Context) (uuid.UUID, error) {
	routerID, exists := c.Get("router_id")
	if !exists {
		return uuid.Nil, nil
	}
	return routerID.(uuid.UUID), nil
}

// ListSalesReports handles listing sales reports
func (h *ReportHandler) ListSalesReports(c *gin.Context) {
	routerID, err := getRouterIDReport(c)
	if err != nil {
		response.BadRequest(c, "router ID not found")
		return
	}

	owner := c.Query("owner")
	reports, err := h.service.GetSalesReports(c.Request.Context(), routerID, owner)
	if err != nil {
		response.InternalServerError(c, err.Error())
		return
	}
	response.OK(c, reports)
}

// GetSalesReportsByDay handles getting sales reports by day
func (h *ReportHandler) GetSalesReportsByDay(c *gin.Context) {
	routerID, err := getRouterIDReport(c)
	if err != nil {
		response.BadRequest(c, "router ID not found")
		return
	}

	day := c.Query("day")
	if day == "" {
		response.BadRequest(c, "day query parameter is required")
		return
	}

	reports, err := h.service.GetSalesReportsByDay(c.Request.Context(), routerID, day)
	if err != nil {
		response.InternalServerError(c, err.Error())
		return
	}
	response.OK(c, reports)
}

// CreateSalesReport handles creating a sales report
func (h *ReportHandler) CreateSalesReport(c *gin.Context) {
	routerID, err := getRouterIDReport(c)
	if err != nil {
		response.BadRequest(c, "router ID not found")
		return
	}

	var report domain.SalesReport
	if err := c.ShouldBindJSON(&report); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	if err := h.service.AddSalesReport(c.Request.Context(), routerID, &report); err != nil {
		response.InternalServerError(c, err.Error())
		return
	}

	response.Created(c, report)
}

// GetReportSummary handles getting a summary of sales reports
func (h *ReportHandler) GetReportSummary(c *gin.Context) {
	routerID, err := getRouterIDReport(c)
	if err != nil {
		response.BadRequest(c, "router ID not found")
		return
	}

	owner := c.Query("owner")
	reports, err := h.service.GetSalesReports(c.Request.Context(), routerID, owner)
	if err != nil {
		response.InternalServerError(c, err.Error())
		return
	}

	// Calculate summary
	summary := &domain.ReportSummary{
		ByProfile: make(map[string]domain.ProfileSummary),
	}

	for _, r := range reports {
		summary.TotalVouchers++
		summary.TotalAmount += r.Price

		profile := r.Profile
		if profile == "" {
			profile = "default"
		}

		ps := summary.ByProfile[profile]
		ps.Count++
		ps.Total += r.Price
		summary.ByProfile[profile] = ps
	}

	response.OK(c, summary)
}
