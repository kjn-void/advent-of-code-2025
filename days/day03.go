package days

import (
	"strconv"
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
		digits := make([]int, len(line))
		for i, ch := range line {
			digits[i] = int(ch - '0') // 1–9
		}

		d.banks = append(d.banks, digits)
	}
}

// -------------------------
// Part 1
// -------------------------

func (d *Day03) SolvePart1() string {
	return d.maxJoltage(2)
}

// -------------------------
// Part 2
// -------------------------

func (d *Day03) SolvePart2() string {
	return d.maxJoltage(12)
}

func (d *Day03) maxJoltage(pick int) string {
	total := int64(0)

	for _, bank := range d.banks {
		n := len(bank)

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

		total += stackToNumber(stack)
	}

	return strconv.FormatInt(total, 10)
}

func stackToNumber(stack []int) int64 {
	var val int64 = 0
	for _, dgt := range stack {
		val = val*10 + int64(dgt)
	}
	return val
}
