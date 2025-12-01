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
		if delta > 0 {
			// Right: pos -> pos + delta
			for i := 1; i <= delta; i++ {
				if mod100(pos+i) == 0 {
					countZero++
				}
			}
		} else if delta < 0 {
			// Left: pos -> pos + delta (delta is negative)
			step := -delta
			for i := 1; i <= step; i++ {
				if mod100(pos-i) == 0 {
					countZero++
				}
			}
		}
		pos = mod100(pos + delta)
	}

	return strconv.Itoa(countZero)
}
