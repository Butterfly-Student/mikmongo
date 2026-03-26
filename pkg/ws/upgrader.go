package ws

import (
	"net/http"
	"strings"
	"time"

	"github.com/gorilla/websocket"
)

const (
	// MaxMessageSize is the maximum size of incoming WebSocket messages (4KB).
	MaxMessageSize = 4096

	// WriteWait is the time allowed to write a message to the peer.
	WriteWait = 10 * time.Second

	// PongWait is the time allowed to read the next pong message from the peer.
	PongWait = 60 * time.Second

	// ReadWait is the time allowed to read the initial message from the peer.
	ReadWait = 30 * time.Second
)

// allowedOrigins stores the list of allowed origins for WebSocket connections.
// Empty means only same-origin is allowed.
var allowedOrigins []string

// SetAllowedOrigins configures the allowed origins for WebSocket connections.
// Pass "*" to allow all origins (development only).
func SetAllowedOrigins(origins []string) {
	allowedOrigins = origins
}

// Upgrader is the default WebSocket upgrader for the application.
var Upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin:     checkOrigin,
}

func checkOrigin(r *http.Request) bool {
	origin := r.Header.Get("Origin")
	if origin == "" {
		return true // same-origin or non-browser client
	}

	for _, allowed := range allowedOrigins {
		if allowed == "*" {
			return true
		}
		if strings.EqualFold(origin, allowed) {
			return true
		}
	}
	return false
}

// ConfigureConn sets read limits and deadlines on a WebSocket connection.
func ConfigureConn(conn *websocket.Conn) {
	conn.SetReadLimit(MaxMessageSize)
	conn.SetReadDeadline(time.Now().Add(PongWait))
	conn.SetPongHandler(func(string) error {
		conn.SetReadDeadline(time.Now().Add(PongWait))
		return nil
	})
}
