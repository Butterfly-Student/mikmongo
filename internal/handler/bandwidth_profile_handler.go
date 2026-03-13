package handler

import (
	"encoding/json"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"mikmongo/internal/model"
	"mikmongo/internal/service"
	"mikmongo/pkg/response"
)

// BandwidthProfileHandler handles bandwidth profile HTTP requests
type BandwidthProfileHandler struct {
	service *service.BandwidthProfileService
}

// NewBandwidthProfileHandler creates a new bandwidth profile handler
func NewBandwidthProfileHandler(svc *service.BandwidthProfileService) *BandwidthProfileHandler {
	return &BandwidthProfileHandler{service: svc}
}

// CreateRequest represents the request body for creating a bandwidth profile
type CreateProfileRequest struct {
	ProfileCode    string          `json:"profile_code" binding:"required"`
	Name           string          `json:"name" binding:"required"`
	Description    string          `json:"description,omitempty"`
	PriceMonthly   float64         `json:"price_monthly" binding:"required"`
	DownloadSpeed  int64           `json:"download_speed" binding:"required"`
	UploadSpeed    int64           `json:"upload_speed" binding:"required"`
	MikrotikConfig json.RawMessage `json:"mikrotik_config,omitempty"`
}

// Create handles bandwidth profile creation for a specific router
func (h *BandwidthProfileHandler) Create(c *gin.Context) {
	routerID, err := uuid.Parse(c.Param("router_id"))
	if err != nil {
		response.BadRequest(c, "invalid router_id")
		return
	}

	var req CreateProfileRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	profile := model.BandwidthProfile{
		RouterID:       routerID.String(),
		ProfileCode:    req.ProfileCode,
		Name:           req.Name,
		PriceMonthly:   req.PriceMonthly,
		DownloadSpeed:  req.DownloadSpeed,
		UploadSpeed:    req.UploadSpeed,
		MikrotikConfig: req.MikrotikConfig,
	}

	if req.Description != "" {
		profile.Description = &req.Description
	}

	if err := h.service.Create(c.Request.Context(), &profile); err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}
	response.Created(c, profile)
}

// Get handles getting a bandwidth profile by ID
func (h *BandwidthProfileHandler) Get(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		response.BadRequest(c, "invalid id")
		return
	}
	profile, err := h.service.GetByID(c.Request.Context(), id)
	if err != nil {
		response.NotFound(c, err.Error())
		return
	}
	response.OK(c, profile)
}

// List handles listing bandwidth profiles for a specific router
func (h *BandwidthProfileHandler) List(c *gin.Context) {
	routerID, err := uuid.Parse(c.Param("router_id"))
	if err != nil {
		response.BadRequest(c, "invalid router_id")
		return
	}

	limit, offset := getPagination(c)
	profiles, count, err := h.service.ListByRouterID(c.Request.Context(), routerID, limit, offset)
	if err != nil {
		response.InternalServerError(c, err.Error())
		return
	}
	response.WithMeta(c, http.StatusOK, profiles, &response.Meta{Total: count, Limit: limit, Offset: offset})
}

// Update handles updating a bandwidth profile
func (h *BandwidthProfileHandler) Update(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		response.BadRequest(c, "invalid id")
		return
	}
	profile, err := h.service.GetByID(c.Request.Context(), id)
	if err != nil {
		response.NotFound(c, err.Error())
		return
	}
	if err := c.ShouldBindJSON(profile); err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	if err := h.service.Update(c.Request.Context(), profile); err != nil {
		response.InternalServerError(c, err.Error())
		return
	}
	response.OK(c, profile)
}

// Delete handles deleting a bandwidth profile
func (h *BandwidthProfileHandler) Delete(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		response.BadRequest(c, "invalid id")
		return
	}
	if err := h.service.Delete(c.Request.Context(), id); err != nil {
		response.InternalServerError(c, err.Error())
		return
	}
	response.OK(c, gin.H{"message": "deleted"})
}
