package handler

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"mikmongo/internal/dto"
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

// Create handles customer creation, with optional subscription
func (h *CustomerHandler) Create(c *gin.Context) {
	var req dto.CreateCustomerRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	customer := req.ToCustomerModel()

	// If plan_id is provided, create customer with subscription
	if req.PlanID != "" {
		username := req.Username
		if username == "" {
			username = generateUsernameFromFullName(req.FullName)
		}
		password := req.Password
		if password == "" {
			password = generatePasswordFromFullName(req.FullName)
		}

		subscription := &model.Subscription{
			PlanID:   req.PlanID,
			RouterID: req.RouterID,
			Username: username,
			Password: password,
			StaticIP: req.StaticIP,
			Status:   "pending",
		}

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
			"customer":     dto.CustomerToResponse(createdCustomer),
			"subscription": dto.SubscriptionToResponse(createdSubscription, nil),
		})
		return
	}

	// Create customer only (no subscription)
	if err := h.service.Create(c.Request.Context(), customer); err != nil {
		response.Error(c, http.StatusBadRequest, err.Error())
		return
	}
	response.Created(c, gin.H{"customer": dto.CustomerToResponse(customer)})
}

// generateUsernameFromFullName generates username from full name
func generateUsernameFromFullName(fullName string) string {
	username := strings.ToLower(fullName)
	username = strings.ReplaceAll(username, " ", "-")
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
	response.OK(c, dto.CustomerToResponse(customer))
}

// List handles listing customers
func (h *CustomerHandler) List(c *gin.Context) {
	limit, offset := getPagination(c)
	customers, total, err := h.service.List(c.Request.Context(), limit, offset)
	if err != nil {
		response.InternalServerError(c, err.Error())
		return
	}
	response.WithMeta(c, http.StatusOK, dto.CustomersToResponse(customers), &response.Meta{
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
	customer, err := h.service.GetByID(c.Request.Context(), id)
	if err != nil {
		response.NotFound(c, err.Error())
		return
	}

	var req dto.UpdateCustomerRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	req.ApplyTo(customer)

	if err := h.service.Update(c.Request.Context(), customer); err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}
	response.OK(c, dto.CustomerToResponse(customer))
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
