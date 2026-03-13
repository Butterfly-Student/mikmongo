package handler

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"mikmongo/internal/model"
	"mikmongo/internal/service"
	"mikmongo/pkg/jwt"
	"mikmongo/pkg/response"
)

// CustomerPortalHandler handles customer self-service portal endpoints
type CustomerPortalHandler struct {
	customerSvc     *service.CustomerService
	subscriptionSvc *service.SubscriptionService
	billingSvc      *service.BillingService
	paymentSvc      *service.PaymentService
	jwtService      *jwt.Service
}

// NewCustomerPortalHandler creates a new customer portal handler
func NewCustomerPortalHandler(
	customerSvc *service.CustomerService,
	subscriptionSvc *service.SubscriptionService,
	billingSvc *service.BillingService,
	paymentSvc *service.PaymentService,
	jwtService *jwt.Service,
) *CustomerPortalHandler {
	return &CustomerPortalHandler{
		customerSvc:     customerSvc,
		subscriptionSvc: subscriptionSvc,
		billingSvc:      billingSvc,
		paymentSvc:      paymentSvc,
		jwtService:      jwtService,
	}
}

// Login authenticates a customer for portal access
func (h *CustomerPortalHandler) Login(c *gin.Context) {
	var req struct {
		CustomerCode string `json:"customer_code" binding:"required"`
		Password     string `json:"password" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	customer, err := h.customerSvc.AuthPortal(c.Request.Context(), req.CustomerCode, req.Password)
	if err != nil {
		response.Error(c, http.StatusUnauthorized, err.Error())
		return
	}
	token, err := h.jwtService.GeneratePortal(customer.ID, customer.CustomerCode)
	if err != nil {
		response.InternalServerError(c, "failed to generate token")
		return
	}
	customer.PortalPasswordHash = nil
	response.OK(c, gin.H{"token": token, "customer": customer})
}

// GetProfile returns the customer's profile
func (h *CustomerPortalHandler) GetProfile(c *gin.Context) {
	customerID := c.MustGet("customer_id").(string)
	id, _ := uuid.Parse(customerID)
	customer, err := h.customerSvc.GetByID(c.Request.Context(), id)
	if err != nil {
		response.NotFound(c, "customer not found")
		return
	}
	customer.PortalPasswordHash = nil
	response.OK(c, customer)
}

// ChangePortalPassword changes the customer's portal password
func (h *CustomerPortalHandler) ChangePortalPassword(c *gin.Context) {
	customerID := c.MustGet("customer_id").(string)
	id, _ := uuid.Parse(customerID)
	var req struct {
		Password string `json:"password" binding:"required,min=6"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	if err := h.customerSvc.SetPortalPassword(c.Request.Context(), id, req.Password); err != nil {
		response.InternalServerError(c, err.Error())
		return
	}
	response.OK(c, gin.H{"message": "password changed"})
}

// GetSubscriptions returns the customer's subscriptions
func (h *CustomerPortalHandler) GetSubscriptions(c *gin.Context) {
	customerID := c.MustGet("customer_id").(string)
	id, _ := uuid.Parse(customerID)
	subs, err := h.subscriptionSvc.GetByCustomerID(c.Request.Context(), id)
	if err != nil {
		response.InternalServerError(c, err.Error())
		return
	}
	response.OK(c, subs)
}

// GetInvoices returns the customer's invoices
func (h *CustomerPortalHandler) GetInvoices(c *gin.Context) {
	customerID := c.MustGet("customer_id").(string)
	id, _ := uuid.Parse(customerID)
	invoices, err := h.billingSvc.GetByCustomer(c.Request.Context(), id)
	if err != nil {
		response.InternalServerError(c, err.Error())
		return
	}
	response.OK(c, invoices)
}

// GetInvoice returns a specific invoice (verified ownership)
func (h *CustomerPortalHandler) GetInvoice(c *gin.Context) {
	invoiceID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		response.BadRequest(c, "invalid id")
		return
	}
	invoice, err := h.billingSvc.GetInvoice(c.Request.Context(), invoiceID)
	if err != nil {
		response.NotFound(c, err.Error())
		return
	}
	customerID := c.MustGet("customer_id").(string)
	if invoice.CustomerID != customerID {
		response.Forbidden(c, "access denied")
		return
	}
	response.OK(c, invoice)
}

// CreatePayment creates a payment from the portal
func (h *CustomerPortalHandler) CreatePayment(c *gin.Context) {
	customerID := c.MustGet("customer_id").(string)
	var req struct {
		Amount        float64 `json:"amount" binding:"required"`
		PaymentMethod string  `json:"payment_method" binding:"required"`
		Notes         *string `json:"notes"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	payment := &model.Payment{
		CustomerID:    customerID,
		Amount:        req.Amount,
		PaymentMethod: req.PaymentMethod,
		PaymentDate:   time.Now(),
		Notes:         req.Notes,
	}
	if err := h.paymentSvc.Create(c.Request.Context(), payment); err != nil {
		response.Error(c, http.StatusBadRequest, err.Error())
		return
	}
	response.Created(c, payment)
}

// GetPayments returns the customer's payment history
func (h *CustomerPortalHandler) GetPayments(c *gin.Context) {
	customerID := c.MustGet("customer_id").(string)
	id, _ := uuid.Parse(customerID)
	payments, err := h.paymentSvc.GetByCustomer(c.Request.Context(), id)
	if err != nil {
		response.InternalServerError(c, err.Error())
		return
	}
	response.OK(c, payments)
}
