package mikhmon

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"mikmongo/internal/dto/mikrotik/mikhmon"
	mikhmonSvc "mikmongo/internal/service/mikrotik/mikhmon"
	"mikmongo/pkg/response"
)

type ReportHandler struct {
	reportSvc *mikhmonSvc.MikhmonReportService
}

func NewReportHandler(reportSvc *mikhmonSvc.MikhmonReportService) *ReportHandler {
	return &ReportHandler{
		reportSvc: reportSvc,
	}
}

func (h *ReportHandler) Add(c *gin.Context) {
	routerID, err := uuid.Parse(c.Param("router_id"))
	if err != nil {
		response.BadRequest(c, "invalid router_id")
		return
	}

	var req mikhmon.CreateReportRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	if err := h.reportSvc.AddReport(c.Request.Context(), routerID, req.ToDomain()); err != nil {
		response.InternalServerError(c, err.Error())
		return
	}
	response.Created(c, gin.H{"message": "report added"})
}

func (h *ReportHandler) GetReports(c *gin.Context) {
	routerID, err := uuid.Parse(c.Param("router_id"))
	if err != nil {
		response.BadRequest(c, "invalid router_id")
		return
	}

	var query mikhmon.ReportFilterQuery
	if err := c.ShouldBindQuery(&query); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	if query.Owner != "" {
		reports, err := h.reportSvc.GetReportsByOwner(c.Request.Context(), routerID, query.Owner)
		if err != nil {
			response.InternalServerError(c, err.Error())
			return
		}
		response.OK(c, mikhmon.SalesReportsToResponse(reports))
		return
	}

	if query.Day != "" {
		reports, err := h.reportSvc.GetReportsByDay(c.Request.Context(), routerID, query.Day)
		if err != nil {
			response.InternalServerError(c, err.Error())
			return
		}
		response.OK(c, mikhmon.SalesReportsToResponse(reports))
		return
	}

	summary, err := h.reportSvc.GetReportSummary(c.Request.Context(), routerID, query.ToDomain())
	if err != nil {
		response.InternalServerError(c, err.Error())
		return
	}
	response.OK(c, mikhmon.ReportSummaryToResponse(summary))
}

func (h *ReportHandler) GetSummary(c *gin.Context) {
	routerID, err := uuid.Parse(c.Param("router_id"))
	if err != nil {
		response.BadRequest(c, "invalid router_id")
		return
	}

	var query mikhmon.ReportFilterQuery
	if err := c.ShouldBindQuery(&query); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	summary, err := h.reportSvc.GetReportSummary(c.Request.Context(), routerID, query.ToDomain())
	if err != nil {
		response.InternalServerError(c, err.Error())
		return
	}
	response.OK(c, mikhmon.ReportSummaryToResponse(summary))
}
