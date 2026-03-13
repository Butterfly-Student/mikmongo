package mikrotik

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	mikrotiksvc "mikmongo/internal/service/mikrotik"
	"mikmongo/pkg/mikrotik/domain"
	"mikmongo/pkg/response"
)

// HotspotHandler handles MikroTik Hotspot HTTP requests
type HotspotHandler struct {
	service *mikrotiksvc.HotspotService
}

// NewHotspotHandler creates a new Hotspot handler
func NewHotspotHandler(service *mikrotiksvc.HotspotService) *HotspotHandler {
	return &HotspotHandler{service: service}
}

// getRouterID extracts router ID from context
func getRouterID(c *gin.Context) (uuid.UUID, error) {
	routerID, exists := c.Get("router_id")
	if !exists {
		return uuid.Nil, nil
	}
	return routerID.(uuid.UUID), nil
}

// ─── Users ────────────────────────────────────────────────────────────────────

// ListUsers handles listing hotspot users
func (h *HotspotHandler) ListUsers(c *gin.Context) {
	routerID, err := getRouterID(c)
	if err != nil {
		response.BadRequest(c, "router ID not found")
		return
	}

	profile := c.Query("profile")
	users, err := h.service.GetUsers(c.Request.Context(), routerID, profile)
	if err != nil {
		response.InternalServerError(c, err.Error())
		return
	}
	response.OK(c, users)
}

// GetUser handles getting a hotspot user by ID
func (h *HotspotHandler) GetUser(c *gin.Context) {
	routerID, err := getRouterID(c)
	if err != nil {
		response.BadRequest(c, "router ID not found")
		return
	}

	id := c.Param("id")
	user, err := h.service.GetUserByID(c.Request.Context(), routerID, id)
	if err != nil {
		response.NotFound(c, err.Error())
		return
	}
	response.OK(c, user)
}

// GetUserByName handles getting a hotspot user by name
func (h *HotspotHandler) GetUserByName(c *gin.Context) {
	routerID, err := getRouterID(c)
	if err != nil {
		response.BadRequest(c, "router ID not found")
		return
	}

	name := c.Query("name")
	if name == "" {
		response.BadRequest(c, "name query parameter is required")
		return
	}

	user, err := h.service.GetUserByName(c.Request.Context(), routerID, name)
	if err != nil {
		response.NotFound(c, err.Error())
		return
	}
	response.OK(c, user)
}

// GetUsersByComment handles getting hotspot users by comment
func (h *HotspotHandler) GetUsersByComment(c *gin.Context) {
	routerID, err := getRouterID(c)
	if err != nil {
		response.BadRequest(c, "router ID not found")
		return
	}

	comment := c.Query("comment")
	if comment == "" {
		response.BadRequest(c, "comment query parameter is required")
		return
	}

	users, err := h.service.GetUsersByComment(c.Request.Context(), routerID, comment)
	if err != nil {
		response.InternalServerError(c, err.Error())
		return
	}
	response.OK(c, users)
}

// CreateUser handles creating a new hotspot user
func (h *HotspotHandler) CreateUser(c *gin.Context) {
	routerID, err := getRouterID(c)
	if err != nil {
		response.BadRequest(c, "router ID not found")
		return
	}

	var user domain.HotspotUser
	if err := c.ShouldBindJSON(&user); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	id, err := h.service.AddUser(c.Request.Context(), routerID, &user)
	if err != nil {
		response.InternalServerError(c, err.Error())
		return
	}

	user.ID = id
	response.Created(c, user)
}

// UpdateUser handles updating a hotspot user
func (h *HotspotHandler) UpdateUser(c *gin.Context) {
	routerID, err := getRouterID(c)
	if err != nil {
		response.BadRequest(c, "router ID not found")
		return
	}

	id := c.Param("id")
	var user domain.HotspotUser
	if err := c.ShouldBindJSON(&user); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	if err := h.service.UpdateUser(c.Request.Context(), routerID, id, &user); err != nil {
		response.InternalServerError(c, err.Error())
		return
	}

	response.OK(c, gin.H{"message": "user updated"})
}

// DeleteUser handles deleting a hotspot user
func (h *HotspotHandler) DeleteUser(c *gin.Context) {
	routerID, err := getRouterID(c)
	if err != nil {
		response.BadRequest(c, "router ID not found")
		return
	}

	id := c.Param("id")
	if err := h.service.RemoveUser(c.Request.Context(), routerID, id); err != nil {
		response.InternalServerError(c, err.Error())
		return
	}

	response.OK(c, gin.H{"message": "user deleted"})
}

// DeleteUsersByComment handles deleting users by comment
func (h *HotspotHandler) DeleteUsersByComment(c *gin.Context) {
	routerID, err := getRouterID(c)
	if err != nil {
		response.BadRequest(c, "router ID not found")
		return
	}

	var req struct {
		Comment string `json:"comment" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	if err := h.service.RemoveUsersByComment(c.Request.Context(), routerID, req.Comment); err != nil {
		response.InternalServerError(c, err.Error())
		return
	}

	response.OK(c, gin.H{"message": "users deleted"})
}

// BatchDeleteUsers handles batch deleting users
func (h *HotspotHandler) BatchDeleteUsers(c *gin.Context) {
	routerID, err := getRouterID(c)
	if err != nil {
		response.BadRequest(c, "router ID not found")
		return
	}

	var req domain.BatchIDsRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	if err := h.service.RemoveUsers(c.Request.Context(), routerID, req.IDs); err != nil {
		response.InternalServerError(c, err.Error())
		return
	}

	response.OK(c, gin.H{"message": "users deleted"})
}

// EnableUser handles enabling a hotspot user
func (h *HotspotHandler) EnableUser(c *gin.Context) {
	routerID, err := getRouterID(c)
	if err != nil {
		response.BadRequest(c, "router ID not found")
		return
	}

	id := c.Param("id")
	if err := h.service.EnableUser(c.Request.Context(), routerID, id); err != nil {
		response.InternalServerError(c, err.Error())
		return
	}

	response.OK(c, gin.H{"message": "user enabled"})
}

// DisableUser handles disabling a hotspot user
func (h *HotspotHandler) DisableUser(c *gin.Context) {
	routerID, err := getRouterID(c)
	if err != nil {
		response.BadRequest(c, "router ID not found")
		return
	}

	id := c.Param("id")
	if err := h.service.DisableUser(c.Request.Context(), routerID, id); err != nil {
		response.InternalServerError(c, err.Error())
		return
	}

	response.OK(c, gin.H{"message": "user disabled"})
}

// BatchEnableUsers handles batch enabling users
func (h *HotspotHandler) BatchEnableUsers(c *gin.Context) {
	routerID, err := getRouterID(c)
	if err != nil {
		response.BadRequest(c, "router ID not found")
		return
	}

	var req domain.BatchIDsRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	if err := h.service.EnableUsers(c.Request.Context(), routerID, req.IDs); err != nil {
		response.InternalServerError(c, err.Error())
		return
	}

	response.OK(c, gin.H{"message": "users enabled"})
}

// BatchDisableUsers handles batch disabling users
func (h *HotspotHandler) BatchDisableUsers(c *gin.Context) {
	routerID, err := getRouterID(c)
	if err != nil {
		response.BadRequest(c, "router ID not found")
		return
	}

	var req domain.BatchIDsRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	if err := h.service.DisableUsers(c.Request.Context(), routerID, req.IDs); err != nil {
		response.InternalServerError(c, err.Error())
		return
	}

	response.OK(c, gin.H{"message": "users disabled"})
}

// ResetUserCounters handles resetting user counters
func (h *HotspotHandler) ResetUserCounters(c *gin.Context) {
	routerID, err := getRouterID(c)
	if err != nil {
		response.BadRequest(c, "router ID not found")
		return
	}

	id := c.Param("id")
	if err := h.service.ResetUserCounters(c.Request.Context(), routerID, id); err != nil {
		response.InternalServerError(c, err.Error())
		return
	}

	response.OK(c, gin.H{"message": "counters reset"})
}

// BatchResetUserCounters handles batch resetting user counters
func (h *HotspotHandler) BatchResetUserCounters(c *gin.Context) {
	routerID, err := getRouterID(c)
	if err != nil {
		response.BadRequest(c, "router ID not found")
		return
	}

	var req domain.BatchIDsRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	if err := h.service.ResetUserCountersMultiple(c.Request.Context(), routerID, req.IDs); err != nil {
		response.InternalServerError(c, err.Error())
		return
	}

	response.OK(c, gin.H{"message": "counters reset"})
}

// GetUsersCount handles getting the count of hotspot users
func (h *HotspotHandler) GetUsersCount(c *gin.Context) {
	routerID, err := getRouterID(c)
	if err != nil {
		response.BadRequest(c, "router ID not found")
		return
	}

	count, err := h.service.GetUsersCount(c.Request.Context(), routerID)
	if err != nil {
		response.InternalServerError(c, err.Error())
		return
	}

	response.OK(c, gin.H{"count": count})
}

// ─── Profiles ─────────────────────────────────────────────────────────────────

// ListProfiles handles listing hotspot profiles
func (h *HotspotHandler) ListProfiles(c *gin.Context) {
	routerID, err := getRouterID(c)
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

// GetProfile handles getting a hotspot profile by ID
func (h *HotspotHandler) GetProfile(c *gin.Context) {
	routerID, err := getRouterID(c)
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

// GetProfileByName handles getting a hotspot profile by name
func (h *HotspotHandler) GetProfileByName(c *gin.Context) {
	routerID, err := getRouterID(c)
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

// CreateProfile handles creating a new hotspot profile
func (h *HotspotHandler) CreateProfile(c *gin.Context) {
	routerID, err := getRouterID(c)
	if err != nil {
		response.BadRequest(c, "router ID not found")
		return
	}

	var profile domain.UserProfile
	if err := c.ShouldBindJSON(&profile); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	id, err := h.service.AddProfile(c.Request.Context(), routerID, &profile)
	if err != nil {
		response.InternalServerError(c, err.Error())
		return
	}

	profile.ID = id
	response.Created(c, profile)
}

// UpdateProfile handles updating a hotspot profile
func (h *HotspotHandler) UpdateProfile(c *gin.Context) {
	routerID, err := getRouterID(c)
	if err != nil {
		response.BadRequest(c, "router ID not found")
		return
	}

	id := c.Param("id")
	var profile domain.UserProfile
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

// DeleteProfile handles deleting a hotspot profile
func (h *HotspotHandler) DeleteProfile(c *gin.Context) {
	routerID, err := getRouterID(c)
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

// EnableProfile handles enabling a hotspot profile
func (h *HotspotHandler) EnableProfile(c *gin.Context) {
	routerID, err := getRouterID(c)
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

// DisableProfile handles disabling a hotspot profile
func (h *HotspotHandler) DisableProfile(c *gin.Context) {
	routerID, err := getRouterID(c)
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
func (h *HotspotHandler) BatchDeleteProfiles(c *gin.Context) {
	routerID, err := getRouterID(c)
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
func (h *HotspotHandler) BatchEnableProfiles(c *gin.Context) {
	routerID, err := getRouterID(c)
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
func (h *HotspotHandler) BatchDisableProfiles(c *gin.Context) {
	routerID, err := getRouterID(c)
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

// ListActive handles listing active hotspot sessions
func (h *HotspotHandler) ListActive(c *gin.Context) {
	routerID, err := getRouterID(c)
	if err != nil {
		response.BadRequest(c, "router ID not found")
		return
	}

	active, err := h.service.GetActive(c.Request.Context(), routerID)
	if err != nil {
		response.InternalServerError(c, err.Error())
		return
	}
	response.OK(c, active)
}

// GetActiveCount handles getting the count of active sessions
func (h *HotspotHandler) GetActiveCount(c *gin.Context) {
	routerID, err := getRouterID(c)
	if err != nil {
		response.BadRequest(c, "router ID not found")
		return
	}

	count, err := h.service.GetActiveCount(c.Request.Context(), routerID)
	if err != nil {
		response.InternalServerError(c, err.Error())
		return
	}

	response.OK(c, gin.H{"count": count})
}

// DeleteActive handles removing an active session
func (h *HotspotHandler) DeleteActive(c *gin.Context) {
	routerID, err := getRouterID(c)
	if err != nil {
		response.BadRequest(c, "router ID not found")
		return
	}

	id := c.Param("id")
	if err := h.service.RemoveActive(c.Request.Context(), routerID, id); err != nil {
		response.InternalServerError(c, err.Error())
		return
	}

	response.OK(c, gin.H{"message": "session removed"})
}

// BatchDeleteActive handles batch removing active sessions
func (h *HotspotHandler) BatchDeleteActive(c *gin.Context) {
	routerID, err := getRouterID(c)
	if err != nil {
		response.BadRequest(c, "router ID not found")
		return
	}

	var req domain.BatchIDsRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	if err := h.service.RemoveActives(c.Request.Context(), routerID, req.IDs); err != nil {
		response.InternalServerError(c, err.Error())
		return
	}

	response.OK(c, gin.H{"message": "sessions removed"})
}

// ─── Hosts ────────────────────────────────────────────────────────────────────

// ListHosts handles listing hotspot hosts
func (h *HotspotHandler) ListHosts(c *gin.Context) {
	routerID, err := getRouterID(c)
	if err != nil {
		response.BadRequest(c, "router ID not found")
		return
	}

	hosts, err := h.service.GetHosts(c.Request.Context(), routerID)
	if err != nil {
		response.InternalServerError(c, err.Error())
		return
	}
	response.OK(c, hosts)
}

// DeleteHost handles removing a hotspot host
func (h *HotspotHandler) DeleteHost(c *gin.Context) {
	routerID, err := getRouterID(c)
	if err != nil {
		response.BadRequest(c, "router ID not found")
		return
	}

	id := c.Param("id")
	if err := h.service.RemoveHost(c.Request.Context(), routerID, id); err != nil {
		response.InternalServerError(c, err.Error())
		return
	}

	response.OK(c, gin.H{"message": "host removed"})
}

// ─── Servers ──────────────────────────────────────────────────────────────────

// ListServers handles listing hotspot servers
func (h *HotspotHandler) ListServers(c *gin.Context) {
	routerID, err := getRouterID(c)
	if err != nil {
		response.BadRequest(c, "router ID not found")
		return
	}

	servers, err := h.service.GetServers(c.Request.Context(), routerID)
	if err != nil {
		response.InternalServerError(c, err.Error())
		return
	}
	response.OK(c, servers)
}

// ─── IP Bindings ──────────────────────────────────────────────────────────────

// ListIPBindings handles listing hotspot IP bindings
func (h *HotspotHandler) ListIPBindings(c *gin.Context) {
	routerID, err := getRouterID(c)
	if err != nil {
		response.BadRequest(c, "router ID not found")
		return
	}

	bindings, err := h.service.GetIPBindings(c.Request.Context(), routerID)
	if err != nil {
		response.InternalServerError(c, err.Error())
		return
	}
	response.OK(c, bindings)
}

// CreateIPBinding handles creating a new IP binding
func (h *HotspotHandler) CreateIPBinding(c *gin.Context) {
	routerID, err := getRouterID(c)
	if err != nil {
		response.BadRequest(c, "router ID not found")
		return
	}

	var binding domain.HotspotIPBinding
	if err := c.ShouldBindJSON(&binding); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	id, err := h.service.AddIPBinding(c.Request.Context(), routerID, &binding)
	if err != nil {
		response.InternalServerError(c, err.Error())
		return
	}

	binding.ID = id
	response.Created(c, binding)
}

// DeleteIPBinding handles deleting an IP binding
func (h *HotspotHandler) DeleteIPBinding(c *gin.Context) {
	routerID, err := getRouterID(c)
	if err != nil {
		response.BadRequest(c, "router ID not found")
		return
	}

	id := c.Param("id")
	if err := h.service.RemoveIPBinding(c.Request.Context(), routerID, id); err != nil {
		response.InternalServerError(c, err.Error())
		return
	}

	response.OK(c, gin.H{"message": "binding deleted"})
}

// EnableIPBinding handles enabling an IP binding
func (h *HotspotHandler) EnableIPBinding(c *gin.Context) {
	routerID, err := getRouterID(c)
	if err != nil {
		response.BadRequest(c, "router ID not found")
		return
	}

	id := c.Param("id")
	if err := h.service.EnableIPBinding(c.Request.Context(), routerID, id); err != nil {
		response.InternalServerError(c, err.Error())
		return
	}

	response.OK(c, gin.H{"message": "binding enabled"})
}

// DisableIPBinding handles disabling an IP binding
func (h *HotspotHandler) DisableIPBinding(c *gin.Context) {
	routerID, err := getRouterID(c)
	if err != nil {
		response.BadRequest(c, "router ID not found")
		return
	}

	id := c.Param("id")
	if err := h.service.DisableIPBinding(c.Request.Context(), routerID, id); err != nil {
		response.InternalServerError(c, err.Error())
		return
	}

	response.OK(c, gin.H{"message": "binding disabled"})
}
