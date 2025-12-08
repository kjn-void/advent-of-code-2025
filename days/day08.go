package days

import (
	"sort"
	"strconv"
	"strings"
)

type Day08 struct {
	points []vec3
}

type vec3 struct {
	x, y, z int64
}

func init() {
	Register(8, func() Solution { return &Day08{} })
}

// -----------------------------------------------------------
// Parsing
// -----------------------------------------------------------

func parseVec3(line string) vec3 {
	parts := strings.Split(line, ",")
	// AoC input always valid integers
	x, _ := strconv.ParseInt(parts[0], 10, 64)
	y, _ := strconv.ParseInt(parts[1], 10, 64)
	z, _ := strconv.ParseInt(parts[2], 10, 64)
	return vec3{x, y, z}
}

func (d *Day08) SetInput(lines []string) {
	d.points = d.points[:0]
	for _, ln := range lines {
		ln = strings.TrimSpace(ln)
		if ln == "" {
			continue
		}
		d.points = append(d.points, parseVec3(ln))
	}
}

// -----------------------------------------------------------
// Distance & Edge Preparation
// -----------------------------------------------------------

type edge struct {
	dist2 int64 // squared distance
	i, j  int   // indices of points
}

func squaredDist(a, b vec3) int64 {
	dx := a.x - b.x
	dy := a.y - b.y
	dz := a.z - b.z
	return dx*dx + dy*dy + dz*dz
}

// Generate all edges sorted by ascending squared distance.
func buildSortedEdges(points []vec3) []edge {
	n := len(points)
	edges := make([]edge, 0, n*(n-1)/2)

	for i := 0; i < n; i++ {
		pi := points[i]
		for j := i + 1; j < n; j++ {
			pj := points[j]
			edges = append(edges, edge{
				dist2: squaredDist(pi, pj),
				i:     i,
				j:     j,
			})
		}
	}

	sort.Slice(edges, func(a, b int) bool {
		return edges[a].dist2 < edges[b].dist2
	})

	return edges
}

// -----------------------------------------------------------
// DSU (Union-Find) with Component Sizes
// -----------------------------------------------------------

type dsu struct {
	parent []int
	size   []int
}

func newDSU(n int) *dsu {
	p := make([]int, n)
	s := make([]int, n)
	for i := range p {
		p[i] = i
		s[i] = 1
	}
	return &dsu{p, s}
}

func (d *dsu) find(x int) int {
	for d.parent[x] != x {
		d.parent[x] = d.parent[d.parent[x]]
		x = d.parent[x]
	}
	return x
}

func (d *dsu) union(a, b int) bool {
	ra := d.find(a)
	rb := d.find(b)
	if ra == rb {
		return false
	}
	// union by size
	if d.size[ra] < d.size[rb] {
		ra, rb = rb, ra
	}
	d.parent[rb] = ra
	d.size[ra] += d.size[rb]
	return true
}

// -----------------------------------------------------------
// Core Solver Logic (testable internal helpers)
// -----------------------------------------------------------

// runConnections performs exactly k connection attempts
// using the shortest k edges in the complete graph.
// Returns component sizes sorted descending.
func runConnections(points []vec3, k int) []int {
	n := len(points)
	if n == 0 {
		return nil
	}

	edges := buildSortedEdges(points)
	if k > len(edges) {
		k = len(edges)
	}

	uf := newDSU(n)

	// Process the first k edges
	for idx := 0; idx < k; idx++ {
		e := edges[idx]
		// Try to union; if already connected, nothing happens
		uf.union(e.i, e.j)
	}

	// Count component sizes
	compMap := make(map[int]int)
	for i := 0; i < n; i++ {
		r := uf.find(i)
		compMap[r] = uf.size[r]
	}

	// Extract sizes
	sizes := make([]int, 0, len(compMap))
	for _, sz := range compMap {
		sizes = append(sizes, sz)
	}

	// Sort descending
	sort.Slice(sizes, func(a, b int) bool { return sizes[a] > sizes[b] })

	return sizes
}

// runUntilSingleCircuit keeps connecting shortest edges until all
// points are in a single connected component. It returns the indices
// of the last pair that actually merged two different components.
func runUntilSingleCircuit(points []vec3) (int, int) {
	n := len(points)
	if n <= 1 {
		return 0, 0
	}

	edges := buildSortedEdges(points)
	uf := newDSU(n)
	components := n

	lastI, lastJ := 0, 0

	for _, e := range edges {
		if uf.union(e.i, e.j) {
			components--
			lastI, lastJ = e.i, e.j
			if components == 1 {
				break
			}
		}
	}

	return lastI, lastJ
}

// -----------------------------------------------------------
// Solve Part 1 (k = 1000) and Part 2
// -----------------------------------------------------------

func (d *Day08) SolvePart1() string {
	sizes := runConnections(d.points, 1000)
	if len(sizes) < 3 {
		return "0"
	}
	result := sizes[0] * sizes[1] * sizes[2]
	return strconv.Itoa(result)
}

func (d *Day08) SolvePart2() string {
	if len(d.points) < 2 {
		return "0"
	}

	i, j := runUntilSingleCircuit(d.points)
	// Multiply X coordinates of the two junction boxes
	xa := d.points[i].x
	xb := d.points[j].x
	product := xa * xb
	return strconv.FormatInt(product, 10)
}
