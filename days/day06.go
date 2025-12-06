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
		d.grid = append(d.grid, line)
	}

	// Normalize row widths so all rows have identical length.
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
// Helpers
// -----------------------------------------------------------

type day06Block struct {
	start, end int
}

// findBlocks finds contiguous column ranges containing non-space characters,
// separated by fully blank columns.
func (d *Day06) findBlocks() []day06Block {
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

	blocks := make([]day06Block, 0, 64)
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
				blocks = append(blocks, day06Block{start, c - 1})
			}
		}
	}
	if inBlock {
		blocks = append(blocks, day06Block{start, d.C - 1})
	}

	return blocks
}

// getOperator reads the operator ('+' or '*') from the bottom row within a block.
func (d *Day06) getOperator(b day06Block) byte {
	opRow := d.grid[d.R-1][b.start : b.end+1]
	for i := range opRow {
		if opRow[i] == '+' || opRow[i] == '*' {
			return opRow[i]
		}
	}
	return '*' // fallback (AoC guarantees this won't happen)
}

// -----------------------------------------------------------
// Number extractors
// -----------------------------------------------------------

// Part 1: Each row forms a number (vertical block)
func (d *Day06) extractNumbersPart1(b day06Block) []int64 {
	nums := make([]int64, 0, d.R)

	for r := 0; r < d.R-1; r++ { // last row is operator
		s := strings.TrimSpace(d.grid[r][b.start : b.end+1])
		v, _ := strconv.ParseInt(s, 10, 64)
		nums = append(nums, v)
	}
	return nums
}

// Part 2: Each column forms a number (cephalopod rules)
func (d *Day06) extractNumbersPart2(b day06Block) []int64 {
	width := b.end - b.start + 1
	nums := make([]int64, 0, width)

	for c := b.start; c <= b.end; c++ {
		var sb strings.Builder
		sb.Grow(d.R)

		for r := 0; r < d.R-1; r++ { // last row is operator
			ch := d.grid[r][c]
			if ch != ' ' {
				sb.WriteByte(ch)
			}
		}

		v, _ := strconv.ParseInt(sb.String(), 10, 64)
		nums = append(nums, v)
	}
	return nums
}

// -----------------------------------------------------------
// Shared block evaluation
// -----------------------------------------------------------

func (d *Day06) evaluateBlocks(extractor func(day06Block) []int64) int64 {
	blocks := d.findBlocks()
	var total int64

	for _, b := range blocks {
		nums := extractor(b)
		op := d.getOperator(b)
		total += evalNumbers(nums, op)
	}

	return total
}

// evalNumbers computes either sum or product of nums depending on op ('+' or '*').
func evalNumbers(nums []int64, op byte) int64 {
	if op == '+' {
		var sum int64
		for _, n := range nums {
			sum += n
		}
		return sum
	}
	prod := int64(1)
	for _, n := range nums {
		prod *= n
	}
	return prod
}

// -----------------------------------------------------------
// Part 1
// -----------------------------------------------------------

func (d *Day06) SolvePart1() string {
	total := d.evaluateBlocks(d.extractNumbersPart1)
	return strconv.FormatInt(total, 10)
}

// -----------------------------------------------------------
// Part 2
// -----------------------------------------------------------

func (d *Day06) SolvePart2() string {
	total := d.evaluateBlocks(d.extractNumbersPart2)
	return strconv.FormatInt(total, 10)
}
