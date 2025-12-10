package days

import (
	"fmt"
	"math"
	"strconv"
	"strings"
)

// MachineData holds the parsed information for a single machine
type MachineData struct {
	TargetLights  []int   // Part 1 target vector (0/1)
	TargetJoltage []int   // Part 2 target vector (integers)
	Buttons       [][]int // For each button: list of affected indices (defines matrix A)
}

type Day10 struct {
	Machines []MachineData
}

func init() {
	Register(10, func() Solution { return &Day10{} })
}

// ------------------------------------------------------------
// Parsing
// ------------------------------------------------------------

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

func (d *Day10) SetInput(lines []string) {
	d.Machines = d.Machines[:0]

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

		d.Machines = append(d.Machines, MachineData{
			TargetLights:  lights,
			TargetJoltage: joltage,
			Buttons:       buttons,
		})
	}
}

// ------------------------------------------------------------
// Part 1: GF(2) linear system, minimal Hamming weight
// ------------------------------------------------------------

func solveLights10(m MachineData) (int, error) {
	N := len(m.TargetLights)
	M := len(m.Buttons)
	if N == 0 || M == 0 {
		return 0, nil
	}

	// mat is N x (M+1) augmented matrix over GF(2)
	mat := make([][]int, N)
	for i := range N {
		row := make([]int, M+1)
		row[M] = m.TargetLights[i]
		mat[i] = row
	}

	// Fill A part
	for j, btn := range m.Buttons {
		for _, idx := range btn {
			if idx < N {
				mat[idx][j] = 1
			}
		}
	}

	pivotRow := 0
	pivotCols := make(map[int]int) // col -> row

	// Gaussian elimination over GF(2)
	for col := 0; col < M && pivotRow < N; col++ {
		sel := -1
		for r := pivotRow; r < N; r++ {
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

		for r := range N {
			if r != pivotRow && mat[r][col] == 1 {
				for k := col; k <= M; k++ {
					mat[r][k] ^= mat[pivotRow][k]
				}
			}
		}
		pivotRow++
	}

	// Identify free variables
	freeVars := make([]int, 0)
	for c := range M {
		if _, ok := pivotCols[c]; !ok {
			freeVars = append(freeVars, c)
		}
	}

	minPresses := M + 1
	count := 1 << len(freeVars)

	// Brute force over free variables; deduce pivot vars
	for mask := range count {
		x := make([]int, M)

		// Assign free variables
		for i, fIdx := range freeVars {
			if (mask>>i)&1 == 1 {
				x[fIdx] = 1
			}
		}

		// Compute pivot variables by back-substitution
		for c := M - 1; c >= 0; c-- {
			if r, isPivot := pivotCols[c]; isPivot {
				val := mat[r][M]
				for k := c + 1; k < M; k++ {
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

func solveJoltage10(m MachineData) (int, error) {
	N := len(m.TargetJoltage)
	M := len(m.Buttons)
	if N == 0 || M == 0 {
		return 0, nil
	}

	// Build augmented matrix mat (N x (M+1)) over R
	mat := make([][]float64, N)
	for i := range N {
		row := make([]float64, M+1)
		row[M] = float64(m.TargetJoltage[i])
		mat[i] = row
	}

	for j, btn := range m.Buttons {
		for _, idx := range btn {
			if idx < N {
				mat[idx][j] = 1.0
			}
		}
	}

	// RREF over R
	pivotRow := 0
	pivotCols := make(map[int]int)

	for col := 0; col < M && pivotRow < N; col++ {
		sel := -1
		for r := pivotRow; r < N; r++ {
			if math.Abs(mat[r][col]) > 1e-9 {
				sel = r
				break
			}
		}
		if sel == -1 {
			continue
		}

		mat[pivotRow], mat[sel] = mat[sel], mat[pivotRow]
		pivotCols[col] = pivotRow

		// Normalize pivot to 1
		div := mat[pivotRow][col]
		for k := col; k <= M; k++ {
			mat[pivotRow][k] /= div
		}

		// Eliminate column in all other rows
		for r := 0; r < N; r++ {
			if r == pivotRow {
				continue
			}
			f := mat[r][col]
			if math.Abs(f) < 1e-9 {
				continue
			}
			for k := col; k <= M; k++ {
				mat[r][k] -= f * mat[pivotRow][k]
			}
		}

		pivotRow++
	}

	// Check consistency
	for r := pivotRow; r < N; r++ {
		if math.Abs(mat[r][M]) > 1e-9 {
			return 0, fmt.Errorf("inconsistent real system")
		}
	}

	// Identify free variables
	freeVars := make([]int, 0)
	isPivotCol := make([]bool, M)
	for c, r := range pivotCols {
		_ = r
		isPivotCol[c] = true
	}
	for c := range M {
		if !isPivotCol[c] {
			freeVars = append(freeVars, c)
		}
	}

	// Very conservative bound: no button can be pressed more than max(target_i)+1
	maxTarget := 0.0
	for _, v := range m.TargetJoltage {
		if float64(v) > maxTarget {
			maxTarget = float64(v)
		}
	}
	searchBound := int(maxTarget) + 1

	minTotal := math.MaxInt64

	// DFS over free variables with simple branch-and-bound on total presses
	var dfs func(idx int, x []float64, currentSum int)

	dfs = func(idx int, x []float64, currentSum int) {
		if currentSum >= minTotal {
			return
		}
		if idx == len(freeVars) {
			// All free vars set -> compute pivot vars
			total := currentSum
			valid := true

			for c := 0; c < M; c++ {
				r, isPivot := pivotCols[c]
				if !isPivot {
					continue
				}

				val := mat[r][M] // RHS
				for k := c + 1; k < M; k++ {
					if math.Abs(mat[r][k]) > 1e-9 {
						val -= mat[r][k] * x[k]
					}
				}

				if val < -1e-5 {
					valid = false
					break
				}
				rounded := math.Round(val)
				if math.Abs(val-rounded) > 1e-5 {
					valid = false
					break
				}
				iVal := int(rounded)
				x[c] = float64(iVal)
				total += iVal
			}

			if valid && total < minTotal {
				minTotal = total
			}
			return
		}

		fc := freeVars[idx]
		for v := 0; v <= searchBound; v++ {
			x[fc] = float64(v)
			dfs(idx+1, x, currentSum+v)
			if currentSum+v >= minTotal {
				break
			}
		}
	}

	x := make([]float64, M)
	dfs(0, x, 0)

	if minTotal == math.MaxInt64 {
		return 0, fmt.Errorf("no integer solution")
	}
	return minTotal, nil
}

// ------------------------------------------------------------
// Day interface
// ------------------------------------------------------------

func (d *Day10) SolvePart1() string {
	total := 0
	for _, m := range d.Machines {
		if len(m.TargetLights) == 0 {
			continue
		}
		res, err := solveLights10(m)
		if err != nil {
			// With AoC input we expect solutions; panic if not.
			panic(fmt.Sprintf("Day10 Part1: %v", err))
		}
		total += res
	}
	return strconv.Itoa(total)
}

func (d *Day10) SolvePart2() string {
	total := 0
	for _, m := range d.Machines {
		if len(m.TargetJoltage) == 0 {
			continue
		}
		res, err := solveJoltage10(m)
		if err != nil {
			panic(fmt.Sprintf("Day10 Part2: %v", err))
		}
		total += res
	}
	return strconv.Itoa(total)
}
