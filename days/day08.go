package days

import (
	"sort"
	"strconv"
	"strings"
)

type Day08 struct {
	points []vec3
	edges  []edge // all pairwise edges, sorted by distance ascending
}

type vec3 struct {
	x, y, z int64
}

type edge struct {
	dist2 int64
	i, j  int
}

func init() {
	Register(8, func() Solution { return &Day08{} })
}

// -----------------------------------------------------------
// Parsing
// -----------------------------------------------------------

func parseVec3(line string) vec3 {
	parts := strings.Split(line, ",")
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

	// Build and sort all pairwise edges once; reuse in both parts.
	d.edges = buildSortedEdges(d.points)
}

// -----------------------------------------------------------
// Distance & Edge Preparation
// -----------------------------------------------------------

func squaredDist(a, b vec3) int64 {
	dx := a.x - b.x
	dy := a.y - b.y
	dz := a.z - b.z
	return dx*dx + dy*dy + dz*dz
}

// Generate all edges sorted by ascending squared distance.
func buildSortedEdges(points []vec3) []edge {
	n := len(points)
	if n < 2 {
		return nil
	}

	edges := make([]edge, n*(n-1)/2)
	idx := 0
	for i := 0; i < n-1; i++ {
		pi := points[i]
		for j := i + 1; j < n; j++ {
			pj := points[j]
			edges[idx] = edge{
				dist2: squaredDist(pi, pj),
				i:     i,
				j:     j,
			}
			idx++
		}
	}

	radixSortEdges(edges)
	return edges
}

func radixSortEdges(edges []edge) {
	if len(edges) < 2 {
		return
	}

	const (
		radixBits = 16
		buckets   = 1 << radixBits
		mask      = buckets - 1
	)

	tmp := make([]edge, len(edges))
	counts := make([]int, buckets)
	src := edges
	dst := tmp

	for shift := uint(0); shift < 64; shift += radixBits {
		for i := range counts {
			counts[i] = 0
		}
		for _, e := range src {
			counts[(uint64(e.dist2)>>shift)&mask]++
		}

		sum := 0
		for i, count := range counts {
			counts[i] = sum
			sum += count
		}

		for _, e := range src {
			bucket := (uint64(e.dist2) >> shift) & mask
			dst[counts[bucket]] = e
			counts[bucket]++
		}

		src, dst = dst, src
	}
}

// -----------------------------------------------------------
// DSU (Union-Find) with component sizes
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

// union returns true if it actually merged two different components.
func (d *dsu) union(a, b int) bool {
	ra := d.find(a)
	rb := d.find(b)
	if ra == rb {
		return false
	}
	if d.size[ra] < d.size[rb] {
		ra, rb = rb, ra
	}
	d.parent[rb] = ra
	d.size[ra] += d.size[rb]
	return true
}

// -----------------------------------------------------------
// Core solver helpers (internal, testable)
// -----------------------------------------------------------

// runConnections performs exactly k connection attempts using the
// shortest k edges in the sorted edge list.
// Returns component sizes sorted descending.
func runConnections(points []vec3, edges []edge, k int) []int {
	n := len(points)
	if n == 0 {
		return nil
	}
	if k > len(edges) {
		k = len(edges)
	}

	uf := newDSU(n)

	for idx := 0; idx < k; idx++ {
		e := edges[idx]
		// If already in same component, nothing happens (but still counts as an attempt)
		uf.union(e.i, e.j)
	}

	seen := make([]bool, n)
	sizes := make([]int, 0, n)
	for i := range n {
		r := uf.find(i)
		if !seen[r] {
			seen[r] = true
			sizes = append(sizes, uf.size[r])
		}
	}

	// Sort descending
	sort.Slice(sizes, func(a, b int) bool { return sizes[a] > sizes[b] })
	return sizes
}

// runUntilSingleCircuit keeps connecting shortest edges until all
// points are in a single connected component. It returns the indices
// of the last pair that actually merged two different components.
func runUntilSingleCircuit(points []vec3, edges []edge) (int, int) {
	n := len(points)
	if n <= 1 {
		return 0, 0
	}

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
// Solve Part 1 & Part 2
// -----------------------------------------------------------

func (d *Day08) SolvePart1() string {
	sizes := runConnections(d.points, d.edges, 1000)
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
	i, j := runUntilSingleCircuit(d.points, d.edges)
	xa := d.points[i].x
	xb := d.points[j].x
	return strconv.FormatInt(xa*xb, 10)
}
