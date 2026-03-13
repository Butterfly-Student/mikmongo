package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"mikmongo/internal/service"
	"mikmongo/pkg/response"
)

// WebhookHandler handles webhook HTTP requests
type WebhookHandler struct {
	paymentService *service.PaymentService
}

// NewWebhookHandler creates a new webhook handler
func NewWebhookHandler(paymentService *service.PaymentService) *WebhookHandler {
	return &WebhookHandler{paymentService: paymentService}
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
