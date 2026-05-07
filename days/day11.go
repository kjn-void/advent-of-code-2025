package days

import (
	"strconv"
	"strings"
)

type day11 struct {
	outputs map[string][]string
}

func init() {
	Register(11, func() Solution { return &day11{} })
}

// SetInput parses device output lines into a directed graph from each device to
// the devices receiving its outputs.
func (d *day11) SetInput(lines []string) {
	d.outputs = make(map[string][]string)

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
		d.outputs[from] = outs
	}
}

// -----------------------------------------------------------
// Part 1 — count all paths from "you" to "out"
// -----------------------------------------------------------

// countPathsFrom counts all directed paths from node to "out" using memoization;
// visiting guards against accidental cycles in malformed input.
func (d *day11) countPathsFrom(node string, memo map[string]int64, visiting map[string]bool) int64 {
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
	for _, next := range d.outputs[node] {
		total += d.countPathsFrom(next, memo, visiting)
	}

	memo[node] = total
	return total
}

// SolvePart1 returns the number of directed paths from "you" to "out".
func (d *day11) SolvePart1() string {
	if d.outputs == nil {
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

// countPathsWithRequired counts directed paths from start to end that visit both
// required nodes, carrying a bitmask of visited requirements through the DFS.
func (d *day11) countPathsWithRequired(start, end string, need1, need2 string) int64 {
	if d.outputs == nil {
		return 0
	}

	memo := make(map[day11State]int64)

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

		var total int64

		for _, nxt := range d.outputs[node] {
			nextMask := mask
			if nxt == need1 {
				nextMask |= 1
			}
			if nxt == need2 {
				nextMask |= 2
			}
			total += dfs(nxt, nextMask)
		}

		memo[st] = total
		return total
	}

	return dfs(start, mask)
}

// SolvePart2 returns the number of paths from "svr" to "out" that visit both
// required diagnostic devices "dac" and "fft".
func (d *day11) SolvePart2() string {
	// Count paths from "svr" to "out" that visit both "dac" and "fft"
	total := d.countPathsWithRequired("svr", "out", "dac", "fft")
	return strconv.FormatInt(total, 10)
}
