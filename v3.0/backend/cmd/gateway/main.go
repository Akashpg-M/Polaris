package main

import (
	"context"
	"encoding/json"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/Akashpg-M/polaris/backend/algo_/logger"
	"github.com/Akashpg-M/polaris/backend/internal/adapter/handler"
	"github.com/Akashpg-M/polaris/backend/internal/adapter/repository"
	"github.com/Akashpg-M/polaris/backend/internal/config"
	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
)

func main() {
	cfg := config.Load()
	logger.Init()
	slog.Info("Booting Polaris v3.0 API Gateway...")

	kafkaBroker := getEnvFallback("KAFKA_BROKER_URL", "localhost:9092")
	kafkaAdapter := repository.NewKafkaStreamAdapter(kafkaBroker)
	defer kafkaAdapter.Close()

	// 1. Initialized Registries
	registry := handler.NewConnectionRegistry()
	dashboardRegistry := handler.NewDashboardRegistry() // FIXED

	// 2. Start Background Async Topic Consumers
	go startCommandSubscriber(cfg.Redis.URL, registry)
	go startDashboardSubscriber(cfg.Redis.URL, dashboardRegistry) // FIXED

	gin.SetMode(gin.ReleaseMode)

	router := gin.New()
	router.Use(gin.Recovery())

	router.Use(func(c *gin.Context) {
    c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
    c.Writer.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
    c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
    
    // Handle preflight OPTIONS requests instantly
    if c.Request.Method == "OPTIONS" {
        c.AbortWithStatus(200)
        return
    }
    c.Next()
	})
	// 3. Wired Handlers
	ingestionHandler := handler.NewIngestionHandler(kafkaAdapter, registry)
	dashboardHandler := handler.NewDashboardHandler(dashboardRegistry) // FIXED

	// Ingress endpoint (Binary inputs from edge hardware devices)
	router.GET("/ws/telemetry", ingestionHandler.HandleIoTConnection)

	// Egress streaming endpoint (JSON outputs out to browser tracking dashboards)
	router.GET("/ws/dashboard", dashboardHandler.HandleWebConnection) // FIXED

	api := router.Group("/api/v1")
	{
		api.GET("/metrics/connections", func(c *gin.Context) {
			c.JSON(200, gin.H{"active_uplinks": registry.GetActiveCount()})
		})
	}

	port := ":" + cfg.Server.GatewayPort
	srv := &http.Server{Addr: port, Handler: router}

	go func() {
		slog.Info("Gateway active", "port", port)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			slog.Error("Server failed", "error", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	slog.Warn("Shutdown signal received. Draining WebSockets...")
	ctxShutdown, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctxShutdown); err != nil {
		slog.Error("Gateway forced to shutdown", "error", err)
	}
	slog.Info("Gateway safely terminated.")
}

func startCommandSubscriber(redisURL string, registry *handler.ConnectionRegistry) {
	opts, _ := redis.ParseURL(redisURL)
	client := redis.NewClient(opts)
	pubsub := client.Subscribe(context.Background(), "telemetry:commands")
	defer pubsub.Close()

	for msg := range pubsub.Channel() {
		var payload struct {
			NodeID  string      `json:"node_id"`
			Command interface{} `json:"command"`
		}
		if err := json.Unmarshal([]byte(msg.Payload), &payload); err == nil {
			registry.SendCommand(payload.NodeID, payload.Command)
		}
	}
}

func startDashboardSubscriber(redisURL string, dashboardRegistry *handler.DashboardRegistry) {
	opts, _ := redis.ParseURL(redisURL)
	client := redis.NewClient(opts)
	pubsub := client.Subscribe(context.Background(), "spatial:updates")
	defer pubsub.Close()

	for msg := range pubsub.Channel() {
		// Broadcast the JSON string instantly to all open web dashboard browsers
		dashboardRegistry.BroadcastToUIs(msg.Payload)
	}
}

func getEnvFallback(key, fallback string) string {
	if val, exists := os.LookupEnv(key); exists && val != "" {
		return val
	}
	return fallback
}