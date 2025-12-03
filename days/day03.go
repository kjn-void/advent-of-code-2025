package days

import (
	"strconv"
	"strings"
)

type Day03 struct {
	banks [][]int // each bank is a slice of digits
}

func init() {
	Register(3, func() Solution { return &Day03{} })
}

func (d *Day03) SetInput(lines []string) {
	d.banks = d.banks[:0]

	for _, line := range lines {
		s := strings.TrimSpace(line)
		if s == "" {
			continue
		}

		digits := make([]int, len(s))
		for i, ch := range s {
			digits[i] = int(ch - '0') // 1–9
		}

		d.banks = append(d.banks, digits)
	}
}

// -------------------------
// Part 1
// -------------------------

func (d *Day03) SolvePart1() string {
	total := 0

	for _, bank := range d.banks {
		if len(bank) < 2 {
			continue
		}

		maxPrev := -1
		best := 0

		for _, digit := range bank {
			if maxPrev != -1 {
				val := maxPrev*10 + digit
				if val > best {
					best = val
				}
			}
			if digit > maxPrev {
				maxPrev = digit
			}
		}

		total += best
	}

	return strconv.Itoa(total)
}

// -------------------------
// Part 2
// -------------------------

func (d *Day03) SolvePart2() string {
	total := int64(0)
	const pick = 12

	for _, bank := range d.banks {
		n := len(bank)
		if n < pick {
			continue
		}

		// Monotonic stack for lexicographically max subsequence of length 12
		need := pick
		stack := make([]int, 0, pick)

		for i := range n {
			dig := bank[i]

			remaining := n - i
			canPop := len(stack) > 0 && remaining > need

			for canPop && stack[len(stack)-1] < dig {
				stack = stack[:len(stack)-1]
				need++
				canPop = len(stack) > 0 && remaining > need
			}

			if need > 0 {
				stack = append(stack, dig)
				need--
			}
		}

		// Convert to number
		var val int64 = 0
		for _, dgt := range stack {
			val = val*10 + int64(dgt)
		}
		total += val
	}

	return strconv.FormatInt(total, 10)
}
