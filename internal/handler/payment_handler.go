package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"mikmongo/internal/model"
	"mikmongo/internal/service"
	"mikmongo/pkg/response"
)

// PaymentHandler handles payment HTTP requests
type PaymentHandler struct {
	service *service.PaymentService
}

// NewPaymentHandler creates a new payment handler
func NewPaymentHandler(service *service.PaymentService) *PaymentHandler {
	return &PaymentHandler{service: service}
}

// List handles listing payments
func (h *PaymentHandler) List(c *gin.Context) {
	limit, offset := getPagination(c)
	payments, err := h.service.List(c.Request.Context(), limit, offset)
	if err != nil {
		response.InternalServerError(c, err.Error())
		return
	}
	response.OK(c, payments)
}

// Create handles creating a payment
func (h *PaymentHandler) Create(c *gin.Context) {
	var payment model.Payment
	if err := c.ShouldBindJSON(&payment); err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	if err := h.service.Create(c.Request.Context(), &payment); err != nil {
		response.Error(c, http.StatusBadRequest, err.Error())
		return
	}
	response.Created(c, payment)
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
	response.OK(c, payment)
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
	var req struct {
		Reason string `json:"reason" binding:"required"`
	}
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
	var req struct {
		Amount float64 `json:"amount" binding:"required"`
		Reason string  `json:"reason" binding:"required"`
	}
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
