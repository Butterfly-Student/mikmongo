package mikhmon

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	mikhdto "mikmongo/internal/dto/mikrotik/mikhmon"
	mikhmonservice "mikmongo/internal/service/mikrotik/mikhmon"
	"mikmongo/pkg/response"
)

type ExpireHandler struct {
	expireSvc *mikhmonservice.MikhmonExpireService
}

func NewExpireHandler(expireSvc *mikhmonservice.MikhmonExpireService) *ExpireHandler {
	return &ExpireHandler{
		expireSvc: expireSvc,
	}
}

func (h *ExpireHandler) Setup(c *gin.Context) {
	routerID, err := uuid.Parse(c.Param("router_id"))
	if err != nil {
		response.BadRequest(c, "invalid router_id")
		return
	}

	if err := h.expireSvc.SetupExpireMonitor(c.Request.Context(), routerID); err != nil {
		response.InternalServerError(c, err.Error())
		return
	}

	response.OK(c, mikhdto.ExpireMonitorStatusResponse{
		Enabled: true,
	})
}

func (h *ExpireHandler) Disable(c *gin.Context) {
	routerID, err := uuid.Parse(c.Param("router_id"))
	if err != nil {
		response.BadRequest(c, "invalid router_id")
		return
	}

	if err := h.expireSvc.DisableExpireMonitor(c.Request.Context(), routerID); err != nil {
		response.InternalServerError(c, err.Error())
		return
	}

	response.OK(c, mikhdto.ExpireMonitorStatusResponse{
		Enabled: false,
	})
}

func (h *ExpireHandler) GetStatus(c *gin.Context) {
	routerID, err := uuid.Parse(c.Param("router_id"))
	if err != nil {
		response.BadRequest(c, "invalid router_id")
		return
	}

	enabled, err := h.expireSvc.IsExpireMonitorEnabled(c.Request.Context(), routerID)
	if err != nil {
		response.InternalServerError(c, err.Error())
		return
	}

	response.OK(c, mikhdto.ExpireMonitorStatusResponse{
		Enabled: enabled,
	})
}

func (h *ExpireHandler) GenerateScript(c *gin.Context) {
	script := h.expireSvc.GenerateExpireMonitorScript()
	response.OK(c, mikhdto.ScriptResponse{
		Name:    "expire-monitor",
		Content: script,
	})
}
