package mikrotik

import (
	mkdomain "github.com/Butterfly-Student/go-ros/domain"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	svcmikrotik "mikmongo/internal/service/mikrotik"
	"mikmongo/pkg/response"
)

// IPHandler handles IP Pool and IP Address REST endpoints.
type IPHandler struct {
	poolSvc *svcmikrotik.IPPoolService
	addrSvc *svcmikrotik.IPAddressService
}

func NewIPHandler(poolSvc *svcmikrotik.IPPoolService, addrSvc *svcmikrotik.IPAddressService) *IPHandler {
	return &IPHandler{poolSvc: poolSvc, addrSvc: addrSvc}
}

func (h *IPHandler) GetPools(c *gin.Context) {
	routerID, err := uuid.Parse(c.Param("router_id"))
	if err != nil {
		response.BadRequest(c, "invalid router_id")
		return
	}
	pools, err := h.poolSvc.GetPools(c.Request.Context(), routerID)
	if err != nil {
		response.InternalServerError(c, err.Error())
		return
	}
	response.OK(c, pools)
}

func (h *IPHandler) AddPool(c *gin.Context) {
	routerID, err := uuid.Parse(c.Param("router_id"))
	if err != nil {
		response.BadRequest(c, "invalid router_id")
		return
	}
	var req mkdomain.IPPool
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	id, err := h.poolSvc.AddPool(c.Request.Context(), routerID, &req)
	if err != nil {
		response.InternalServerError(c, err.Error())
		return
	}
	response.Created(c, gin.H{"id": id})
}

func (h *IPHandler) RemovePool(c *gin.Context) {
	routerID, err := uuid.Parse(c.Param("router_id"))
	if err != nil {
		response.BadRequest(c, "invalid router_id")
		return
	}
	id := c.Param("id")
	if err := h.poolSvc.RemovePool(c.Request.Context(), routerID, id); err != nil {
		response.InternalServerError(c, err.Error())
		return
	}
	response.OK(c, gin.H{"message": "pool removed"})
}

func (h *IPHandler) GetAddresses(c *gin.Context) {
	routerID, err := uuid.Parse(c.Param("router_id"))
	if err != nil {
		response.BadRequest(c, "invalid router_id")
		return
	}
	addresses, err := h.addrSvc.GetAddresses(c.Request.Context(), routerID)
	if err != nil {
		response.InternalServerError(c, err.Error())
		return
	}
	response.OK(c, addresses)
}
