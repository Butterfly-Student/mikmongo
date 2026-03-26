package mikrotik

import (
	mkdomain "github.com/Butterfly-Student/go-ros/domain"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	svcmikrotik "mikmongo/internal/service/mikrotik"
	"mikmongo/pkg/response"
)

// PPPHandler handles PPP-related REST endpoints.
type PPPHandler struct {
	pppSvc *svcmikrotik.PPPService
}

func NewPPPHandler(pppSvc *svcmikrotik.PPPService) *PPPHandler {
	return &PPPHandler{pppSvc: pppSvc}
}

func (h *PPPHandler) GetProfiles(c *gin.Context) {
	routerID, err := uuid.Parse(c.Param("router_id"))
	if err != nil {
		response.BadRequest(c, "invalid router_id")
		return
	}
	profiles, err := h.pppSvc.GetProfiles(c.Request.Context(), routerID)
	if err != nil {
		response.InternalServerError(c, err.Error())
		return
	}
	response.OK(c, profiles)
}

func (h *PPPHandler) AddProfile(c *gin.Context) {
	routerID, err := uuid.Parse(c.Param("router_id"))
	if err != nil {
		response.BadRequest(c, "invalid router_id")
		return
	}
	var req mkdomain.PPPProfile
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	if err := h.pppSvc.AddProfile(c.Request.Context(), routerID, &req); err != nil {
		response.InternalServerError(c, err.Error())
		return
	}
	response.Created(c, gin.H{"message": "profile created"})
}

func (h *PPPHandler) GetProfileByName(c *gin.Context) {
	routerID, err := uuid.Parse(c.Param("router_id"))
	if err != nil {
		response.BadRequest(c, "invalid router_id")
		return
	}
	name := c.Param("name")
	profile, err := h.pppSvc.GetProfileByName(c.Request.Context(), routerID, name)
	if err != nil {
		response.InternalServerError(c, err.Error())
		return
	}
	if profile == nil {
		response.NotFound(c, "profile not found")
		return
	}
	response.OK(c, profile)
}

func (h *PPPHandler) RemoveProfile(c *gin.Context) {
	routerID, err := uuid.Parse(c.Param("router_id"))
	if err != nil {
		response.BadRequest(c, "invalid router_id")
		return
	}
	id := c.Param("id")
	if err := h.pppSvc.RemoveProfile(c.Request.Context(), routerID, id); err != nil {
		response.InternalServerError(c, err.Error())
		return
	}
	response.OK(c, gin.H{"message": "profile removed"})
}

func (h *PPPHandler) GetSecrets(c *gin.Context) {
	routerID, err := uuid.Parse(c.Param("router_id"))
	if err != nil {
		response.BadRequest(c, "invalid router_id")
		return
	}
	profile := c.Query("profile")
	secrets, err := h.pppSvc.GetSecrets(c.Request.Context(), routerID, profile)
	if err != nil {
		response.InternalServerError(c, err.Error())
		return
	}
	response.OK(c, secrets)
}

func (h *PPPHandler) AddSecret(c *gin.Context) {
	routerID, err := uuid.Parse(c.Param("router_id"))
	if err != nil {
		response.BadRequest(c, "invalid router_id")
		return
	}
	var req mkdomain.PPPSecret
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	if err := h.pppSvc.AddSecret(c.Request.Context(), routerID, &req); err != nil {
		response.InternalServerError(c, err.Error())
		return
	}
	response.Created(c, gin.H{"message": "secret created"})
}

func (h *PPPHandler) GetSecretByName(c *gin.Context) {
	routerID, err := uuid.Parse(c.Param("router_id"))
	if err != nil {
		response.BadRequest(c, "invalid router_id")
		return
	}
	name := c.Param("name")
	secret, err := h.pppSvc.GetSecretByName(c.Request.Context(), routerID, name)
	if err != nil {
		response.InternalServerError(c, err.Error())
		return
	}
	if secret == nil {
		response.NotFound(c, "secret not found")
		return
	}
	response.OK(c, secret)
}

func (h *PPPHandler) RemoveSecret(c *gin.Context) {
	routerID, err := uuid.Parse(c.Param("router_id"))
	if err != nil {
		response.BadRequest(c, "invalid router_id")
		return
	}
	id := c.Param("id")
	if err := h.pppSvc.RemoveSecret(c.Request.Context(), routerID, id); err != nil {
		response.InternalServerError(c, err.Error())
		return
	}
	response.OK(c, gin.H{"message": "secret removed"})
}

func (h *PPPHandler) GetActive(c *gin.Context) {
	routerID, err := uuid.Parse(c.Param("router_id"))
	if err != nil {
		response.BadRequest(c, "invalid router_id")
		return
	}
	service := c.Query("service")
	active, err := h.pppSvc.GetActiveUsers(c.Request.Context(), routerID, service)
	if err != nil {
		response.InternalServerError(c, err.Error())
		return
	}
	response.OK(c, active)
}
