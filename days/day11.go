package days

import (
	"strconv"
	"strings"
)

type Day11 struct {
	adj map[string][]string
}

func init() {
	Register(11, func() Solution { return &Day11{} })
}

func (d *Day11) SetInput(lines []string) {
	d.adj = make(map[string][]string)

	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}

		// Format: "aaa: you hhh"
		parts := strings.SplitN(line, ":", 2)
		if len(parts) == 0 {
			continue
		}
		from := strings.TrimSpace(parts[0])

		var outs []string
		if len(parts) == 2 {
			right := strings.TrimSpace(parts[1])
			if right != "" {
				for _, tok := range strings.Fields(right) {
					outs = append(outs, tok)
				}
			}
		}
		d.adj[from] = outs
	}
}

// -----------------------------------------------------------
// Part 1 — count all paths from "you" to "out"
// -----------------------------------------------------------

func (d *Day11) countPathsFrom(node string, memo map[string]int64, visiting map[string]bool) int64 {
	// Base case: reaching "out" is one valid path.
	if node == "out" {
		return 1
	}

	if v, ok := memo[node]; ok {
		return v
	}

	// Simple cycle guard (shouldn't happen in valid input).
	if visiting[node] {
		return 0
	}
	visiting[node] = true
	defer func() {
		delete(visiting, node)
	}()

	var total int64
	for _, next := range d.adj[node] {
		total += d.countPathsFrom(next, memo, visiting)
	}

	memo[node] = total
	return total
}

func (d *Day11) SolvePart1() string {
	if d.adj == nil {
		return "0"
	}

	memo := make(map[string]int64)
	visiting := make(map[string]bool)

	total := d.countPathsFrom("you", memo, visiting)
	return strconv.FormatInt(total, 10)
}

// -----------------------------------------------------------
// Part 2 — paths from "svr" to "out" that visit both "dac" and "fft"
// -----------------------------------------------------------

type day11State struct {
	node string
	mask int // bit0: visited dac, bit1: visited fft
}

func (d *Day11) countPathsWithRequired(start, end string, need1, need2 string) int64 {
	if d.adj == nil {
		return 0
	}

	memo := make(map[day11State]int64)
	visiting := make(map[day11State]bool)

	// initial mask (in case start is one of the required nodes)
	mask := 0
	if start == need1 {
		mask |= 1
	}
	if start == need2 {
		mask |= 2
	}

	var dfs func(node string, mask int) int64
	dfs = func(node string, mask int) int64 {
		st := day11State{node: node, mask: mask}

		// Cycle guard
		if visiting[st] {
			return 0
		}
		if v, ok := memo[st]; ok {
			return v
		}

		// Base case: at end, count path only if both required nodes seen.
		if node == end {
			if mask == 3 {
				return 1
			}
			return 0
		}

		visiting[st] = true
		var total int64

		for _, nxt := range d.adj[node] {
			nextMask := mask
			if nxt == need1 {
				nextMask |= 1
			}
			if nxt == need2 {
				nextMask |= 2
			}
			total += dfs(nxt, nextMask)
		}

		visiting[st] = false
		memo[st] = total
		return total
	}

	return dfs(start, mask)
}

func (d *Day11) SolvePart2() string {
	// Count paths from "svr" to "out" that visit both "dac" and "fft"
	total := d.countPathsWithRequired("svr", "out", "dac", "fft")
	return strconv.FormatInt(total, 10)
}
