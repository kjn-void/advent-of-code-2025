package days

import (
	"strings"
	"testing"
)

const day10Example = `
[.##.] (3) (1,3) (2) (2,3) (0,2) (0,1) {3,5,4,7}
[...#.] (0,2,3,4) (2,3) (0,4) (0,1,2) (1,2,3,4) {7,5,12,7,2}
[.###.#] (0,1,2,3,4) (0,3,4) (0,1,2,4,5) (1,2) {10,11,11,5,10,5}
`

func TestDay10_Part1_Example(t *testing.T) {
	var d Day10
	lines := splitLines(day10Example)
	d.SetInput(lines)

	got := d.SolvePart1()
	want := "7"
	if got != want {
		t.Fatalf("Part1 example: got %s, want %s", got, want)
	}
}

func TestDay10_Part2_Example(t *testing.T) {
	var d Day10
	lines := splitLines(day10Example)
	d.SetInput(lines)

	got := d.SolvePart2()
	want := "33"
	if got != want {
		t.Fatalf("Part2 example: got %s, want %s", got, want)
	}
}

// splitLines is a tiny helper shared by tests in this package.
func splitLines(s string) []string {
	raw := strings.Split(s, "\n")
	out := make([]string, 0, len(raw))
	for _, line := range raw {
		line = strings.TrimRight(line, "\r")
		if line == "" {
			continue
		}
		out = append(out, line)
	}
	return out
}
