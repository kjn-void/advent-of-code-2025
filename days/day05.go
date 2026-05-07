package days

import (
	"cmp"
	"slices"
	"strconv"
	"strings"
)

type day05 struct {
	freshRanges   []freshRange
	ingredientIDs []int64
}

type freshRange struct {
	start, end int64
}

func init() {
	Register(5, func() Solution { return &day05{} })
}

// SetInput parses fresh ingredient ranges and available ingredient IDs, then
// merges overlapping fresh ranges for efficient membership checks.
func (d *day05) SetInput(lines []string) {
	d.freshRanges = d.freshRanges[:0]
	d.ingredientIDs = d.ingredientIDs[:0]

	// Split into two blocks: ranges, blank line, then available IDs
	section := 0
	for _, line := range lines {
		s := strings.TrimSpace(line)
		if s == "" {
			section++
			continue
		}

		if section == 0 {
			// fresh ranges
			parts := strings.Split(s, "-")
			start, _ := strconv.ParseInt(parts[0], 10, 64)
			end, _ := strconv.ParseInt(parts[1], 10, 64)
			d.freshRanges = append(d.freshRanges, freshRange{start: start, end: end})
		} else {
			// available ingredient IDs (used only in part 1)
			id, _ := strconv.ParseInt(s, 10, 64)
			d.ingredientIDs = append(d.ingredientIDs, id)
		}
	}

	if len(d.freshRanges) == 0 {
		return
	}

	// Merge overlapping ranges for efficient lookup.
	slices.SortFunc(d.freshRanges, func(a, b freshRange) int {
		return cmp.Compare(a.start, b.start)
	})

	merged := d.freshRanges[:0]
	curStart, curEnd := d.freshRanges[0].start, d.freshRanges[0].end

	for i := 1; i < len(d.freshRanges); i++ {
		s, e := d.freshRanges[i].start, d.freshRanges[i].end
		if s <= curEnd { // overlapping
			if e > curEnd {
				curEnd = e
			}
		} else {
			merged = append(merged, freshRange{start: curStart, end: curEnd})
			curStart, curEnd = s, e
		}
	}
	merged = append(merged, freshRange{start: curStart, end: curEnd})
	d.freshRanges = merged
}

// SolvePart1 counts available ingredient IDs that fall inside any merged fresh
// range and returns that count.
func (d *day05) SolvePart1() string {
	count := 0
	for _, id := range d.ingredientIDs {
		if d.isFresh(id) {
			count++
		}
	}
	return strconv.Itoa(count)
}

// isFresh checks whether id is contained in the sorted, merged fresh ranges and
// returns true when the ingredient is fresh.
func (d *day05) isFresh(id int64) bool {
	// binary search in merged ranges
	lo, hi := 0, len(d.freshRanges)-1
	for lo <= hi {
		mid := (lo + hi) / 2
		r := d.freshRanges[mid]
		if id < r.start {
			hi = mid - 1
		} else if id > r.end {
			lo = mid + 1
		} else {
			return true
		}
	}
	return false
}

// SolvePart2 returns the total number of distinct ingredient IDs covered by the
// merged fresh ranges.
func (d *day05) SolvePart2() string {
	var total int64 = 0
	for _, r := range d.freshRanges {
		total += r.end - r.start + 1
	}

	return strconv.FormatInt(total, 10)
}
