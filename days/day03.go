package days

import (
	"strconv"
)

type day03 struct {
	batteryBanks [][]int
}

func init() {
	Register(3, func() Solution { return &day03{} })
}

// SetInput converts each battery-bank line into digits while preserving order,
// which matters because selected batteries cannot be rearranged.
func (d *day03) SetInput(lines []string) {
	d.batteryBanks = d.batteryBanks[:0]

	for _, line := range lines {
		digits := make([]int, len(line))
		for i, ch := range line {
			digits[i] = int(ch - '0') // 1–9
		}

		d.batteryBanks = append(d.batteryBanks, digits)
	}
}

// -------------------------
// Part 1
// -------------------------

// SolvePart1 finds the best two-battery joltage for each bank and returns their
// total as a decimal string.
func (d *day03) SolvePart1() string {
	return d.maxJoltage(2)
}

// -------------------------
// Part 2
// -------------------------

// SolvePart2 applies the same ordered digit-selection algorithm using twelve
// batteries per bank and returns the total joltage.
func (d *day03) SolvePart2() string {
	return d.maxJoltage(12)
}

// maxJoltage selects pick digits from each bank to form the largest possible
// ordered number, sums those numbers, and returns the total as a string.
func (d *day03) maxJoltage(pick int) string {
	total := int64(0)

	for _, bank := range d.batteryBanks {
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

// stackToNumber converts an ordered digit stack into the integer represented by
// those digits.
func stackToNumber(stack []int) int64 {
	var val int64 = 0
	for _, dgt := range stack {
		val = val*10 + int64(dgt)
	}
	return val
}
