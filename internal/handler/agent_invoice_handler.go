package handler

import (
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"mikmongo/internal/repository"
	"mikmongo/internal/service"
	"mikmongo/pkg/response"
)

// AgentInvoiceHandler handles agent invoice management.
type AgentInvoiceHandler struct {
	svc *service.AgentInvoiceService
}

// NewAgentInvoiceHandler creates a new AgentInvoiceHandler.
func NewAgentInvoiceHandler(svc *service.AgentInvoiceService) *AgentInvoiceHandler {
	return &AgentInvoiceHandler{svc: svc}
}

type generateAgentInvoiceRequest struct {
	PeriodStart string `json:"period_start" binding:"required"` // YYYY-MM-DD
	PeriodEnd   string `json:"period_end" binding:"required"`   // YYYY-MM-DD
}

type markPaidRequest struct {
	PaidAmount float64 `json:"paid_amount" binding:"required,min=0"`
}

// List handles listing agent invoices with optional filters.
// Query params: agent_id, router_id, status, billing_cycle, billing_year, billing_month, billing_week
func (h *AgentInvoiceHandler) List(c *gin.Context) {
	filter := repository.AgentInvoiceFilter{}

	if v := c.Query("agent_id"); v != "" {
		id, err := uuid.Parse(v)
		if err != nil {
			response.BadRequest(c, "invalid agent_id")
			return
		}
		filter.AgentID = &id
	}
	if v := c.Query("router_id"); v != "" {
		id, err := uuid.Parse(v)
		if err != nil {
			response.BadRequest(c, "invalid router_id")
			return
		}
		filter.RouterID = &id
	}
	filter.Status = c.Query("status")
	filter.BillingCycle = c.Query("billing_cycle")
	if v := c.Query("billing_year"); v != "" {
		var y int
		if _, err := toInt(v, &y); err != nil {
			response.BadRequest(c, "invalid billing_year")
			return
		}
		filter.BillingYear = &y
	}
	if v := c.Query("billing_month"); v != "" {
		var m int
		if _, err := toInt(v, &m); err != nil {
			response.BadRequest(c, "invalid billing_month")
			return
		}
		filter.BillingMonth = &m
	}
	if v := c.Query("billing_week"); v != "" {
		var w int
		if _, err := toInt(v, &w); err != nil {
			response.BadRequest(c, "invalid billing_week")
			return
		}
		filter.BillingWeek = &w
	}

	limit, offset := getPagination(c)
	invs, count, err := h.svc.ListInvoices(c.Request.Context(), filter, limit, offset)
	if err != nil {
		response.InternalServerError(c, err.Error())
		return
	}
	response.WithMeta(c, http.StatusOK, invs, &response.Meta{Total: count, Limit: limit, Offset: offset})
}

// ListByAgent handles listing invoices for a specific agent (path param :id).
func (h *AgentInvoiceHandler) ListByAgent(c *gin.Context) {
	agentID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		response.BadRequest(c, "invalid agent id")
		return
	}
	limit, offset := getPagination(c)
	invs, count, err := h.svc.ListByAgent(c.Request.Context(), agentID, limit, offset)
	if err != nil {
		response.InternalServerError(c, err.Error())
		return
	}
	response.WithMeta(c, http.StatusOK, invs, &response.Meta{Total: count, Limit: limit, Offset: offset})
}

// Get handles fetching a single invoice by ID.
func (h *AgentInvoiceHandler) Get(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		response.BadRequest(c, "invalid id")
		return
	}
	inv, err := h.svc.GetInvoice(c.Request.Context(), id)
	if err != nil {
		response.NotFound(c, err.Error())
		return
	}
	response.OK(c, inv)
}

// Generate handles manual invoice generation for a specific agent.
// Path param :id = agent ID. Body: period_start, period_end (YYYY-MM-DD).
func (h *AgentInvoiceHandler) Generate(c *gin.Context) {
	agentID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		response.BadRequest(c, "invalid agent id")
		return
	}

	var req generateAgentInvoiceRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	periodStart, err := time.Parse("2006-01-02", req.PeriodStart)
	if err != nil {
		response.BadRequest(c, "invalid period_start, use YYYY-MM-DD")
		return
	}
	periodEnd, err := time.Parse("2006-01-02", req.PeriodEnd)
	if err != nil {
		response.BadRequest(c, "invalid period_end, use YYYY-MM-DD")
		return
	}
	// period_end is exclusive: add one day so "2026-03-31" covers the full last day
	periodEnd = periodEnd.AddDate(0, 0, 1)

	inv, err := h.svc.GenerateManual(c.Request.Context(), agentID, periodStart, periodEnd)
	if err != nil {
		response.InternalServerError(c, err.Error())
		return
	}
	response.Created(c, inv)
}

// MarkPaid handles marking an invoice as paid.
func (h *AgentInvoiceHandler) MarkPaid(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		response.BadRequest(c, "invalid id")
		return
	}

	var req markPaidRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	inv, err := h.svc.MarkPaid(c.Request.Context(), id, req.PaidAmount)
	if err != nil {
		response.InternalServerError(c, err.Error())
		return
	}
	response.OK(c, inv)
}

// Cancel handles cancelling an invoice.
func (h *AgentInvoiceHandler) Cancel(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		response.BadRequest(c, "invalid id")
		return
	}
	if err := h.svc.Cancel(c.Request.Context(), id); err != nil {
		response.InternalServerError(c, err.Error())
		return
	}
	response.OK(c, gin.H{"message": "cancelled"})
}

// ProcessScheduled handles manually triggering the scheduled invoice generation.
func (h *AgentInvoiceHandler) ProcessScheduled(c *gin.Context) {
	if err := h.svc.ProcessScheduled(c.Request.Context()); err != nil {
		response.InternalServerError(c, err.Error())
		return
	}
	response.OK(c, gin.H{"message": "processing complete"})
}

// toInt parses a string into an int pointer target and returns an error if invalid.
func toInt(s string, out *int) (int, error) {
	var n int
	_, err := fmt.Sscanf(s, "%d", &n)
	if err != nil {
		return 0, err
	}
	*out = n
	return n, nil
}
