package days

import (
	"strconv"
	"strings"
)

type day06 struct {
	grid       []string
	rows, cols int
}

func init() {
	Register(6, func() Solution { return &day06{} })
}

// SetInput stores the worksheet rows and pads them to equal width so column
// scans can safely index every row.
func (d *day06) SetInput(lines []string) {
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

	d.rows = len(d.grid)
	d.cols = maxC
}

// -----------------------------------------------------------
// Helpers
// -----------------------------------------------------------

type worksheetProblem struct {
	start, end int
}

// findProblems finds contiguous column ranges containing non-space characters,
// separated by fully blank columns.
// It returns one worksheetProblem per horizontally arranged math problem.
func (d *day06) findProblems() []worksheetProblem {
	isBlank := make([]bool, d.cols)

	for c := 0; c < d.cols; c++ {
		allSpace := true
		for r := 0; r < d.rows; r++ {
			if d.grid[r][c] != ' ' {
				allSpace = false
				break
			}
		}
		isBlank[c] = allSpace
	}

	problems := make([]worksheetProblem, 0, 64)
	inBlock := false
	start := 0

	for c := 0; c < d.cols; c++ {
		if !isBlank[c] {
			if !inBlock {
				inBlock = true
				start = c
			}
		} else {
			if inBlock {
				inBlock = false
				problems = append(problems, worksheetProblem{start, c - 1})
			}
		}
	}
	if inBlock {
		problems = append(problems, worksheetProblem{start, d.cols - 1})
	}

	return problems
}

// getOperator reads the operator ('+' or '*') from the bottom row within a
// worksheet problem and returns '*' only as a defensive fallback for bad input.
func (d *day06) getOperator(problem worksheetProblem) byte {
	opRow := d.grid[d.rows-1][problem.start : problem.end+1]
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

// extractNumbersPart1 reads a problem using the part-one layout, where each row
// segment forms one operand, and returns those operands.
func (d *day06) extractNumbersPart1(problem worksheetProblem) []int64 {
	nums := make([]int64, 0, d.rows)

	for r := 0; r < d.rows-1; r++ { // last row is operator
		s := strings.TrimSpace(d.grid[r][problem.start : problem.end+1])
		v, _ := strconv.ParseInt(s, 10, 64)
		nums = append(nums, v)
	}
	return nums
}

// extractNumbersPart2 reads a problem using the part-two layout, where each
// column forms one operand after spaces are removed, and returns those operands.
func (d *day06) extractNumbersPart2(problem worksheetProblem) []int64 {
	width := problem.end - problem.start + 1
	nums := make([]int64, 0, width)

	for c := problem.start; c <= problem.end; c++ {
		var sb strings.Builder
		sb.Grow(d.rows)

		for r := 0; r < d.rows-1; r++ { // last row is operator
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

// evaluateProblems applies extractor to each parsed worksheet problem, evaluates
// the operands with that problem's operator, and returns the grand total.
func (d *day06) evaluateProblems(extractor func(worksheetProblem) []int64) int64 {
	problems := d.findProblems()
	var total int64

	for _, problem := range problems {
		nums := extractor(problem)
		op := d.getOperator(problem)
		total += evalNumbers(nums, op)
	}

	return total
}

// evalNumbers computes either the sum or product of nums depending on op ('+' or
// '*') and returns the problem result.
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

// SolvePart1 evaluates the worksheet with row-oriented operands and returns the
// grand total.
func (d *day06) SolvePart1() string {
	total := d.evaluateProblems(d.extractNumbersPart1)
	return strconv.FormatInt(total, 10)
}

// -----------------------------------------------------------
// Part 2
// -----------------------------------------------------------

// SolvePart2 evaluates the worksheet with column-oriented operands and returns
// the grand total.
func (d *day06) SolvePart2() string {
	total := d.evaluateProblems(d.extractNumbersPart2)
	return strconv.FormatInt(total, 10)
}
