package mikrotik

import (
	"context"
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/gorilla/websocket"
	mikrotiksvc "mikmongo/internal/service/mikrotik"
	"mikmongo/pkg/mikrotik/domain"
)

// WebSocketHandler handles WebSocket connections for MikroTik streaming
type WebSocketHandler struct {
	hotspotService *mikrotiksvc.HotspotService
	pppService     *mikrotiksvc.PPPService
	queueService   *mikrotiksvc.QueueService
	monitorService *mikrotiksvc.MonitorService
	upgrader       websocket.Upgrader
	internalKey    string
}

// NewWebSocketHandler creates a new WebSocket handler
func NewWebSocketHandler(hotspotService *mikrotiksvc.HotspotService, pppService *mikrotiksvc.PPPService, queueService *mikrotiksvc.QueueService, monitorService *mikrotiksvc.MonitorService) *WebSocketHandler {
	return &WebSocketHandler{
		hotspotService: hotspotService,
		pppService:     pppService,
		queueService:   queueService,
		monitorService: monitorService,
		upgrader: websocket.Upgrader{
			CheckOrigin: func(r *http.Request) bool {
				return true // Allow all origins for now
			},
		},
		internalKey: os.Getenv("INTERNAL_KEY"),
	}
}

// validateInternalKey validates the internal key from header or query parameter.
// Browser WebSocket API does not support custom headers, so we also accept
// the key as a "key" query parameter.
func (h *WebSocketHandler) validateInternalKey(c *gin.Context) bool {
	if h.internalKey == "" {
		return false
	}
	key := c.GetHeader("X-Internal-Key")
	if key == "" {
		key = c.Query("key")
	}
	return key == h.internalKey
}

// getRouterID extracts router ID from context
func getRouterIDWS(c *gin.Context) uuid.UUID {
	routerID, exists := c.Get("router_id")
	if !exists {
		return uuid.Nil
	}
	return routerID.(uuid.UUID)
}

// StreamHotspotActive streams active hotspot sessions via WebSocket
func (h *WebSocketHandler) StreamHotspotActive(c *gin.Context) {
	if !h.validateInternalKey(c) {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	routerID := getRouterIDWS(c)
	if routerID == uuid.Nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "router ID not found"})
		return
	}

	conn, err := h.upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		return
	}
	defer conn.Close()

	ctx, cancel := context.WithCancel(c.Request.Context())
	defer cancel()

	resultChan := make(chan []*domain.HotspotActive, 10)
	cleanup, err := h.hotspotService.ListenActive(ctx, routerID, resultChan)
	if err != nil {
		conn.WriteJSON(gin.H{"error": err.Error()})
		return
	}
	defer cleanup()

	// Handle client disconnect
	go func() {
		for {
			_, _, err := conn.ReadMessage()
			if err != nil {
				cancel()
				return
			}
		}
	}()

	// Stream data to client
	for {
		select {
		case data := <-resultChan:
			if err := conn.WriteJSON(data); err != nil {
				return
			}
		case <-ctx.Done():
			return
		}
	}
}

// StreamHotspotInactive streams inactive hotspot users via WebSocket
func (h *WebSocketHandler) StreamHotspotInactive(c *gin.Context) {
	if !h.validateInternalKey(c) {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	routerID := getRouterIDWS(c)
	if routerID == uuid.Nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "router ID not found"})
		return
	}

	conn, err := h.upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		return
	}
	defer conn.Close()

	ctx, cancel := context.WithCancel(c.Request.Context())
	defer cancel()

	resultChan := make(chan []*domain.HotspotUser, 10)
	cleanup, err := h.hotspotService.ListenInactive(ctx, routerID, resultChan)
	if err != nil {
		conn.WriteJSON(gin.H{"error": err.Error()})
		return
	}
	defer cleanup()

	go func() {
		for {
			_, _, err := conn.ReadMessage()
			if err != nil {
				cancel()
				return
			}
		}
	}()

	for {
		select {
		case data := <-resultChan:
			if err := conn.WriteJSON(data); err != nil {
				return
			}
		case <-ctx.Done():
			return
		}
	}
}

// StreamPPPActive streams active PPP sessions via WebSocket
func (h *WebSocketHandler) StreamPPPActive(c *gin.Context) {
	if !h.validateInternalKey(c) {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	routerID := getRouterIDWS(c)
	if routerID == uuid.Nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "router ID not found"})
		return
	}

	conn, err := h.upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		return
	}
	defer conn.Close()

	ctx, cancel := context.WithCancel(c.Request.Context())
	defer cancel()

	resultChan := make(chan []*domain.PPPActive, 10)
	cleanup, err := h.pppService.ListenActive(ctx, routerID, resultChan)
	if err != nil {
		conn.WriteJSON(gin.H{"error": err.Error()})
		return
	}
	defer cleanup()

	go func() {
		for {
			_, _, err := conn.ReadMessage()
			if err != nil {
				cancel()
				return
			}
		}
	}()

	for {
		select {
		case data := <-resultChan:
			if err := conn.WriteJSON(data); err != nil {
				return
			}
		case <-ctx.Done():
			return
		}
	}
}

// StreamPPPInactive streams inactive PPP secrets via WebSocket
func (h *WebSocketHandler) StreamPPPInactive(c *gin.Context) {
	if !h.validateInternalKey(c) {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	routerID := getRouterIDWS(c)
	if routerID == uuid.Nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "router ID not found"})
		return
	}

	conn, err := h.upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		return
	}
	defer conn.Close()

	ctx, cancel := context.WithCancel(c.Request.Context())
	defer cancel()

	resultChan := make(chan []*domain.PPPSecret, 10)
	cleanup, err := h.pppService.ListenInactive(ctx, routerID, resultChan)
	if err != nil {
		conn.WriteJSON(gin.H{"error": err.Error()})
		return
	}
	defer cleanup()

	go func() {
		for {
			_, _, err := conn.ReadMessage()
			if err != nil {
				cancel()
				return
			}
		}
	}()

	for {
		select {
		case data := <-resultChan:
			if err := conn.WriteJSON(data); err != nil {
				return
			}
		case <-ctx.Done():
			return
		}
	}
}

// StreamQueueStats streams queue statistics via WebSocket
func (h *WebSocketHandler) StreamQueueStats(c *gin.Context) {
	if !h.validateInternalKey(c) {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	routerID := getRouterIDWS(c)
	if routerID == uuid.Nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "router ID not found"})
		return
	}

	queueName := c.Query("name")
	if queueName == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "name query parameter is required"})
		return
	}

	conn, err := h.upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		return
	}
	defer conn.Close()

	ctx, cancel := context.WithCancel(c.Request.Context())
	defer cancel()

	cfg := domain.QueueStatsConfig{Name: queueName}
	resultChan := make(chan domain.QueueStats, 10)
	cleanup, err := h.queueService.StartQueueStatsListen(ctx, routerID, cfg, resultChan)
	if err != nil {
		conn.WriteJSON(gin.H{"error": err.Error()})
		return
	}
	defer cleanup()

	go func() {
		for {
			_, _, err := conn.ReadMessage()
			if err != nil {
				cancel()
				return
			}
		}
	}()

	for {
		select {
		case data := <-resultChan:
			if err := conn.WriteJSON(data); err != nil {
				return
			}
		case <-ctx.Done():
			return
		}
	}
}

// StreamSystemResource streams system resource stats via WebSocket
func (h *WebSocketHandler) StreamSystemResource(c *gin.Context) {
	if !h.validateInternalKey(c) {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	routerID := getRouterIDWS(c)
	if routerID == uuid.Nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "router ID not found"})
		return
	}

	conn, err := h.upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		return
	}
	defer conn.Close()

	ctx, cancel := context.WithCancel(c.Request.Context())
	defer cancel()

	resultChan := make(chan domain.SystemResourceMonitorStats, 10)
	cleanup, err := h.monitorService.StartSystemResourceMonitorListen(ctx, routerID, resultChan)
	if err != nil {
		conn.WriteJSON(gin.H{"error": err.Error()})
		return
	}
	defer cleanup()

	go func() {
		for {
			_, _, err := conn.ReadMessage()
			if err != nil {
				cancel()
				return
			}
		}
	}()

	for {
		select {
		case data := <-resultChan:
			if err := conn.WriteJSON(data); err != nil {
				return
			}
		case <-ctx.Done():
			return
		}
	}
}

// StreamTraffic streams interface traffic stats via WebSocket
func (h *WebSocketHandler) StreamTraffic(c *gin.Context) {
	if !h.validateInternalKey(c) {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	routerID := getRouterIDWS(c)
	if routerID == uuid.Nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "router ID not found"})
		return
	}

	iface := c.Query("interface")
	if iface == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "interface query parameter is required"})
		return
	}

	conn, err := h.upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		return
	}
	defer conn.Close()

	ctx, cancel := context.WithCancel(c.Request.Context())
	defer cancel()

	resultChan := make(chan domain.TrafficMonitorStats, 10)
	cleanup, err := h.monitorService.StartTrafficMonitorListen(ctx, routerID, iface, resultChan)
	if err != nil {
		conn.WriteJSON(gin.H{"error": err.Error()})
		return
	}
	defer cleanup()

	go func() {
		for {
			_, _, err := conn.ReadMessage()
			if err != nil {
				cancel()
				return
			}
		}
	}()

	for {
		select {
		case data := <-resultChan:
			if err := conn.WriteJSON(data); err != nil {
				return
			}
		case <-ctx.Done():
			return
		}
	}
}

// StreamPing streams ping results via WebSocket
func (h *WebSocketHandler) StreamPing(c *gin.Context) {
	if !h.validateInternalKey(c) {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	routerID := getRouterIDWS(c)
	if routerID == uuid.Nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "router ID not found"})
		return
	}

	address := c.Query("address")
	if address == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "address query parameter is required"})
		return
	}

	conn, err := h.upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		return
	}
	defer conn.Close()

	ctx, cancel := context.WithCancel(c.Request.Context())
	defer cancel()

	cfg := domain.DefaultPingConfig(address)
	resultChan := make(chan domain.PingResult, 10)
	cleanup, err := h.monitorService.StartPingListen(ctx, routerID, cfg, resultChan)
	if err != nil {
		conn.WriteJSON(gin.H{"error": err.Error()})
		return
	}
	defer cleanup()

	go func() {
		for {
			_, _, err := conn.ReadMessage()
			if err != nil {
				cancel()
				return
			}
		}
	}()

	for {
		select {
		case data := <-resultChan:
			if err := conn.WriteJSON(data); err != nil {
				return
			}
		case <-ctx.Done():
			return
		}
	}
}

// StreamLogs streams logs via WebSocket
func (h *WebSocketHandler) StreamLogs(c *gin.Context) {
	if !h.validateInternalKey(c) {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	routerID := getRouterIDWS(c)
	if routerID == uuid.Nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "router ID not found"})
		return
	}

	topics := c.Query("topics")

	conn, err := h.upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		return
	}
	defer conn.Close()

	ctx, cancel := context.WithCancel(c.Request.Context())
	defer cancel()

	resultChan := make(chan *domain.LogEntry, 10)
	cleanup, err := h.monitorService.ListenLogs(ctx, routerID, topics, resultChan)
	if err != nil {
		conn.WriteJSON(gin.H{"error": err.Error()})
		return
	}
	defer cleanup()

	go func() {
		for {
			_, _, err := conn.ReadMessage()
			if err != nil {
				cancel()
				return
			}
		}
	}()

	for {
		select {
		case data := <-resultChan:
			if err := conn.WriteJSON(data); err != nil {
				return
			}
		case <-ctx.Done():
			return
		}
	}
}

// writeJSON writes JSON data to WebSocket connection with timeout
func writeJSON(conn *websocket.Conn, v interface{}) error {
	conn.SetWriteDeadline(time.Now().Add(10 * time.Second))
	return conn.WriteJSON(v)
}

// readJSON reads JSON data from WebSocket connection with timeout
func readJSON(conn *websocket.Conn, v interface{}) error {
	conn.SetReadDeadline(time.Now().Add(60 * time.Second))
	return conn.ReadJSON(v)
}
