package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"mikmongo/internal/dto"
	"mikmongo/internal/service"
	"mikmongo/pkg/response"
)

// BillingHandler handles billing HTTP requests
type BillingHandler struct {
	service *service.BillingService
}

// NewBillingHandler creates a new billing handler
func NewBillingHandler(service *service.BillingService) *BillingHandler {
	return &BillingHandler{service: service}
}

// ListInvoices handles listing all invoices
func (h *BillingHandler) ListInvoices(c *gin.Context) {
	limit, offset := getPagination(c)
	invoices, err := h.service.List(c.Request.Context(), limit, offset)
	if err != nil {
		response.InternalServerError(c, err.Error())
		return
	}
	response.OK(c, dto.InvoicesToResponse(invoices))
}

// GetInvoice handles getting an invoice by ID
func (h *BillingHandler) GetInvoice(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		response.BadRequest(c, "invalid id")
		return
	}
	invoice, err := h.service.GetInvoice(c.Request.Context(), id)
	if err != nil {
		response.NotFound(c, err.Error())
		return
	}
	response.OK(c, dto.InvoiceToResponse(invoice))
}

// GetOverdue handles getting overdue invoices
func (h *BillingHandler) GetOverdue(c *gin.Context) {
	invoices, err := h.service.GetOverdueInvoices(c.Request.Context())
	if err != nil {
		response.InternalServerError(c, err.Error())
		return
	}
	response.OK(c, dto.InvoicesToResponse(invoices))
}

// CancelInvoice handles cancelling an invoice
func (h *BillingHandler) CancelInvoice(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		response.BadRequest(c, "invalid id")
		return
	}
	if err := h.service.Cancel(c.Request.Context(), id); err != nil {
		response.Error(c, http.StatusBadRequest, err.Error())
		return
	}
	response.OK(c, gin.H{"message": "invoice cancelled"})
}

// TriggerMonthlyBilling triggers monthly billing process
func (h *BillingHandler) TriggerMonthlyBilling(c *gin.Context) {
	if err := h.service.ProcessMonthlyBilling(c.Request.Context()); err != nil {
		response.InternalServerError(c, err.Error())
		return
	}
	response.OK(c, gin.H{"message": "billing process triggered"})
}
