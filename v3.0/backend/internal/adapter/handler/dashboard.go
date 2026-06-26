package handler

import (
	"log/slog"
	"net/http"
	"sync"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

// DashboardRegistry tracks all active UI dashboard connections
type DashboardRegistry struct {
	mu          sync.RWMutex
	connections map[*websocket.Conn]bool
}

// NewDashboardRegistry initializes an empty thread-safe connection tracker
func NewDashboardRegistry() *DashboardRegistry {
	return &DashboardRegistry{
		connections: make(map[*websocket.Conn]bool),
	}
}

// Register adds a new UI dashboard connection to the active broadcast list
func (r *DashboardRegistry) Register(conn *websocket.Conn) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.connections[conn] = true
	slog.Info("[DashboardRegistry] New web client connected to telemetry stream", "active_dashboards", len(r.connections))
}

// Unregister safely drops a connection when a user closes the browser tab
func (r *DashboardRegistry) Unregister(conn *websocket.Conn) {
	r.mu.Lock()
	defer r.mu.Unlock()
	if _, exists := r.connections[conn]; exists {
		delete(r.connections, conn)
		conn.Close()
		slog.Info("[DashboardRegistry] Web client disconnected", "active_dashboards", len(r.connections))
	}
}

// BroadcastToUIs pumps a raw message string out to every single open dashboard browser concurrently
func (r *DashboardRegistry) BroadcastToUIs(payload string) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	for conn := range r.connections {
		// Send standard Text message to UIs (JSON strings)
		err := conn.WriteMessage(websocket.TextMessage, []byte(payload))
		if err != nil {
			slog.Warn("[DashboardRegistry] Failed to push frame down streaming channel, breaking pipe", "err", err)
			// Schedule cleanup asynchronously to prevent deadlocking the write lock
			go r.Unregister(conn)
		}
	}
}

// DashboardHandler provides the REST-to-WS upgrade entrypoint for web clients
type DashboardHandler struct {
	registry *DashboardRegistry
	upgrader websocket.Upgrader
}

// NewDashboardHandler constructs the gateway handler for web clients
func NewDashboardHandler(registry *DashboardRegistry) *DashboardHandler {
	return &DashboardHandler{
		registry: registry,
		upgrader: websocket.Upgrader{
			ReadBufferSize:  1024,
			WriteBufferSize: 1024,
			// Allow cross-origin requests so your local frontend can connect easily
			CheckOrigin: func(r *http.Request) bool {
				return true
			},
		},
	}
}

// HandleWebConnection converts incoming HTTP requests into an asynchronous JSON stream
func (h *DashboardHandler) HandleWebConnection(c *gin.Context) {
	conn, err := h.upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		slog.Error("[DashboardGateway] Handshake Upgrade Error", "error", err)
		return
	}

	h.registry.Register(conn)

	// Keep connection alive, listen for client-side closures
	go func() {
		defer h.registry.Unregister(conn)
		for {
			// Dashboards are consumer-only; if they send messages or close, clean up the pipe
			if _, _, err := conn.ReadMessage(); err != nil {
				break
			}
		}
	}()
}