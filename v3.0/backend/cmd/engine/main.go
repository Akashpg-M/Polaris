package main

import (
	"context"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"runtime"
	"strconv"
	"syscall"
	"time"

	"github.com/Akashpg-M/polaris/backend/algo_/logger"
	"github.com/Akashpg-M/polaris/backend/internal/adapter/handler"
	"github.com/Akashpg-M/polaris/backend/internal/adapter/repository"
	"github.com/Akashpg-M/polaris/backend/internal/application/dispatch"
	"github.com/Akashpg-M/polaris/backend/internal/application/orchestration"
	"github.com/Akashpg-M/polaris/backend/internal/application/orchestrator"
	"github.com/Akashpg-M/polaris/backend/internal/application/outbox"
	"github.com/Akashpg-M/polaris/backend/internal/application/reconciliation"
	"github.com/Akashpg-M/polaris/backend/internal/application/spatial"
	"github.com/Akashpg-M/polaris/backend/internal/application/stream"
	"github.com/Akashpg-M/polaris/backend/internal/application/twin"
	"github.com/Akashpg-M/polaris/backend/internal/config"
	"github.com/Akashpg-M/polaris/backend/internal/core/extension"
	redisinfra "github.com/Akashpg-M/polaris/backend/internal/infra/redis"
	mobilitymodule "github.com/Akashpg-M/polaris/backend/internal/modules/mobility"
	"github.com/Akashpg-M/polaris/backend/internal/modules/mobility/matching"
	"github.com/Akashpg-M/polaris/backend/internal/modules/mobility/planning"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

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

	redisClient, err := redisinfra.NewClient(cfg.Redis.URL)
	if err != nil {
		panic("Cannot start engine without Redis: " + err.Error())
	}
	registryStore, err := repository.NewRegistryStore(cfg.DB.URL)
	if err != nil {
		panic("Cannot start registry: " + err.Error())
	}
	defer registryStore.Close()
	if err = registryStore.BootstrapPlatformAdmin(ctx, os.Getenv("DEV_PLATFORM_ADMIN_TOKEN")); err != nil {
		panic("Cannot bootstrap development operator: " + err.Error())
	}
	extensionRegistry := extension.NewRegistry()
	mobilityCfg, err := mobilitymodule.LoadConfig()
	if err != nil {
		panic("Invalid Mobility configuration: " + err.Error())
	}
	var mobilityModule *mobilitymodule.Module
	stateFanout := &stream.StateFanout{Primary: engine}
	if mobilityCfg.Enabled {
		mobilityModule = mobilitymodule.New(mobilityCfg, mobilityRebuildLoader(redisClient, registryStore))
		extensionRegistry.RegisterModule(mobilityModule)
		if mobilityCfg.SpatialEnabled {
			extensionRegistry.RegisterCandidateProvider(&matching.Provider{Spatial: mobilityModule.Spatial, Routing: mobilityModule, RawLimit: mobilityCfg.MaxRawCandidates, RoutedLimit: mobilityCfg.MaxRoutedCandidates, MaxRadius: mobilityCfg.MaxSearchRadiusMeters})
		}
		extensionRegistry.RegisterTaskPlanner(&planning.Planner{SpatialState: mobilityModule.Spatial.Get, Routing: mobilityModule, MaxPlanAge: 2 * time.Minute})
		if mobilityCfg.SpatialEnabled {
			stateFanout.Projections = append(stateFanout.Projections, &mobilitymodule.TelemetryProjector{Manager: mobilityModule.Spatial})
		}
	}
	extensionRegistry.RegisterTaskPlanner(extension.DefaultTaskPlanner{})
	if err = extensionRegistry.Start(ctx); err != nil {
		panic("Cannot start capability modules: " + err.Error())
	}
	var mobilityTraffic *mobilitymodule.TrafficConsumer
	if mobilityModule != nil && mobilityModule.Traffic() != nil {
		mobilityTraffic = mobilitymodule.NewTrafficConsumer(kafkaBroker, mobilityModule.Traffic(), mobilityCfg.TrafficRefreshInterval)
		go mobilityTraffic.Start(ctx)
	}
	kafkaConsumer := stream.NewKafkaConsumer(kafkaBroker, stateFanout, redisClient)
	go kafkaConsumer.Start(ctx, "engine-node-1")
	archiver, err := stream.NewKafkaPostgresArchiver(kafkaBroker, cfg.DB.URL)
	if err != nil {
		slog.Warn("Kafka/PostgreSQL Archiver offline", "error", err)
	} else {
		go archiver.Start(ctx)
	}
	outboxRelay := outbox.New(registryStore, kafkaBroker, envInt("OUTBOX_BATCH_SIZE", 100), envDuration("OUTBOX_POLL_INTERVAL", 500*time.Millisecond))
	go outboxRelay.Start(ctx)
	orchestrationMetrics := orchestration.NewMetrics()
	orchestrationService := orchestration.NewServiceWithRegistry(registryStore, redisClient, envInt("COMMAND_MAX_ATTEMPTS", 5), orchestrationMetrics, extensionRegistry)
	ownershipStore := repository.NewConnectionOwnershipStore(redisClient, envDuration("CONNECTION_LEASE_TTL", 30*time.Second))
	commandDispatcher := dispatch.New(kafkaBroker, redisClient, ownershipStore)
	go commandDispatcher.Start(ctx)
	commandReconciler := reconciliation.New(registryStore, orchestrationService, ownershipStore, envDuration("COMMAND_RECONCILE_INTERVAL", time.Second), envDuration("COMMAND_ACK_TIMEOUT", 5*time.Second))
	go commandReconciler.Start(ctx)
	connectivityDetector := twin.NewDetector(redisClient, kafkaBroker, envDuration("DEVICE_STALE_AFTER", 30*time.Second), envDuration("DEVICE_OFFLINE_AFTER", 90*time.Second), envDuration("OFFLINE_SCAN_INTERVAL", 10*time.Second))
	if mobilityModule != nil {
		connectivityDetector.SetTransitionHandler(func(tenant, device, status string) {
			_ = mobilityModule.Spatial.EvictInactive(tenant, device, "ACTIVE", status)
		})
	}
	go connectivityDetector.Start(ctx)

	predictiveStrategy, err := orchestrator.NewPredictiveZoneStrategy(cfg.DB.URL)
	if err != nil {
		slog.Warn("Predictive zone view offline", "error", err)
	} else {
		defer predictiveStrategy.Close()
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

	api := router.Group("/api/v1")
	{
		registryAPI := handler.NewRegistryAPI(registryStore, redisClient, envDuration("DEVICE_STALE_AFTER", 30*time.Second), envDuration("DEVICE_OFFLINE_AFTER", 90*time.Second), envDuration("CONNECTION_TICKET_TTL", 30*time.Second))
		if mobilityModule != nil && mobilityCfg.SpatialEnabled {
			registryAPI.SetLifecycleHook(func(tenant, device, status string) {
				if device == "" {
					if status != "ACTIVE" {
						_ = mobilityModule.Spatial.RemoveTenant(tenant)
					}
				} else {
					_ = mobilityModule.Spatial.EvictInactive(tenant, device, status, "ONLINE")
				}
			})
		}
		registryAPI.Register(api)
		handler.NewOrchestrationAPI(registryStore, orchestrationService, orchestrationMetrics).Register(api, registryAPI)
		protected := api.Group("")
		protected.Use(registryAPI.Middleware("read"))
		protected.GET("/nodes/match", matchHandler.GetNearestNodes)
		if mobilityModule != nil {
			handler.NewMobilityAPI(mobilityModule.Spatial, mobilityModule, mobilityCfg.MaxRawCandidates).Register(protected)
		}
		protected.GET("/zones/predicted", func(c *gin.Context) {
			if predictiveStrategy != nil {
				c.JSON(200, gin.H{"status": "success", "data": predictiveStrategy.GetTargetZones(context.Background())})
			} else {
				c.JSON(200, gin.H{"status": "success", "data": []interface{}{}})
			}
		})
	}
	router.GET("/healthz", func(c *gin.Context) {
		if !kafkaConsumer.Healthy() || !commandDispatcher.Healthy() || (archiver != nil && !archiver.Healthy()) {
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
		dbStats := registryStore.DB.Stats()
		c.JSON(http.StatusOK, gin.H{"status": "ready", "core": "READY", "modules": extensionRegistry.Status(probeCtx), "runtime": gin.H{"goroutines": runtime.NumGoroutine(), "db_open_connections": dbStats.OpenConnections, "db_in_use": dbStats.InUse, "db_idle": dbStats.Idle, "db_wait_count": dbStats.WaitCount}})
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
	if err := extensionRegistry.Close(ctxShutdown); err != nil {
		slog.Error("Capability module shutdown failed", "error", err)
	}
	if err := kafkaConsumer.Wait(ctxShutdown); err != nil {
		slog.Error("Spatial consumer shutdown timed out", "error", err)
	}
	if archiver != nil {
		if err := archiver.Wait(ctxShutdown); err != nil {
			slog.Error("Archive consumer shutdown timed out", "error", err)
		}
	}
	if err := commandDispatcher.Wait(ctxShutdown); err != nil {
		slog.Error("Command dispatcher shutdown timed out", "error", err)
	}
	if mobilityTraffic != nil {
		if err := mobilityTraffic.Wait(ctxShutdown); err != nil {
			slog.Error("Mobility traffic consumer shutdown timed out", "error", err)
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
