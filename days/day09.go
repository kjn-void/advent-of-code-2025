package days

import (
	"strconv"
	"strings"
)

type Day09 struct {
	reds      []pt9
	edges     []edge9
	vertEdges []edge9 // cached vertical edges for fast inside test
	horEdges  []edge9
}

type pt9 struct {
	x, y int
}

type edge9 struct {
	x1, y1 int
	x2, y2 int
	hor    bool
}

func init() {
	Register(9, func() Solution { return &Day09{} })
}

func (d *Day09) SetInput(lines []string) {
	d.reds = d.reds[:0]
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		parts := strings.Split(line, ",")
		x, _ := strconv.Atoi(parts[0])
		y, _ := strconv.Atoi(parts[1])
		d.reds = append(d.reds, pt9{x, y})
	}
	d.edges = nil
	d.vertEdges = nil
	d.horEdges = nil
}

// ----------------------------------------------------------
// Part 1 - largest rectangle from any two red tiles
// ----------------------------------------------------------

func (d *Day09) SolvePart1() string {
	best := maxAreaInclusive(d.reds)
	return strconv.Itoa(best)
}

func maxAreaInclusive(points []pt9) int {
	n := len(points)
	if n < 2 {
		return 0
	}
	best := 0
	for i := 0; i < n; i++ {
		for j := i + 1; j < n; j++ {
			dx := absInt(points[i].x-points[j].x) + 1
			dy := absInt(points[i].y-points[j].y) + 1
			area := dx * dy
			if area > best {
				best = area
			}
		}
	}
	return best
}

// ----------------------------------------------------------
// Part 2 - rectangles fully inside orthogonal polygon
// ----------------------------------------------------------

func (d *Day09) SolvePart2() string {
	n := len(d.reds)
	if n < 2 {
		return "0"
	}
	if d.edges == nil {
		d.buildEdges()
	}

	candidates := make([]rectCandidate9, 0, n*(n-1)/2)

	for i := 0; i < n; i++ {
		a := d.reds[i]
		for j := i + 1; j < n; j++ {
			b := d.reds[j]

			x1 := minInt(a.x, b.x)
			x2 := maxInt(a.x, b.x)
			y1 := minInt(a.y, b.y)
			y2 := maxInt(a.y, b.y)

			dx := x2 - x1 + 1
			dy := y2 - y1 + 1
			candidates = append(candidates, rectCandidate9{
				x1:   x1,
				y1:   y1,
				x2:   x2,
				y2:   y2,
				area: dx * dy,
			})
		}
	}

	radixSortRectCandidatesByArea(candidates)

	insideCache := make(map[pt9]bool, n*n/2)
	for i := len(candidates) - 1; i >= 0; i-- {
		c := candidates[i]
		// The other two corners must be inside or on the polygon.
		if !d.pointInsideOrOnCached(pt9{c.x1, c.y2}, insideCache) ||
			!d.pointInsideOrOnCached(pt9{c.x2, c.y1}, insideCache) {
			continue
		}

		// Ensure no polygon edge cuts through the interior of this rectangle.
		if d.rectangleCutByPolygon(c.x1, c.y1, c.x2, c.y2) {
			continue
		}

		return strconv.Itoa(c.area)
	}

	return "0"
}

type rectCandidate9 struct {
	x1, y1 int
	x2, y2 int
	area   int
}

func radixSortRectCandidatesByArea(candidates []rectCandidate9) {
	if len(candidates) < 2 {
		return
	}

	const (
		radixBits = 16
		buckets   = 1 << radixBits
		mask      = buckets - 1
	)

	tmp := make([]rectCandidate9, len(candidates))
	counts := make([]int, buckets)
	src := candidates
	dst := tmp

	for shift := uint(0); shift < 64; shift += radixBits {
		for i := range counts {
			counts[i] = 0
		}
		for _, c := range src {
			counts[(uint64(c.area)>>shift)&mask]++
		}

		sum := 0
		for i, count := range counts {
			counts[i] = sum
			sum += count
		}

		for _, c := range src {
			bucket := (uint64(c.area) >> shift) & mask
			dst[counts[bucket]] = c
			counts[bucket]++
		}

		src, dst = dst, src
	}
}

// ----------------------------------------------------------
// Polygon edges
// ----------------------------------------------------------

func (d *Day09) buildEdges() {
	n := len(d.reds)
	edges := make([]edge9, 0, n)
	verts := make([]edge9, 0, n)
	hors := make([]edge9, 0, n)

	for i := 0; i < n; i++ {
		a := d.reds[i]
		b := d.reds[(i+1)%n]
		e := edge9{x1: a.x, y1: a.y, x2: b.x, y2: b.y}

		if a.y == b.y {
			e.hor = true
			if e.x1 > e.x2 {
				e.x1, e.x2 = e.x2, e.x1
			}
			// y1 == y2 already
			hors = append(hors, e)
		} else {
			e.hor = false
			// normalize to y1 < y2 for half-open tests
			if e.y1 > e.y2 {
				e.y1, e.y2 = e.y2, e.y1
			}
			// x1 == x2 already
			verts = append(verts, e)
		}

		edges = append(edges, e)
	}

	d.edges = edges
	d.vertEdges = verts
	d.horEdges = hors
}

// ----------------------------------------------------------
// Point in polygon (inside or on boundary)
// ----------------------------------------------------------

func (d *Day09) pointInsideOrOn(p pt9) bool {
	// Boundary check (allowed):
	for _, e := range d.horEdges {
		if p.y == e.y1 && p.x >= e.x1 && p.x <= e.x2 {
			return true
		}
	}
	for _, e := range d.vertEdges {
		if p.x == e.x1 && p.y >= e.y1 && p.y <= e.y2 {
			return true
		}
	}

	// Strict interior test using integer-only vertical-edge crossing.
	return d.pointInsideStrict(p)
}

func (d *Day09) pointInsideOrOnCached(p pt9, cache map[pt9]bool) bool {
	if inside, ok := cache[p]; ok {
		return inside
	}
	inside := d.pointInsideOrOn(p)
	cache[p] = inside
	return inside
}

// Integer-only odd-even rule specialized for orthogonal polygons.
// Count vertical edges strictly to the right of p that cross scanline y = p.y.
// Use half-open interval [y1, y2) to avoid double counting at vertices.
func (d *Day09) pointInsideStrict(p pt9) bool {
	crossings := 0
	for _, e := range d.vertEdges {
		// vertical edge at x = e.x1, y in [e.y1, e.y2)
		if e.x1 > p.x && e.y1 <= p.y && p.y < e.y2 {
			crossings++
		}
	}
	return (crossings & 1) == 1
}

// ----------------------------------------------------------
// Check if polygon cuts through rectangle interior
// ----------------------------------------------------------
//
// Rectangle is [x1,x2] × [y1,y2], inclusive corner tiles.
// We treat the *interior* as (x1,x2) × (y1,y2) (open intervals).
func (d *Day09) rectangleCutByPolygon(x1, y1, x2, y2 int) bool {
	if x1 == x2 || y1 == y2 {
		// Degenerate rectangles have no interior.
		return false
	}

	for _, e := range d.horEdges {
		y0 := e.y1
		if y0 <= y1 || y0 >= y2 {
			continue
		}
		// overlap with (x1,x2)
		if maxInt(e.x1, x1) < minInt(e.x2, x2) {
			return true
		}
	}
	for _, e := range d.vertEdges {
		x0 := e.x1
		if x0 <= x1 || x0 >= x2 {
			continue
		}
		// overlap with (y1,y2)
		if maxInt(e.y1, y1) < minInt(e.y2, y2) {
			return true
		}
	}
	return false
}

// ----------------------------------------------------------
// Helpers
// ----------------------------------------------------------

func absInt(a int) int {
	if a < 0 {
		return -a
	}
	return a
}

func minInt(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func maxInt(a, b int) int {
	if a > b {
		return a
	}
	return b
}
