package days

import (
	"strconv"
	"strings"
)

type Day02 struct {
	ranges [][2]int64 // inclusive ranges
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

		// AoC input is always well-formed
		start, _ := strconv.ParseInt(b[0], 10, 64)
		end, _ := strconv.ParseInt(b[1], 10, 64)

		d.ranges = append(d.ranges, [2]int64{start, end})
	}
}

func pow10(n int) int64 {
	v := int64(1)
	for range n {
		v *= 10
	}
	return v
}

// Part 1: sequence repeated exactly twice
func (d *Day02) SolvePart1() string {
	var sum int64 = 0

	// For each range, sequential approach
	for _, r := range d.ranges {
		L := r[0]
		R := r[1]

		// Max digits of any possible number in this range
		maxDigits := len(strconv.FormatInt(R, 10))

		// For each possible half-length k: total length = 2k (must be even)
		for k := 1; 2*k <= maxDigits; k++ {
			base := pow10(k)
			repeatFactor := base + 1 // N = d*(10^k+1)

			// Valid d range: [10^(k-1), 10^k - 1]
			dLo := pow10(k - 1)
			dHi := base - 1

			// L <= d*repeatFactor <= R
			candMin := (L + repeatFactor - 1) / repeatFactor // ceil(L / repeatFactor)
			candMax := R / repeatFactor                      // floor(R / repeatFactor)

			// Intersect with valid d bounds
			if candMin < dLo {
				candMin = dLo
			}
			if candMax > dHi {
				candMax = dHi
			}
			if candMin > candMax {
				continue
			}

			// All d in [candMin..candMax] produce invalid IDs in this range
			for dd := candMin; dd <= candMax; dd++ {
				N := dd * repeatFactor
				sum += N
			}
		}
	}

	return strconv.FormatInt(sum, 10)
}

// Part 2: sequence repeated at least twice (2, 3, 4, ... times)
func (d *Day02) SolvePart2() string {
	var sum int64 = 0
	seen := make(map[int64]struct{})

	for _, r := range d.ranges {
		L := r[0]
		R := r[1]

		maxDigits := len(strconv.FormatInt(R, 10))

		// Total digit length of N
		for lenDigits := 2; lenDigits <= maxDigits; lenDigits++ {
			tenLen := pow10(lenDigits) // 10^lenDigits

			// m = number of repetitions, must divide lenDigits, and m >= 2
			for m := 2; m <= lenDigits; m++ {
				if lenDigits%m != 0 {
					continue
				}

				k := lenDigits / m // digit length of the repeating block

				baseK := pow10(k)
				// geometric series: 1 + 10^k + ... + 10^{(m-1)k}
				// = (10^{lenDigits} - 1) / (10^k - 1)
				repFactor := (tenLen - 1) / (baseK - 1)

				// Valid d range: exactly k-digit numbers
				dLo := pow10(k - 1)
				dHi := baseK - 1

				// From range: L <= d*repFactor <= R
				candMin := (L + repFactor - 1) / repFactor // ceil(L / repFactor)
				candMax := R / repFactor                   // floor(R / repFactor)

				// Intersect with valid d bounds
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
					N := dd * repFactor

					// Avoid double-counting numbers that can be represented
					// with multiple (k, m) combinations.
					if _, ok := seen[N]; ok {
						continue
					}
					seen[N] = struct{}{}
					sum += N
				}
			}
		}
	}

	return strconv.FormatInt(sum, 10)
}
