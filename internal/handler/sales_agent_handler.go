package handler

import (
	"context"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
	"mikmongo/internal/model"
	"mikmongo/internal/repository"
	"mikmongo/pkg/response"
)

// SalesAgentHandler handles sales agent CRUD and profile price management.
type SalesAgentHandler struct {
	agentRepo   repository.SalesAgentRepository
	settingRepo repository.SystemSettingRepository
}

// NewSalesAgentHandler creates a new SalesAgentHandler.
func NewSalesAgentHandler(agentRepo repository.SalesAgentRepository, settingRepo repository.SystemSettingRepository) *SalesAgentHandler {
	return &SalesAgentHandler{agentRepo: agentRepo, settingRepo: settingRepo}
}

// getSetting reads a system setting value, returning fallback if not found.
func (h *SalesAgentHandler) getSetting(ctx context.Context, group, key, fallback string) string {
	s, err := h.settingRepo.GetByGroupAndKey(ctx, group, key)
	if err != nil || s.Value == nil || *s.Value == "" {
		return fallback
	}
	return *s.Value
}

// getSettingInt reads a system setting as int, returning fallback if not found.
func (h *SalesAgentHandler) getSettingInt(ctx context.Context, group, key string, fallback int) int {
	s, err := h.settingRepo.GetByGroupAndKey(ctx, group, key)
	if err != nil || s.Value == nil || *s.Value == "" {
		return fallback
	}
	var v int
	if _, err := fmt.Sscanf(*s.Value, "%d", &v); err != nil {
		return fallback
	}
	return v
}

type createSalesAgentRequest struct {
	RouterID      string  `json:"router_id" binding:"required,uuid"`
	Name          string  `json:"name" binding:"required"`
	Phone         *string `json:"phone"`
	Username      string  `json:"username" binding:"required"`
	Password      string  `json:"password" binding:"required,min=6"`
	Status        string  `json:"status"`
	VoucherMode   string  `json:"voucher_mode"`
	VoucherLength int     `json:"voucher_length"`
	VoucherType   string  `json:"voucher_type"`
	BillDiscount  float64 `json:"bill_discount"`
	BillingCycle  string  `json:"billing_cycle" binding:"omitempty,oneof=weekly monthly"`
	BillingDay    int     `json:"billing_day" binding:"omitempty,min=1,max=31"`
}

type updateSalesAgentRequest struct {
	Name          *string  `json:"name"`
	Phone         *string  `json:"phone"`
	Password      *string  `json:"password"`
	Status        *string  `json:"status"`
	VoucherMode   *string  `json:"voucher_mode"`
	VoucherLength *int     `json:"voucher_length"`
	VoucherType   *string  `json:"voucher_type"`
	BillDiscount  *float64 `json:"bill_discount"`
	BillingCycle  *string  `json:"billing_cycle" binding:"omitempty,oneof=weekly monthly"`
	BillingDay    *int     `json:"billing_day" binding:"omitempty,min=1,max=31"`
}

type upsertProfilePriceRequest struct {
	BasePrice     float64 `json:"base_price"`
	SellingPrice  float64 `json:"selling_price"`
	VoucherLength *int    `json:"voucher_length"`
	IsActive      *bool   `json:"is_active"`
}

// Create handles creating a new sales agent.
func (h *SalesAgentHandler) Create(c *gin.Context) {
	var req createSalesAgentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		response.InternalServerError(c, "failed to hash password")
		return
	}

	// Apply system-setting defaults for billing cycle and day when not provided
	ctx := c.Request.Context()
	billingCycle := req.BillingCycle
	if billingCycle == "" {
		billingCycle = h.getSetting(ctx, "billing", "agent_default_billing_cycle", "monthly")
	}
	billingDay := req.BillingDay
	if billingDay == 0 {
		key := "agent_default_billing_day_monthly"
		if billingCycle == "weekly" {
			key = "agent_default_billing_day_weekly"
		}
		billingDay = h.getSettingInt(ctx, "billing", key, 1)
	}

	agent := &model.SalesAgent{
		RouterID:     req.RouterID,
		Name:         req.Name,
		Phone:        req.Phone,
		Username:     req.Username,
		PasswordHash: string(hash),
		Status:       coalesce(req.Status, "active"),
		VoucherMode:  coalesce(req.VoucherMode, "mix"),
		VoucherType:  coalesce(req.VoucherType, "upp"),
		BillDiscount: req.BillDiscount,
		BillingCycle: billingCycle,
		BillingDay:   billingDay,
	}
	if req.VoucherLength > 0 {
		agent.VoucherLength = req.VoucherLength
	} else {
		agent.VoucherLength = 6
	}

	if err := h.agentRepo.Create(c.Request.Context(), agent); err != nil {
		response.InternalServerError(c, err.Error())
		return
	}
	agent.PasswordHash = "" // never expose hash
	response.Created(c, agent)
}

// Get handles fetching a single sales agent.
func (h *SalesAgentHandler) Get(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		response.BadRequest(c, "invalid id")
		return
	}
	agent, err := h.agentRepo.GetByID(c.Request.Context(), id)
	if err != nil {
		response.NotFound(c, err.Error())
		return
	}
	agent.PasswordHash = ""
	response.OK(c, agent)
}

// List handles listing sales agents with optional router_id filter.
func (h *SalesAgentHandler) List(c *gin.Context) {
	var routerID *uuid.UUID
	if rid := c.Query("router_id"); rid != "" {
		id, err := uuid.Parse(rid)
		if err != nil {
			response.BadRequest(c, "invalid router_id")
			return
		}
		routerID = &id
	}

	limit, offset := getPagination(c)
	agents, count, err := func() ([]model.SalesAgent, int64, error) {
		a, err := h.agentRepo.List(c.Request.Context(), routerID, limit, offset)
		if err != nil {
			return nil, 0, err
		}
		cnt, err := h.agentRepo.Count(c.Request.Context(), routerID)
		return a, cnt, err
	}()
	if err != nil {
		response.InternalServerError(c, err.Error())
		return
	}
	for i := range agents {
		agents[i].PasswordHash = ""
	}
	response.WithMeta(c, http.StatusOK, agents, &response.Meta{Total: count, Limit: limit, Offset: offset})
}

// Update handles updating a sales agent.
func (h *SalesAgentHandler) Update(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		response.BadRequest(c, "invalid id")
		return
	}

	agent, err := h.agentRepo.GetByID(c.Request.Context(), id)
	if err != nil {
		response.NotFound(c, err.Error())
		return
	}

	var req updateSalesAgentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	if req.Name != nil {
		agent.Name = *req.Name
	}
	agent.Phone = req.Phone
	if req.Status != nil {
		agent.Status = *req.Status
	}
	if req.VoucherMode != nil {
		agent.VoucherMode = *req.VoucherMode
	}
	if req.VoucherLength != nil {
		agent.VoucherLength = *req.VoucherLength
	}
	if req.VoucherType != nil {
		agent.VoucherType = *req.VoucherType
	}
	if req.BillDiscount != nil {
		agent.BillDiscount = *req.BillDiscount
	}
	if req.BillingCycle != nil {
		agent.BillingCycle = *req.BillingCycle
	}
	if req.BillingDay != nil {
		agent.BillingDay = *req.BillingDay
	}
	if req.Password != nil && *req.Password != "" {
		hash, err := bcrypt.GenerateFromPassword([]byte(*req.Password), bcrypt.DefaultCost)
		if err != nil {
			response.InternalServerError(c, "failed to hash password")
			return
		}
		agent.PasswordHash = string(hash)
	}

	if err := h.agentRepo.Update(c.Request.Context(), agent); err != nil {
		response.InternalServerError(c, err.Error())
		return
	}
	agent.PasswordHash = ""
	response.OK(c, agent)
}

// Delete handles soft-deleting a sales agent.
func (h *SalesAgentHandler) Delete(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		response.BadRequest(c, "invalid id")
		return
	}
	if err := h.agentRepo.Delete(c.Request.Context(), id); err != nil {
		response.InternalServerError(c, err.Error())
		return
	}
	response.OK(c, gin.H{"message": "deleted"})
}

// ListProfilePrices handles listing profile price overrides for a sales agent.
func (h *SalesAgentHandler) ListProfilePrices(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		response.BadRequest(c, "invalid id")
		return
	}
	prices, err := h.agentRepo.ListProfilePrices(c.Request.Context(), id)
	if err != nil {
		response.InternalServerError(c, err.Error())
		return
	}
	response.OK(c, prices)
}

// UpsertProfilePrice handles creating or updating a profile price override.
func (h *SalesAgentHandler) UpsertProfilePrice(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		response.BadRequest(c, "invalid id")
		return
	}
	profileName := c.Param("profile")
	if profileName == "" {
		response.BadRequest(c, "profile name required")
		return
	}

	var req upsertProfilePriceRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	isActive := true
	if req.IsActive != nil {
		isActive = *req.IsActive
	}

	price := &model.SalesProfilePrice{
		SalesAgentID:  id.String(),
		ProfileName:   profileName,
		BasePrice:     req.BasePrice,
		SellingPrice:  req.SellingPrice,
		VoucherLength: req.VoucherLength,
		IsActive:      isActive,
	}

	if err := h.agentRepo.UpsertProfilePrice(c.Request.Context(), price); err != nil {
		response.InternalServerError(c, err.Error())
		return
	}
	response.OK(c, price)
}

// coalesce returns the first non-empty string.
func coalesce(a, b string) string {
	if a != "" {
		return a
	}
	return b
}
