package days

import (
	"fmt"
	"math"
	"strconv"
	"strings"
)

// MachineData holds the parsed information for a single machine
type MachineData struct {
	TargetLights  []int   // Part 1 target vector b (0/1)
	TargetJoltage []int   // Part 2 target vector b (integers)
	Buttons       [][]int // List of affected indices for each button (defines Matrix A)
}

// Day10 implements the Solution interface.
type Day10 struct {
	Machines []MachineData
}

func init() {
	Register(10, func() Solution { return &Day10{} })
}

// parseList parses strings like "(0,2,3,4)" or "{3,5,4,7}" into an []int.
func parseList(s string) []int {
	s = strings.TrimSpace(s)
	if len(s) < 2 {
		return nil
	}
	// Remove outer brackets/braces
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

		// Input format: [.##.] (3) (1,3) ... {3,5,4,7}
		// Split by spaces to get parts
		// We need to be careful not to split inside the brackets if spaces existed there,
		// but the problem description implies standard spacing.

		// 1. Extract Lights [ ... ]
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

		// 2. Extract Joltage { ... }
		startBrace := strings.Index(line, "{")
		endBrace := strings.Index(line, "}")
		var joltage []int
		if startBrace != -1 && endBrace != -1 {
			joltage = parseList(line[startBrace : endBrace+1])
		}

		// 3. Extract Buttons (...)
		// We scan between the light diagram and the joltage braces
		midSection := line[endBracket+1:]
		if startBrace != -1 {
			midSection = line[endBracket+1 : startBrace]
		}

		buttons := make([][]int, 0)
		// Naive split might fail if nested commas, but (1,2,3) is simple.
		// Let's iterate finding (...) pairs.
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

// --- Solver Logic ---

// solveSystem solves the linear system for the given machine.
// isPart2: false for GF(2) (lights), true for Real Integers (joltage).
func solveSystem(data MachineData, isPart2 bool) (int, error) {
	// 1. Construct Matrix A (N x M) and Vector b (N)
	var N, M int
	M = len(data.Buttons)

	var b []float64
	if isPart2 {
		N = len(data.TargetJoltage)
		b = make([]float64, N)
		for i, v := range data.TargetJoltage {
			b[i] = float64(v)
		}
	} else {
		N = len(data.TargetLights)
		b = make([]float64, N)
		for i, v := range data.TargetLights {
			b[i] = float64(v)
		}
	}

	if N == 0 || M == 0 {
		return 0, nil
	}

	// Matrix A: rows=counters, cols=buttons
	A := make([][]float64, N)
	for i := 0; i < N; i++ {
		A[i] = make([]float64, M)
	}
	for j, btn := range data.Buttons {
		for _, affected := range btn {
			if affected < N {
				A[affected][j] = 1.0
			}
		}
	}

	// 2. Gaussian Elimination to RREF
	// We use float64 for both parts to share logic, manually handling mod 2 for Part 1 later if needed,
	// or we can implement distinct logic. Distinct is safer for Part 1 GF(2).

	if !isPart2 {
		return solveGF2(N, M, A, b)
	}
	return solveReal(N, M, A, b)
}

// solveGF2 solves Ax = b over GF(2) minimizing Hamming weight.
func solveGF2(N, M int, A [][]float64, b []float64) (int, error) {
	// Convert to int for easier bitwise ops
	mat := make([][]int, N)
	for i := 0; i < N; i++ {
		row := make([]int, M+1) // Augmented
		for j := 0; j < M; j++ {
			row[j] = int(A[i][j])
		}
		row[M] = int(b[i])
		mat[i] = row
	}

	pivotRow := 0
	pivotCols := make(map[int]int) // col -> row

	for col := 0; col < M && pivotRow < N; col++ {
		// Find pivot
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

		// Swap
		mat[pivotRow], mat[sel] = mat[sel], mat[pivotRow]
		pivotCols[col] = pivotRow

		// Eliminate
		for r := 0; r < N; r++ {
			if r != pivotRow && mat[r][col] == 1 {
				for k := col; k <= M; k++ {
					mat[r][k] ^= mat[pivotRow][k]
				}
			}
		}
		pivotRow++
	}

	// Check consistency
	for r := pivotRow; r < N; r++ {
		if mat[r][M] == 1 {
			return 0, fmt.Errorf("inconsistent system")
		}
	}

	// Identify free variables
	var freeVars []int
	for c := 0; c < M; c++ {
		if _, ok := pivotCols[c]; !ok {
			freeVars = append(freeVars, c)
		}
	}

	minPresses := M + 1

	// Brute force free variables (2^F)
	count := 1 << len(freeVars)
	for mask := 0; mask < count; mask++ {
		x := make([]int, M)
		// Set free vars
		for i, fIdx := range freeVars {
			if (mask>>i)&1 == 1 {
				x[fIdx] = 1
			}
		}
		// Solve dependent vars
		// In RREF, pivot row `r` for col `c` looks like: 1*x_c + sum(other*x_other) = b_r
		// x_c = b_r - sum(...)  (subtraction is XOR in GF2)
		// Since we did full Jordan elimination (above and below), the row only contains the pivot and free vars to the right.

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

		// Count presses
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

// solveReal solves Ax = b over Reals, restricting to non-negative integers.
func solveReal(N, M int, A [][]float64, b []float64) (int, error) {
	// Create Augmented Matrix
	mat := make([][]float64, N)
	for i := 0; i < N; i++ {
		row := make([]float64, M+1)
		copy(row, A[i])
		row[M] = b[i]
		mat[i] = row
	}

	pivotRow := 0
	pivotCols := make(map[int]int)

	for col := 0; col < M && pivotRow < N; col++ {
		// Find pivot
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

		// Swap
		mat[pivotRow], mat[sel] = mat[sel], mat[pivotRow]
		pivotCols[col] = pivotRow

		// Normalize pivot to 1
		div := mat[pivotRow][col]
		for k := col; k <= M; k++ {
			mat[pivotRow][k] /= div
		}

		// Eliminate
		for r := 0; r < N; r++ {
			if r != pivotRow {
				factor := mat[r][col]
				if math.Abs(factor) > 1e-9 {
					for k := col; k <= M; k++ {
						mat[r][k] -= factor * mat[pivotRow][k]
					}
				}
			}
		}
		pivotRow++
	}

	// Check consistency
	for r := pivotRow; r < N; r++ {
		if math.Abs(mat[r][M]) > 1e-9 {
			return 0, fmt.Errorf("inconsistent system")
		}
	}

	// Identify Free Vars
	var freeVars []int
	for c := 0; c < M; c++ {
		if _, ok := pivotCols[c]; !ok {
			freeVars = append(freeVars, c)
		}
	}

	// Determine Max Bound for Search
	// Since A has only 0/1 and b >= 0, no x_i can exceed Max(b)
	// (unless a column is all zeros, in which case x_i=0 is optimal).
	maxTarget := 0.0
	for _, val := range b {
		if val > maxTarget {
			maxTarget = val
		}
	}
	searchBound := int(maxTarget) + 1

	minTotal := math.MaxInt64

	// Recursive Search over Free Variables
	// Using Branch and Bound
	var search func(idx int, currentX []float64, currentSum int)
	search = func(idx int, currentX []float64, currentSum int) {
		// Optimization: Pruning
		if currentSum >= minTotal {
			return
		}

		if idx == len(freeVars) {
			// All free vars set, calculate dependent vars
			total := currentSum
			valid := true

			// We need to determine dependent vars.
			// Thanks to RREF, for pivot col `c` at row `r`:
			// 1*x_c + sum(mat[r][k]*x_k for k > c) = mat[r][M]
			// x_c = mat[r][M] - sum(...)

			// We must calculate in reverse order of columns to handle dependencies correctly?
			// Actually, in RREF, the only non-zero entries in row `r` other than the pivot `c`
			// are free variable columns to the right.

			// Let's iterate all pivot columns.
			for c := 0; c < M; c++ {
				if r, isPivot := pivotCols[c]; isPivot {
					val := mat[r][M]
					for k := c + 1; k < M; k++ {
						if math.Abs(mat[r][k]) > 1e-9 {
							val -= mat[r][k] * currentX[k]
						}
					}

					// Check Non-Negative Integer
					if val < -1e-5 {
						valid = false
						break
					}
					// Round to nearest int
					rounded := math.Round(val)
					if math.Abs(val-rounded) > 1e-5 {
						valid = false // Not an integer
						break
					}

					iVal := int(rounded)
					currentX[c] = float64(iVal)
					total += iVal
				}
			}

			if valid {
				if total < minTotal {
					minTotal = total
				}
			}
			return
		}

		// Try values for free var freeVars[idx]
		// Optimization: Check if this free var column is all zeros in the original matrix A.
		// If so, it contributes nothing to targets, only cost. Best value is 0.
		// (Gaussian elimination handles this: column of zeros => free var).
		// We can check the coefficient sum for this column in the *original* matrix or just search.
		// Searching 0 is first, so it handles it naturally.

		fIdx := freeVars[idx]
		for v := 0; v <= searchBound; v++ {
			currentX[fIdx] = float64(v)
			search(idx+1, currentX, currentSum+v)

			// Pruning inside loop: if simple sum of free vars already exceeds best, stop
			if currentSum+v >= minTotal {
				break
			}
		}
	}

	initialX := make([]float64, M)
	search(0, initialX, 0)

	if minTotal == math.MaxInt64 {
		return 0, fmt.Errorf("no solution found")
	}
	return minTotal, nil
}

func (d *Day10) SolvePart1() string {
	total := 0
	for _, m := range d.Machines {
		res, err := solveSystem(m, false)
		if err != nil {
			// fmt.Println("Part 1 Error:", err)
			// Return error string or ignore?
			// For AOC, usually we assume valid inputs.
			continue
		}
		total += res
	}
	return strconv.Itoa(total)
}

func (d *Day10) SolvePart2() string {
	total := 0
	// Sort free variables or optimizing the matrix could speed this up,
	// but given N,M <= 13, this is sufficient.
	for _, m := range d.Machines {
		res, err := solveSystem(m, true)
		if err != nil {
			// fmt.Println("Part 2 Error:", err)
			return "Error: " + err.Error()
		}
		total += res
	}
	return strconv.Itoa(total)
}
