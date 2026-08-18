package main

import (
	"context"
	"encoding/json"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"

	"github.com/Akashpg-M/polaris/backend/algo_/graph"
	"github.com/Akashpg-M/polaris/backend/algo_/logger"
	"github.com/Akashpg-M/polaris/backend/internal/adapter/handler"
	"github.com/Akashpg-M/polaris/backend/internal/adapter/repository"
	"github.com/Akashpg-M/polaris/backend/internal/application/orchestrator"
	"github.com/Akashpg-M/polaris/backend/internal/application/outbox"
	"github.com/Akashpg-M/polaris/backend/internal/application/spatial"
	"github.com/Akashpg-M/polaris/backend/internal/application/stream"
	"github.com/Akashpg-M/polaris/backend/internal/application/twin"
	"github.com/Akashpg-M/polaris/backend/internal/config"
	"github.com/Akashpg-M/polaris/backend/internal/infra/osm"
	redisinfra "github.com/Akashpg-M/polaris/backend/internal/infra/redis"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
)

type RedisCommander struct {
	client *redis.Client
}

func (r *RedisCommander) SendCommand(nodeID string, payload interface{}) error {
	msg := map[string]interface{}{"node_id": nodeID, "command": payload}
	data, _ := json.Marshal(msg)
	return r.client.Publish(context.Background(), "telemetry:commands", string(data)).Err()
}

func main() {
	// 1. Initialize Config & Logger
	cfg := config.Load()
	logger.Init()
	slog.Info("Booting Polaris v3.0 Spatial Engine...", "env", cfg.App.Env)

	roadNetwork, err := osm.LoadRoadNetwork("data/chennai-metro.osm.pbf")
	if err != nil {
		slog.Warn("OSM Graph initialization failed. Routing API will be offline.", "error", err)
		// If you don't have the file yet, we initialize an empty graph so the server doesn't crash
		roadNetwork = graph.NewRoadNetwork()
	}

	engine := spatial.NewEngine()
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// 2. Initialize Dependencies (Using the nested Config structs)
	// redisConsumer, _ := stream.NewRedisConsumer(cfg.Redis.URL, engine)
	// go redisConsumer.Start(ctx, "engine-node-1")

	kafkaBroker := os.Getenv("KAFKA_BROKER_URL")
	if kafkaBroker == "" {
		kafkaBroker = "localhost:9092"
	}

	redisClient, err := redisinfra.NewClient(cfg.Redis.URL)
	if err != nil {
		panic("Cannot start engine without Redis: " + err.Error())
	}
	commander := &RedisCommander{client: redisClient}
	registryStore, err := repository.NewRegistryStore(cfg.DB.URL)
	if err != nil {
		panic("Cannot start registry: " + err.Error())
	}
	defer registryStore.Close()
	if err = registryStore.BootstrapPlatformAdmin(ctx, os.Getenv("DEV_PLATFORM_ADMIN_TOKEN")); err != nil {
		panic("Cannot bootstrap development operator: " + err.Error())
	}
	kafkaConsumer := stream.NewKafkaConsumer(kafkaBroker, engine, redisClient)
	go kafkaConsumer.Start(ctx, "engine-node-1")
	archiver, err := stream.NewKafkaPostgresArchiver(kafkaBroker, cfg.DB.URL)
	if err != nil {
		slog.Warn("Kafka/PostgreSQL Archiver offline", "error", err)
	} else {
		go archiver.Start(ctx)
	}
	if roadNetwork != nil {
		go stream.NewTrafficAnalyzer(kafkaBroker, roadNetwork).Start(ctx)
	}
	outboxRelay := outbox.New(registryStore, kafkaBroker, envInt("OUTBOX_BATCH_SIZE", 100), envDuration("OUTBOX_POLL_INTERVAL", 500*time.Millisecond))
	go outboxRelay.Start(ctx)
	connectivityDetector := twin.NewDetector(redisClient, kafkaBroker, envDuration("DEVICE_STALE_AFTER", 30*time.Second), envDuration("DEVICE_OFFLINE_AFTER", 90*time.Second), envDuration("OFFLINE_SCAN_INTERVAL", 10*time.Second))
	go connectivityDetector.Start(ctx)

	predictiveStrategy, err := orchestrator.NewPredictiveZoneStrategy(cfg.DB.URL)
	if err != nil {
		slog.Warn("Predictive Strategy offline. Falling back to Static Zones.")
		rebalancer := orchestrator.NewRebalancer(engine, commander, &orchestrator.StaticZoneStrategy{})
		go rebalancer.StartAutonomousLoop(ctx)
	} else {
		rebalancer := orchestrator.NewRebalancer(engine, commander, predictiveStrategy)
		go rebalancer.StartAutonomousLoop(ctx)
	}

	// 3. Setup Router
	gin.SetMode(gin.ReleaseMode)
	router := gin.New()
	router.Use(gin.Recovery(), cors.New(cors.Config{
		AllowOrigins: []string{"*"},
		AllowMethods: []string{http.MethodGet, http.MethodPost, http.MethodPatch, http.MethodPut, http.MethodDelete, http.MethodOptions},
		AllowHeaders: []string{"Origin", "Content-Type", "Authorization", "X-Tenant-ID"},
	}))

	matchHandler := handler.NewMatchHandler(engine)
	routingHandler := handler.NewRoutingHandler(roadNetwork)

	api := router.Group("/api/v1")
	{
		registryAPI := handler.NewRegistryAPI(registryStore, redisClient, envDuration("DEVICE_STALE_AFTER", 30*time.Second), envDuration("DEVICE_OFFLINE_AFTER", 90*time.Second), envDuration("CONNECTION_TICKET_TTL", 30*time.Second))
		registryAPI.Register(api)
		protected := api.Group("")
		protected.Use(registryAPI.Middleware("read"))
		protected.GET("/nodes/match", matchHandler.GetNearestNodes)
		protected.GET("/routes/calculate", routingHandler.CalculateRoute)
		protected.GET("/zones/predicted", func(c *gin.Context) {
			if predictiveStrategy != nil {
				c.JSON(200, gin.H{"status": "success", "data": predictiveStrategy.GetTargetZones(context.Background())})
			} else {
				c.JSON(200, gin.H{"status": "success", "data": []interface{}{}})
			}
		})
	}
	router.GET("/healthz", func(c *gin.Context) {
		if !kafkaConsumer.Healthy() || (archiver != nil && !archiver.Healthy()) {
			c.JSON(http.StatusServiceUnavailable, gin.H{"status": "unhealthy"})
			return
		}
		c.JSON(http.StatusOK, gin.H{"status": "live"})
	})
	router.GET("/readyz", func(c *gin.Context) {
		probeCtx, cancel := context.WithTimeout(c.Request.Context(), 2*time.Second)
		defer cancel()
		if err := kafkaConsumer.Ready(probeCtx); err != nil {
			c.JSON(http.StatusServiceUnavailable, gin.H{"status": "not_ready", "dependency": "kafka_or_redis", "error": err.Error()})
			return
		}
		if archiver == nil {
			c.JSON(http.StatusServiceUnavailable, gin.H{"status": "not_ready", "dependency": "postgres_archiver"})
			return
		}
		if err := registryStore.DB.PingContext(probeCtx); err != nil {
			c.JSON(http.StatusServiceUnavailable, gin.H{"status": "not_ready", "dependency": "registry", "error": err.Error()})
			return
		}
		if err := archiver.Ready(probeCtx); err != nil {
			c.JSON(http.StatusServiceUnavailable, gin.H{"status": "not_ready", "dependency": "postgres", "error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, gin.H{"status": "ready"})
	})

	// 4. Start Server with Graceful Shutdown
	port := ":" + cfg.Server.EnginePort
	srv := &http.Server{Addr: port, Handler: router}

	go func() {
		slog.Info("Engine REST API active", "port", port)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			slog.Error("Server failed", "error", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	slog.Warn("Shutdown signal received...")
	ctxShutdown, cancelShutdown := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancelShutdown()

	srv.Shutdown(ctxShutdown)
	cancel() // Stops background context (workers)
	if err := kafkaConsumer.Wait(ctxShutdown); err != nil {
		slog.Error("Spatial consumer shutdown timed out", "error", err)
	}
	if archiver != nil {
		if err := archiver.Wait(ctxShutdown); err != nil {
			slog.Error("Archive consumer shutdown timed out", "error", err)
		}
	}
	redisClient.Close()
	slog.Info("Engine safely terminated.")
}

func envDuration(key string, fallback time.Duration) time.Duration {
	if raw := os.Getenv(key); raw != "" {
		if v, err := time.ParseDuration(raw); err == nil {
			return v
		}
	}
	return fallback
}
func envInt(key string, fallback int) int {
	if raw := os.Getenv(key); raw != "" {
		if v, err := strconv.Atoi(raw); err == nil {
			return v
		}
	}
	return fallback
}
