package osm

import (
	// "context"
	"io"
	"log/slog"
	"os"
	"runtime"

	"github.com/Akashpg-M/polaris/backend/algo_/geo"
	"github.com/Akashpg-M/polaris/backend/algo_/graph"
	"github.com/qedus/osmpbf"
)

// LoadRoadNetwork reads a .pbf file and populates the graph topology
func LoadRoadNetwork(pbfPath string) (*graph.RoadNetwork, error) {
	slog.Info("Parsing OSM Geographic Topology...", "file", pbfPath)

	f, err := os.Open(pbfPath)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	d := osmpbf.NewDecoder(f)
	d.SetBufferSize(osmpbf.MaxBlobSize)

	// Use all available CPU cores for concurrent decompression
	err = d.Start(runtime.GOMAXPROCS(-1))
	if err != nil {
		return nil, err
	}

	network := graph.NewRoadNetwork()

	// Step 1: We must temporarily store nodes in a standard map because
	// OSM Ways reference nodes by ID, but we need their Lat/Lon to calculate distance.
	tempNodes := make(map[int64]osmpbf.Node)

	var nc, wc uint64
	for {
		if v, err := d.Decode(); err == io.EOF {
			break
		} else if err != nil {
			return nil, err
		} else {
			switch v := v.(type) {
			case *osmpbf.Node:
				tempNodes[v.ID] = *v
				network.AddIntersection(v.ID, v.Lat, v.Lon)
				nc++
			case *osmpbf.Way:
				// Filter: We only care about drivable roads
				if isDrivableRoad(v.Tags) {
					processWay(network, tempNodes, v)
					wc++
				}
			}
		}
	}

	slog.Info("OSM Graph Construction Complete", "nodes_loaded", nc, "drivable_ways", wc)
	return network, nil
}

// isDrivableRoad filters out buildings, footpaths, and rivers
func isDrivableRoad(tags map[string]string) bool {
	if val, ok := tags["highway"]; ok {
		switch val {
		case "motorway", "trunk", "primary", "secondary", "tertiary", "unclassified", "residential":
			return true
		}
	}
	return false
}

// processWay calculates distances between intersections and binds the edges
func processWay(network *graph.RoadNetwork, nodes map[int64]osmpbf.Node, way *osmpbf.Way) {
	// Determine if the street is one-way
	isOneWay := false
	if val, ok := way.Tags["oneway"]; ok && val == "yes" {
		isOneWay = true
	}

	// Connect each sequential node in the street
	for i := 0; i < len(way.NodeIDs)-1; i++ {
		sourceID := way.NodeIDs[i]
		targetID := way.NodeIDs[i+1]

		sourceNode, sourceFound := nodes[sourceID]
		targetNode, targetFound := nodes[targetID]

		if sourceFound && targetFound {
			// Calculate real-world physical distance using Haversine from our geo package
			distKm := geo.Haversine(sourceNode.Lat, sourceNode.Lon, targetNode.Lat, targetNode.Lon)
			network.AddRoadSegment(sourceID, targetID, way.ID, distKm, isOneWay)
		}
	}
}