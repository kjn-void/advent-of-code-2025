package days

import (
	"strconv"
	"strings"
)

type day02 struct {
	productIDRanges []idRange
}

type idRange struct {
	first, last int64
}

func init() {
	Register(2, func() Solution { return &day02{} })
}

// SetInput parses comma-separated inclusive product ID ranges into named range
// structs used by both invalid-ID summations.
func (d *day02) SetInput(lines []string) {
	d.productIDRanges = d.productIDRanges[:0]

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

		d.productIDRanges = append(d.productIDRanges, idRange{first: start, last: end})
	}
}

// pow10Table builds powers of ten used to construct repeated numeric patterns
// without repeatedly calling slower math/string helpers.
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

// smallestBlock returns the length of the smallest digit block that can be
// repeated to form s; it returns len(s) when s has no internal repetition.
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

// SolvePart1 sums all product IDs in the configured ranges whose decimal form
// is exactly two copies of the same digit sequence.
func (d *day02) SolvePart1() string {
	sum := int64(0)

	for _, productIDs := range d.productIDRanges {
		maxDigits := len(strconv.FormatInt(productIDs.last, 10))

		for k := 1; 2*k <= maxDigits; k++ {
			base := p10[k]
			repFactor := base + 1

			dLo := p10[k-1]
			dHi := base - 1

			candMin := (productIDs.first + repFactor - 1) / repFactor
			candMax := productIDs.last / repFactor

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

// SolvePart2 sums all product IDs in the configured ranges whose decimal form
// is two or more copies of a primitive digit sequence.
func (d *day02) SolvePart2() string {
	total := int64(0)

	for _, productIDs := range d.productIDRanges {
		maxDigits := len(strconv.FormatInt(productIDs.last, 10))

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

				candMin := (productIDs.first + repFactor - 1) / repFactor
				candMax := productIDs.last / repFactor

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
