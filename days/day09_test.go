package days

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

var exampleDay09 = []string{
	"7,1",
	"11,1",
	"11,7",
	"9,7",
	"9,5",
	"2,5",
	"2,3",
	"7,3",
}

func TestDay09ExamplePart1(t *testing.T) {
	d := &Day09{}
	d.SetInput(exampleDay09)

	got := d.SolvePart1()
	want := "50"

	if got != want {
		t.Fatalf("Day09 Part1 got %s want %s", got, want)
	}
}

func TestDay09ExamplePart2(t *testing.T) {
	d := &Day09{}
	d.SetInput(exampleDay09)

	got := d.SolvePart2()
	want := "24"

	if got != want {
		t.Fatalf("Day09 Part2 got %s want %s", got, want)
	}
}

// loadRealInputDay09 loads the actual AoC input from input/day09.txt
func loadRealInputDay09(b *testing.B) []string {
	path := filepath.Join("..", "input", "day09.txt")

	data, err := os.ReadFile(path)
	if err != nil {
		b.Fatalf("Missing input file: %v", err)
	}
	raw := strings.Split(strings.TrimRight(string(data), "\n"), "\n")
	return raw
}

// ------------------------
// Benchmarks
// ------------------------

// 1. Benchmark SetInput() — parsing only
func BenchmarkDay09_SetInput(b *testing.B) {
	for b.Loop() {
		lines := loadRealInputDay09(b)
		s := &Day09{}
		s.SetInput(lines)
	}
}

// 2. Benchmark SolvePart1 only
func BenchmarkDay09_SolvePart1(b *testing.B) {
	lines := loadRealInputDay09(b)

	s := &Day09{}
	s.SetInput(lines)

	for b.Loop() {
		_ = s.SolvePart1()
	}
}

// 3. Benchmark SolvePart2 only
func BenchmarkDay09_SolvePart2(b *testing.B) {
	lines := loadRealInputDay09(b)

	s := &Day09{}
	s.SetInput(lines)

	for b.Loop() {
		_ = s.SolvePart2()
	}
}

// 4. Benchmark full pipeline (SetInput + Part1 + Part2)
func BenchmarkDay09_FullPipeline(b *testing.B) {
	for b.Loop() {
		lines := loadRealInputDay09(b)
		s := &Day09{}
		s.SetInput(lines)
		_ = s.SolvePart1()
		_ = s.SolvePart2()
	}
}
