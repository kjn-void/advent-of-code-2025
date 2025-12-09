package days

import (
	"strconv"
	"strings"
)

type Day09 struct {
	reds  []pt9
	edges []edge9
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

	best := 0

	for i := range n {
		a := d.reds[i]
		for j := i + 1; j < n; j++ {
			b := d.reds[j]

			x1 := minInt(a.x, b.x)
			x2 := maxInt(a.x, b.x)
			y1 := minInt(a.y, b.y)
			y2 := maxInt(a.y, b.y)

			dx := x2 - x1 + 1
			dy := y2 - y1 + 1
			area := dx * dy

			if area <= best {
				continue
			}

			// The other two corners must be inside or on the polygon.
			c3 := pt9{x1, y2}
			c4 := pt9{x2, y1}

			if !d.pointInsideOrOn(c3) || !d.pointInsideOrOn(c4) {
				continue
			}

			// Ensure no polygon edge cuts through the interior of this rectangle.
			if d.rectangleCutByPolygon(x1, y1, x2, y2) {
				continue
			}

			if area > best {
				best = area
			}
		}
	}

	return strconv.Itoa(best)
}

// ----------------------------------------------------------
// Polygon edges
// ----------------------------------------------------------

func (d *Day09) buildEdges() {
	n := len(d.reds)
	edges := make([]edge9, 0, n)
	for i := range n {
		a := d.reds[i]
		b := d.reds[(i+1)%n]
		e := edge9{x1: a.x, y1: a.y, x2: b.x, y2: b.y}
		if a.y == b.y {
			e.hor = true
			if e.x1 > e.x2 {
				e.x1, e.x2 = e.x2, e.x1
			}
		} else {
			e.hor = false
			if e.y1 > e.y2 {
				e.y1, e.y2 = e.y2, e.y1
			}
		}
		edges = append(edges, e)
	}
	d.edges = edges
}

// ----------------------------------------------------------
// Point in polygon (inside or on boundary)
// ----------------------------------------------------------

func (d *Day09) pointInsideOrOn(p pt9) bool {
	// First, check if point lies exactly on any edge (boundary is allowed).
	for _, e := range d.edges {
		if e.hor {
			if p.y == e.y1 && p.x >= e.x1 && p.x <= e.x2 {
				return true
			}
		} else {
			if p.x == e.x1 && p.y >= e.y1 && p.y <= e.y2 {
				return true
			}
		}
	}
	// Standard ray-casting (odd-even) for inside test.
	return pointInPolygonRayCast(p, d.reds)
}

func pointInPolygonRayCast(p pt9, poly []pt9) bool {
	inside := false
	n := len(poly)
	if n < 3 {
		return false
	}

	j := n - 1
	for i := 0; i < n; i++ {
		pi := poly[i]
		pj := poly[j]

		// Check if edge (pj -> pi) crosses a horizontal ray to the right of p.
		if (pi.y > p.y) != (pj.y > p.y) {
			// Compute intersection x coordinate of the edge with horizontal line y = p.y
			// x_intersect = xj + (y - yj)*(xi - xj)/(yi - yj)
			xIntersect := float64(pj.x) + float64(p.y-pj.y)*float64(pi.x-pj.x)/float64(pi.y-pj.y)
			if float64(p.x) < xIntersect {
				inside = !inside
			}
		}
		j = i
	}
	return inside
}

// ----------------------------------------------------------
// Check if polygon cuts through rectangle interior
// ----------------------------------------------------------
//
// Rectangle is [x1,x2] × [y1,y2], inclusive corner tiles.
// We treat the *interior* as (x1,x2) × (y1,y2) (open intervals).
// Any polygon segment strictly crossing this interior means the
// rectangle is not entirely inside the polygon.
//
// Horizontal segment at y0:
//
//	if y1 < y0 < y2 and (x-range ∩ (x1,x2)) non-empty => interior cut.
//
// Vertical segment at x0:
//
//	if x1 < x0 < x2 and (y-range ∩ (y1,y2)) non-empty => interior cut.
func (d *Day09) rectangleCutByPolygon(x1, y1, x2, y2 int) bool {
	if x1 == x2 || y1 == y2 {
		// Degenerate (line) rectangles don't really have interior.
		// For AoC purposes, they are fine if corners are inside.
		return false
	}

	for _, e := range d.edges {
		if e.hor {
			y0 := e.y1
			if y0 <= y1 || y0 >= y2 {
				continue
			}
			// e.x1..e.x2 vs (x1,x2)
			if maxInt(e.x1, x1) < minInt(e.x2, x2) {
				return true
			}
		} else {
			x0 := e.x1
			if x0 <= x1 || x0 >= x2 {
				continue
			}
			// e.y1..e.y2 vs (y1,y2)
			if maxInt(e.y1, y1) < minInt(e.y2, y2) {
				return true
			}
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
