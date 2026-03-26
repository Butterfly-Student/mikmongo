package mikrotik

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	svcmikrotik "mikmongo/internal/service/mikrotik"
	"mikmongo/pkg/response"
)

// QueueHandler handles Queue-related REST endpoints.
type QueueHandler struct {
	queueSvc *svcmikrotik.QueueService
}

func NewQueueHandler(queueSvc *svcmikrotik.QueueService) *QueueHandler {
	return &QueueHandler{queueSvc: queueSvc}
}

func (h *QueueHandler) GetSimpleQueues(c *gin.Context) {
	routerID, err := uuid.Parse(c.Param("router_id"))
	if err != nil {
		response.BadRequest(c, "invalid router_id")
		return
	}
	queues, err := h.queueSvc.GetSimpleQueues(c.Request.Context(), routerID)
	if err != nil {
		response.InternalServerError(c, err.Error())
		return
	}
	response.OK(c, queues)
}
