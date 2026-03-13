package handler

import (
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"mikmongo/internal/service"
	"mikmongo/pkg/response"
)

// ReportHandler handles report HTTP requests
type ReportHandler struct {
	service *service.ReportService
}

// NewReportHandler creates a new report handler
func NewReportHandler(service *service.ReportService) *ReportHandler {
	return &ReportHandler{service: service}
}

// GetSummary returns a summary report
func (h *ReportHandler) GetSummary(c *gin.Context) {
	from, err := time.Parse("2006-01-02", c.DefaultQuery("from", time.Now().AddDate(0, -1, 0).Format("2006-01-02")))
	if err != nil {
		response.BadRequest(c, "invalid 'from' date format, use YYYY-MM-DD")
		return
	}
	to, err := time.Parse("2006-01-02", c.DefaultQuery("to", time.Now().Format("2006-01-02")))
	if err != nil {
		response.BadRequest(c, "invalid 'to' date format, use YYYY-MM-DD")
		return
	}
	// Set to end of day
	to = to.Add(24*time.Hour - time.Second)

	summary, err := h.service.GetSummary(c.Request.Context(), from, to)
	if err != nil {
		response.InternalServerError(c, err.Error())
		return
	}
	response.OK(c, summary)
}

// GetSubscriptions returns all subscriptions report
func (h *ReportHandler) GetSubscriptions(c *gin.Context) {
	h.getSubscriptionReport(c)
}

func (h *ReportHandler) getSubscriptionReport(c *gin.Context) {
	var from, to time.Time
	var err error

	if fromStr := c.Query("from"); fromStr != "" {
		from, err = time.Parse("2006-01-02", fromStr)
		if err != nil {
			response.BadRequest(c, "invalid 'from' date format, use YYYY-MM-DD")
			return
		}
	}
	if toStr := c.Query("to"); toStr != "" {
		to, err = time.Parse("2006-01-02", toStr)
		if err != nil {
			response.BadRequest(c, "invalid 'to' date format, use YYYY-MM-DD")
			return
		}
		to = to.Add(24*time.Hour - time.Second)
	}

	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "50"))
	offset, _ := strconv.Atoi(c.DefaultQuery("offset", "0"))
	if limit <= 0 || limit > 500 {
		limit = 50
	}
	if offset < 0 {
		offset = 0
	}

	items, total, err := h.service.GetSubscriptionReport(c.Request.Context(), from, to, limit, offset)
	if err != nil {
		response.InternalServerError(c, err.Error())
		return
	}

	response.WithMeta(c, http.StatusOK, items, &response.Meta{
		Total:  total,
		Limit:  limit,
		Offset: offset,
	})
}
