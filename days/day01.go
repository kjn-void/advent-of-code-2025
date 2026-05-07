package days

import (
	"strconv"
	"strings"
)

type day01 struct {
	// Each rotation is a signed dial delta:
	//   Rn => +n (right / increasing)
	//   Ln => -n (left / decreasing)
	rotations []int
}

func init() {
	Register(1, func() Solution { return &day01{} })
}

// SetInput parses rotation instructions like "L68" and "R8" into signed dial
// movements stored on the solver for both parts.
func (d *day01) SetInput(lines []string) {
	d.rotations = d.rotations[:0]

	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}

		dir := line[0]
		val, _ := strconv.Atoi(line[1:])

		if dir == 'L' {
			d.rotations = append(d.rotations, -val)
		} else { // 'R'
			d.rotations = append(d.rotations, val)
		}
	}
}

// dialPosition wraps n onto the safe dial's 0..99 range and returns the
// normalized position.
func dialPosition(n int) int {
	n %= 100
	if n < 0 {
		n += 100
	}
	return n
}

// SolvePart1 follows each full rotation from the starting dial position and
// returns how many rotations leave the dial pointing at zero.
func (d *day01) SolvePart1() string {
	dial := 50
	zeroStops := 0

	for _, rotation := range d.rotations {
		dial = dialPosition(dial + rotation)
		if dial == 0 {
			zeroStops++
		}
	}

	return strconv.Itoa(zeroStops)
}

// SolvePart2 walks every individual click in each rotation and returns how many
// times the dial crosses or lands on zero during the full instruction list.
func (d *day01) SolvePart2() string {
	dial := 50
	zeroClicks := 0

	for _, rotation := range d.rotations {
		step := 1
		if rotation < 0 {
			step = -1
		}

		for moved := 0; moved != rotation; moved += step {
			dial += step
			if dial < 0 {
				dial += 100
			} else if dial >= 100 {
				dial -= 100
			}
			if dial == 0 {
				zeroClicks++
			}
		}
	}

	return strconv.Itoa(zeroClicks)
}
