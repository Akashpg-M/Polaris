package routing

import (
	"errors"
	"sort"
)

type RoadGraph struct {
	version     string
	nodes       map[int64]RoadNode
	edges       []RoadEdge
	adjacency   map[int64][]int
	incident    map[int64][]int
	nodeIndex   RoadNodeIndex
	maxSpeedMPS float64
}

func (g *RoadGraph) Version() string                { return g.version }
func (g *RoadGraph) NodeCount() int                 { return len(g.nodes) }
func (g *RoadGraph) EdgeCount() int                 { return len(g.edges) }
func (g *RoadGraph) Node(id int64) (RoadNode, bool) { v, ok := g.nodes[id]; return v, ok }
func (g *RoadGraph) Edge(index int) RoadEdge        { return g.edges[index] }
func (g *RoadGraph) Outgoing(id int64) []int        { return g.adjacency[id] }
func (g *RoadGraph) NodeIndex() RoadNodeIndex       { return g.nodeIndex }
func (g *RoadGraph) MaxSpeedMPS() float64           { return g.maxSpeedMPS }

type GraphBuilder struct {
	version string
	nodes   map[int64]RoadNode
	edges   []RoadEdge
	seen    map[[3]int64]struct{}
}

func NewGraphBuilder(version string) *GraphBuilder {
	return &GraphBuilder{version: version, nodes: map[int64]RoadNode{}, seen: map[[3]int64]struct{}{}}
}
func (b *GraphBuilder) AddNode(n RoadNode) error {
	if n.ID == 0 {
		return errors.New("invalid road node id")
	}
	if _, ok := b.nodes[n.ID]; ok {
		return nil
	}
	b.nodes[n.ID] = n
	return nil
}
func (b *GraphBuilder) AddEdge(e RoadEdge) error {
	if e.FromID == e.ToID || e.DistanceM <= 0 || e.BaseTravelTime <= 0 {
		return errors.New("invalid road edge")
	}
	if _, ok := b.nodes[e.FromID]; !ok {
		return ErrNoRoadNode
	}
	if _, ok := b.nodes[e.ToID]; !ok {
		return ErrNoRoadNode
	}
	k := [3]int64{e.FromID, e.ToID, int64(e.ID)}
	if _, ok := b.seen[k]; ok {
		return nil
	}
	b.seen[k] = struct{}{}
	b.edges = append(b.edges, e)
	return nil
}
func (b *GraphBuilder) Build() (*RoadGraph, error) {
	if len(b.nodes) == 0 {
		return nil, ErrNoRoadNode
	}
	sort.Slice(b.edges, func(i, j int) bool {
		if b.edges[i].FromID != b.edges[j].FromID {
			return b.edges[i].FromID < b.edges[j].FromID
		}
		if b.edges[i].ToID != b.edges[j].ToID {
			return b.edges[i].ToID < b.edges[j].ToID
		}
		return b.edges[i].ID < b.edges[j].ID
	})
	g := &RoadGraph{version: b.version, nodes: map[int64]RoadNode{}, edges: append([]RoadEdge(nil), b.edges...), adjacency: map[int64][]int{}, incident: map[int64][]int{}}
	nodes := make([]RoadNode, 0, len(b.nodes))
	for id, n := range b.nodes {
		g.nodes[id] = n
		nodes = append(nodes, n)
	}
	for i, e := range g.edges {
		g.adjacency[e.FromID] = append(g.adjacency[e.FromID], i)
		g.incident[e.FromID] = append(g.incident[e.FromID], i)
		g.incident[e.ToID] = append(g.incident[e.ToID], i)
		speed := e.DistanceM / e.BaseTravelTime.Seconds()
		if speed > g.maxSpeedMPS {
			g.maxSpeedMPS = speed
		}
	}
	for id, edges := range g.incident {
		if len(edges) > 2 {
			n := g.nodes[id]
			n.Type = NodeIntersection
			g.nodes[id] = n
		}
	}
	if g.maxSpeedMPS <= 0 {
		g.maxSpeedMPS = 40
	}
	g.nodeIndex = NewKDTree(nodes)
	return g, nil
}
