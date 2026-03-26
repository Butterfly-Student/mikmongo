package ws

import (
	"context"
	"regexp"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

// validIfaceName matches valid MikroTik interface names (alphanumeric, dash, dot, underscore).
var validIfaceName = regexp.MustCompile(`^[a-zA-Z0-9_\-\.]{1,64}$`)

// ValidateInterfaceName checks if a MikroTik interface name is safe.
func ValidateInterfaceName(name string) bool {
	return validIfaceName.MatchString(name)
}

// StreamConfig holds the configuration for a WebSocket stream.
type StreamConfig struct {
	Conn   *websocket.Conn
	Ctx    context.Context
	Cancel context.CancelFunc
}

// UpgradeAndConfigure upgrades the HTTP connection to WebSocket and sets up
// read limits, deadlines, and a read-pump goroutine for disconnect detection.
// Returns nil StreamConfig if upgrade fails.
func UpgradeAndConfigure(c *gin.Context) *StreamConfig {
	conn, err := Upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		return nil
	}

	ConfigureConn(conn)

	ctx, cancel := context.WithCancel(c.Request.Context())

	// Read pump goroutine — detect client disconnect and pong responses
	go func() {
		defer cancel()
		for {
			if _, _, err := conn.ReadMessage(); err != nil {
				return
			}
		}
	}()

	return &StreamConfig{
		Conn:   conn,
		Ctx:    ctx,
		Cancel: cancel,
	}
}

// ForwardChannel reads from a channel and writes each item as JSON to the WebSocket.
// Blocks until ctx is done or channel is closed. Caller must defer conn.Close() and cancel().
func ForwardChannel[T any](sc *StreamConfig, ch <-chan T) {
	for {
		select {
		case <-sc.Ctx.Done():
			return
		case data, ok := <-ch:
			if !ok {
				return
			}
			sc.Conn.SetWriteDeadline(time.Now().Add(WriteWait))
			if err := sc.Conn.WriteJSON(data); err != nil {
				return
			}
		}
	}
}
