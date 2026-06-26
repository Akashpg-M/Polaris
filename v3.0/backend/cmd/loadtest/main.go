package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"math/rand"
	"os"
	"os/signal"
	"sync"
	"sync/atomic"
	"syscall"
	"time"

	"github.com/gorilla/websocket"
  pb "github.com/Akashpg-M/polaris/backend/api/proto/v1"
  "google.golang.org/protobuf/proto"
)

// Atomic counters for thread-safe, high-speed metrics tracking
var (
	activeConnections int64
	messagesSent      int64
	connectionErrors  int64
)

// Payload matches the Polaris domain model
type Payload struct {
	TenantID string  `json:"tenant_id"`
	NodeID   string  `json:"node_id"`
	Class    uint16  `json:"asset_class"`
	Lat      float64 `json:"lat"`
	Lon      float64 `json:"lon"`
	Status   string  `json:"status"`
	Battery  int     `json:"battery"`
}

func main() {
	// 1. Configurable Flags
	targetNodes := flag.Int("nodes", 1000, "Number of concurrent drones to simulate")
	serverURL := flag.String("url", "ws://localhost:6080/ws/telemetry", "Gateway WebSocket URL")
	rampRate := flag.Int("ramp", 100, "How many new connections to open per second")
	flag.Parse()

	log.Printf("🚀 Initiating Polaris Stress Test...")
	log.Printf("Targeting %d concurrent drones on %s", *targetNodes, *serverURL)

	var wg sync.WaitGroup
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// 2. Start the live metrics dashboard in the background
	go printMetricsDashboard(ctx)

	// 3. Smooth Connection Ramping
	// We calculate the delay between each dial to hit the requested ramp rate
	delayStr := fmt.Sprintf("%dms", 1000/(*rampRate))
	delay, _ := time.ParseDuration(delayStr)
	ticker := time.NewTicker(delay)
	defer ticker.Stop()

	// 4. Listen for Ctrl+C to stop the test gracefully
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	// Spawn the drones
SpawnLoop:
	for i := 1; i <= *targetNodes; i++ {
		select {
		case <-quit:
			log.Println("\n⚠️ Aborting launch sequence early...")
			break SpawnLoop
		case <-ticker.C:
			wg.Add(1)
			go simulateDrone(ctx, i, *serverURL, &wg)
		}
	}

	log.Println("\n✅ All requested drones deployed. Press Ctrl+C to terminate test.")
	<-quit // Wait here until user kills the script
	
	cancel() // Tell all drones to shut down
	wg.Wait() // Wait for them to actually close
	fmt.Println("\nStress test concluded cleanly.")
}

func simulateDrone(ctx context.Context, id int, wsURL string, wg *sync.WaitGroup) {
	defer wg.Done()

	conn, _, err := websocket.DefaultDialer.Dial(wsURL, nil)
	if err != nil {
		atomic.AddInt64(&connectionErrors, 1)
		return
	}
	defer conn.Close()

	atomic.AddInt64(&activeConnections, 1)
	defer atomic.AddInt64(&activeConnections, -1)

	nodeID := fmt.Sprintf("STRESS-DRONE-%d", id)
	lat := 13.04 + (rand.Float64() * 0.1)
	lon := 80.24 + (rand.Float64() * 0.1)

	go func() {
		for {
			if _, _, err := conn.ReadMessage(); err != nil {
				return
			}
		}
	}()

	ticker := time.NewTicker(1 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			lat += (rand.Float64() - 0.5) * 0.001
			lon += (rand.Float64() - 0.5) * 0.001

			// 1. Construct the Protobuf payload
			payload := &pb.SpatialObject{
				TenantId:      "alpha_logistics",
				Id:            nodeID,
				Type:          pb.NodeType_NODE_TYPE_DRONE,
				Lat:           lat,
				Lon:           lon,
				Status:        pb.NodeStatus_NODE_STATUS_ACTIVE,
				EnergyPercent: uint32(rand.Intn(100)),
			}

			// 2. Marshal to raw bytes using proto
			msgBytes, err := proto.Marshal(payload)
			if err != nil {
				atomic.AddInt64(&connectionErrors, 1)
				return
			}

			// 3. Send as a BinaryMessage
			if err := conn.WriteMessage(websocket.BinaryMessage, msgBytes); err != nil {
				atomic.AddInt64(&connectionErrors, 1)
				return
			}

			atomic.AddInt64(&messagesSent, 1)
		}
	}
}

// printMetricsDashboard clears the console line and prints live stats
func printMetricsDashboard(ctx context.Context) {
	ticker := time.NewTicker(1 * time.Second)
	defer ticker.Stop()
	var lastMessagesSent int64

	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			currentMessages := atomic.LoadInt64(&messagesSent)
			throughput := currentMessages - lastMessagesSent
			lastMessagesSent = currentMessages

			// \033[2K clears the current terminal line, \r returns cursor to the start
			fmt.Printf("\033[2K\r📡 ACTIVE UPLINKS: %d | ⚡ THROUGHPUT: %d msgs/sec | ❌ ERRORS: %d | 📦 TOTAL SENT: %d",
				atomic.LoadInt64(&activeConnections),
				throughput,
				atomic.LoadInt64(&connectionErrors),
				currentMessages,
			)
		}
	}
}