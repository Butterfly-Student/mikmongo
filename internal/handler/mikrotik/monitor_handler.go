package mikrotik

import (
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	mikrotiksvc "mikmongo/internal/service/mikrotik"
	"mikmongo/pkg/mikrotik/domain"
	"mikmongo/pkg/response"
)

// MonitorHandler handles MikroTik Monitor HTTP requests
type MonitorHandler struct {
	service *mikrotiksvc.MonitorService
}

// NewMonitorHandler creates a new Monitor handler
func NewMonitorHandler(service *mikrotiksvc.MonitorService) *MonitorHandler {
	return &MonitorHandler{service: service}
}

// getRouterIDMonitor extracts router ID from context
func getRouterIDMonitor(c *gin.Context) (uuid.UUID, error) {
	routerID, exists := c.Get("router_id")
	if !exists {
		return uuid.Nil, nil
	}
	return routerID.(uuid.UUID), nil
}

// GetSystemResource handles getting system resource information
func (h *MonitorHandler) GetSystemResource(c *gin.Context) {
	routerID, err := getRouterIDMonitor(c)
	if err != nil {
		response.BadRequest(c, "router ID not found")
		return
	}

	resource, err := h.service.GetSystemResource(c.Request.Context(), routerID)
	if err != nil {
		response.InternalServerError(c, err.Error())
		return
	}
	response.OK(c, resource)
}

// GetSystemHealth handles getting system health information
func (h *MonitorHandler) GetSystemHealth(c *gin.Context) {
	routerID, err := getRouterIDMonitor(c)
	if err != nil {
		response.BadRequest(c, "router ID not found")
		return
	}

	health, err := h.service.GetSystemHealth(c.Request.Context(), routerID)
	if err != nil {
		response.InternalServerError(c, err.Error())
		return
	}
	response.OK(c, health)
}

// GetSystemIdentity handles getting system identity
func (h *MonitorHandler) GetSystemIdentity(c *gin.Context) {
	routerID, err := getRouterIDMonitor(c)
	if err != nil {
		response.BadRequest(c, "router ID not found")
		return
	}

	identity, err := h.service.GetSystemIdentity(c.Request.Context(), routerID)
	if err != nil {
		response.InternalServerError(c, err.Error())
		return
	}
	response.OK(c, identity)
}

// GetSystemClock handles getting system clock
func (h *MonitorHandler) GetSystemClock(c *gin.Context) {
	routerID, err := getRouterIDMonitor(c)
	if err != nil {
		response.BadRequest(c, "router ID not found")
		return
	}

	clock, err := h.service.GetSystemClock(c.Request.Context(), routerID)
	if err != nil {
		response.InternalServerError(c, err.Error())
		return
	}
	response.OK(c, clock)
}

// GetRouterBoardInfo handles getting routerboard information
func (h *MonitorHandler) GetRouterBoardInfo(c *gin.Context) {
	routerID, err := getRouterIDMonitor(c)
	if err != nil {
		response.BadRequest(c, "router ID not found")
		return
	}

	info, err := h.service.GetRouterBoardInfo(c.Request.Context(), routerID)
	if err != nil {
		response.InternalServerError(c, err.Error())
		return
	}
	response.OK(c, info)
}

// ListInterfaces handles listing network interfaces
func (h *MonitorHandler) ListInterfaces(c *gin.Context) {
	routerID, err := getRouterIDMonitor(c)
	if err != nil {
		response.BadRequest(c, "router ID not found")
		return
	}

	interfaces, err := h.service.GetInterfaces(c.Request.Context(), routerID)
	if err != nil {
		response.InternalServerError(c, err.Error())
		return
	}
	response.OK(c, interfaces)
}

// GetLogs handles getting logs
func (h *MonitorHandler) GetLogs(c *gin.Context) {
	routerID, err := getRouterIDMonitor(c)
	if err != nil {
		response.BadRequest(c, "router ID not found")
		return
	}

	topics := c.Query("topics")
	limitStr := c.DefaultQuery("limit", "100")
	limit, _ := strconv.Atoi(limitStr)

	logs, err := h.service.GetLogs(c.Request.Context(), routerID, topics, limit)
	if err != nil {
		response.InternalServerError(c, err.Error())
		return
	}
	response.OK(c, logs)
}

// GetHotspotLogs handles getting hotspot logs
func (h *MonitorHandler) GetHotspotLogs(c *gin.Context) {
	routerID, err := getRouterIDMonitor(c)
	if err != nil {
		response.BadRequest(c, "router ID not found")
		return
	}

	limitStr := c.DefaultQuery("limit", "100")
	limit, _ := strconv.Atoi(limitStr)

	logs, err := h.service.GetHotspotLogs(c.Request.Context(), routerID, limit)
	if err != nil {
		response.InternalServerError(c, err.Error())
		return
	}
	response.OK(c, logs)
}

// GetPPPLogs handles getting PPP logs
func (h *MonitorHandler) GetPPPLogs(c *gin.Context) {
	routerID, err := getRouterIDMonitor(c)
	if err != nil {
		response.BadRequest(c, "router ID not found")
		return
	}

	limitStr := c.DefaultQuery("limit", "100")
	limit, _ := strconv.Atoi(limitStr)

	logs, err := h.service.GetPPPLogs(c.Request.Context(), routerID, limit)
	if err != nil {
		response.InternalServerError(c, err.Error())
		return
	}
	response.OK(c, logs)
}

// EnableHotspotLogging handles enabling hotspot logging
func (h *MonitorHandler) EnableHotspotLogging(c *gin.Context) {
	routerID, err := getRouterIDMonitor(c)
	if err != nil {
		response.BadRequest(c, "router ID not found")
		return
	}

	if err := h.service.EnableHotspotLogging(c.Request.Context(), routerID); err != nil {
		response.InternalServerError(c, err.Error())
		return
	}

	response.OK(c, gin.H{"message": "hotspot logging enabled"})
}

// EnablePPPLogging handles enabling PPP logging
func (h *MonitorHandler) EnablePPPLogging(c *gin.Context) {
	routerID, err := getRouterIDMonitor(c)
	if err != nil {
		response.BadRequest(c, "router ID not found")
		return
	}

	if err := h.service.EnablePPPLogging(c.Request.Context(), routerID); err != nil {
		response.InternalServerError(c, err.Error())
		return
	}

	response.OK(c, gin.H{"message": "ppp logging enabled"})
}

// Ping performs a one-shot ping from the router
func (h *MonitorHandler) Ping(c *gin.Context) {
	routerID, err := getRouterIDMonitor(c)
	if err != nil {
		response.BadRequest(c, "router ID not found")
		return
	}

	var req struct {
		Address  string `json:"address" binding:"required"`
		Count    int    `json:"count"`
		Size     int    `json:"size"`
		Interval int    `json:"interval"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	cfg := domain.PingConfig{
		Address: req.Address,
		Count:   req.Count,
		Size:    req.Size,
	}
	if cfg.Count <= 0 {
		cfg.Count = 4
	}
	if cfg.Size <= 0 {
		cfg.Size = 64
	}
	if req.Interval > 0 {
		cfg.Interval = time.Duration(req.Interval) * time.Second
	} else {
		cfg.Interval = 1 * time.Second
	}

	results, err := h.service.Ping(c.Request.Context(), routerID, cfg)
	if err != nil {
		response.InternalServerError(c, err.Error())
		return
	}

	response.OK(c, results)
}

// Stre