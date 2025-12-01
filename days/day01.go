package days

import (
	"strconv"
	"strings"
)

type Day01 struct {
	// Each move is a signed delta:
	//   Rn => +n (right / increasing)
	//   Ln => -n (left / decreasing)
	moves []int
}

func init() {
	Register(1, func() Solution { return &Day01{} })
}

func (d *Day01) SetInput(lines []string) {
	d.moves = d.moves[:0]

	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}

		dir := line[0]
		val, _ := strconv.Atoi(line[1:])

		if dir == 'L' {
			d.moves = append(d.moves, -val)
		} else { // 'R'
			d.moves = append(d.moves, val)
		}
	}
}

func mod100(n int) int {
	n %= 100
	if n < 0 {
		n += 100
	}
	return n
}

func (d *Day01) SolvePart1() string {
	pos := 50
	countZero := 0

	for _, delta := range d.moves {
		pos = mod100(pos + delta)
		if pos == 0 {
			countZero++
		}
	}

	return strconv.Itoa(countZero)
}

func (d *Day01) SolvePart2() string {
	pos := 50
	countZero := 0

	for _, delta := range d.moves {
		step := 1
		if delta < 0 {
			step = -1
		}

		for moved := 0; moved != delta; moved += step {
			pos = pos + step
			if pos < 0 {
				pos = pos + 100
			} else if pos >= 100 {
				pos = pos - 100
			}
			if pos == 0 {
				countZero++
			}
		}
	}

	return strconv.Itoa(countZero)
}
