package days

import (
	"fmt"
	"math"
	"strconv"
	"strings"
)

type machine struct {
	targetLights  []int
	targetJoltage []int
	buttons       [][]int
}

type day10 struct {
	machines []machine
}

func init() {
	Register(10, func() Solution { return &day10{} })
}

// ------------------------------------------------------------
// Parsing
// ------------------------------------------------------------

// parseList parses a bracketed, parenthesized, or braced comma-separated list of
// integers and returns the values inside it.
func parseList(s string) []int {
	s = strings.TrimSpace(s)
	if len(s) < 2 {
		return nil
	}
	// Remove outer brackets/parens/braces
	s = s[1 : len(s)-1]
	if s == "" {
		return nil
	}
	parts := strings.Split(s, ",")
	result := make([]int, 0, len(parts))
	for _, p := range parts {
		val, err := strconv.Atoi(strings.TrimSpace(p))
		if err == nil {
			result = append(result, val)
		}
	}
	return result
}

// SetInput parses each machine manual line into target indicator lights,
// button wiring, and joltage requirements for the two solvers.
func (d *day10) SetInput(lines []string) {
	d.machines = d.machines[:0]

	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}

		// Format: [.##.] (3) (1,3) ... {3,5,4,7}

		// 1. Extract lights [ ... ]
		startBracket := strings.Index(line, "[")
		endBracket := strings.Index(line, "]")
		if startBracket == -1 || endBracket == -1 {
			continue
		}
		lightStr := line[startBracket+1 : endBracket]
		lights := make([]int, len(lightStr))
		for i, char := range lightStr {
			if char == '#' {
				lights[i] = 1
			} else {
				lights[i] = 0
			}
		}

		// 2. Extract joltage { ... }
		startBrace := strings.Index(line, "{")
		endBrace := strings.Index(line, "}")
		var joltage []int
		if startBrace != -1 && endBrace != -1 {
			joltage = parseList(line[startBrace : endBrace+1])
		}

		// 3. Extract buttons (...) between ']' and '{' (if present)
		midSection := line[endBracket+1:]
		if startBrace != -1 {
			midSection = line[endBracket+1 : startBrace]
		}

		buttons := make([][]int, 0)
		for {
			pStart := strings.Index(midSection, "(")
			if pStart == -1 {
				break
			}
			pEnd := strings.Index(midSection, ")")
			if pEnd == -1 {
				break
			}
			buttons = append(buttons, parseList(midSection[pStart:pEnd+1]))
			midSection = midSection[pEnd+1:]
		}

		d.machines = append(d.machines, machine{
			targetLights:  lights,
			targetJoltage: joltage,
			buttons:       buttons,
		})
	}
}

// ------------------------------------------------------------
// Part 1: GF(2) linear system, minimal Hamming weight
// ------------------------------------------------------------

// solveIndicatorLights solves the button-toggle system over GF(2) for one
// machine and returns the minimum number of button presses needed for lights.
func solveIndicatorLights(m machine) (int, error) {
	nLights := len(m.targetLights)
	nButtons := len(m.buttons)
	if nLights == 0 || nButtons == 0 {
		return 0, nil
	}

	// mat is nLights x (nButtons+1) augmented matrix over GF(2)
	mat := make([][]int, nLights)
	for i := range nLights {
		row := make([]int, nButtons+1)
		row[nButtons] = m.targetLights[i]
		mat[i] = row
	}

	// Fill A part
	for j, btn := range m.buttons {
		for _, idx := range btn {
			if idx < nLights {
				mat[idx][j] = 1
			}
		}
	}

	pivotRow := 0
	pivotCols := make(map[int]int) // col -> row

	// Gaussian elimination over GF(2)
	for col := 0; col < nButtons && pivotRow < nLights; col++ {
		sel := -1
		for r := pivotRow; r < nLights; r++ {
			if mat[r][col] == 1 {
				sel = r
				break
			}
		}
		if sel == -1 {
			continue
		}

		mat[pivotRow], mat[sel] = mat[sel], mat[pivotRow]
		pivotCols[col] = pivotRow

		for r := range nLights {
			if r != pivotRow && mat[r][col] == 1 {
				for k := col; k <= nButtons; k++ {
					mat[r][k] ^= mat[pivotRow][k]
				}
			}
		}
		pivotRow++
	}

	// Identify free variables
	freeVars := make([]int, 0)
	for c := range nButtons {
		if _, ok := pivotCols[c]; !ok {
			freeVars = append(freeVars, c)
		}
	}

	minPresses := nButtons + 1
	count := 1 << len(freeVars)

	// Brute force over free variables; deduce pivot vars
	for mask := range count {
		x := make([]int, nButtons)

		// Assign free variables
		for i, fIdx := range freeVars {
			if (mask>>i)&1 == 1 {
				x[fIdx] = 1
			}
		}

		// Compute pivot variables by back-substitution
		for c := nButtons - 1; c >= 0; c-- {
			if r, isPivot := pivotCols[c]; isPivot {
				val := mat[r][nButtons]
				for k := c + 1; k < nButtons; k++ {
					if mat[r][k] == 1 {
						val ^= x[k]
					}
				}
				x[c] = val
			}
		}

		presses := 0
		for _, v := range x {
			presses += v
		}
		if presses < minPresses {
			minPresses = presses
		}
	}

	return minPresses, nil
}

// ------------------------------------------------------------
// Part 2: Real RREF + integer search over free vars
// ------------------------------------------------------------

// solveJoltageRequirements solves the integer button-press system for one
// machine's joltage requirements and returns the minimum total press count.
func solveJoltageRequirements(m machine) (int, error) {
	nLights := len(m.targetJoltage)
	nButtons := len(m.buttons)
	if nLights == 0 || nButtons == 0 {
		return 0, nil
	}

	cols := nButtons + 1
	mat := make([]float64, nLights*cols)
	for i, target := range m.targetJoltage {
		mat[i*cols+nButtons] = float64(target)
	}

	for j, btn := range m.buttons {
		for _, idx := range btn {
			if idx < nLights {
				mat[idx*cols+j] = 1.0
			}
		}
	}

	// RREF over R
	pivotRow := 0
	pivotRowForCol := make([]int, nButtons)
	for i := range pivotRowForCol {
		pivotRowForCol[i] = -1
	}
	pivotCap := nLights
	if nButtons < pivotCap {
		pivotCap = nButtons
	}
	pivotCols := make([]int, 0, pivotCap)

	for col := 0; col < nButtons && pivotRow < nLights; col++ {
		sel := -1
		for r := pivotRow; r < nLights; r++ {
			if math.Abs(mat[r*cols+col]) > 1e-9 {
				sel = r
				break
			}
		}
		if sel == -1 {
			continue
		}

		if sel != pivotRow {
			pivotBase := pivotRow * cols
			selBase := sel * cols
			for k := 0; k <= nButtons; k++ {
				mat[pivotBase+k], mat[selBase+k] = mat[selBase+k], mat[pivotBase+k]
			}
		}

		pivotBase := pivotRow * cols
		pivotRowForCol[col] = pivotRow
		pivotCols = append(pivotCols, col)

		// Normalize pivot to 1
		div := mat[pivotBase+col]
		for k := col; k <= nButtons; k++ {
			mat[pivotBase+k] /= div
		}

		// Eliminate column in all other rows
		for r := 0; r < nLights; r++ {
			if r == pivotRow {
				continue
			}
			rowBase := r * cols
			f := mat[rowBase+col]
			if math.Abs(f) < 1e-9 {
				continue
			}
			for k := col; k <= nButtons; k++ {
				mat[rowBase+k] -= f * mat[pivotBase+k]
			}
		}

		pivotRow++
	}

	// Check consistency
	for r := pivotRow; r < nLights; r++ {
		if math.Abs(mat[r*cols+nButtons]) > 1e-9 {
			return 0, fmt.Errorf("inconsistent real system")
		}
	}

	// Identify free variables
	freeVars := make([]int, 0, nButtons-len(pivotCols))
	for c := range nButtons {
		if pivotRowForCol[c] == -1 {
			freeVars = append(freeVars, c)
		}
	}

	freeBounds := make([]int, len(freeVars))
	for i, col := range freeVars {
		bound := math.MaxInt
		for _, idx := range m.buttons[col] {
			if idx < nLights && m.targetJoltage[idx] < bound {
				bound = m.targetJoltage[idx]
			}
		}
		if bound == math.MaxInt {
			bound = 0
		}
		freeBounds[i] = bound
	}

	// Try smaller domains first so branch-and-bound finds a good total early.
	for i := 1; i < len(freeVars); i++ {
		col := freeVars[i]
		bound := freeBounds[i]
		j := i - 1
		for j >= 0 && freeBounds[j] > bound {
			freeVars[j+1] = freeVars[j]
			freeBounds[j+1] = freeBounds[j]
			j--
		}
		freeVars[j+1] = col
		freeBounds[j+1] = bound
	}

	pivotCount := len(pivotCols)
	pivotRHS := make([]float64, pivotCount)
	pivotFreeCoeff := make([]float64, pivotCount*len(freeVars))
	for i, col := range pivotCols {
		rowBase := pivotRowForCol[col] * cols
		pivotRHS[i] = mat[rowBase+nButtons]
		coeffBase := i * len(freeVars)
		for j, freeCol := range freeVars {
			pivotFreeCoeff[coeffBase+j] = mat[rowBase+freeCol]
		}
	}

	freeValues := make([]int, len(freeVars))
	minTotal := math.MaxInt

	// DFS over free variables with simple branch-and-bound on total presses
	var dfs func(idx int, currentSum int)

	dfs = func(idx int, currentSum int) {
		if currentSum >= minTotal {
			return
		}
		if idx == len(freeVars) {
			// All free vars set -> compute pivot vars
			total := currentSum

			for p := 0; p < pivotCount; p++ {
				val := pivotRHS[p]
				coeffBase := p * len(freeVars)
				for f, x := range freeValues {
					coeff := pivotFreeCoeff[coeffBase+f]
					if math.Abs(coeff) > 1e-9 {
						val -= coeff * float64(x)
					}
				}

				if val < -1e-5 {
					return
				}
				rounded := math.Round(val)
				if math.Abs(val-rounded) > 1e-5 {
					return
				}
				iVal := int(rounded)
				total += iVal
				if total >= minTotal {
					return
				}
			}

			if total < minTotal {
				minTotal = total
			}
			return
		}

		for v := 0; v <= freeBounds[idx]; v++ {
			freeValues[idx] = v
			dfs(idx+1, currentSum+v)
			if currentSum+v >= minTotal {
				break
			}
		}
	}

	dfs(0, 0)

	if minTotal == math.MaxInt {
		return 0, fmt.Errorf("no integer solution")
	}
	return minTotal, nil
}

// ------------------------------------------------------------
// Day interface
// ------------------------------------------------------------

// SolvePart1 sums the minimum indicator-light button presses across all parsed
// machines and returns the total.
func (d *day10) SolvePart1() string {
	total := 0
	for _, m := range d.machines {
		if len(m.targetLights) == 0 {
			continue
		}
		res, err := solveIndicatorLights(m)
		if err != nil {
			// With AoC input we expect solutions; panic if not.
			panic(fmt.Sprintf("Day10 Part1: %v", err))
		}
		total += res
	}
	return strconv.Itoa(total)
}

// SolvePart2 solves each machine's joltage system concurrently, sums the minimum
// press counts, and returns the total.
func (d *day10) SolvePart2() string {
	total := 0
	resultCh := make(chan int, len(d.machines))
	for _, m := range d.machines {
		go func(m machine) {
			res, err := solveJoltageRequirements(m)
			if err != nil {
				panic(fmt.Sprintf("Day10 Part2: %v", err))
			}
			resultCh <- res
		}(m)
	}
	for range d.machines {
		total += <-resultCh
	}
	return strconv.Itoa(total)
}
