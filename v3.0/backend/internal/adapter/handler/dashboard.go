package handler

import (
	"context"
	"encoding/json"
	"log/slog"
	"net/http"
	"strings"
	"sync"

	"github.com/Akashpg-M/polaris/backend/internal/core/auth"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

// DashboardRegistry tracks all active UI dashboard connections
type DashboardRegistry struct {
	mu          sync.RWMutex
	connections map[*websocket.Conn]string
}

// NewDashboardRegistry initializes an empty thread-safe connection tracker
func NewDashboardRegistry() *DashboardRegistry {
	return &DashboardRegistry{
		connections: make(map[*websocket.Conn]string),
	}
}

// Register adds a new UI dashboard connection to the active broadcast list
func (r *DashboardRegistry) Register(conn *websocket.Conn, tenantID string) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.connections[conn] = tenantID
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

	var event struct {
		TenantID string `json:"tenant_id"`
	}
	if json.Unmarshal([]byte(payload), &event) != nil || event.TenantID == "" {
		return
	}
	for conn, tenantID := range r.connections {
		if tenantID != event.TenantID {
			continue
		}
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
	registry      *DashboardRegistry
	upgrader      websocket.Upgrader
	authenticator DashboardAuthenticator
}
type DashboardAuthenticator interface {
	ResolveOperator(context.Context, string) (auth.OperatorPrincipal, error)
	ConsumeOperatorTicket(context.Context, string) (auth.OperatorPrincipal, error)
}

// NewDashboardHandler constructs the gateway handler for web clients
func NewDashboardHandler(registry *DashboardRegistry, authenticator DashboardAuthenticator) *DashboardHandler {
	return &DashboardHandler{
		registry:      registry,
		authenticator: authenticator,
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
	principal, err := h.authenticate(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": gin.H{"code": "UNAUTHENTICATED", "message": "Dashboard authorization is required"}})
		return
	}
	tenantID := principal.TenantID
	if tenantID == "" {
		tenantID = c.GetHeader("X-Tenant-ID")
		if tenantID == "" {
			tenantID = c.Query("tenant_id")
		}
	}
	if tenantID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": gin.H{"code": "TENANT_REQUIRED", "message": "Tenant scope is required"}})
		return
	}
	conn, err := h.upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		slog.Error("[DashboardGateway] Handshake Upgrade Error", "error", err)
		return
	}

	h.registry.Register(conn, tenantID)

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
func (h *DashboardHandler) authenticate(c *gin.Context) (auth.OperatorPrincipal, error) {
	if ticket := c.Query("ticket"); ticket != "" {
		return h.authenticator.ConsumeOperatorTicket(c.Request.Context(), ticket)
	}
	header := c.GetHeader("Authorization")
	if !strings.HasPrefix(header, "Bearer ") {
		return auth.OperatorPrincipal{}, auth.ErrInvalidCredential
	}
	return h.authenticator.ResolveOperator(c.Request.Context(), strings.TrimSpace(strings.TrimPrefix(header, "Bearer ")))
}
