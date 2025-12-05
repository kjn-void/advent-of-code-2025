package days

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

var day05ExampleInput = []string{
	"3-5",
	"10-14",
	"16-20",
	"12-18",
	"",
	"1",
	"5",
	"8",
	"11",
	"17",
	"32",
}

func TestDay05Part1Example(t *testing.T) {
	s := &Day05{}
	s.SetInput(day05ExampleInput)

	got := s.SolvePart1()
	want := "3"

	if got != want {
		t.Fatalf("Day05 Part1 example failed: got %s, want %s", got, want)
	}
}

func TestDay05Part2Example(t *testing.T) {
	s := &Day05{}
	s.SetInput(day05ExampleInput)

	got := s.SolvePart2()
	want := "14"

	if got != want {
		t.Fatalf("Day05 Part2 example failed: got %s, want %s", got, want)
	}
}

// ------------------------
// Benchmarks
// ------------------------

func loadRealInputDay05(b *testing.B) []string {
	path := filepath.Join("..", "input", "day05.txt")

	data, err := os.ReadFile(path)
	if err != nil {
		b.Fatalf("Missing input file: %v", err)
	}
	return strings.Split(strings.TrimSpace(string(data)), "\n")
}

// 1. Benchmark SetInput() — parsing ranges + IDs
func BenchmarkDay05_SetInput(b *testing.B) {
	for b.Loop() {
		lines := loadRealInputDay05(b)
		s := &Day05{}
		s.SetInput(lines)
	}
}

// 2. Benchmark SolvePart1 only
func BenchmarkDay05_SolvePart1(b *testing.B) {
	lines := loadRealInputDay05(b)

	s := &Day05{}
	s.SetInput(lines)

	for b.Loop() {
		_ = s.SolvePart1()
	}
}

// 3. Benchmark SolvePart2 only
func BenchmarkDay05_SolvePart2(b *testing.B) {
	lines := loadRealInputDay05(b)

	s := &Day05{}
	s.SetInput(lines)

	for b.Loop() {
		_ = s.SolvePart2()
	}
}

// 4. Benchmark full pipeline (SetInput + Part1 + Part2)
func BenchmarkDay05_FullPipeline(b *testing.B) {
	for b.Loop() {
		lines := loadRealInputDay05(b)
		s := &Day05{}
		s.SetInput(lines)
		_ = s.SolvePart1()
		_ = s.SolvePart2()
	}
}
