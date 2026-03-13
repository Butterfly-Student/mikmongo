package mikrotik

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	mikrotiksvc "mikmongo/internal/service/mikrotik"
	"mikmongo/pkg/mikrotik/domain"
	"mikmongo/pkg/response"
)

// QueueHandler handles MikroTik Queue HTTP requests
type QueueHandler struct {
	service *mikrotiksvc.QueueService
}

// NewQueueHandler creates a new Queue handler
func NewQueueHandler(service *mikrotiksvc.QueueService) *QueueHandler {
	return &QueueHandler{service: service}
}

// getRouterIDQueue extracts router ID from context
func getRouterIDQueue(c *gin.Context) (uuid.UUID, error) {
	routerID, exists := c.Get("router_id")
	if !exists {
		return uuid.Nil, nil
	}
	return routerID.(uuid.UUID), nil
}

// ListSimpleQueues handles listing simple queues
func (h *QueueHandler) ListSimpleQueues(c *gin.Context) {
	routerID, err := getRouterIDQueue(c)
	if err != nil {
		response.BadRequest(c, "router ID not found")
		return
	}

	queues, err := h.service.GetSimpleQueues(c.Request.Context(), routerID)
	if err != nil {
		response.InternalServerError(c, err.Error())
		return
	}
	response.OK(c, queues)
}

// GetSimpleQueue handles getting a simple queue by ID
func (h *QueueHandler) GetSimpleQueue(c *gin.Context) {
	routerID, err := getRouterIDQueue(c)
	if err != nil {
		response.BadRequest(c, "router ID not found")
		return
	}

	id := c.Param("id")
	queue, err := h.service.GetSimpleQueueByID(c.Request.Context(), routerID, id)
	if err != nil {
		response.NotFound(c, err.Error())
		return
	}
	response.OK(c, queue)
}

// GetSimpleQueueByName handles getting a simple queue by name
func (h *QueueHandler) GetSimpleQueueByName(c *gin.Context) {
	routerID, err := getRouterIDQueue(c)
	if err != nil {
		response.BadRequest(c, "router ID not found")
		return
	}

	name := c.Query("name")
	if name == "" {
		response.BadRequest(c, "name query parameter is required")
		return
	}

	queue, err := h.service.GetSimpleQueueByName(c.Request.Context(), routerID, name)
	if err != nil {
		response.NotFound(c, err.Error())
		return
	}
	response.OK(c, queue)
}

// GetAllQueues handles getting all queue names
func (h *QueueHandler) GetAllQueues(c *gin.Context) {
	routerID, err := getRouterIDQueue(c)
	if err != nil {
		response.BadRequest(c, "router ID not found")
		return
	}

	queues, err := h.service.GetAllQueues(c.Request.Context(), routerID)
	if err != nil {
		response.InternalServerError(c, err.Error())
		return
	}
	response.OK(c, queues)
}

// GetAllParentQueues handles getting all parent queue names
func (h *QueueHandler) GetAllParentQueues(c *gin.Context) {
	routerID, err := getRouterIDQueue(c)
	if err != nil {
		response.BadRequest(c, "router ID not found")
		return
	}

	queues, err := h.service.GetAllParentQueues(c.Request.Context(), routerID)
	if err != nil {
		response.InternalServerError(c, err.Error())
		return
	}
	response.OK(c, queues)
}

// CreateSimpleQueue handles creating a new simple queue
func (h *QueueHandler) CreateSimpleQueue(c *gin.Context) {
	routerID, err := getRouterIDQueue(c)
	if err != nil {
		response.BadRequest(c, "router ID not found")
		return
	}

	var queue domain.SimpleQueue
	if err := c.ShouldBindJSON(&queue); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	id, err := h.service.AddSimpleQueue(c.Request.Context(), routerID, &queue)
	if err != nil {
		response.InternalServerError(c, err.Error())
		return
	}

	queue.ID = id
	response.Created(c, queue)
}

// UpdateSimpleQueue handles updating a simple queue
func (h *QueueHandler) UpdateSimpleQueue(c *gin.Context) {
	routerID, err := getRouterIDQueue(c)
	if err != nil {
		response.BadRequest(c, "router ID not found")
		return
	}

	id := c.Param("id")
	var queue domain.SimpleQueue
	if err := c.ShouldBindJSON(&queue); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	if err := h.service.UpdateSimpleQueue(c.Request.Context(), routerID, id, &queue); err != nil {
		response.InternalServerError(c, err.Error())
		return
	}

	response.OK(c, gin.H{"message": "queue updated"})
}

// DeleteSimpleQueue handles deleting a simple queue
func (h *QueueHandler) DeleteSimpleQueue(c *gin.Context) {
	routerID, err := getRouterIDQueue(c)
	if err != nil {
		response.BadRequest(c, "router ID not found")
		return
	}

	id := c.Param("id")
	if err := h.service.RemoveSimpleQueue(c.Request.Context(), routerID, id); err != nil {
		response.InternalServerError(c, err.Error())
		return
	}

	response.OK(c, gin.H{"message": "queue deleted"})
}

// EnableSimpleQueue handles enabling a simple queue
func (h *QueueHandler) EnableSimpleQueue(c *gin.Context) {
	routerID, err := getRouterIDQueue(c)
	if err != nil {
		response.BadRequest(c, "router ID not found")
		return
	}

	id := c.Param("id")
	if err := h.service.EnableSimpleQueue(c.Request.Context(), routerID, id); err != nil {
		response.InternalServerError(c, err.Error())
		return
	}

	response.OK(c, gin.H{"message": "queue enabled"})
}

// DisableSimpleQueue handles disabling a simple queue
func (h *QueueHandler) DisableSimpleQueue(c *gin.Context) {
	routerID, err := getRouterIDQueue(c)
	if err != nil {
		response.BadRequest(c, "router ID not found")
		return
	}

	id := c.Param("id")
	if err := h.service.DisableSimpleQueue(c.Request.Context(), routerID, id); err != nil {
		response.InternalServerError(c, err.Error())
		return
	}

	response.OK(c, gin.H{"message": "queue disabled"})
}

// ResetQueueCounters handles resetting queue counters
func (h *QueueHandler) ResetQueueCounters(c *gin.Context) {
	routerID, err := getRouterIDQueue(c)
	if err != nil {
		response.BadRequest(c, "router ID not found")
		return
	}

	id := c.Param("id")
	if err := h.service.ResetQueueCounters(c.Request.Context(), routerID, id); err != nil {
		response.InternalServerError(c, err.Error())
		return
	}

	response.OK(c, gin.H{"message": "counters reset"})
}

// ResetAllQueueCounters handles resetting all queue counters
func (h *QueueHandler) ResetAllQueueCounters(c *gin.Context) {
	routerID, err := getRouterIDQueue(c)
	if err != nil {
		response.BadRequest(c, "router ID not found")
		return
	}

	if err := h.service.ResetAllQueueCounters(c.Request.Context(), routerID); err != nil {
		response.InternalServerError(c, err.Error())
		return
	}

	response.OK(c, gin.H{"message": "all counters reset"})
}
