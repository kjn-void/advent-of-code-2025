package days

import (
	"container/heap"
	"strconv"
	"strings"
)

type Day10 struct {
	machines []machine10
}

func init() {
	Register(10, func() Solution { return &Day10{} })
}

type machine10 struct {
	targetLights []int   // Part 1: target bitmask bits
	buttons1     []int   // Part 1: toggles as bitmask
	targetJolt   []int   // Part 2: jolt vector
	buttons2     [][]int // Part 2: list of counters incremented per button
}

// ------------------------------------------------------------
// Parse input (shared for Part 1 + Part 2)
// ------------------------------------------------------------

func (d *Day10) SetInput(lines []string) {
	d.machines = d.machines[:0]

	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}

		// Extract segments: [lights] (...) (...) ... {jolts}
		// Example:
		// [.##.] (3) (1,3) (2) (2,3) (0,2) (0,1) {3,5,4,7}

		// --- Extract lights ---
		i1 := strings.Index(line, "[")
		i2 := strings.Index(line, "]")
		lightsStr := line[i1+1 : i2]

		targetLights := make([]int, len(lightsStr))
		for i, ch := range lightsStr {
			if ch == '#' {
				targetLights[i] = 1
			} else {
				targetLights[i] = 0
			}
		}

		rest := line[i2+1:]

		// --- Extract joltage { } ---
		j1 := strings.Index(rest, "{")
		j2 := strings.Index(rest, "}")
		joltStr := rest[j1+1 : j2]
		joltParts := strings.Split(joltStr, ",")
		targetJolt := make([]int, len(joltParts))
		for i, p := range joltParts {
			v, _ := strconv.Atoi(strings.TrimSpace(p))
			targetJolt[i] = v
		}

		// --- Extract button tuples between ] and { ---
		buttonSection := strings.TrimSpace(rest[:j1])
		parts := strings.Fields(buttonSection)

		var buttons1 []int
		var buttons2 [][]int

		for _, p := range parts {
			if !strings.HasPrefix(p, "(") {
				continue
			}
			s := strings.Trim(p, "()")
			if s == "" {
				continue
			}
			idxStrs := strings.Split(s, ",")
			idxs := make([]int, len(idxStrs))
			for i, q := range idxStrs {
				v, _ := strconv.Atoi(q)
				idxs[i] = v
			}

			// Part 1: convert clicked indices to toggle bitmask
			mask := 0
			for _, i := range idxs {
				mask ^= (1 << i)
			}
			buttons1 = append(buttons1, mask)

			// Part 2: keep as slice of counters affected
			buttons2 = append(buttons2, idxs)
		}

		d.machines = append(d.machines, machine10{
			targetLights: targetLights,
			buttons1:     buttons1,
			targetJolt:   targetJolt,
			buttons2:     buttons2,
		})
	}
}

// ------------------------------------------------------------
// PART 1 — BFS on GF(2)
// ------------------------------------------------------------

func (d *Day10) SolvePart1() string {
	total := 0
	for _, m := range d.machines {
		total += solveLights(m)
	}
	return strconv.Itoa(total)
}

func solveLights(m machine10) int {
	n := len(m.targetLights)
	target := 0
	for i, v := range m.targetLights {
		if v == 1 {
			target |= 1 << i
		}
	}

	type node struct {
		val int
		d   int
	}

	visited := map[int]bool{}
	queue := []node{{0, 0}}
	visited[0] = true

	for len(queue) > 0 {
		cur := queue[0]
		queue = queue[1:]

		if cur.val == target {
			return cur.d
		}

		for _, b := range m.buttons1 {
			nv := cur.val ^ b
			if !visited[nv] {
				visited[nv] = true
				queue = append(queue, node{nv, cur.d + 1})
			}
		}
	}
	return 0
}

// ------------------------------------------------------------
// PART 2 — Multi-dimensional Dijkstra/BFS on bounded counter grid
// ------------------------------------------------------------

func (d *Day10) SolvePart2() string {
	total := 0
	for _, m := range d.machines {
		total += solveJolts(m)
	}
	return strconv.Itoa(total)
}

type state10 struct {
	cost int
	id   int
}

type pq10 []state10

func (p pq10) Len() int            { return len(p) }
func (p pq10) Less(i, j int) bool  { return p[i].cost < p[j].cost }
func (p pq10) Swap(i, j int)       { p[i], p[j] = p[j], p[i] }
func (p *pq10) Push(x interface{}) { *p = append(*p, x.(state10)) }
func (p *pq10) Pop() interface{} {
	old := *p
	n := len(old)
	v := old[n-1]
	*p = old[:n-1]
	return v
}

// solveJolts uses Dijkstra (all edges weight=1, so BFS-ish)
// State is tuple of counters, encoded into single int state ID.
func solveJolts(m machine10) int {
	K := len(m.targetJolt)
	T := m.targetJolt

	// Precompute dimension size multipliers for indexing
	mults := make([]int, K)
	mults[0] = 1
	for i := 1; i < K; i++ {
		mults[i] = mults[i-1] * (T[i-1] + 1)
	}

	// Encode & decode helpers
	encode := func(vec []int) int {
		id := 0
		for i := 0; i < K; i++ {
			id += vec[i] * mults[i]
		}
		return id
	}

	targetID := encode(T)

	// Priority queue for Dijkstra
	pq := &pq10{}
	heap.Init(pq)
	heap.Push(pq, state10{cost: 0, id: 0}) // start at all zeros

	visited := map[int]int{} // id → best cost so far
	visited[0] = 0

	buf := make([]int, K)

	for pq.Len() > 0 {
		cur := heap.Pop(pq).(state10)

		if cur.id == targetID {
			return cur.cost
		}

		// Decode cur.id into counters
		tmp := cur.id
		for i := K - 1; i >= 0; i-- {
			den := mults[i]
			buf[i] = tmp / den
			tmp %= den
		}

		// Try pressing each button
		for _, button := range m.buttons2 {
			next := buf
			// copy current counters
			// (reuse buf slice ensures no allocations)
			for xx := 0; xx < K; xx++ {
				next[xx] = buf[xx]
			}

			valid := true
			for _, idx := range button {
				next[idx]++
				if next[idx] > T[idx] {
					valid = false
					break
				}
			}
			if !valid {
				continue
			}

			id2 := encode(next)
			newCost := cur.cost + 1
			old, ok := visited[id2]
			if !ok || newCost < old {
				visited[id2] = newCost
				heap.Push(pq, state10{cost: newCost, id: id2})
			}
		}
	}

	return 0
}
