package days

import (
	"strconv"
	"strings"
)

type Day07 struct {
	grid   []string
	R, C   int
	sr, sc int // start row / column
}

func init() {
	Register(7, func() Solution { return &Day07{} })
}

func (d *Day07) SetInput(lines []string) {
	d.grid = d.grid[:0]

	for _, line := range lines {
		// Keep layout exactly as in the input
		d.grid = append(d.grid, line)
	}

	// Normalize row widths to avoid out-of-bounds on ragged lines
	maxC := 0
	for _, row := range d.grid {
		if len(row) > maxC {
			maxC = len(row)
		}
	}
	for i := range d.grid {
		if len(d.grid[i]) < maxC {
			d.grid[i] += strings.Repeat(" ", maxC-len(d.grid[i]))
		}
	}

	d.R = len(d.grid)
	d.C = maxC

	// Locate the start 'S'
	for r := 0; r < d.R; r++ {
		for c := 0; c < d.C; c++ {
			if d.grid[r][c] == 'S' {
				d.sr, d.sc = r, c
				return
			}
		}
	}
}

// -----------------------------------------------------------
// Part 1 — count total number of splits
// -----------------------------------------------------------

func (d *Day07) SolvePart1() string {
	active := make([]bool, d.C)
	active[d.sc] = true

	splitCount := 0

	for r := d.sr + 1; r < d.R; r++ {
		nextActive := make([]bool, d.C)

		for c := 0; c < d.C; c++ {
			if !active[c] {
				continue
			}

			switch d.grid[r][c] {
			case '.':
				// Beam continues straight down
				nextActive[c] = true

			case '^':
				// Splitter: beam splits; original path stops here
				splitCount++

				// Left beam
				if c > 0 {
					nextActive[c-1] = true
				}
				// Right beam
				if c+1 < d.C {
					nextActive[c+1] = true
				}

			default:
				// Treat any unexpected character as empty space
				nextActive[c] = true
			}
		}

		active = nextActive
	}

	return strconv.Itoa(splitCount)
}

// -----------------------------------------------------------
// Part 2 — count number of distinct timelines
// -----------------------------------------------------------

func (d *Day07) SolvePart2() string {
	// timelines[c] = number of timelines with the particle at column c
	// just above row r (we will process row by row).
	timelines := make([]int64, d.C)
	timelines[d.sc] = 1

	for r := d.sr + 1; r < d.R; r++ {
		next := make([]int64, d.C)

		for c := 0; c < d.C; c++ {
			count := timelines[c]
			if count == 0 {
				continue
			}

			switch d.grid[r][c] {
			case '.':
				// All those timelines just move straight down.
				next[c] += count

			case '^':
				// Splitter: each timeline forks into (at most) 2.
				if c > 0 {
					next[c-1] += count
				}
				if c+1 < d.C {
					next[c+1] += count
				}

			default:
				// Any other cell acts like empty space.
				next[c] += count
			}
		}

		timelines = next
	}

	var total int64
	for _, cnt := range timelines {
		total += cnt
	}

	return strconv.FormatInt(total, 10)
}
