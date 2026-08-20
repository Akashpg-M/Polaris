package routing

import (
	"context"
	"math"
	"sort"

	"github.com/Akashpg-M/polaris/backend/internal/modules/mobility/model"
)

type point3 struct {
	x, y, z float64
	n       RoadNode
}
type kdNode struct {
	p           point3
	axis        int
	left, right *kdNode
}
type KDTree struct{ root *kdNode }

func sphere(p model.Position) (x, y, z float64) {
	lat, lon := p.Latitude*math.Pi/180, p.Longitude*math.Pi/180
	return math.Cos(lat) * math.Cos(lon), math.Cos(lat) * math.Sin(lon), math.Sin(lat)
}
func NewKDTree(nodes []RoadNode) *KDTree {
	pts := make([]point3, len(nodes))
	for i, n := range nodes {
		pts[i].x, pts[i].y, pts[i].z = sphere(n.Position)
		pts[i].n = n
	}
	return &KDTree{root: buildKD(pts, 0)}
}
func buildKD(p []point3, depth int) *kdNode {
	if len(p) == 0 {
		return nil
	}
	axis := depth % 3
	value := func(v point3) float64 {
		if axis == 0 {
			return v.x
		}
		if axis == 1 {
			return v.y
		}
		return v.z
	}
	sort.Slice(p, func(i, j int) bool { return value(p[i]) < value(p[j]) })
	m := len(p) / 2
	return &kdNode{p: p[m], axis: axis, left: buildKD(append([]point3(nil), p[:m]...), depth+1), right: buildKD(append([]point3(nil), p[m+1:]...), depth+1)}
}
func (t *KDTree) Nearest(ctx context.Context, p model.Position) (RoadNode, error) {
	if t.root == nil {
		return RoadNode{}, ErrNoRoadNode
	}
	x, y, z := sphere(p)
	target := [3]float64{x, y, z}
	best := t.root.p
	bestD := math.Inf(1)
	var walk func(*kdNode)
	walk = func(n *kdNode) {
		if n == nil || ctx.Err() != nil {
			return
		}
		v := [3]float64{n.p.x, n.p.y, n.p.z}
		d := (v[0]-x)*(v[0]-x) + (v[1]-y)*(v[1]-y) + (v[2]-z)*(v[2]-z)
		if d < bestD {
			bestD, best = d, n.p
		}
		delta := target[n.axis] - v[n.axis]
		first, second := n.left, n.right
		if delta > 0 {
			first, second = second, first
		}
		walk(first)
		if delta*delta < bestD {
			walk(second)
		}
	}
	walk(t.root)
	if ctx.Err() != nil {
		return RoadNode{}, ctx.Err()
	}
	return best.n, nil
}
