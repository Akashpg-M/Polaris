package routing

import (
	"context"
	"io"
	"os"
	"runtime"
	"strconv"
	"strings"
	"time"

	"github.com/Akashpg-M/polaris/backend/internal/modules/mobility/model"
	"github.com/Akashpg-M/polaris/backend/internal/modules/mobility/spatial"
	"github.com/qedus/osmpbf"
)

type osmWay struct {
	id    int64
	nodes []int64
	tags  map[string]string
}

func LoadOSMPBF(ctx context.Context, path, version string) (*RoadGraph, error) {
	f, e := os.Open(path)
	if e != nil {
		return nil, e
	}
	defer f.Close()
	d := osmpbf.NewDecoder(f)
	if e = d.Start(runtime.GOMAXPROCS(0)); e != nil {
		return nil, e
	}
	nodes := map[int64]RoadNode{}
	ways := []osmWay{}
	for {
		if ctx.Err() != nil {
			return nil, ctx.Err()
		}
		v, e := d.Decode()
		if e == io.EOF {
			break
		}
		if e != nil {
			return nil, e
		}
		switch x := v.(type) {
		case *osmpbf.Node:
			nodes[x.ID] = RoadNode{ID: x.ID, Position: model.Position{Latitude: x.Lat, Longitude: x.Lon}, Type: NodeRoad}
		case *osmpbf.Way:
			if roadClass(x.Tags["highway"]) != "" {
				ways = append(ways, osmWay{x.ID, append([]int64(nil), x.NodeIDs...), x.Tags})
			}
		}
	}
	b := NewGraphBuilder(version)
	usedNodes := make(map[int64]struct{})
	for _, way := range ways {
		for _, id := range way.nodes {
			usedNodes[id] = struct{}{}
		}
	}
	for id := range usedNodes {
		if n, ok := nodes[id]; ok {
			_ = b.AddNode(n)
		}
	}
	var edgeID int64 = 1
	for _, w := range ways {
		class := roadClass(w.tags["highway"])
		speed := parseSpeed(w.tags["maxspeed"], defaultSpeed(class))
		reverseOnly := w.tags["oneway"] == "-1"
		oneway := reverseOnly || w.tags["oneway"] == "yes" || w.tags["oneway"] == "1" || w.tags["junction"] == "roundabout"
		for i := 0; i+1 < len(w.nodes); i++ {
			a, aok := nodes[w.nodes[i]]
			z, zok := nodes[w.nodes[i+1]]
			if !aok || !zok {
				continue
			}
			distance := spatial.DistanceMeters(a.Position, z.Position)
			from, to := a.ID, z.ID
			if reverseOnly {
				from, to = to, from
			}
			_ = b.AddEdge(RoadEdge{ID: EdgeID(edgeID), FromID: from, ToID: to, DistanceM: distance, BaseTravelTime: time.Duration(distance / speed * float64(time.Second)), RoadClass: class})
			edgeID++
			if !oneway {
				_ = b.AddEdge(RoadEdge{ID: EdgeID(edgeID), FromID: to, ToID: from, DistanceM: distance, BaseTravelTime: time.Duration(distance / speed * float64(time.Second)), RoadClass: class})
				edgeID++
			}
		}
	}
	return b.Build()
}
func roadClass(v string) RoadClass {
	switch v {
	case "motorway", "motorway_link", "trunk", "trunk_link":
		return RoadMotorway
	case "primary", "primary_link":
		return RoadPrimary
	case "secondary", "secondary_link":
		return RoadSecondary
	case "tertiary", "tertiary_link":
		return RoadTertiary
	case "residential", "living_street":
		return RoadResidential
	case "unclassified", "service":
		return RoadUnclassified
	}
	return ""
}
func defaultSpeed(c RoadClass) float64 {
	switch c {
	case RoadMotorway:
		return 27.8
	case RoadPrimary:
		return 19.4
	case RoadSecondary:
		return 16.7
	case RoadTertiary:
		return 13.9
	default:
		return 8.3
	}
}
func parseSpeed(raw string, fallback float64) float64 {
	if raw == "" {
		return fallback
	}
	value := strings.Fields(strings.Split(raw, ";")[0])
	if len(value) == 0 {
		return fallback
	}
	n, e := strconv.ParseFloat(value[0], 64)
	if e != nil || n <= 0 {
		return fallback
	}
	if strings.Contains(strings.ToLower(raw), "mph") {
		return n * .44704
	}
	return n / 3.6
}
