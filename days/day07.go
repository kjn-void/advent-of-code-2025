package days

import (
	"strconv"
	"strings"
)

type Day07 struct {
	grid   []string
	R, C   int
	sr, sc int // start row/column for S
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

	// Locate S exactly once
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
// Part 1 — Linear beam simulation counting split events
// -----------------------------------------------------------

func (d *Day07) SolvePart1() string {
	active := make([]bool, d.C)
	active[d.sc] = true
	splitCount := 0

	for r := d.sr + 1; r < d.R; r++ {
		row := d.grid[r]
		next := make([]bool, d.C)

		for c, hasBeam := range active {
			if !hasBeam {
				continue
			}

			if row[c] == '^' {
				// Splitter: original beam ends, two new beams start
				splitCount++
				if c > 0 {
					next[c-1] = true
				}
				if c+1 < d.C {
					next[c+1] = true
				}
			} else {
				// Otherwise beam continues straight
				next[c] = true
			}
		}

		active = next
	}

	return strconv.Itoa(splitCount)
}

// -----------------------------------------------------------
// Part 2 — Count all timelines (many-worlds interpretation)
// -----------------------------------------------------------

func (d *Day07) SolvePart2() string {
	// timelines[c] = how many timelines have reached column c
	timelines := make([]int64, d.C)
	timelines[d.sc] = 1

	for r := d.sr + 1; r < d.R; r++ {
		row := d.grid[r]
		next := make([]int64, d.C)

		for c, count := range timelines {
			if count == 0 {
				continue
			}

			if row[c] == '^' {
				// Split to left and right
				if c > 0 {
					next[c-1] += count
				}
				if c+1 < d.C {
					next[c+1] += count
				}
			} else {
				// Continue straight down
				next[c] += count
			}
		}

		timelines = next
	}

	var total int64
	for _, t := range timelines {
		total += t
	}

	return strconv.FormatInt(total, 10)
}
