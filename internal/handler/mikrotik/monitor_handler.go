package mikrotik

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	svcmikrotik "mikmongo/internal/service/mikrotik"
	"mikmongo/pkg/response"
)

// MonitorHandler handles Monitor REST endpoints (non-streaming).
type MonitorHandler struct {
	monitorSvc *svcmikrotik.MonitorService
}

func NewMonitorHandler(monitorSvc *svcmikrotik.MonitorService) *MonitorHandler {
	return &MonitorHandler{monitorSvc: monitorSvc}
}

func (h *MonitorHandler) GetSystemResource(c *gin.Context) {
	routerID, err := uuid.Parse(c.Param("router_id"))
	if err != nil {
		response.BadRequest(c, "invalid router_id")
		return
	}
	resource, err := h.monitorSvc.GetSystemResource(c.Request.Context(), routerID)
	if err != nil {
		response.InternalServerError(c, err.Error())
		return
	}
	response.OK(c, resource)
}

func (h *MonitorHandler) GetInterfaces(c *gin.Context) {
	routerID, err := uuid.Parse(c.Param("router_id"))
	if err != nil {
		response.BadRequest(c, "invalid router_id")
		return
	}
	ifaces, err := h.monitorSvc.GetInterfaces(c.Request.Context(), routerID)
	if err != nil {
		response.InternalServerError(c, err.Error())
		return
	}
	response.OK(c, ifaces)
}
