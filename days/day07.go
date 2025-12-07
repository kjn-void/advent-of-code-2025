package days

import (
	"strconv"
	"strings"
)

type Day07 struct {
	grid       []string
	Rows, Cols int
	startCol   int // start column for S
}

func init() {
	Register(7, func() Solution { return &Day07{} })
}

func (d *Day07) SetInput(lines []string) {
	d.grid = d.grid[:0]

	// Keep layout exactly; AoC never gives malformed lines
	for _, line := range lines {
		d.grid = append(d.grid, line)
	}

	// Normalize width so all rows have same length
	maxC := len(d.grid[0])
	for i := range d.grid {
		if len(d.grid[i]) < maxC {
			d.grid[i] += strings.Repeat(" ", maxC-len(d.grid[i]))
		}
	}

	d.Rows = len(d.grid)
	d.Cols = maxC

	// Locate S on first row
	for c := 0; c < d.Cols; c++ {
		if d.grid[0][c] == 'S' {
			d.startCol = c
			return
		}
	}
}

// -----------------------------------------------------------
// Part 1 — Linear beam simulation counting split events
// -----------------------------------------------------------

func (d *Day07) SolvePart1() string {
	// Double buffer: two rows of bools we alternate between
	bufA := make([]bool, d.Cols)
	bufB := make([]bool, d.Cols)

	active := bufA
	next := bufB

	active[d.startCol] = true
	splitCount := 0

	for r := 1; r < d.Rows; r++ {
		row := d.grid[r]

		clear(next)

		for c, hasBeam := range active {
			if !hasBeam {
				continue
			}

			if row[c] == '^' {
				// Splitter: original beam ends, two new beams start
				splitCount++
				next[c-1] = true
				next[c+1] = true
			} else {
				// Otherwise beam continues straight
				next[c] = true
			}
		}

		// Swap buffers
		active, next = next, active
	}

	return strconv.Itoa(splitCount)
}

// -----------------------------------------------------------
// Part 2 — Count all timelines (many-worlds interpretation)
// -----------------------------------------------------------

func (d *Day07) SolvePart2() string {
	// Double buffer: two rows of int64 we alternate between
	bufA := make([]int64, d.Cols)
	bufB := make([]int64, d.Cols)

	timelines := bufA
	next := bufB

	timelines[d.startCol] = 1

	for r := 1; r < d.Rows; r++ {
		row := d.grid[r]

		clear(next)

		for c, count := range timelines {
			if count == 0 {
				continue
			}

			if row[c] == '^' {
				// Split to left and right
				next[c-1] += count
				next[c+1] += count
			} else {
				// Continue straight down
				next[c] += count
			}
		}

		// Swap buffers
		timelines, next = next, timelines
	}

	var totalCount int64
	for _, count := range timelines {
		totalCount += count
	}

	return strconv.FormatInt(totalCount, 10)
}
