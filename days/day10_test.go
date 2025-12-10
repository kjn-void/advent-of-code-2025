package days

import (
	"strings"
	"testing"
)

var day10Example = `
[.##.] (3) (1,3) (2) (2,3) (0,2) (0,1) {3,5,4,7}
[...#.] (0,2,3,4) (2,3) (0,4) (0,1,2) (1,2,3,4) {7,5,12,7,2}
[.###.#] (0,1,2,3,4) (0,3,4) (0,1,2,4,5) (1,2) {10,11,11,5,10,5}
`

func TestDay10_Part1(t *testing.T) {
	lines := strings.Split(strings.TrimSpace(day10Example), "\n")
	d := &Day10{}
	d.SetInput(lines)

	got := d.SolvePart1()
	want := "7"

	if got != want {
		t.Fatalf("Part1: got %s want %s", got, want)
	}
}

func TestDay10_Part2(t *testing.T) {
	lines := strings.Split(strings.TrimSpace(day10Example), "\n")
	d := &Day10{}
	d.SetInput(lines)

	got := d.SolvePart2()
	want := "33"

	if got != want {
		t.Fatalf("Part2: got %s want %s", got, want)
	}
}
