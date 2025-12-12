package days

import (
	"strconv"
	"strings"
)

// --- Data types -------------------------------------------------------------

type point struct {
	x, y int
}

type variant struct {
	width, height int
	cells         []point // relative positions of '#' cells
}

type shape struct {
	area     int
	variants []variant
}

type region struct {
	width, height int
	counts        []int // number of presents of each shape index
}

type Day12 struct {
	shapes  []shape
	regions []region
}

func init() {
	Register(12, func() Solution { return &Day12{} })
}

// --- Parsing ---------------------------------------------------------------

func (d *Day12) SetInput(lines []string) {
	d.shapes = d.shapes[:0]
	d.regions = d.regions[:0]

	// Parse shapes first, then regions.
	i := 0
	for i < len(lines) && strings.TrimSpace(lines[i]) == "" {
		i++
	}

	// --- Parse shape blocks of form:
	// 0:
	// ###
	// ..#
	// ###
	for i < len(lines) {
		line := strings.TrimSpace(lines[i])
		if line == "" {
			i++
			continue
		}
		if isRegionLine(line) {
			// We've reached the region section.
			break
		}

		// Expect a header like "0:" or "5:"
		if !strings.HasSuffix(line, ":") {
			i++
			continue
		}
		i++

		var rows []string
		for i < len(lines) {
			s := strings.TrimRight(strings.TrimRight(lines[i], "\r"), "\n")
			if strings.TrimSpace(s) == "" {
				i++
				break
			}
			trimmed := strings.TrimSpace(s)

			// Stop if we hit the next shape header or a region line.
			if strings.HasSuffix(trimmed, ":") || isRegionLine(trimmed) {
				break
			}
			rows = append(rows, trimmed)
			i++
		}
		if len(rows) > 0 {
			d.shapes = append(d.shapes, buildShape(rows))
		}
	}

	// --- Parse regions: "WxH: c0 c1 c2 ..."
	for i < len(lines) {
		line := strings.TrimSpace(lines[i])
		i++
		if line == "" {
			continue
		}
		if !isRegionLine(line) {
			continue
		}
		parts := strings.SplitN(line, ":", 2)
		if len(parts) != 2 {
			continue
		}
		dimPart := strings.TrimSpace(parts[0])
		cntPart := strings.TrimSpace(parts[1])

		wh := strings.Split(dimPart, "x")
		if len(wh) != 2 {
			continue
		}
		w, err1 := strconv.Atoi(strings.TrimSpace(wh[0]))
		h, err2 := strconv.Atoi(strings.TrimSpace(wh[1]))
		if err1 != nil || err2 != nil {
			continue
		}

		countFields := strings.Fields(cntPart)
		if len(countFields) == 0 {
			continue
		}
		counts := make([]int, len(countFields))
		for idx, s := range countFields {
			val, err := strconv.Atoi(s)
			if err != nil {
				val = 0
			}
			counts[idx] = val
		}

		d.regions = append(d.regions, region{
			width:  w,
			height: h,
			counts: counts,
		})
	}
}

func isRegionLine(s string) bool {
	s = strings.TrimSpace(s)
	colonIdx := strings.IndexByte(s, ':')
	if colonIdx <= 0 {
		return false
	}
	head := s[:colonIdx]
	parts := strings.Split(head, "x")
	if len(parts) != 2 {
		return false
	}
	if _, err := strconv.Atoi(strings.TrimSpace(parts[0])); err != nil {
		return false
	}
	if _, err := strconv.Atoi(strings.TrimSpace(parts[1])); err != nil {
		return false
	}
	return true
}

// --- Shape construction (variants: rotations + flips) ----------------------

func buildShape(rows []string) shape {
	// Build initial grid as [][]bool
	h0 := len(rows)
	w0 := 0
	for _, row := range rows {
		if len(row) > w0 {
			w0 = len(row)
		}
	}
	grid := make([][]bool, h0)
	for y, row := range rows {
		grid[y] = make([]bool, w0)
		for x := 0; x < len(row); x++ {
			if row[x] == '#' {
				grid[y][x] = true
			}
		}
	}

	var variants []variant
	seen := make(map[string]bool)

	g := grid
	for r := range 4 {
		if r > 0 {
			g = rotateGrid(g)
		}
		for f := range 2 {
			var gf [][]bool
			if f == 0 {
				gf = g
			} else {
				gf = flipGridH(g)
			}
			v := gridToVariant(gf)
			if len(v.cells) == 0 {
				continue
			}
			key := variantKey(v)
			if !seen[key] {
				seen[key] = true
				variants = append(variants, v)
			}
		}
	}

	area := 0
	if len(variants) > 0 {
		area = len(variants[0].cells)
	}

	return shape{
		area:     area,
		variants: variants,
	}
}

func rotateGrid(grid [][]bool) [][]bool {
	height := len(grid)
	if height == 0 {
		return [][]bool{}
	}
	width := len(grid[0])
	res := make([][]bool, width)
	for y := range width {
		res[y] = make([]bool, height)
	}

	for y := range height {
		for x := range width {
			res[x][height-1-y] = grid[y][x]
		}
	}
	return res
}

func flipGridH(grid [][]bool) [][]bool {
	height := len(grid)
	if height == 0 {
		return [][]bool{}
	}
	width := len(grid[0])
	res := make([][]bool, height)
	for y := range height {
		res[y] = make([]bool, width)
		for x := range width {
			res[y][width-1-x] = grid[y][x]
		}
	}
	return res
}

func gridToVariant(grid [][]bool) variant {
	height := len(grid)
	if height == 0 {
		return variant{}
	}
	width := len(grid[0])

	minX, minY := width, height
	maxX, maxY := -1, -1

	for y := range height {
		for x := range width {
			if grid[y][x] {
				if x < minX {
					minX = x
				}
				if x > maxX {
					maxX = x
				}
				if y < minY {
					minY = y
				}
				if y > maxY {
					maxY = y
				}
			}
		}
	}
	if maxX < minX || maxY < minY {
		return variant{}
	}

	vw := maxX - minX + 1
	vh := maxY - minY + 1
	var cells []point

	for y := minY; y <= maxY; y++ {
		for x := minX; x <= maxX; x++ {
			if grid[y][x] {
				cells = append(cells, point{x: x - minX, y: y - minY})
			}
		}
	}

	return variant{
		width:  vw,
		height: vh,
		cells:  cells,
	}
}

func variantKey(v variant) string {
	var sb strings.Builder
	sb.WriteString(strconv.Itoa(v.width))
	sb.WriteByte('x')
	sb.WriteString(strconv.Itoa(v.height))
	sb.WriteByte(':')
	for _, c := range v.cells {
		sb.WriteString(strconv.Itoa(c.x))
		sb.WriteByte(',')
		sb.WriteString(strconv.Itoa(c.y))
		sb.WriteByte(';')
	}
	return sb.String()
}

// --- Solver core -----------------------------------------------------------

const smallBoardMaxArea12 = 15 * 15 // full tiling search only if w*h <= this

func (d *Day12) SolvePart1() string {
	valid := 0
	for _, region := range d.regions {
		if d.regionCanFit12(region) {
			valid++
		}
	}
	return strconv.Itoa(valid)
}

// For now, part two is unknown; implement a stub.
func (d *Day12) SolvePart2() string {
	return "0"
}

func (d *Day12) regionCanFit12(r region) bool {
	if len(d.shapes) == 0 {
		return false
	}

	// --- Area check (hard necessary condition) ---
	totalArea := 0
	for i, cnt := range r.counts {
		if i >= len(d.shapes) {
			break
		}
		if cnt < 0 {
			return false
		}
		totalArea += cnt * d.shapes[i].area
	}
	if totalArea > r.width*r.height {
		return false
	}

	// Small boards: do an actual tiling search (geometric fit).
	if r.width*r.height <= smallBoardMaxArea12 {
		return d.canTileRegionSmall(r)
	}

	// Large boards: assume area is sufficient (fast heuristic).
	// If needed, we can refine this later with more constraints.
	return true
}

// --- Exact tiling search for small regions ---------------------------------

func (d *Day12) canTileRegionSmall(r region) bool {
	w, h := r.width, r.height
	numShapes := len(d.shapes)
	if numShapes == 0 {
		return false
	}

	// Precompute all possible placements for each shape type on this board.
	allPlacements := make([][][]int, numShapes) // shapeIdx -> list of placements -> []boardIndex

	for si := range numShapes {
		shape := d.shapes[si]
		for _, v := range shape.variants {
			if v.width == 0 || v.height == 0 {
				continue
			}
			for by := 0; by <= h-v.height; by++ {
				for bx := 0; bx <= w-v.width; bx++ {
					var cells []int
					ok := true
					for _, c := range v.cells {
						x := bx + c.x
						y := by + c.y
						if x < 0 || x >= w || y < 0 || y >= h {
							ok = false
							break
						}
						cells = append(cells, y*w+x)
					}
					if ok && len(cells) > 0 {
						allPlacements[si] = append(allPlacements[si], cells)
					}
				}
			}
		}
	}

	board := make([]bool, w*h)
	counts := make([]int, len(r.counts))
	copy(counts, r.counts)

	return d.btTile(board, w, h, counts, allPlacements)
}

func (d *Day12) btTile(board []bool, w, h int, counts []int, placements [][][]int) bool {
	// Check if all counts are zero (all presents placed).
	done := true
	totalRemainingArea := 0
	for i, c := range counts {
		if i >= len(placements) {
			break
		}
		if c > 0 {
			done = false
			totalRemainingArea += c * d.shapes[i].area
		}
	}
	if done {
		return true
	}

	// Quick prune: not enough free cells to fit remaining area.
	freeCells := 0
	for _, occupied := range board {
		if !occupied {
			freeCells++
		}
	}
	if totalRemainingArea > freeCells {
		return false
	}

	// Choose the shape type that is currently most constrained:
	// the one with the smallest number of feasible placements.
	bestShape := -1
	bestCount := 1 << 30

	for si, c := range counts {
		if c <= 0 {
			continue
		}
		plList := placements[si]
		if len(plList) == 0 {
			// No way to place this shape at all.
			return false
		}
		feasible := 0
		for _, pl := range plList {
			ok := true
			for _, idx := range pl {
				if board[idx] {
					ok = false
					break
				}
			}
			if ok {
				feasible++
				if feasible >= bestCount {
					break
				}
			}
		}
		if feasible == 0 {
			return false
		}
		if feasible < bestCount {
			bestCount = feasible
			bestShape = si
		}
	}

	if bestShape == -1 {
		// Shouldn't happen if totalRemainingArea > 0, but guard anyway.
		return false
	}

	// Try placing one copy of bestShape in each feasible way.
	counts[bestShape]--
	for _, pl := range placements[bestShape] {
		ok := true
		for _, idx := range pl {
			if board[idx] {
				ok = false
				break
			}
		}
		if !ok {
			continue
		}

		// Place shape
		for _, idx := range pl {
			board[idx] = true
		}

		if d.btTile(board, w, h, counts, placements) {
			return true
		}

		// Undo
		for _, idx := range pl {
			board[idx] = false
		}
	}
	counts[bestShape]++
	return false
}
