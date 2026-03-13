package mikrotik

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	mikrotiksvc "mikmongo/internal/service/mikrotik"
	"mikmongo/pkg/mikrotik/domain"
	"mikmongo/pkg/response"
)

// PPPHandler handles MikroTik PPP HTTP requests
type PPPHandler struct {
	service *mikrotiksvc.PPPService
}

// NewPPPHandler creates a new PPP handler
func NewPPPHandler(service *mikrotiksvc.PPPService) *PPPHandler {
	return &PPPHandler{service: service}
}

// getRouterIDPPP extracts router ID from context
func getRouterIDPPP(c *gin.Context) (uuid.UUID, error) {
	routerID, exists := c.Get("router_id")
	if !exists {
		return uuid.Nil, nil
	}
	return routerID.(uuid.UUID), nil
}

// ─── Secrets ──────────────────────────────────────────────────────────────────

// ListSecrets handles listing PPP secrets
func (h *PPPHandler) ListSecrets(c *gin.Context) {
	routerID, err := getRouterIDPPP(c)
	if err != nil {
		response.BadRequest(c, "router ID not found")
		return
	}

	profile := c.Query("profile")
	secrets, err := h.service.GetSecrets(c.Request.Context(), routerID, profile)
	if err != nil {
		response.InternalServerError(c, err.Error())
		return
	}
	response.OK(c, secrets)
}

// GetSecret handles getting a PPP secret by ID
func (h *PPPHandler) GetSecret(c *gin.Context) {
	routerID, err := getRouterIDPPP(c)
	if err != nil {
		response.BadRequest(c, "router ID not found")
		return
	}

	id := c.Param("id")
	secret, err := h.service.GetSecretByID(c.Request.Context(), routerID, id)
	if err != nil {
		response.NotFound(c, err.Error())
		return
	}
	response.OK(c, secret)
}

// GetSecretByName handles getting a PPP secret by name
func (h *PPPHandler) GetSecretByName(c *gin.Context) {
	routerID, err := getRouterIDPPP(c)
	if err != nil {
		response.BadRequest(c, "router ID not found")
		return
	}

	name := c.Query("name")
	if name == "" {
		response.BadRequest(c, "name query parameter is required")
		return
	}

	secret, err := h.service.GetSecretByName(c.Request.Context(), routerID, name)
	if err != nil {
		response.NotFound(c, err.Error())
		return
	}
	response.OK(c, secret)
}

// CreateSecret handles creating a new PPP secret
func (h *PPPHandler) CreateSecret(c *gin.Context) {
	routerID, err := getRouterIDPPP(c)
	if err != nil {
		response.BadRequest(c, "router ID not found")
		return
	}

	var secret domain.PPPSecret
	if err := c.ShouldBindJSON(&secret); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	if err := h.service.AddSecret(c.Request.Context(), routerID, &secret); err != nil {
		response.InternalServerError(c, err.Error())
		return
	}

	response.Created(c, secret)
}

// UpdateSecret handles updating a PPP secret
func (h *PPPHandler) UpdateSecret(c *gin.Context) {
	routerID, err := getRouterIDPPP(c)
	if err != nil {
		response.BadRequest(c, "router ID not found")
		return
	}

	id := c.Param("id")
	var secret domain.PPPSecret
	if err := c.ShouldBindJSON(&secret); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	if err := h.service.UpdateSecret(c.Request.Context(), routerID, id, &secret); err != nil {
		response.InternalServerError(c, err.Error())
		return
	}

	response.OK(c, gin.H{"message": "secret updated"})
}

// DeleteSecret handles deleting a PPP secret
func (h *PPPHandler) DeleteSecret(c *gin.Context) {
	routerID, err := getRouterIDPPP(c)
	if err != nil {
		response.BadRequest(c, "router ID not found")
		return
	}

	id := c.Param("id")
	if err := h.service.RemoveSecret(c.Request.Context(), routerID, id); err != nil {
		response.InternalServerError(c, err.Error())
		return
	}

	response.OK(c, gin.H{"message": "secret deleted"})
}

// EnableSecret handles enabling a PPP secret
func (h *PPPHandler) EnableSecret(c *gin.Context) {
	routerID, err := getRouterIDPPP(c)
	if err != nil {
		response.BadRequest(c, "router ID not found")
		return
	}

	id := c.Param("id")
	if err := h.service.EnableSecret(c.Request.Context(), routerID, id); err != nil {
		response.InternalServerError(c, err.Error())
		return
	}

	response.OK(c, gin.H{"message": "secret enabled"})
}

// DisableSecret handles disabling a PPP secret
func (h *PPPHandler) DisableSecret(c *gin.Context) {
	routerID, err := getRouterIDPPP(c)
	if err != nil {
		response.BadRequest(c, "router ID not found")
		return
	}

	id := c.Param("id")
	if err := h.service.DisableSecret(c.Request.Context(), routerID, id); err != nil {
		response.InternalServerError(c, err.Error())
		return
	}

	response.OK(c, gin.H{"message": "secret disabled"})
}

// BatchDeleteSecrets handles batch deleting secrets
func (h *PPPHandler) BatchDeleteSecrets(c *gin.Context) {
	routerID, err := getRouterIDPPP(c)
	if err != nil {
		response.BadRequest(c, "router ID not found")
		return
	}

	var req domain.BatchIDsRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	if err := h.service.RemoveSecrets(c.Request.Context(), routerID, req.IDs); err != nil {
		response.InternalServerError(c, err.Error())
		return
	}

	response.OK(c, gin.H{"message": "secrets deleted"})
}

// BatchEnableSecrets handles batch enabling secrets
func (h *PPPHandler) BatchEnableSecrets(c *gin.Context) {
	routerID, err := getRouterIDPPP(c)
	if err != nil {
		response.BadRequest(c, "router ID not found")
		return
	}

	var req domain.BatchIDsRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	if err := h.service.EnableSecrets(c.Request.Context(), routerID, req.IDs); err != nil {
		response.InternalServerError(c, err.Error())
		return
	}

	response.OK(c, gin.H{"message": "secrets enabled"})
}

// BatchDisableSecrets handles batch disabling secrets
func (h *PPPHandler) BatchDisableSecrets(c *gin.Context) {
	routerID, err := getRouterIDPPP(c)
	if err != nil {
		response.BadRequest(c, "router ID not found")
		return
	}

	var req domain.BatchIDsRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	if err := h.service.DisableSecrets(c.Request.Context(), routerID, req.IDs); err != nil {
		response.InternalServerError(c, err.Error())
		return
	}

	response.OK(c, gin.H{"message": "secrets disabled"})
}

// ─── Profiles ─────────────────────────────────────────────────────────────────

// ListProfiles handles listing PPP profiles
func (h *PPPHandler) ListProfiles(c *gin.Context) {
	routerID, err := getRouterIDPPP(c)
	if err != nil {
		response.BadRequest(c, "router ID not found")
		return
	}

	profiles, err := h.service.GetProfiles(c.Request.Context(), routerID)
	if err != nil {
		response.InternalServerError(c, err.Error())
		return
	}
	response.OK(c, profiles)
}

// GetProfile handles getting a PPP profile by ID
func (h *PPPHandler) GetProfile(c *gin.Context) {
	routerID, err := getRouterIDPPP(c)
	if err != nil {
		response.BadRequest(c, "router ID not found")
		return
	}

	id := c.Param("id")
	profile, err := h.service.GetProfileByID(c.Request.Context(), routerID, id)
	if err != nil {
		response.NotFound(c, err.Error())
		return
	}
	response.OK(c, profile)
}

// GetProfileByName handles getting a PPP profile by name
func (h *PPPHandler) GetProfileByName(c *gin.Context) {
	routerID, err := getRouterIDPPP(c)
	if err != nil {
		response.BadRequest(c, "router ID not found")
		return
	}

	name := c.Query("name")
	if name == "" {
		response.BadRequest(c, "name query parameter is required")
		return
	}

	profile, err := h.service.GetProfileByName(c.Request.Context(), routerID, name)
	if err != nil {
		response.NotFound(c, err.Error())
		return
	}
	response.OK(c, profile)
}

// CreateProfile handles creating a new PPP profile
func (h *PPPHandler) CreateProfile(c *gin.Context) {
	routerID, err := getRouterIDPPP(c)
	if err != nil {
		response.BadRequest(c, "router ID not found")
		return
	}

	var profile domain.PPPProfile
	if err := c.ShouldBindJSON(&profile); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	if err := h.service.AddProfile(c.Request.Context(), routerID, &profile); err != nil {
		response.InternalServerError(c, err.Error())
		return
	}

	response.Created(c, profile)
}

// UpdateProfile handles updating a PPP profile
func (h *PPPHandler) UpdateProfile(c *gin.Context) {
	routerID, err := getRouterIDPPP(c)
	if err != nil {
		response.BadRequest(c, "router ID not found")
		return
	}

	id := c.Param("id")
	var profile domain.PPPProfile
	if err := c.ShouldBindJSON(&profile); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	if err := h.service.UpdateProfile(c.Request.Context(), routerID, id, &profile); err != nil {
		response.InternalServerError(c, err.Error())
		return
	}

	response.OK(c, gin.H{"message": "profile updated"})
}

// DeleteProfile handles deleting a PPP profile
func (h *PPPHandler) DeleteProfile(c *gin.Context) {
	routerID, err := getRouterIDPPP(c)
	if err != nil {
		response.BadRequest(c, "router ID not found")
		return
	}

	id := c.Param("id")
	if err := h.service.RemoveProfile(c.Request.Context(), routerID, id); err != nil {
		response.InternalServerError(c, err.Error())
		return
	}

	response.OK(c, gin.H{"message": "profile deleted"})
}

// EnableProfile handles enabling a PPP profile
func (h *PPPHandler) EnableProfile(c *gin.Context) {
	routerID, err := getRouterIDPPP(c)
	if err != nil {
		response.BadRequest(c, "router ID not found")
		return
	}

	id := c.Param("id")
	if err := h.service.EnableProfile(c.Request.Context(), routerID, id); err != nil {
		response.InternalServerError(c, err.Error())
		return
	}

	response.OK(c, gin.H{"message": "profile enabled"})
}

// DisableProfile handles disabling a PPP profile
func (h *PPPHandler) DisableProfile(c *gin.Context) {
	routerID, err := getRouterIDPPP(c)
	if err != nil {
		response.BadRequest(c, "router ID not found")
		return
	}

	id := c.Param("id")
	if err := h.service.DisableProfile(c.Request.Context(), routerID, id); err != nil {
		response.InternalServerError(c, err.Error())
		return
	}

	response.OK(c, gin.H{"message": "profile disabled"})
}

// BatchDeleteProfiles handles batch deleting profiles
func (h *PPPHandler) BatchDeleteProfiles(c *gin.Context) {
	routerID, err := getRouterIDPPP(c)
	if err != nil {
		response.BadRequest(c, "router ID not found")
		return
	}

	var req domain.BatchIDsRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	if err := h.service.RemoveProfiles(c.Request.Context(), routerID, req.IDs); err != nil {
		response.InternalServerError(c, err.Error())
		return
	}

	response.OK(c, gin.H{"message": "profiles deleted"})
}

// BatchEnableProfiles handles batch enabling profiles
func (h *PPPHandler) BatchEnableProfiles(c *gin.Context) {
	routerID, err := getRouterIDPPP(c)
	if err != nil {
		response.BadRequest(c, "router ID not found")
		return
	}

	var req domain.BatchIDsRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	if err := h.service.EnableProfiles(c.Request.Context(), routerID, req.IDs); err != nil {
		response.InternalServerError(c, err.Error())
		return
	}

	response.OK(c, gin.H{"message": "profiles enabled"})
}

// BatchDisableProfiles handles batch disabling profiles
func (h *PPPHandler) BatchDisableProfiles(c *gin.Context) {
	routerID, err := getRouterIDPPP(c)
	if err != nil {
		response.BadRequest(c, "router ID not found")
		return
	}

	var req domain.BatchIDsRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	if err := h.service.DisableProfiles(c.Request.Context(), routerID, req.IDs); err != nil {
		response.InternalServerError(c, err.Error())
		return
	}

	response.OK(c, gin.H{"message": "profiles disabled"})
}

// ─── Active Sessions ──────────────────────────────────────────────────────────

// ListActive handles listing active PPP sessions
func (h *PPPHandler) ListActive(c *gin.Context) {
	routerID, err := getRouterIDPPP(c)
	if err != nil {
		response.BadRequest(c, "router ID not found")
		return
	}

	service := c.Query("service")
	active, err := h.service.GetActiveUsers(c.Request.Context(), routerID, service)
	if err != nil {
		response.InternalServerError(c, err.Error())
		return
	}
	response.OK(c, active)
}

// GetActive handles getting an active PPP session by ID
func (h *PPPHandler) GetActive(c *gin.Context) {
	routerID, err := getRouterIDPPP(c)
	if err != nil {
		response.BadRequest(c, "router ID not found")
		return
	}

	id := c.Param("id")
	active, err := h.service.GetActiveByID(c.Request.Context(), routerID, id)
	if err != nil {
		response.NotFound(c, err.Error())
		return
	}
	response.OK(c, active)
}

// DisconnectActive handles disconnecting an active PPP session
func (h *PPPHandler) DisconnectActive(c *gin.Context) {
	routerID, err := getRouterIDPPP(c)
	if err != nil {
		response.BadRequest(c, "router ID not found")
		return
	}

	id := c.Param("id")
	if err := h.service.DisconnectActive(c.Request.Context(), routerID, id); err != nil {
		response.InternalServerError(c, err.Error())
		return
	}

	response.OK(c, gin.H{"message": "session disconnected"})
}

// BatchDisconnectActive handles batch disconnecting active sessions
func (h *PPPHandler) BatchDisconnectActive(c *gin.Context) {
	routerID, err := getRouterIDPPP(c)
	if err != nil {
		response.BadRequest(c, "router ID not found")
		return
	}

	var req domain.BatchIDsRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	if err := h.service.DisconnectActives(c.Request.Context(), routerID, req.IDs); err != nil {
		response.InternalServerError(c, err.Error())
		return
	}

	response.OK(c, gin.H{"message": "sessions disconnected"})
}
