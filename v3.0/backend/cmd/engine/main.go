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
	"github.com/Akashpg-M/polaris/backend/internal/application/orchestrator"
	"github.com/Akashpg-M/polaris/backend/internal/application/spatial"
	"github.com/Akashpg-M/polaris/backend/internal/application/stream"
	"github.com/Akashpg-M/polaris/backend/internal/config"
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

	kafkaConsumer := stream.NewKafkaConsumer(kafkaBroker, engine)
	go kafkaConsumer.Start(ctx, "engine-node-1")

	archiver, err := stream.NewPostgresArchiver(cfg.Redis.URL, cfg.DB.URL)
	if err != nil {
		slog.Warn("PostgreSQL unreachable...")
	} else {
		go archiver.Start(ctx)
	}

	redisClient, err := redisinfra.NewClient(cfg.Redis.URL)
	if err != nil {
		panic("Cannot start engine without Redis: " + err.Error())
	}
	commander := &RedisCommander{client: redisClient}

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
	router.Use(gin.Recovery(), cors.Default())

	matchHandler := handler.NewMatchHandler(engine)
	api := router.Group("/api/v1")
	{
		api.GET("/nodes/match", matchHandler.GetNearestNodes)
		api.GET("/zones/predicted", func(c *gin.Context) {
			if predictiveStrategy != nil {
				c.JSON(200, gin.H{"status": "success", "data": predictiveStrategy.GetTargetZones(context.Background())})
			} else {
				c.JSON(200, gin.H{"status": "success", "data": []interface{}{}})
			}
		})
	}

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
	redisClient.Close()
	slog.Info("Engine safely terminated.")
}