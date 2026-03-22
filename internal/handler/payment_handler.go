package handler

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"mikmongo/internal/dto"
	"mikmongo/internal/service"
	gateway "mikmongo/pkg/payment"
	"mikmongo/pkg/response"
)

// PaymentHandler handles payment HTTP requests
type PaymentHandler struct {
	service   *service.PaymentService
	providers map[string]gateway.Provider
}

// NewPaymentHandler creates a new payment handler
func NewPaymentHandler(svc *service.PaymentService) *PaymentHandler {
	return &PaymentHandler{service: svc, providers: map[string]gateway.Provider{}}
}

// SetProvider registers a gateway provider (e.g. "xendit").
func (h *PaymentHandler) SetProvider(name string, p gateway.Provider) {
	h.providers[name] = p
}

// List handles listing payments
func (h *PaymentHandler) List(c *gin.Context) {
	limit, offset := getPagination(c)
	payments, err := h.service.List(c.Request.Context(), limit, offset)
	if err != nil {
		response.InternalServerError(c, err.Error())
		return
	}
	response.OK(c, dto.PaymentsToResponse(payments))
}

// Create handles creating a payment
func (h *PaymentHandler) Create(c *gin.Context) {
	var req dto.CreatePaymentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	payment := req.ToModel()
	if err := h.service.Create(c.Request.Context(), payment); err != nil {
		response.Error(c, http.StatusBadRequest, err.Error())
		return
	}
	response.Created(c, dto.PaymentToResponse(payment))
}

// Get handles getting a payment by ID
func (h *PaymentHandler) Get(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		response.BadRequest(c, "invalid id")
		return
	}
	payment, err := h.service.GetPayment(c.Request.Context(), id)
	if err != nil {
		response.NotFound(c, err.Error())
		return
	}
	response.OK(c, dto.PaymentToResponse(payment))
}

// Confirm handles confirming a payment
func (h *PaymentHandler) Confirm(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		response.BadRequest(c, "invalid id")
		return
	}
	userID, _ := c.Get("user_id")
	processedBy := ""
	if userID != nil {
		processedBy = userID.(string)
	}
	if err := h.service.Confirm(c.Request.Context(), id, processedBy); err != nil {
		response.Error(c, http.StatusBadRequest, err.Error())
		return
	}
	response.OK(c, gin.H{"message": "payment confirmed"})
}

// Reject handles rejecting a payment
func (h *PaymentHandler) Reject(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		response.BadRequest(c, "invalid id")
		return
	}
	var req dto.RejectPaymentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	if err := h.service.Reject(c.Request.Context(), id, req.Reason); err != nil {
		response.Error(c, http.StatusBadRequest, err.Error())
		return
	}
	response.OK(c, gin.H{"message": "payment rejected"})
}

// Refund handles refunding a payment
func (h *PaymentHandler) Refund(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		response.BadRequest(c, "invalid id")
		return
	}
	var req dto.RefundPaymentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	if err := h.service.Refund(c.Request.Context(), id, req.Amount, req.Reason); err != nil {
		response.Error(c, http.StatusBadRequest, err.Error())
		return
	}
	response.OK(c, gin.H{"message": "payment refunded"})
}

// InitiateGateway creates a gateway invoice for a payment.
// POST /api/v1/payments/:id/initiate-gateway?gateway=xendit
func (h *PaymentHandler) InitiateGateway(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		response.BadRequest(c, "invalid id")
		return
	}

	gatewayName := c.Query("gateway")
	if gatewayName == "" {
		response.BadRequest(c, "gateway query param is required")
		return
	}

	provider, ok := h.providers[gatewayName]
	if !ok {
		response.Error(c, http.StatusBadRequest, "unsupported gateway: "+gatewayName)
		return
	}

	payment, err := h.service.GetPayment(c.Request.Context(), id)
	if err != nil {
		response.NotFound(c, err.Error())
		return
	}

	// Guard: payment.ID must be a valid UUID before sending to external gateway
	if _, err := uuid.Parse(payment.ID); err != nil {
		response.InternalServerError(c, "payment has invalid ID")
		return
	}

	// Build customer info for the invoice
	var customerEmail, customerName string
	if payment.Customer.Email != nil {
		customerEmail = *payment.Customer.Email
	}
	customerName = payment.Customer.FullName

	result, err := provider.CreateInvoice(c.Request.Context(), gateway.CreateInvoiceRequest{
		ExternalID:    payment.ID,
		Amount:        payment.Amount,
		Description:   "Payment " + payment.PaymentNumber,
		CustomerEmail: customerEmail,
		CustomerName:  customerName,
	})
	if err != nil {
		log.Printf("ERROR InitiateGateway gateway=%s payment=%s: %v", gatewayName, payment.ID, err)
		response.Error(c, http.StatusBadGateway, "payment gateway is unavailable, please try again later")
		return
	}

	if err := h.service.SetGatewayInfo(
		c.Request.Context(), id,
		provider.Name(), result.GatewayID, result.PaymentURL, result.RawJSON,
	); err != nil {
		response.InternalServerError(c, err.Error())
		return
	}

	response.OK(c, gin.H{
		"payment_url": result.PaymentURL,
		"expires_at":  result.ExpiresAt,
		"gateway_id":  result.GatewayID,
	})
}
