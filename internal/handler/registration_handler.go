package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"mikmongo/internal/dto"
	"mikmongo/internal/service"
	"mikmongo/pkg/response"
)

// RegistrationHandler handles customer registration HTTP requests
type RegistrationHandler struct {
	service *service.RegistrationService
}

// NewRegistrationHandler creates a new registration handler
func NewRegistrationHandler(svc *service.RegistrationService) *RegistrationHandler {
	return &RegistrationHandler{service: svc}
}

// Create handles creating a registration request
func (h *RegistrationHandler) Create(c *gin.Context) {
	var req dto.CreateRegistrationRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	reg := req.ToModel()
	if err := h.service.Create(c.Request.Context(), reg); err != nil {
		response.Error(c, http.StatusBadRequest, err.Error())
		return
	}
	response.Created(c, dto.RegistrationToResponse(reg))
}

// Get handles getting a registration by ID
func (h *RegistrationHandler) Get(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		response.BadRequest(c, "invalid id")
		return
	}
	reg, err := h.service.GetByID(c.Request.Context(), id)
	if err != nil {
		response.NotFound(c, err.Error())
		return
	}
	response.OK(c, dto.RegistrationToResponse(reg))
}

// List handles listing registrations
func (h *RegistrationHandler) List(c *gin.Context) {
	limit, offset := getPagination(c)
	regs, count, err := h.service.List(c.Request.Context(), limit, offset)
	if err != nil {
		response.InternalServerError(c, err.Error())
		return
	}
	response.WithMeta(c, http.StatusOK, dto.RegistrationsToResponse(regs), &response.Meta{Total: count, Limit: limit, Offset: offset})
}

// Approve handles approving a registration
func (h *RegistrationHandler) Approve(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		response.BadRequest(c, "invalid id")
		return
	}
	var req dto.ApproveRegistrationRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	approvedBy, _ := c.Get("user_id")
	approvedByStr, _ := approvedBy.(string)
	if err := h.service.Approve(c.Request.Context(), id, approvedByStr, req.RouterID, req.ProfileID); err != nil {
		response.Error(c, http.StatusBadRequest, err.Error())
		return
	}
	response.OK(c, gin.H{"message": "approved"})
}

// Reject handles rejecting a registration
func (h *RegistrationHandler) Reject(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		response.BadRequest(c, "invalid id")
		return
	}
	var req dto.RejectRegistrationRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	rejectedBy, _ := c.Get("user_id")
	rejectedByStr, _ := rejectedBy.(string)
	if err := h.service.Reject(c.Request.Context(), id, req.Reason, rejectedByStr); err != nil {
		response.Error(c, http.StatusBadRequest, err.Error())
		return
	}
	response.OK(c, gin.H{"message": "rejected"})
}
