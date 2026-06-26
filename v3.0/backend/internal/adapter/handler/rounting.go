package handler

import (
	"net/http"
	"strconv"

	"github.com/Akashpg-M/polaris/backend/algo_/graph"
	"github.com/gin-gonic/gin"
)

type RoutingHandler struct {
	network *graph.RoadNetwork
}

func NewRoutingHandler(network *graph.RoadNetwork) *RoutingHandler {
	return &RoutingHandler{network: network}
}

// Coordinate represents a point for the frontend map renderer
type Coordinate struct {
	Lat float64 `json:"lat"`
	Lon float64 `json:"lon"`
}

// CalculateRoute handles GET /api/v1/routes/calculate
func (h *RoutingHandler) CalculateRoute(c *gin.Context) {
	srcLat, err1 := strconv.ParseFloat(c.Query("src_lat"), 64)
	srcLon, err2 := strconv.ParseFloat(c.Query("src_lon"), 64)
	tgtLat, err3 := strconv.ParseFloat(c.Query("tgt_lat"), 64)
	tgtLon, err4 := strconv.ParseFloat(c.Query("tgt_lon"), 64)

	if err1 != nil || err2 != nil || err3 != nil || err4 != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid coordinate parameters"})
		return
	}

	// 1. Snap GPS to Graph Nodes
	srcNode, err := h.network.GetNearestIntersection(srcLat, srcLon)
	tgtNode, errTgt := h.network.GetNearestIntersection(tgtLat, tgtLon)
	
	if err != nil || errTgt != nil {
		c.JSON(http.StatusServiceUnavailable, gin.H{"error": "road network offline"})
		return
	}

	// 2. Execute Dijkstra's Algorithm
	pathResult, err := h.network.FindShortestPath(srcNode, tgtNode)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "no drivable route found between these locations"})
		return
	}

	// 3. Convert abstract Node IDs back to plottable Lat/Lon coordinates
	var polyline []Coordinate
	for _, nodeID := range pathResult.RouteNodes {
		lat, lon, ok := h.network.GetNodeCoordinates(nodeID)
		if ok {
			polyline = append(polyline, Coordinate{Lat: lat, Lon: lon})
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"status":          "success",
		"total_dist_km":   pathResult.TotalDistKm,
		"congestion_cost": pathResult.TotalCost,
		"route":           polyline,
	})
}