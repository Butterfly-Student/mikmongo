package mikrotik

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	mikrotiksvc "mikmongo/internal/service/mikrotik"
	"mikmongo/pkg/mikrotik/domain"
	"mikmongo/pkg/response"
)

// IPAddressHandler handles MikroTik IP Address HTTP requests
type IPAddressHandler struct {
	service *mikrotiksvc.IPAddressService
}

// NewIPAddressHandler creates a new IP Address handler
func NewIPAddressHandler(service *mikrotiksvc.IPAddressService) *IPAddressHandler {
	return &IPAddressHandler{service: service}
}

// getRouterIDIPAddress extracts router ID from context
func getRouterIDIPAddress(c *gin.Context) (uuid.UUID, error) {
	routerID, exists := c.Get("router_id")
	if !exists {
		return uuid.Nil, nil
	}
	return routerID.(uuid.UUID), nil
}

// ListAddresses handles listing IP addresses
func (h *IPAddressHandler) ListAddresses(c *gin.Context) {
	routerID, err := getRouterIDIPAddress(c)
	if err != nil {
		response.BadRequest(c, "router ID not found")
		return
	}

	addresses, err := h.service.GetAddresses(c.Request.Context(), routerID)
	if err != nil {
		response.InternalServerError(c, err.Error())
		return
	}
	response.OK(c, addresses)
}

// GetAddress handles getting an IP address by ID
func (h *IPAddressHandler) GetAddress(c *gin.Context) {
	routerID, err := getRouterIDIPAddress(c)
	if err != nil {
		response.BadRequest(c, "router ID not found")
		return
	}

	id := c.Param("id")
	addr, err := h.service.GetAddressByID(c.Request.Context(), routerID, id)
	if err != nil {
		response.NotFound(c, err.Error())
		return
	}
	response.OK(c, addr)
}

// GetAddressesByInterface handles getting IP addresses by interface
func (h *IPAddressHandler) GetAddressesByInterface(c *gin.Context) {
	routerID, err := getRouterIDIPAddress(c)
	if err != nil {
		response.BadRequest(c, "router ID not found")
		return
	}

	iface := c.Param("interface")
	addresses, err := h.service.GetAddressesByInterface(c.Request.Context(), routerID, iface)
	if err != nil {
		response.InternalServerError(c, err.Error())
		return
	}
	response.OK(c, addresses)
}

// CreateAddress handles creating a new IP address
func (h *IPAddressHandler) CreateAddress(c *gin.Context) {
	routerID, err := getRouterIDIPAddress(c)
	if err != nil {
		response.BadRequest(c, "router ID not found")
		return
	}

	var addr domain.IPAddress
	if err := c.ShouldBindJSON(&addr); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	id, err := h.service.AddAddress(c.Request.Context(), routerID, &addr)
	if err != nil {
		response.InternalServerError(c, err.Error())
		return
	}

	addr.ID = id
	response.Created(c, addr)
}

// UpdateAddress handles updating an IP address
func (h *IPAddressHandler) UpdateAddress(c *gin.Context) {
	routerID, err := getRouterIDIPAddress(c)
	if err != nil {
		response.BadRequest(c, "router ID not found")
		return
	}

	id := c.Param("id")
	var addr domain.IPAddress
	if err := c.ShouldBindJSON(&addr); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	if err := h.service.UpdateAddress(c.Request.Context(), routerID, id, &addr); err != nil {
		response.InternalServerError(c, err.Error())
		return
	}

	response.OK(c, gin.H{"message": "address updated"})
}

// DeleteAddress handles deleting an IP address
func (h *IPAddressHandler) DeleteAddress(c *gin.Context) {
	routerID, err := getRouterIDIPAddress(c)
	if err != nil {
		response.BadRequest(c, "router ID not found")
		return
	}

	id := c.Param("id")
	if err := h.service.RemoveAddress(c.Request.Context(), routerID, id); err != nil {
		response.InternalServerError(c, err.Error())
		return
	}

	response.OK(c, gin.H{"message": "address deleted"})
}

// EnableAddress handles enabling an IP address
func (h *IPAddressHandler) EnableAddress(c *gin.Context) {
	routerID, err := getRouterIDIPAddress(c)
	if err != nil {
		response.BadRequest(c, "router ID not found")
		return
	}

	id := c.Param("id")
	if err := h.service.EnableAddress(c.Request.Context(), routerID, id); err != nil {
		response.InternalServerError(c, err.Error())
		return
	}

	response.OK(c, gin.H{"message": "address enabled"})
}

// DisableAddress handles disabling an IP address
func (h *IPAddressHandler) DisableAddress(c *gin.Context) {
	routerID, err := getRouterIDIPAddress(c)
	if err != nil {
		response.BadRequest(c, "router ID not found")
		return
	}

	id := c.Param("id")
	if err := h.service.DisableAddress(c.Request.Context(), routerID, id); err != nil {
		response.InternalServerError(c, err.Error())
		return
	}

	response.OK(c, gin.H{"message": "address disabled"})
}
