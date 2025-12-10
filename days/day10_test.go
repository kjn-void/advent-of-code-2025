package days

import (
	"os"
	"path/filepath"
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

// loadRealInputDay10 loads the actual AoC input from input/day09.txt
func loadRealInputDay10(b *testing.B) []string {
	path := filepath.Join("..", "input", "day10.txt")

	data, err := os.ReadFile(path)
	if err != nil {
		b.Fatalf("Missing input file: %v", err)
	}
	raw := strings.Split(strings.TrimRight(string(data), "\n"), "\n")
	return raw
}

// day10_bench_test.go
func BenchmarkDay10_SetInput(b *testing.B) {
	lines := loadRealInputDay10(b)

	var d Day10
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		d.SetInput(lines)
	}
}

func BenchmarkDay10_SolvePart1(b *testing.B) {
	lines := loadRealInputDay10(b)

	var d Day10
	d.SetInput(lines)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = d.SolvePart1()
	}
}

func BenchmarkDay10_SolvePart2(b *testing.B) {
	lines := loadRealInputDay10(b)

	var d Day10
	d.SetInput(lines)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = d.SolvePart2()
	}
}

func BenchmarkDay10_FullPipeline(b *testing.B) {
	lines := loadRealInputDay10(b)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		var d Day10
		d.SetInput(lines)
		_ = d.SolvePart1()
		_ = d.SolvePart2()
	}
}
