package handler

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"mikmongo/internal/repository"
	"mikmongo/internal/service"
	"mikmongo/pkg/response"
)

// HotspotSaleHandler handles hotspot sales record requests.
type HotspotSaleHandler struct {
	svc *service.HotspotSaleService
}

// NewHotspotSaleHandler creates a new HotspotSaleHandler.
func NewHotspotSaleHandler(svc *service.HotspotSaleService) *HotspotSaleHandler {
	return &HotspotSaleHandler{svc: svc}
}

// List handles listing hotspot sales with optional filters.
// Query params: router_id, agent_id, profile, batch_code, date_from, date_to, limit, offset
func (h *HotspotSaleHandler) List(c *gin.Context) {
	filter := repository.HotspotSaleFilter{}

	if rid := c.Query("router_id"); rid != "" {
		id, err := uuid.Parse(rid)
		if err != nil {
			response.BadRequest(c, "invalid router_id")
			return
		}
		filter.RouterID = &id
	}

	if aid := c.Query("agent_id"); aid != "" {
		id, err := uuid.Parse(aid)
		if err != nil {
			response.BadRequest(c, "invalid agent_id")
			return
		}
		filter.SalesAgentID = &id
	}

	filter.Profile = c.Query("profile")
	filter.BatchCode = c.Query("batch_code")

	if df := c.Query("date_from"); df != "" {
		t, err := time.Parse("2006-01-02", df)
		if err != nil {
			response.BadRequest(c, "invalid date_from, use YYYY-MM-DD")
			return
		}
		filter.DateFrom = &t
	}

	if dt := c.Query("date_to"); dt != "" {
		t, err := time.Parse("2006-01-02", dt)
		if err != nil {
			response.BadRequest(c, "invalid date_to, use YYYY-MM-DD")
			return
		}
		filter.DateTo = &t
	}

	limit, offset := getPagination(c)
	sales, count, err := h.svc.ListSales(c.Request.Context(), filter, limit, offset)
	if err != nil {
		response.InternalServerError(c, err.Error())
		return
	}
	response.WithMeta(c, http.StatusOK, sales, &response.Meta{Total: count, Limit: limit, Offset: offset})
}

// ListByRouter handles listing hotspot sales scoped to a specific router.
func (h *HotspotSaleHandler) ListByRouter(c *gin.Context) {
	routerID, err := uuid.Parse(c.Param("router_id"))
	if err != nil {
		response.BadRequest(c, "invalid router_id")
		return
	}

	filter := repository.HotspotSaleFilter{RouterID: &routerID}
	filter.Profile = c.Query("profile")
	filter.BatchCode = c.Query("batch_code")

	if aid := c.Query("agent_id"); aid != "" {
		id, err := uuid.Parse(aid)
		if err != nil {
			response.BadRequest(c, "invalid agent_id")
			return
		}
		filter.SalesAgentID = &id
	}

	limit, offset := getPagination(c)
	sales, count, err := h.svc.ListSales(c.Request.Context(), filter, limit, offset)
	if err != nil {
		response.InternalServerError(c, err.Error())
		return
	}
	response.WithMeta(c, http.StatusOK, sales, &response.Meta{Total: count, Limit: limit, Offset: offset})
}
