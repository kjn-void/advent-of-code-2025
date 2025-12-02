package days

import (
	"strconv"
	"strings"
)

type Day02 struct {
	ranges [][2]int64 // list of [Low, High] ranges, inclusive
}

func init() {
	Register(2, func() Solution { return &Day02{} })
}

func (d *Day02) SetInput(lines []string) {
	d.ranges = d.ranges[:0]

	if len(lines) == 0 {
		return
	}

	parts := strings.SplitSeq(strings.TrimSpace(lines[0]), ",")
	for part := range parts {
		if part == "" {
			continue
		}
		b := strings.Split(part, "-")

		start, _ := strconv.ParseInt(b[0], 10, 64)
		end, _ := strconv.ParseInt(b[1], 10, 64)

		d.ranges = append(d.ranges, [2]int64{start, end})
	}
}

func pow10Table() []int64 {
	t := make([]int64, 20)
	x := int64(1)
	for i := 0; i < 20; i++ {
		t[i] = x
		x *= 10
	}
	return t
}

var p10 = pow10Table()

// smallest repeating block size of numeric string s
func smallestBlock(s string) int {
	n := len(s)
	for k := 1; k <= n/2; k++ {
		if n%k != 0 {
			continue
		}
		block := s[:k]
		ok := true
		for i := k; i < n; i += k {
			if s[i:i+k] != block {
				ok = false
				break
			}
		}
		if ok {
			return k
		}
	}
	return n
}

// ----- Part 1 -----

func (d *Day02) SolvePart1() string {
	sum := int64(0)

	for _, r := range d.ranges {
		L, R := r[0], r[1]
		maxDigits := len(strconv.FormatInt(R, 10))

		for k := 1; 2*k <= maxDigits; k++ {
			base := p10[k]
			repFactor := base + 1

			dLo := p10[k-1]
			dHi := base - 1

			candMin := (L + repFactor - 1) / repFactor
			candMax := R / repFactor

			if candMin < dLo {
				candMin = dLo
			}
			if candMax > dHi {
				candMax = dHi
			}
			if candMin > candMax {
				continue
			}

			for dd := candMin; dd <= candMax; dd++ {
				sum += dd * repFactor
			}
		}
	}

	return strconv.FormatInt(sum, 10)
}

// ----- Part 2 -----

func (d *Day02) SolvePart2() string {
	total := int64(0)

	for _, r := range d.ranges {
		L, R := r[0], r[1]
		maxDigits := len(strconv.FormatInt(R, 10))

		for totalDigits := 2; totalDigits <= maxDigits; totalDigits++ {
			tenLen := p10[totalDigits]

			for m := 2; m <= totalDigits; m++ {
				if totalDigits%m != 0 {
					continue
				}

				k := totalDigits / m

				baseK := p10[k]
				repFactor := (tenLen - 1) / (baseK - 1)

				dLo := p10[k-1]
				dHi := baseK - 1

				candMin := (L + repFactor - 1) / repFactor
				candMax := R / repFactor

				if candMin < dLo {
					candMin = dLo
				}
				if candMax > dHi {
					candMax = dHi
				}
				if candMin > candMax {
					continue
				}

				for dd := candMin; dd <= candMax; dd++ {
					ds := strconv.FormatInt(dd, 10)

					// uniqueness: dd must not have internal repetition
					if smallestBlock(ds) != len(ds) {
						continue
					}

					total += dd * repFactor
				}
			}
		}
	}

	return strconv.FormatInt(total, 10)
}
