package days

import "testing"

const day12Example = `
0:
###
##.
##.

1:
###
##.
.##

2:
.##
###
##.

3:
##.
###
##.

4:
###
#..
###

5:
###
.#.
###

4x4: 0 0 0 0 2 0
12x5: 1 0 1 0 2 2
12x5: 1 0 1 0 3 2
`

func TestDay12_Part1_Example(t *testing.T) {
	var d day12
	lines := splitLines(day12Example) // shared helper from day10_test.go
	d.SetInput(lines)

	got := d.SolvePart1()
	want := "2"
	if got != want {
		t.Fatalf("Part1 example: got %s, want %s", got, want)
	}
}

// --- Benchmarks ------------------------------------------------------------

func BenchmarkDay12(b *testing.B) {
	benchmarkDay(b, 12, func() Solution { return &day12{} })
}
