package days

import (
	"cmp"
	"slices"
	"strconv"
	"strings"
)

type day08 struct {
	junctionBoxes []vec3
	connections   []connection // all pairwise connections, sorted by distance ascending
}

type vec3 struct {
	x, y, z int64
}

type connection struct {
	dist2 int64
	i, j  int
}

func init() {
	Register(8, func() Solution { return &day08{} })
}

// -----------------------------------------------------------
// Parsing
// -----------------------------------------------------------

// parseVec3 parses one X,Y,Z junction-box coordinate line and returns the 3D
// point used by the distance calculations.
func parseVec3(line string) vec3 {
	parts := strings.Split(line, ",")
	x, _ := strconv.ParseInt(parts[0], 10, 64)
	y, _ := strconv.ParseInt(parts[1], 10, 64)
	z, _ := strconv.ParseInt(parts[2], 10, 64)
	return vec3{x, y, z}
}

// SetInput parses junction-box coordinates and precomputes all sorted pairwise
// connections for the circuit-building algorithms.
func (d *day08) SetInput(lines []string) {
	d.junctionBoxes = d.junctionBoxes[:0]

	for _, ln := range lines {
		ln = strings.TrimSpace(ln)
		if ln == "" {
			continue
		}
		d.junctionBoxes = append(d.junctionBoxes, parseVec3(ln))
	}

	// Build and sort all pairwise connections once; reuse in both parts.
	d.connections = buildSortedConnections(d.junctionBoxes)
}

// -----------------------------------------------------------
// Distance & Edge Preparation
// -----------------------------------------------------------

// squaredDist returns the squared Euclidean distance between two junction boxes,
// avoiding square roots because only relative ordering is needed.
func squaredDist(a, b vec3) int64 {
	dx := a.x - b.x
	dy := a.y - b.y
	dz := a.z - b.z
	return dx*dx + dy*dy + dz*dz
}

// buildSortedConnections generates every pairwise connection among points,
// sorts them by ascending squared distance, and returns the sorted slice.
func buildSortedConnections(points []vec3) []connection {
	n := len(points)
	if n < 2 {
		return nil
	}

	connections := make([]connection, n*(n-1)/2)
	idx := 0
	for i := 0; i < n-1; i++ {
		pi := points[i]
		for j := i + 1; j < n; j++ {
			pj := points[j]
			connections[idx] = connection{
				dist2: squaredDist(pi, pj),
				i:     i,
				j:     j,
			}
			idx++
		}
	}

	radixSortConnections(connections)
	return connections
}

// radixSortConnections sorts connections in place by their non-negative squared
// distances and returns nothing.
func radixSortConnections(connections []connection) {
	if len(connections) < 2 {
		return
	}

	const (
		radixBits = 16
		buckets   = 1 << radixBits
		mask      = buckets - 1
	)

	tmp := make([]connection, len(connections))
	counts := make([]int, buckets)
	src := connections
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

// newDSU creates a union-find structure with n singleton circuits and returns it
// ready for component-size tracking.
func newDSU(n int) *dsu {
	p := make([]int, n)
	s := make([]int, n)
	for i := range p {
		p[i] = i
		s[i] = 1
	}
	return &dsu{p, s}
}

// find returns the representative circuit for x, compressing the path for
// faster future lookups.
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
// shortest k connections in the sorted list.
// Returns component sizes sorted descending.
func runConnections(points []vec3, connections []connection, k int) []int {
	n := len(points)
	if n == 0 {
		return nil
	}
	if k > len(connections) {
		k = len(connections)
	}

	uf := newDSU(n)

	for idx := 0; idx < k; idx++ {
		e := connections[idx]
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

	slices.SortFunc(sizes, func(a, b int) int { return cmp.Compare(b, a) })
	return sizes
}

// runUntilSingleCircuit keeps connecting shortest pairs until all
// junction boxes are in a single connected component. It returns the indices
// of the last pair that actually merged two different components.
func runUntilSingleCircuit(points []vec3, connections []connection) (int, int) {
	n := len(points)
	if n <= 1 {
		return 0, 0
	}

	uf := newDSU(n)
	components := n
	lastI, lastJ := 0, 0

	for _, e := range connections {
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

// SolvePart1 makes the first 1000 shortest connection attempts and returns the
// product of the three largest resulting circuit sizes.
func (d *day08) SolvePart1() string {
	sizes := runConnections(d.junctionBoxes, d.connections, 1000)
	if len(sizes) < 3 {
		return "0"
	}
	result := sizes[0] * sizes[1] * sizes[2]
	return strconv.Itoa(result)
}

// SolvePart2 connects circuits until all junction boxes share one circuit and
// returns the product of the x-coordinates from the final merging pair.
func (d *day08) SolvePart2() string {
	if len(d.junctionBoxes) < 2 {
		return "0"
	}
	i, j := runUntilSingleCircuit(d.junctionBoxes, d.connections)
	xa := d.junctionBoxes[i].x
	xb := d.junctionBoxes[j].x
	return strconv.FormatInt(xa*xb, 10)
}
