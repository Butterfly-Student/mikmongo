package mikrotik

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	svcmikrotik "mikmongo/internal/service/mikrotik"
	"mikmongo/pkg/response"
)

// FirewallHandler handles Firewall-related REST endpoints.
type FirewallHandler struct {
	firewallSvc *svcmikrotik.FirewallService
}

func NewFirewallHandler(firewallSvc *svcmikrotik.FirewallService) *FirewallHandler {
	return &FirewallHandler{firewallSvc: firewallSvc}
}

func (h *FirewallHandler) GetFilterRules(c *gin.Context) {
	routerID, err := uuid.Parse(c.Param("router_id"))
	if err != nil {
		response.BadRequest(c, "invalid router_id")
		return
	}
	rules, err := h.firewallSvc.GetFilterRules(c.Request.Context(), routerID)
	if err != nil {
		response.InternalServerError(c, err.Error())
		return
	}
	response.OK(c, rules)
}

func (h *FirewallHandler) GetNATRules(c *gin.Context) {
	routerID, err := uuid.Parse(c.Param("router_id"))
	if err != nil {
		response.BadRequest(c, "invalid router_id")
		return
	}
	rules, err := h.firewallSvc.GetNATRules(c.Request.Context(), routerID)
	if err != nil {
		response.InternalServerError(c, err.Error())
		return
	}
	response.OK(c, rules)
}

func (h *FirewallHandler) GetAddressLists(c *gin.Context) {
	routerID, err := uuid.Parse(c.Param("router_id"))
	if err != nil {
		response.BadRequest(c, "invalid router_id")
		return
	}
	lists, err := h.firewallSvc.GetAddressLists(c.Request.Context(), routerID)
	if err != nil {
		response.InternalServerError(c, err.Error())
		return
	}
	response.OK(c, lists)
}
