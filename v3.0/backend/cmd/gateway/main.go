package main

import (
	"context"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"runtime"
	"syscall"
	"time"

	"github.com/Akashpg-M/polaris/backend/algo_/logger"
	"github.com/Akashpg-M/polaris/backend/internal/adapter/handler"
	"github.com/Akashpg-M/polaris/backend/internal/adapter/repository"
	"github.com/Akashpg-M/polaris/backend/internal/application/orchestration"
	"github.com/Akashpg-M/polaris/backend/internal/config"
	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
)

func main() {
	// 1. CONFIGURATION & LOGGING
	// Loads environment variables/flags and initializes structured JSON/text logging.
	cfg := config.Load()
	logger.Init()
	slog.Info("Booting Polaris v3.0 API Gateway...")

	// 2. INFRASTRUCTURE CLIENTS SETUP
	// Kafka is the durable ingress buffer where all raw telemetry is published.
	kafkaBroker := getEnvFallback("KAFKA_BROKER_URL", "localhost:9092")
	kafkaPublisher := repository.NewKafkaEventPublisher(kafkaBroker)
	defer kafkaPublisher.Close()

	// Redis client used for health checks and lease tracking.
	redisOptions, err := redis.ParseURL(cfg.Redis.URL)
	if err != nil {
		panic("invalid Redis URL: " + err.Error())
	}
	healthRedis := redis.NewClient(redisOptions)
	defer healthRedis.Close()

	// PostgreSQL Registry store used for authenticating incoming device connection tokens/tenants.
	registryStore, err := repository.NewRegistryStore(cfg.DB.URL)
	if err != nil {
		panic("Cannot start authenticated gateway without registry: " + err.Error())
	}
	defer registryStore.Close()

	// 3. REAL-TIME DASHBOARD BROADCASTER
	// Subscribes to Redis Pub/Sub topic "spatial:updates" and broadcasts live updates to frontend WebSockets.
	dashboardRegistry := handler.NewDashboardRegistry()
	go startDashboardSubscriber(cfg.Redis.URL, dashboardRegistry)

	// 4. DEVICE CONNECTION & LEASE MANAGEMENT
	// Tracks active WebSocket connections per gateway instance using distributed leases in Redis.
	orchestrationMetrics := orchestration.NewMetrics()
	connectionManager := handler.NewDeviceConnectionManager(
		getEnvFallback("GATEWAY_ID", "gateway-1"),
		getDurationFallback("CONNECTION_LEASE_TTL", 30*time.Second),
		healthRedis,
		registryStore,
		orchestrationMetrics,
	)
	subscriberCtx, stopSubscriber := context.WithCancel(context.Background())
	defer stopSubscriber()
	// Listens for cross-gateway routing commands targeting devices held by this gateway instance.
	go connectionManager.StartSubscriber(subscriberCtx)

	// 5. HTTP & WEBSOCKET ROUTING (GIN ENGINE)
	gin.SetMode(gin.ReleaseMode)
	router := gin.New()
	router.Use(gin.Recovery())

	// Global CORS middleware for web clients
	router.Use(func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization, X-Tenant-ID")
		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(http.StatusOK)
			return
		}
		c.Next()
	})

	// Ingestion handler upgrades HTTP to WebSocket for IoT devices and pushes data directly to Kafka.
	ingestionHandler := handler.NewIngestionHandler(kafkaPublisher, registryStore, connectionManager)
	dashboardHandler := handler.NewDashboardHandler(dashboardRegistry, registryStore)

	// Ingress WebSocket routes
	router.GET("/ws/telemetry", ingestionHandler.HandleIoTConnection)
	router.GET("/ws/dashboard", dashboardHandler.HandleWebConnection)

	// Gateway Metrics endpoints
	api := router.Group("/api/v1")
	api.GET("/metrics/connections", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"active_uplinks": ingestionHandler.GetActiveConnectionsCount()})
	})
	api.GET("/metrics/orchestration", func(c *gin.Context) {
		c.JSON(http.StatusOK, orchestrationMetrics.Snapshot())
	})

	// Kubernetes Liveness and Readiness probes
	router.GET("/healthz", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "live"})
	})
	router.GET("/readyz", func(c *gin.Context) {
		probeCtx, cancel := context.WithTimeout(c.Request.Context(), 2*time.Second)
		defer cancel()

		// Verify Kafka, Redis, and Postgres are alive before marking gateway ready for traffic.
		if err := kafkaPublisher.Ready(probeCtx, kafkaBroker); err != nil {
			c.JSON(http.StatusServiceUnavailable, gin.H{"status": "not_ready", "dependency": "kafka", "error": err.Error()})
			return
		}
		if err := healthRedis.Ping(probeCtx).Err(); err != nil {
			c.JSON(http.StatusServiceUnavailable, gin.H{"status": "not_ready", "dependency": "redis", "error": err.Error()})
			return
		}
		if err := registryStore.DB.PingContext(probeCtx); err != nil {
			c.JSON(http.StatusServiceUnavailable, gin.H{"status": "not_ready", "dependency": "registry", "error": err.Error()})
			return
		}

		dbStats := registryStore.DB.Stats()
		c.JSON(http.StatusOK, gin.H{
			"status": "ready",
			"runtime": gin.H{
				"goroutines":          runtime.NumGoroutine(),
				"db_open_connections": dbStats.OpenConnections,
				"db_in_use":           dbStats.InUse,
				"db_idle":             dbStats.Idle,
				"db_wait_count":       dbStats.WaitCount,
			},
		})
	})

	// 6. SERVER START & GRACEFUL SHUTDOWN
	port := ":" + cfg.Server.GatewayPort
	srv := &http.Server{Addr: port, Handler: router}
	go func() {
		slog.Info("Gateway active", "port", port)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			slog.Error("Server failed", "error", err)
		}
	}()

	// Intercept OS termination signals for clean WebSocket teardown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	slog.Warn("Shutdown signal received. Draining WebSockets...")
	stopSubscriber()

	ctxShutdown, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctxShutdown); err != nil {
		slog.Error("Gateway forced to shutdown", "error", err)
	}
	slog.Info("Gateway safely terminated.")
}

// Background Redis subscriber forwarding real-time coordinates to connected web dashboards.
func startDashboardSubscriber(redisURL string, dashboardRegistry *handler.DashboardRegistry) {
	opts, _ := redis.ParseURL(redisURL)
	client := redis.NewClient(opts)
	defer client.Close()
	pubsub := client.Subscribe(context.Background(), "spatial:updates")
	defer pubsub.Close()
	for msg := range pubsub.Channel() {
		dashboardRegistry.BroadcastToUIs(msg.Payload)
	}
}

func getEnvFallback(key, fallback string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return fallback
}

func getDurationFallback(key string, fallback time.Duration) time.Duration {
	if value := os.Getenv(key); value != "" {
		if parsed, err := time.ParseDuration(value); err == nil {
			return parsed
		}
	}
	return fallback
}