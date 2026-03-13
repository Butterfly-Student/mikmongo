package handler

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"mikmongo/internal/model"
	"mikmongo/internal/service"
	"mikmongo/pkg/response"
)

// CustomerHandler handles customer HTTP requests
type CustomerHandler struct {
	service *service.CustomerService
}

// NewCustomerHandler creates a new customer handler
func NewCustomerHandler(service *service.CustomerService) *CustomerHandler {
	return &CustomerHandler{service: service}
}

// CreateCustomerRequest represents the request body for creating a customer with subscription
type CreateCustomerRequest struct {
	// Customer fields
	FullName  string   `json:"full_name" binding:"required"`
	Phone     string   `json:"phone" binding:"required"`
	Email     *string  `json:"email"`
	Address   *string  `json:"address"`
	Latitude  *float64 `json:"latitude"`
	Longitude *float64 `json:"longitude"`

	// Subscription fields (mandatory)
	PlanID   string  `json:"plan_id" binding:"required"`
	RouterID string  `json:"router_id" binding:"required"`
	Username string  `json:"username"`  // optional, default dari FullName
	Password string  `json:"password"`  // optional, default dari FullName
	StaticIP *string `json:"static_ip"` // optional
}

// Create handles customer creation with auto subscription
func (h *CustomerHandler) Create(c *gin.Context) {
	var req CreateCustomerRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	// Generate username from FullName if not provided
	username := req.Username
	if username == "" {
		username = generateUsernameFromFullName(req.FullName)
	}

	// Generate password from FullName if not provided
	password := req.Password
	if password == "" {
		password = generatePasswordFromFullName(req.FullName)
	}

	// Create customer model
	customer := &model.Customer{
		FullName:  req.FullName,
		Phone:     req.Phone,
		Email:     req.Email,
		Address:   req.Address,
		Latitude:  req.Latitude,
		Longitude: req.Longitude,
	}

	// Create subscription model
	subscription := &model.Subscription{
		PlanID:   req.PlanID,
		RouterID: req.RouterID,
		Username: username,
		Password: password,
		StaticIP: req.StaticIP,
		Status:   "pending",
	}

	// Create both customer and subscription
	createdCustomer, createdSubscription, err := h.service.CreateWithSubscription(
		c.Request.Context(),
		customer,
		subscription,
	)
	if err != nil {
		response.Error(c, http.StatusBadRequest, err.Error())
		return
	}

	response.Created(c, gin.H{
		"customer":     createdCustomer,
		"subscription": createdSubscription,
	})
}

// generateUsernameFromFullName generates username from full name
func generateUsernameFromFullName(fullName string) string {
	// Lowercase and replace spaces with -
	username := strings.ToLower(fullName)
	username = strings.ReplaceAll(username, " ", "-")

	// Remove special characters (keep only alphanumeric and -)
	var result strings.Builder
	for _, char := range username {
		if (char >= 'a' && char <= 'z') || (char >= '0' && char <= '9') || char == '-' {
			result.WriteRune(char)
		}
	}

	return result.String()
}

// generatePasswordFromFullName generates password from full name
func generatePasswordFromFullName(fullName string) string {
	// Same logic as username
	return generateUsernameFromFullName(fullName)
}

// Get handles getting a customer by ID
func (h *CustomerHandler) Get(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		response.BadRequest(c, "invalid id")
		return
	}
	customer, err := h.service.GetByID(c.Request.Context(), id)
	if err != nil {
		response.NotFound(c, err.Error())
		return
	}
	response.OK(c, customer)
}

// List handles listing customers
func (h *CustomerHandler) List(c *gin.Context) {
	limit, offset := getPagination(c)
	customers, total, err := h.service.List(c.Request.Context(), limit, offset)
	if err != nil {
		response.InternalServerError(c, err.Error())
		return
	}
	response.WithMeta(c, http.StatusOK, customers, &response.Meta{
		Total:  total,
		Limit:  limit,
		Offset: offset,
	})
}

// Update handles updating a customer
func (h *CustomerHandler) Update(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		response.BadRequest(c, "invalid id")
		return
	}
	var customer model.Customer
	if err := c.ShouldBindJSON(&customer); err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	customer.ID = id.String()
	if err := h.service.Update(c.Request.Context(), &customer); err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}
	response.OK(c, customer)
}

// Delete handles deleting a customer
func (h *CustomerHandler) Delete(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		response.BadRequest(c, "invalid id")
		return
	}
	if err := h.service.Delete(c.Request.Context(), id); err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}
	response.OK(c, gin.H{"message": "customer deleted"})
}

// ActivateAccount handles activating a customer account
func (h *CustomerHandler) ActivateAccount(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		response.BadRequest(c, "invalid id")
		return
	}
	if err := h.service.ActivateAccount(c.Request.Context(), id); err != nil {
		response.Error(c, http.StatusBadRequest, err.Error())
		return
	}
	response.OK(c, gin.H{"message": "customer account activated"})
}

// DeactivateAccount handles deactivating a customer account
func (h *CustomerHandler) DeactivateAccount(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		response.BadRequest(c, "invalid id")
		return
	}
	if err := h.service.DeactivateAccount(c.Request.Context(), id); err != nil {
		response.Error(c, http.StatusBadRequest, err.Error())
		return
	}
	response.OK(c, gin.H{"message": "customer account deactivated"})
}
