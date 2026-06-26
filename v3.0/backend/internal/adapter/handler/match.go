package handler

import (
	"net/http"
	"strconv"

	"github.com/Akashpg-M/polaris/backend/internal/application/spatial"
	pb "github.com/Akashpg-M/polaris/backend/api/proto/v1"
	"github.com/gin-gonic/gin"
)

type MatchHandler struct {
	engine *spatial.Engine
}

func NewMatchHandler(engine *spatial.Engine) *MatchHandler {
	return &MatchHandler{engine: engine}
}

// GetNearestNodes handles GET /api/v1/nodes/match
func (h *MatchHandler) GetNearestNodes(c *gin.Context) {
	tenantID := c.Query("tenant_id") 
	if tenantID == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "missing tenant identity"})
			return
	}
	
	lat, errLat := strconv.ParseFloat(c.Query("lat"), 64)
	lon, errLon := strconv.ParseFloat(c.Query("lon"), 64)
	radius, errRad := strconv.ParseFloat(c.DefaultQuery("radius_km", "10.0"), 64)
	assetClass, errClass := strconv.ParseInt(c.DefaultQuery("class", "2"), 10, 32)

	if errLat != nil || errLon != nil || errRad != nil || errClass != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid query parameters. lat and lon are required."})
		return
	}

	matches := h.engine.FindNearest(tenantID, lat, lon, radius, pb.NodeType(assetClass))	
	c.JSON(http.StatusOK, gin.H{
		"status": "success",
		"count":  len(matches),
		"data":   matches,
	})
}