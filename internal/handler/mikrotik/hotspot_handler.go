package mikrotik

import (
	mkdomain "github.com/Butterfly-Student/go-ros/domain"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	svcmikrotik "mikmongo/internal/service/mikrotik"
	"mikmongo/pkg/response"
)

// HotspotHandler handles Hotspot-related REST endpoints.
type HotspotHandler struct {
	hotspotSvc *svcmikrotik.HotspotService
}

func NewHotspotHandler(hotspotSvc *svcmikrotik.HotspotService) *HotspotHandler {
	return &HotspotHandler{hotspotSvc: hotspotSvc}
}

func (h *HotspotHandler) GetProfiles(c *gin.Context) {
	routerID, err := uuid.Parse(c.Param("router_id"))
	if err != nil {
		response.BadRequest(c, "invalid router_id")
		return
	}
	profiles, err := h.hotspotSvc.GetProfiles(c.Request.Context(), routerID)
	if err != nil {
		response.InternalServerError(c, err.Error())
		return
	}
	response.OK(c, profiles)
}

func (h *HotspotHandler) AddProfile(c *gin.Context) {
	routerID, err := uuid.Parse(c.Param("router_id"))
	if err != nil {
		response.BadRequest(c, "invalid router_id")
		return
	}
	var req mkdomain.UserProfile
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	id, err := h.hotspotSvc.AddProfile(c.Request.Context(), routerID, &req)
	if err != nil {
		response.InternalServerError(c, err.Error())
		return
	}
	response.Created(c, gin.H{"id": id})
}

func (h *HotspotHandler) GetProfileByName(c *gin.Context) {
	routerID, err := uuid.Parse(c.Param("router_id"))
	if err != nil {
		response.BadRequest(c, "invalid router_id")
		return
	}
	name := c.Param("name")
	profile, err := h.hotspotSvc.GetProfileByName(c.Request.Context(), routerID, name)
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

func (h *HotspotHandler) RemoveProfile(c *gin.Context) {
	routerID, err := uuid.Parse(c.Param("router_id"))
	if err != nil {
		response.BadRequest(c, "invalid router_id")
		return
	}
	id := c.Param("id")
	if err := h.hotspotSvc.RemoveProfile(c.Request.Context(), routerID, id); err != nil {
		response.InternalServerError(c, err.Error())
		return
	}
	response.OK(c, gin.H{"message": "profile removed"})
}

func (h *HotspotHandler) GetUsers(c *gin.Context) {
	routerID, err := uuid.Parse(c.Param("router_id"))
	if err != nil {
		response.BadRequest(c, "invalid router_id")
		return
	}
	profile := c.Query("profile")
	users, err := h.hotspotSvc.GetUsers(c.Request.Context(), routerID, profile)
	if err != nil {
		response.InternalServerError(c, err.Error())
		return
	}
	response.OK(c, users)
}

func (h *HotspotHandler) AddUser(c *gin.Context) {
	routerID, err := uuid.Parse(c.Param("router_id"))
	if err != nil {
		response.BadRequest(c, "invalid router_id")
		return
	}
	var req mkdomain.HotspotUser
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	id, err := h.hotspotSvc.AddUser(c.Request.Context(), routerID, &req)
	if err != nil {
		response.InternalServerError(c, err.Error())
		return
	}
	response.Created(c, gin.H{"id": id})
}

func (h *HotspotHandler) GetUserByName(c *gin.Context) {
	routerID, err := uuid.Parse(c.Param("router_id"))
	if err != nil {
		response.BadRequest(c, "invalid router_id")
		return
	}
	name := c.Param("name")
	user, err := h.hotspotSvc.GetUserByName(c.Request.Context(), routerID, name)
	if err != nil {
		response.InternalServerError(c, err.Error())
		return
	}
	if user == nil {
		response.NotFound(c, "user not found")
		return
	}
	response.OK(c, user)
}

func (h *HotspotHandler) RemoveUser(c *gin.Context) {
	routerID, err := uuid.Parse(c.Param("router_id"))
	if err != nil {
		response.BadRequest(c, "invalid router_id")
		return
	}
	id := c.Param("id")
	if err := h.hotspotSvc.RemoveUser(c.Request.Context(), routerID, id); err != nil {
		response.InternalServerError(c, err.Error())
		return
	}
	response.OK(c, gin.H{"message": "user removed"})
}

func (h *HotspotHandler) GetActive(c *gin.Context) {
	routerID, err := uuid.Parse(c.Param("router_id"))
	if err != nil {
		response.BadRequest(c, "invalid router_id")
		return
	}
	active, err := h.hotspotSvc.GetActive(c.Request.Context(), routerID)
	if err != nil {
		response.InternalServerError(c, err.Error())
		return
	}
	response.OK(c, active)
}

func (h *HotspotHandler) GetHosts(c *gin.Context) {
	routerID, err := uuid.Parse(c.Param("router_id"))
	if err != nil {
		response.BadRequest(c, "invalid router_id")
		return
	}
	hosts, err := h.hotspotSvc.GetHosts(c.Request.Context(), routerID)
	if err != nil {
		response.InternalServerError(c, err.Error())
		return
	}
	response.OK(c, hosts)
}

func (h *HotspotHandler) GetServers(c *gin.Context) {
	routerID, err := uuid.Parse(c.Param("router_id"))
	if err != nil {
		response.BadRequest(c, "invalid router_id")
		return
	}
	servers, err := h.hotspotSvc.GetServers(c.Request.Context(), routerID)
	if err != nil {
		response.InternalServerError(c, err.Error())
		return
	}
	response.OK(c, servers)
}
