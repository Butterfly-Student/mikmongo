package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"mikmongo/internal/model"
	"mikmongo/internal/service"
	"mikmongo/pkg/response"
)

// RouterHandler handles router HTTP requests
type RouterHandler struct {
	service *service.RouterService
}

// NewRouterHandler creates a new router handler
func NewRouterHandler(service *service.RouterService) *RouterHandler {
	return &RouterHandler{service: service}
}

// List handles listing routers
func (h *RouterHandler) List(c *gin.Context) {
	limit, offset := getPagination(c)
	routers, err := h.service.List(c.Request.Context(), limit, offset)
	if err != nil {
		response.InternalServerError(c, err.Error())
		return
	}
	response.OK(c, routers)
}

// Create handles creating a router
func (h *RouterHandler) Create(c *gin.Context) {
	var req struct {
		model.MikrotikRouter
		Password string `json:"password" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	if err := h.service.Create(c.Request.Context(), &req.MikrotikRouter, req.Password); err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}
	response.Created(c, req.MikrotikRouter)
}

// GetDevice handles getting a router by ID
func (h *RouterHandler) GetDevice(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		response.BadRequest(c, "invalid id")
		return
	}
	device, err := h.service.GetDevice(c.Request.Context(), id)
	if err != nil {
		response.NotFound(c, err.Error())
		return
	}
	response.OK(c, device)
}

// Update handles updating a router
func (h *RouterHandler) Update(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		response.BadRequest(c, "invalid id")
		return
	}
	var req struct {
		model.MikrotikRouter
		Password string `json:"password"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	req.MikrotikRouter.ID = id.String()
	if err := h.service.Update(c.Request.Context(), &req.MikrotikRouter, req.Password); err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}
	response.OK(c, req.MikrotikRouter)
}

// Delete handles deleting a router
func (h *RouterHandler) Delete(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		response.BadRequest(c, "invalid id")
		return
	}
	if err := h.service.Delete(c.Request.Context(), id); err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}
	response.OK(c, gin.H{"message": "router deleted"})
}

// SyncDevice handles syncing a router device
func (h *RouterHandler) SyncDevice(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		response.BadRequest(c, "invalid id")
		return
	}
	if err := h.service.SyncDevice(c.Request.Context(), id); err != nil {
		response.InternalServerError(c, err.Error())
		return
	}
	response.OK(c, gin.H{"message": "device synced"})
}

// TestConnection handles testing a router connection
func (h *RouterHandler) TestConnection(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		response.BadRequest(c, "invalid id")
		return
	}
	if err := h.service.TestConnection(c.Request.Context(), id); err != nil {
		response.Error(c, http.StatusBadRequest, err.Error())
		return
	}
	response.OK(c, gin.H{"message": "connection successful"})
}

// SyncAll handles syncing all devices
func (h *RouterHandler) SyncAll(c *gin.Context) {
	if err := h.service.SyncAllDevices(c.Request.Context()); err != nil {
		response.InternalServerError(c, err.Error())
		return
	}
	response.OK(c, gin.H{"message": "all devices synced"})
}

// SelectRouter handles selecting an active router for the current user
func (h *RouterHandler) SelectRouter(c *gin.Context) {
	userIDStr, _ := c.Get("user_id")
	routerID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		response.BadRequest(c, "invalid router id")
		return
	}

	router, err := h.service.SelectRouter(c.Request.Context(), userIDStr.(string), routerID)
	if err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	response.OK(c, router)
}

// GetSelectedRouter handles getting the currently selected router
func (h *RouterHandler) GetSelectedRouter(c *gin.Context) {
	userIDStr, _ := c.Get("user_id")

	router, err := h.service.GetSelectedRouter(c.Request.Context(), userIDStr.(string))
	if err != nil {
		response.InternalServerError(c, err.Error())
		return
	}

	if router == nil {
		response.OK(c, nil)
		return
	}

	response.OK(c, router)
}
