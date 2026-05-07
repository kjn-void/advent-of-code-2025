package days

import (
	"strconv"
)

type day04 struct {
	grid []string
	rows int
	cols int
}

func init() {
	Register(4, func() Solution { return &day04{} })
}

// SetInput stores the paper-roll diagram and records its dimensions for the
// adjacency checks used by both parts.
func (d *day04) SetInput(lines []string) {
	d.grid = d.grid[:0]

	for _, line := range lines {
		d.grid = append(d.grid, line)
	}

	d.rows = len(d.grid)
	d.cols = len(d.grid[0])
}

// -----------------------------------------------------------------------------
// Common helpers
// -----------------------------------------------------------------------------

var day04Dirs = [8][2]int{
	{-1, -1}, {-1, 0}, {-1, 1},
	{0, -1}, {0, 1},
	{1, -1}, {1, 0}, {1, 1},
}

// makeBoolGrid converts the original diagram into a mutable occupancy grid and
// returns true for cells containing a paper roll.
func (d *day04) makeBoolGrid() [][]bool {
	out := make([][]bool, d.rows)
	for r := 0; r < d.rows; r++ {
		row := make([]bool, d.cols)
		src := d.grid[r]
		for c := 0; c < d.cols; c++ {
			row[c] = (src[c] == '@')
		}
		out[r] = row
	}
	return out
}

// computeDegrees counts occupied neighboring cells for each occupied roll in on
// and returns a grid of those adjacency counts.
func (d *day04) computeDegrees(on [][]bool) [][]int {
	deg := make([][]int, d.rows)
	for r := 0; r < d.rows; r++ {
		row := make([]int, d.cols)
		for c := 0; c < d.cols; c++ {
			if !on[r][c] {
				continue
			}
			cnt := 0
			for _, dxy := range day04Dirs {
				nr := r + dxy[0]
				nc := c + dxy[1]
				if nr >= 0 && nr < d.rows && nc >= 0 && nc < d.cols && on[nr][nc] {
					cnt++
				}
			}
			row[c] = cnt
		}
		deg[r] = row
	}
	return deg
}

// countAdjacentRolls counts the eight-neighbor paper rolls around grid cell
// (r,c) and returns that count.
func (d *day04) countAdjacentRolls(r, c int) int {
	count := 0
	for _, dxy := range day04Dirs {
		nr := r + dxy[0]
		nc := c + dxy[1]
		if nr >= 0 && nr < d.rows && nc >= 0 && nc < d.cols && d.grid[nr][nc] == '@' {
			count++
		}
	}
	return count
}

// -----------------------------------------------------------------------------
// Part 1
// -----------------------------------------------------------------------------

// SolvePart1 counts rolls that are immediately accessible to forklifts because
// they have fewer than four adjacent rolls.
func (d *day04) SolvePart1() string {
	if d.rows == 0 || d.cols == 0 {
		return "0"
	}

	total := 0

	for r := 0; r < d.rows; r++ {
		for c := 0; c < d.cols; c++ {
			if d.grid[r][c] != '@' {
				continue
			}
			if d.countAdjacentRolls(r, c) < 4 {
				total++
			}
		}
	}

	return strconv.Itoa(total)
}

// -----------------------------------------------------------------------------
// Part 2
// -----------------------------------------------------------------------------

// SolvePart2 repeatedly removes currently accessible rolls and returns the
// total number removed after accessibility cascades through the grid.
func (d *day04) SolvePart2() string {
	if d.rows == 0 || d.cols == 0 {
		return "0"
	}

	on := d.makeBoolGrid()
	deg := d.computeDegrees(on)

	type cell struct{ r, c int }
	queue := make([]cell, 0, d.rows*d.cols)

	for r := 0; r < d.rows; r++ {
		for c := 0; c < d.cols; c++ {
			if on[r][c] && deg[r][c] < 4 {
				queue = append(queue, cell{r, c})
			}
		}
	}

	removed := 0
	qp := 0

	for qp < len(queue) {
		cc := queue[qp]
		qp++
		r := cc.r
		c := cc.c

		if !on[r][c] {
			continue
		}

		on[r][c] = false
		removed++

		for _, dxy := range day04Dirs {
			nr := r + dxy[0]
			nc := c + dxy[1]
			if nr < 0 || nr >= d.rows || nc < 0 || nc >= d.cols {
				continue
			}
			if !on[nr][nc] {
				continue
			}

			deg[nr][nc]--
			if deg[nr][nc] == 3 {
				queue = append(queue, cell{nr, nc})
			}
		}
	}

	return strconv.Itoa(removed)
}
