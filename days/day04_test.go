package days

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

var day04ExampleInput = []string{
	"..@@.@@@@.",
	"@@@.@.@.@@",
	"@@@@@.@.@@",
	"@.@@@@..@.",
	"@@.@@@@.@@",
	".@@@@@@@.@",
	".@.@.@.@@@",
	"@.@@@.@@@@",
	".@@@@@@@@.",
	"@.@.@@@.@.",
}

func TestDay04Part1(t *testing.T) {
	s := &Day04{}
	s.SetInput(day04ExampleInput)

	got := s.SolvePart1()
	want := "13"

	if got != want {
		t.Fatalf("Day04 Part1 example failed: got %s, want %s", got, want)
	}
}

func TestDay04Part2(t *testing.T) {
	s := &Day04{}
	s.SetInput(day04ExampleInput)

	got := s.SolvePart2()
	want := "43"

	if got != want {
		t.Fatalf("Day04 Part2 example failed: got %s, want %s", got, want)
	}
}

// ------------------------
// Benchmarks
// ------------------------

func loadRealInputDay04(b *testing.B) []string {
	path := filepath.Join("..", "input", "day04.txt")

	data, err := os.ReadFile(path)
	if err != nil {
		b.Fatalf("Missing input file: %v", err)
	}
	return strings.Split(strings.TrimSpace(string(data)), "\n")
}

// 1. Benchmark SetInput() — parsing only
func BenchmarkDay04_SetInput(b *testing.B) {
	for b.Loop() {
		lines := loadRealInputDay04(b)
		s := &Day04{}
		s.SetInput(lines)
	}
}

// 2. Benchmark SolvePart1 only
func BenchmarkDay04_SolvePart1(b *testing.B) {
	lines := loadRealInputDay04(b)

	s := &Day04{}
	s.SetInput(lines)

	for b.Loop() {
		_ = s.SolvePart1()
	}
}

// 3. Benchmark SolvePart2 only
func BenchmarkDay04_SolvePart2(b *testing.B) {
	lines := loadRealInputDay04(b)

	s := &Day04{}
	s.SetInput(lines)

	for b.Loop() {
		_ = s.SolvePart2()
	}
}

// 4. Benchmark full pipeline (SetInput + Part1 + Part2)
func BenchmarkDay04_FullPipeline(b *testing.B) {
	for b.Loop() {
		lines := loadRealInputDay04(b)
		s := &Day04{}
		s.SetInput(lines)
		_ = s.SolvePart1()
		_ = s.SolvePart2()
	}
}
