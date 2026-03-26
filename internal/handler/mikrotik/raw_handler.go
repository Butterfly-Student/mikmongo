package mikrotik

import (
	"context"
	"encoding/json"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	dto "mikmongo/internal/dto/mikrotik"
	"mikmongo/internal/service"
	"mikmongo/pkg/response"
	"mikmongo/pkg/ws"
)

// blockedCommands contains RouterOS command paths that should never be executed
// via the raw endpoint. These are destructive or security-sensitive operations.
var blockedCommands = []string{
	"/system/reboot",
	"/system/shutdown",
	"/system/reset-configuration",
	"/system/backup",
	"/system/export",
	"/user/add",
	"/user/set",
	"/user/remove",
	"/password",
	"/certificate",
}

// blockedVerbs contains RouterOS command verbs that are not allowed.
var blockedVerbs = []string{
	"remove",
	"set",
	"add",
	"disable",
	"enable",
	"reset",
	"uninstall",
	"move",
}

// RawHandler handles generic raw RouterOS command execution.
type RawHandler struct {
	routerSvc *service.RouterService
}

func NewRawHandler(routerSvc *service.RouterService) *RawHandler {
	return &RawHandler{routerSvc: routerSvc}
}

// validateRawArgs checks that the command args don't contain blocked commands or verbs.
func validateRawArgs(args []string) string {
	if len(args) == 0 {
		return "args is required"
	}
	if len(args) > 20 {
		return "too many args (max 20)"
	}

	// First arg is the command path
	cmd := strings.ToLower(args[0])

	for _, blocked := range blockedCommands {
		if strings.HasPrefix(cmd, blocked) {
			return "command not allowed: " + blocked
		}
	}

	// Check the last segment of the path for blocked verbs
	parts := strings.Split(cmd, "/")
	if len(parts) > 0 {
		verb := parts[len(parts)-1]
		for _, blocked := range blockedVerbs {
			if verb == blocked {
				return "command verb not allowed: " + blocked
			}
		}
	}

	return ""
}

// Run executes a raw RouterOS command and returns the results.
func (h *RawHandler) Run(c *gin.Context) {
	routerID, err := uuid.Parse(c.Param("router_id"))
	if err != nil {
		response.BadRequest(c, "invalid router_id")
		return
	}

	var req dto.RawCommandRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	if errMsg := validateRawArgs(req.Args); errMsg != "" {
		response.Forbidden(c, errMsg)
		return
	}

	mikClient, err := h.routerSvc.GetMikrotikClient(c.Request.Context(), routerID)
	if err != nil {
		response.InternalServerError(c, err.Error())
		return
	}

	results, err := mikClient.RunRaw(c.Request.Context(), req.Args)
	if err != nil {
		response.InternalServerError(c, err.Error())
		return
	}
	response.OK(c, results)
}

// ListenWS starts a raw RouterOS listen command via WebSocket.
// The client sends a JSON message with the args to listen on.
func (h *RawHandler) ListenWS(c *gin.Context) {
	routerID, err := uuid.Parse(c.Param("router_id"))
	if err != nil {
		return
	}

	conn, err := ws.Upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		return
	}
	defer conn.Close()

	ws.ConfigureConn(conn)

	// Read the initial command message from the client with a deadline
	conn.SetReadDeadline(time.Now().Add(ws.ReadWait))
	_, msg, err := conn.ReadMessage()
	if err != nil {
		return
	}

	var req dto.RawListenRequest
	if err := json.Unmarshal(msg, &req); err != nil {
		conn.WriteJSON(gin.H{"error": "invalid request: " + err.Error()})
		return
	}

	if errMsg := validateRawArgs(req.Args); errMsg != "" {
		conn.WriteJSON(gin.H{"error": errMsg})
		return
	}

	ctx, cancel := context.WithCancel(c.Request.Context())
	defer cancel()

	// Read goroutine — handle client disconnect with pong support
	go func() {
		defer cancel()
		conn.SetReadDeadline(time.Now().Add(ws.PongWait))
		conn.SetPongHandler(func(string) error {
			conn.SetReadDeadline(time.Now().Add(ws.PongWait))
			return nil
		})
		for {
			if _, _, err := conn.ReadMessage(); err != nil {
				return
			}
		}
	}()

	mikClient, err := h.routerSvc.GetMikrotikClient(ctx, routerID)
	if err != nil {
		conn.WriteJSON(gin.H{"error": err.Error()})
		return
	}

	ch := make(chan map[string]string, 50)
	cancelListen, err := mikClient.Conn().ListenRaw(ctx, req.Args, ch)
	if err != nil {
		conn.WriteJSON(gin.H{"error": err.Error()})
		return
	}
	defer cancelListen()

	for {
		select {
		case <-ctx.Done():
			return
		case data, ok := <-ch:
			if !ok {
				return
			}
			conn.SetWriteDeadline(time.Now().Add(ws.WriteWait))
			if err := conn.WriteJSON(data); err != nil {
				return
			}
		}
	}
}
