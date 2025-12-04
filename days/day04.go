package days

import (
	"strconv"
	"strings"
)

type Day04 struct {
	grid []string
	rows int
	cols int
}

func init() {
	Register(4, func() Solution { return &Day04{} })
}

func (d *Day04) SetInput(lines []string) {
	d.grid = d.grid[:0]

	for _, line := range lines {
		s := strings.TrimSpace(line)
		if s == "" {
			continue
		}
		d.grid = append(d.grid, s)
	}

	d.rows = len(d.grid)
	if d.rows > 0 {
		d.cols = len(d.grid[0])
	} else {
		d.cols = 0
	}
}

func (d *Day04) SolvePart1() string {
	if d.rows == 0 || d.cols == 0 {
		return "0"
	}

	accessible := 0

	dirs := [8][2]int{
		{-1, -1}, {-1, 0}, {-1, 1},
		{0, -1}, {0, 1},
		{1, -1}, {1, 0}, {1, 1},
	}

	for r := 0; r < d.rows; r++ {
		row := d.grid[r]
		for c := 0; c < d.cols; c++ {
			if row[c] != '@' {
				continue
			}

			count := 0
			for _, dxy := range dirs {
				nr := r + dxy[0]
				nc := c + dxy[1]
				if nr < 0 || nr >= d.rows || nc < 0 || nc >= d.cols {
					continue
				}
				if d.grid[nr][nc] == '@' {
					count++
				}
			}

			if count < 4 {
				accessible++
			}
		}
	}

	return strconv.Itoa(accessible)
}

// ---------------------------
// Part 2
// ---------------------------

func (d *Day04) SolvePart2() string {
	if d.rows == 0 || d.cols == 0 {
		return "0"
	}

	// Convert grid to mutable form
	on := make([][]bool, d.rows)
	deg := make([][]int, d.rows)

	for r := 0; r < d.rows; r++ {
		on[r] = make([]bool, d.cols)
		deg[r] = make([]int, d.cols)
		row := d.grid[r]
		for c := 0; c < d.cols; c++ {
			if row[c] == '@' {
				on[r][c] = true
			}
		}
	}

	dirs := [8][2]int{
		{-1, -1}, {-1, 0}, {-1, 1},
		{0, -1}, {0, 1},
		{1, -1}, {1, 0}, {1, 1},
	}

	// Compute initial neighbor counts
	for r := 0; r < d.rows; r++ {
		for c := 0; c < d.cols; c++ {
			if !on[r][c] {
				continue
			}
			cnt := 0
			for _, dxy := range dirs {
				nr := r + dxy[0]
				nc := c + dxy[1]
				if nr < 0 || nr >= d.rows || nc < 0 || nc >= d.cols {
					continue
				}
				if on[nr][nc] {
					cnt++
				}
			}
			deg[r][c] = cnt
		}
	}

	// Queue for BFS-like peeling
	type cell struct{ r, c int }
	queue := make([]cell, 0, d.rows*d.cols)

	// Enqueue initial removable rolls
	for r := 0; r < d.rows; r++ {
		for c := 0; c < d.cols; c++ {
			if on[r][c] && deg[r][c] < 4 {
				queue = append(queue, cell{r, c})
			}
		}
	}

	removed := 0

	// Pointer for consuming queue without popping front
	qp := 0

	for qp < len(queue) {
		cc := queue[qp]
		qp++

		r := cc.r
		c := cc.c
		if !on[r][c] {
			continue
		}

		// Remove this roll
		on[r][c] = false
		removed++

		// Reduce neighbor degrees; new candidates may appear
		for _, dxy := range dirs {
			nr := r + dxy[0]
			nc := c + dxy[1]
			if nr < 0 || nr >= d.rows || nc < 0 || nc >= d.cols {
				continue
			}
			if !on[nr][nc] {
				continue
			}

			deg[nr][nc]--
			if deg[nr][nc] == 3 { // became accessible now
				queue = append(queue, cell{nr, nc})
			}
		}
	}

	return strconv.Itoa(removed)
}
