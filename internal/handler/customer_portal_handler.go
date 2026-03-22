package handler

import (
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"mikmongo/internal/model"
	"mikmongo/internal/service"
	"mikmongo/pkg/jwt"
	gateway "mikmongo/pkg/payment"
	"mikmongo/pkg/response"
)

// CustomerPortalHandler handles customer self-service portal endpoints
type CustomerPortalHandler struct {
	customerSvc     *service.CustomerService
	subscriptionSvc *service.SubscriptionService
	billingSvc      *service.BillingService
	paymentSvc      *service.PaymentService
	jwtService      *jwt.Service
	providers       map[string]gateway.Provider
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
		providers:       map[string]gateway.Provider{},
	}
}

// SetProvider registers a gateway provider for the portal.
func (h *CustomerPortalHandler) SetProvider(name string, p gateway.Provider) {
	h.providers[name] = p
}

// Login authenticates a customer for portal access
func (h *CustomerPortalHandler) Login(c *gin.Context) {
	var req struct {
		Identifier string `json:"identifier" binding:"required"` // username or email
		Password   string `json:"password" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	customer, err := h.customerSvc.AuthPortal(c.Request.Context(), req.Identifier, req.Password)
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

// GetPayment returns a specific payment (ownership-verified).
// GET /portal/v1/payments/:id
func (h *CustomerPortalHandler) GetPayment(c *gin.Context) {
	paymentID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		response.BadRequest(c, "invalid id")
		return
	}
	payment, err := h.paymentSvc.GetPayment(c.Request.Context(), paymentID)
	if err != nil {
		response.NotFound(c, "payment not found")
		return
	}
	customerID := c.MustGet("customer_id").(string)
	if payment.CustomerID != customerID {
		response.Forbidden(c, "access denied")
		return
	}
	response.OK(c, payment)
}

// PayWithGateway creates a payment gateway invoice for a pending payment.
// POST /portal/v1/payments/:id/pay?gateway=xendit
func (h *CustomerPortalHandler) PayWithGateway(c *gin.Context) {
	paymentID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		response.BadRequest(c, "invalid id")
		return
	}
	customerID := c.MustGet("customer_id").(string)

	payment, err := h.paymentSvc.GetPayment(c.Request.Context(), paymentID)
	if err != nil {
		response.NotFound(c, "payment not found")
		return
	}
	if payment.CustomerID != customerID {
		response.Forbidden(c, "access denied")
		return
	}
	if payment.Status != "pending" {
		response.BadRequest(c, "payment is not pending")
		return
	}
	// Idempotency: if gateway URL already exists, return it without creating a new invoice.
	if payment.GatewayPaymentURL != nil && *payment.GatewayPaymentURL != "" {
		response.OK(c, gin.H{
			"payment_url": *payment.GatewayPaymentURL,
			"gateway_id":  payment.GatewayTrxID,
		})
		return
	}

	gatewayName := c.DefaultQuery("gateway", "")
	provider, ok := h.providers[gatewayName]
	if !ok {
		response.BadRequest(c, "unsupported gateway: "+gatewayName)
		return
	}

	// Build invoice request — preload customer info if available.
	req := gateway.CreateInvoiceRequest{
		ExternalID:  payment.ID,
		Amount:      payment.Amount,
		Description: "Payment " + payment.PaymentNumber,
	}
	if payment.Customer.Email != nil {
		req.CustomerEmail = *payment.Customer.Email
	}
	req.CustomerName = payment.Customer.FullName

	result, err := provider.CreateInvoice(c.Request.Context(), req)
	if err != nil {
		log.Printf("portal PayWithGateway: provider %s error for payment %s: %v", gatewayName, paymentID, err)
		response.Error(c, http.StatusBadGateway, "payment gateway error")
		return
	}

	if err := h.paymentSvc.SetGatewayInfo(c.Request.Context(), paymentID, gatewayName, result.GatewayID, result.PaymentURL, result.RawJSON); err != nil {
		log.Printf("portal PayWithGateway: SetGatewayInfo error for payment %s: %v", paymentID, err)
		response.InternalServerError(c, "failed to save gateway info")
		return
	}

	response.OK(c, gin.H{
		"payment_url": result.PaymentURL,
		"expires_at":  result.ExpiresAt,
		"gateway_id":  result.GatewayID,
	})
}
