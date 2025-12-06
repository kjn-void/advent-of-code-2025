package days

import (
	"strconv"
	"strings"
)

type Day06 struct {
	grid []string
	R, C int
}

func init() {
	Register(6, func() Solution { return &Day06{} })
}

func (d *Day06) SetInput(lines []string) {
	d.grid = d.grid[:0]

	for _, line := range lines {
		s := strings.TrimRight(line, "\n")
		if len(strings.TrimSpace(s)) == 0 {
			continue
		}
		d.grid = append(d.grid, s)
	}

	// Normalize row widths so all rows have identical length
	maxC := 0
	for _, row := range d.grid {
		if len(row) > maxC {
			maxC = len(row)
		}
	}
	for i := range d.grid {
		if len(d.grid[i]) < maxC {
			missing := maxC - len(d.grid[i])
			d.grid[i] += strings.Repeat(" ", missing)
		}
	}

	d.R = len(d.grid)
	d.C = maxC
}

// -----------------------------------------------------------
// Part 1 (unchanged)
// -----------------------------------------------------------

func (d *Day06) SolvePart1() string {
	// Identify blank columns
	isBlank := make([]bool, d.C)
	for c := 0; c < d.C; c++ {
		allSpace := true
		for r := 0; r < d.R; r++ {
			if d.grid[r][c] != ' ' {
				allSpace = false
				break
			}
		}
		isBlank[c] = allSpace
	}

	// Identify problem-blocks
	type block struct{ start, end int }
	blocks := make([]block, 0, 64)
	inBlock := false
	start := 0

	for c := 0; c < d.C; c++ {
		if !isBlank[c] {
			if !inBlock {
				inBlock = true
				start = c
			}
		} else {
			if inBlock {
				inBlock = false
				blocks = append(blocks, block{start, c - 1})
			}
		}
	}
	if inBlock {
		blocks = append(blocks, block{start, d.C - 1})
	}

	// Evaluate blocks (Part 1 logic)
	var total int64 = 0

	for _, b := range blocks {
		nums := make([]int64, 0, d.R)

		// Collect rows 0..R-2 (numbers)
		for r := 0; r < d.R-1; r++ {
			sub := strings.TrimSpace(d.grid[r][b.start : b.end+1])
			v, _ := strconv.ParseInt(sub, 10, 64)
			nums = append(nums, v)
		}

		// Operator is in last row
		opRow := d.grid[d.R-1][b.start : b.end+1]
		op := byte('*')
		for i := range opRow {
			if opRow[i] == '+' || opRow[i] == '*' {
				op = opRow[i]
				break
			}
		}

		if op == '+' {
			s := int64(0)
			for _, v := range nums {
				s += v
			}
			total += s
		} else {
			p := int64(1)
			for _, v := range nums {
				p *= v
			}
			total += p
		}
	}

	return strconv.FormatInt(total, 10)
}

// -----------------------------------------------------------
// Part 2 — CEPhalopod right-to-left column-wise numbers
// -----------------------------------------------------------

func (d *Day06) SolvePart2() string {
	// Identify blank columns
	isBlank := make([]bool, d.C)
	for c := 0; c < d.C; c++ {
		allSpace := true
		for r := 0; r < d.R; r++ {
			if d.grid[r][c] != ' ' {
				allSpace = false
				break
			}
		}
		isBlank[c] = allSpace
	}

	// Identify problem blocks again
	type block struct{ start, end int }
	blocks := make([]block, 0, 64)
	inBlock := false
	start := 0

	for c := 0; c < d.C; c++ {
		if !isBlank[c] {
			if !inBlock {
				inBlock = true
				start = c
			}
		} else {
			if inBlock {
				inBlock = false
				blocks = append(blocks, block{start, c - 1})
			}
		}
	}
	if inBlock {
		blocks = append(blocks, block{start, d.C - 1})
	}

	var total int64 = 0

	for _, b := range blocks {

		// Operator is still in bottom row
		opRow := d.grid[d.R-1][b.start : b.end+1]
		op := byte('*')
		for i := range opRow {
			if opRow[i] == '+' || opRow[i] == '*' {
				op = opRow[i]
				break
			}
		}

		// Extract numbers: each column is one number
		// Read top-to-bottom rows 0..R-2
		numbers := make([]int64, 0, b.end-b.start+1)

		for c := b.start; c <= b.end; c++ {
			builder := strings.Builder{}
			builder.Grow(d.R)

			for r := 0; r < d.R-1; r++ {
				ch := d.grid[r][c]
				if ch != ' ' {
					builder.WriteByte(ch)
				}
			}

			if builder.Len() == 0 {
				continue
			}

			v, _ := strconv.ParseInt(builder.String(), 10, 64)
			numbers = append(numbers, v)
		}

		// Evaluate right-to-left
		if op == '+' {
			s := int64(0)
			for _, number := range numbers {
				s += number
			}
			total += s
		} else {
			p := int64(1)
			for _, number := range numbers {
				p *= number
			}
			total += p
		}
	}

	return strconv.FormatInt(total, 10)
}
