package days

import (
	"sort"
	"strconv"
	"strings"
)

type Day05 struct {
	ranges [][2]int64
	ids    []int64
}

func init() {
	Register(5, func() Solution { return &Day05{} })
}

func (d *Day05) SetInput(lines []string) {
	d.ranges = d.ranges[:0]
	d.ids = d.ids[:0]

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
			d.ranges = append(d.ranges, [2]int64{start, end})
		} else {
			// available ingredient IDs (used only in part 1)
			id, _ := strconv.ParseInt(s, 10, 64)
			d.ids = append(d.ids, id)
		}
	}

	// Merge overlapping ranges for efficient lookup
	sort.Slice(d.ranges, func(i, j int) bool {
		return d.ranges[i][0] < d.ranges[j][0]
	})

	merged := d.ranges[:0]
	curStart, curEnd := d.ranges[0][0], d.ranges[0][1]

	for i := 1; i < len(d.ranges); i++ {
		s, e := d.ranges[i][0], d.ranges[i][1]
		if s <= curEnd { // overlapping
			if e > curEnd {
				curEnd = e
			}
		} else {
			merged = append(merged, [2]int64{curStart, curEnd})
			curStart, curEnd = s, e
		}
	}
	merged = append(merged, [2]int64{curStart, curEnd})
	d.ranges = merged
}

func (d *Day05) SolvePart1() string {
	count := 0
	for _, id := range d.ids {
		if d.isFresh(id) {
			count++
		}
	}
	return strconv.Itoa(count)
}

func (d *Day05) isFresh(id int64) bool {
	// binary search in merged ranges
	lo, hi := 0, len(d.ranges)-1
	for lo <= hi {
		mid := (lo + hi) / 2
		r := d.ranges[mid]
		if id < r[0] {
			hi = mid - 1
		} else if id > r[1] {
			lo = mid + 1
		} else {
			return true
		}
	}
	return false
}

func (d *Day05) SolvePart2() string {
	var total int64 = 0
	for _, r := range d.ranges {
		// r[0]..r[1] inclusive
		total += (r[1] - r[0] + 1)
	}

	return strconv.FormatInt(total, 10)
}
