package mikrotik

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	mikrotiksvc "mikmongo/internal/service/mikrotik"
	"mikmongo/pkg/response"
)

// FirewallHandler handles MikroTik Firewall HTTP requests
type FirewallHandler struct {
	service *mikrotiksvc.FirewallService
}

// NewFirewallHandler creates a new Firewall handler
func NewFirewallHandler(service *mikrotiksvc.FirewallService) *FirewallHandler {
	return &FirewallHandler{service: service}
}

// getRouterIDFirewall extracts router ID from context
func getRouterIDFirewall(c *gin.Context) (uuid.UUID, error) {
	routerID, exists := c.Get("router_id")
	if !exists {
		return uuid.Nil, nil
	}
	return routerID.(uuid.UUID), nil
}

// ListNATRules handles listing NAT rules
func (h *FirewallHandler) ListNATRules(c *gin.Context) {
	routerID, err := getRouterIDFirewall(c)
	if err != nil {
		response.BadRequest(c, "router ID not found")
		return
	}

	rules, err := h.service.GetNATRules(c.Request.Context(), routerID)
	if err != nil {
		response.InternalServerError(c, err.Error())
		return
	}
	response.OK(c, rules)
}

// ListFilterRules handles listing filter rules
func (h *FirewallHandler) ListFilterRules(c *gin.Context) {
	routerID, err := getRouterIDFirewall(c)
	if err != nil {
		response.BadRequest(c, "router ID not found")
		return
	}

	rules, err := h.service.GetFilterRules(c.Request.Context(), routerID)
	if err != nil {
		response.InternalServerError(c, err.Error())
		return
	}
	response.OK(c, rules)
}

// ListAddressLists handles listing address lists
func (h *FirewallHandler) ListAddressLists(c *gin.Context) {
	routerID, err := getRouterIDFirewall(c)
	if err != nil {
		response.BadRequest(c, "router ID not found")
		return
	}

	lists, err := h.service.GetAddressLists(c.Request.Context(), routerID)
	if err != nil {
		response.InternalServerError(c, err.Error())
		return
	}
	response.OK(c, lists)
}
