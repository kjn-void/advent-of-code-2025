package days

import (
	"os"
	"path/filepath"
	"testing"
)

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
	var d Day12
	lines := splitLines(day12Example) // shared helper from day10_test.go
	d.SetInput(lines)

	got := d.SolvePart1()
	want := "2"
	if got != want {
		t.Fatalf("Part1 example: got %s, want %s", got, want)
	}
}

// --- Benchmarks ------------------------------------------------------------

func loadRealInputDay12(b *testing.B) []string {
	path := filepath.Join("..", "input", "day12.txt")

	data, err := os.ReadFile(path)
	if err != nil {
		b.Fatalf("Missing input file: %v", err)
	}
	// reuse splitLines semantics via simple split helper from other tests
	return splitLines(string(data))
}

func BenchmarkDay12_SetInput(b *testing.B) {
	lines := loadRealInputDay12(b)

	var d Day12
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		d.SetInput(lines)
	}
}

func BenchmarkDay12_SolvePart1(b *testing.B) {
	lines := loadRealInputDay12(b)

	var d Day12
	d.SetInput(lines)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = d.SolvePart1()
	}
}

func BenchmarkDay12_FullPipeline(b *testing.B) {
	lines := loadRealInputDay12(b)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		var d Day12
		d.SetInput(lines)
		_ = d.SolvePart1()
	}
}
