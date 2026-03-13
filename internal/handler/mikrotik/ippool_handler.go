package mikrotik

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	mikrotiksvc "mikmongo/internal/service/mikrotik"
	"mikmongo/pkg/mikrotik/domain"
	"mikmongo/pkg/response"
)

// IPPoolHandler handles MikroTik IP Pool HTTP requests
type IPPoolHandler struct {
	service *mikrotiksvc.IPPoolService
}

// NewIPPoolHandler creates a new IP Pool handler
func NewIPPoolHandler(service *mikrotiksvc.IPPoolService) *IPPoolHandler {
	return &IPPoolHandler{service: service}
}

// getRouterIDIPPool extracts router ID from context
func getRouterIDIPPool(c *gin.Context) (uuid.UUID, error) {
	routerID, exists := c.Get("router_id")
	if !exists {
		return uuid.Nil, nil
	}
	return routerID.(uuid.UUID), nil
}

// ListPools handles listing IP pools
func (h *IPPoolHandler) ListPools(c *gin.Context) {
	routerID, err := getRouterIDIPPool(c)
	if err != nil {
		response.BadRequest(c, "router ID not found")
		return
	}

	pools, err := h.service.GetPools(c.Request.Context(), routerID)
	if err != nil {
		response.InternalServerError(c, err.Error())
		return
	}
	response.OK(c, pools)
}

// GetPool handles getting an IP pool by ID
func (h *IPPoolHandler) GetPool(c *gin.Context) {
	routerID, err := getRouterIDIPPool(c)
	if err != nil {
		response.BadRequest(c, "router ID not found")
		return
	}

	id := c.Param("id")
	pool, err := h.service.GetPoolByID(c.Request.Context(), routerID, id)
	if err != nil {
		response.NotFound(c, err.Error())
		return
	}
	response.OK(c, pool)
}

// GetPoolByName handles getting an IP pool by name
func (h *IPPoolHandler) GetPoolByName(c *gin.Context) {
	routerID, err := getRouterIDIPPool(c)
	if err != nil {
		response.BadRequest(c, "router ID not found")
		return
	}

	name := c.Query("name")
	if name == "" {
		response.BadRequest(c, "name query parameter is required")
		return
	}

	pool, err := h.service.GetPoolByName(c.Request.Context(), routerID, name)
	if err != nil {
		response.NotFound(c, err.Error())
		return
	}
	response.OK(c, pool)
}

// GetPoolNames handles getting all pool names
func (h *IPPoolHandler) GetPoolNames(c *gin.Context) {
	routerID, err := getRouterIDIPPool(c)
	if err != nil {
		response.BadRequest(c, "router ID not found")
		return
	}

	names, err := h.service.GetPoolNames(c.Request.Context(), routerID)
	if err != nil {
		response.InternalServerError(c, err.Error())
		return
	}
	response.OK(c, names)
}

// CreatePool handles creating a new IP pool
func (h *IPPoolHandler) CreatePool(c *gin.Context) {
	routerID, err := getRouterIDIPPool(c)
	if err != nil {
		response.BadRequest(c, "router ID not found")
		return
	}

	var pool domain.IPPool
	if err := c.ShouldBindJSON(&pool); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	id, err := h.service.AddPool(c.Request.Context(), routerID, &pool)
	if err != nil {
		response.InternalServerError(c, err.Error())
		return
	}

	pool.ID = id
	response.Created(c, pool)
}

// UpdatePool handles updating an IP pool
func (h *IPPoolHandler) UpdatePool(c *gin.Context) {
	routerID, err := getRouterIDIPPool(c)
	if err != nil {
		response.BadRequest(c, "router ID not found")
		return
	}

	id := c.Param("id")
	var pool domain.IPPool
	if err := c.ShouldBindJSON(&pool); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	if err := h.service.UpdatePool(c.Request.Context(), routerID, id, &pool); err != nil {
		response.InternalServerError(c, err.Error())
		return
	}

	response.OK(c, gin.H{"message": "pool updated"})
}

// DeletePool handles deleting an IP pool
func (h *IPPoolHandler) DeletePool(c *gin.Context) {
	routerID, err := getRouterIDIPPool(c)
	if err != nil {
		response.BadRequest(c, "router ID not found")
		return
	}

	id := c.Param("id")
	if err := h.service.RemovePool(c.Request.Context(), routerID, id); err != nil {
		response.InternalServerError(c, err.Error())
		return
	}

	response.OK(c, gin.H{"message": "pool deleted"})
}

// GetPoolUsed handles getting used IP allocations
func (h *IPPoolHandler) GetPoolUsed(c *gin.Context) {
	routerID, err := getRouterIDIPPool(c)
	if err != nil {
		response.BadRequest(c, "router ID not found")
		return
	}

	used, err := h.service.GetPoolUsed(c.Request.Context(), routerID)
	if err != nil {
		response.InternalServerError(c, err.Error())
		return
	}
	response.OK(c, used)
}
