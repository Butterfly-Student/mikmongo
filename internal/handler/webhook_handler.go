package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"mikmongo/internal/service"
	gateway "mikmongo/pkg/payment"
	"mikmongo/pkg/response"
)

// WebhookHandler handles webhook HTTP requests
type WebhookHandler struct {
	paymentService *service.PaymentService
	xenditProvider gateway.Provider
}

// NewWebhookHandler creates a new webhook handler
func NewWebhookHandler(paymentService *service.PaymentService) *WebhookHandler {
	return &WebhookHandler{paymentService: paymentService}
}

// SetXenditProvider injects the Xendit gateway provider.
func (h *WebhookHandler) SetXenditProvider(p gateway.Provider) {
	h.xenditProvider = p
}

// MidtransWebhook handles Midtrans payment webhook
func (h *WebhookHandler) MidtransWebhook(c *gin.Context) {
	var payload struct {
		TransactionID   string `json:"transaction_id"`
		TransactionStatus string `json:"transaction_status"`
	}
	
	if err := c.ShouldBindJSON(&payload); err != nil {
		response.Error(c, http.StatusBadRequest, err.Error())
		return
	}
	
	var status string
	switch payload.TransactionStatus {
	case "capture", "settlement":
		status = "confirmed"
	case "deny", "cancel", "expire":
		status = "rejected"
	default:
		status = "pending"
	}
	
	if err := h.paymentService.HandleWebhook(c.Request.Context(), payload.TransactionID, status); err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}
	
	response.Success(c, http.StatusOK, gin.H{"message": "webhook processed"})
}

// XenditWebhook handles Xendit invoice webhook.
// POST /api/v1/webhooks/xendit
func (h *WebhookHandler) XenditWebhook(c *gin.Context) {
	if h.xenditProvider == nil {
		response.Error(c, http.StatusServiceUnavailable, "xendit not configured")
		return
	}

	event, err := h.xenditProvider.VerifyWebhook(c.Request)
	if err != nil {
		response.Error(c, http.StatusUnauthorized, "webhook verification failed: "+err.Error())
		return
	}

	if err := h.paymentService.HandleGatewayWebhook(c.Request.Context(), event); err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}

	response.Success(c, http.StatusOK, gin.H{"message": "webhook processed"})
}
