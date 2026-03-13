package mikrotik

import (
	"github.com/gin-gonic/gin"
	mikrotiksvc "mikmongo/internal/service/mikrotik"
	"mikmongo/pkg/mikrotik/domain"
	"mikmongo/pkg/response"
)

// ScriptHandler handles script generation HTTP requests
type ScriptHandler struct {
	service *mikrotiksvc.ScriptService
}

// NewScriptHandler creates a new Script handler
func NewScriptHandler(service *mikrotiksvc.ScriptService) *ScriptHandler {
	return &ScriptHandler{service: service}
}

// GenerateOnLogin handles generating on-login script
func (h *ScriptHandler) GenerateOnLogin(c *gin.Context) {
	var req domain.ProfileRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	script := h.service.GenerateOnLoginScript(&req)
	response.OK(c, gin.H{
		"script": script,
		"type":   "on-login",
	})
}

// ParseOnLogin handles parsing on-login script
func (h *ScriptHandler) ParseOnLogin(c *gin.Context) {
	var req struct {
		Script string `json:"script" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	parsed := h.service.ParseOnLoginScript(req.Script)
	response.OK(c, parsed)
}

// GenerateExpiredAction handles generating expired action script
func (h *ScriptHandler) GenerateExpiredAction(c *gin.Context) {
	var req struct {
		ExpireMode string `json:"expireMode" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	script := h.service.GenerateExpiredAction(req.ExpireMode)
	response.OK(c, gin.H{
		"script": script,
		"type":   "expired-action",
	})
}

// GenerateExpireMonitor handles generating expire monitor script
func (h *ScriptHandler) GenerateExpireMonitor(c *gin.Context) {
	script := h.service.GenerateExpireMonitorScript()
	response.OK(c, gin.H{
		"script": script,
		"type":   "expire-monitor",
	})
}

// GetExpireModes handles getting available expire modes
func (h *ScriptHandler) GetExpireModes(c *gin.Context) {
	modes := []gin.H{
		{"value": "0", "label": "No Expiration", "description": "User never expires"},
		{"value": "rem", "label": "Remove User", "description": "Remove user when expired"},
		{"value": "remc", "label": "Remove User + Record", "description": "Remove user and record sales when expired"},
		{"value": "ntf", "label": "Notify Only", "description": "Mark user as expired but keep account"},
		{"value": "ntfc", "label": "Notify + Record", "description": "Mark user as expired and record sales"},
	}
	response.OK(c, modes)
}
