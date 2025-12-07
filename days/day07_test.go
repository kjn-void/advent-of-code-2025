package days

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

var day07ExampleInput = []string{
	".......S.......",
	"...............",
	".......^.......",
	"...............",
	"......^.^......",
	"...............",
	".....^.^.^.....",
	"...............",
	"....^.^...^....",
	"...............",
	"...^.^...^.^...",
	"...............",
	"..^...^.....^..",
	"...............",
	".^.^.^.^.^...^.",
	"...............",
}

func TestDay07Part1(t *testing.T) {
	s := &Day07{}
	s.SetInput(day07ExampleInput)

	got := s.SolvePart1()
	want := "21"

	if got != want {
		t.Fatalf("Day07 Part1 example failed: got %s, want %s", got, want)
	}
}

func TestDay07Part2(t *testing.T) {
	s := &Day07{}
	s.SetInput(day07ExampleInput)

	got := s.SolvePart2()
	want := "40"

	if got != want {
		t.Fatalf("Day07 Part2 example failed: got %s, want %s", got, want)
	}
}

// ------------------------
// Benchmarks
// ------------------------

func loadRealInputDay07(b *testing.B) []string {
	path := filepath.Join("..", "input", "day07.txt")

	data, err := os.ReadFile(path)
	if err != nil {
		b.Fatalf("Missing input file: %v", err)
	}
	return strings.Split(strings.TrimRight(string(data), "\n"), "\n")
}

// 1. Benchmark SetInput() — parsing and normalization only
func BenchmarkDay07_SetInput(b *testing.B) {
	for b.Loop() {
		lines := loadRealInputDay07(b)
		s := &Day07{}
		s.SetInput(lines)
	}
}

// 2. Benchmark SolvePart1 only
func BenchmarkDay07_SolvePart1(b *testing.B) {
	lines := loadRealInputDay07(b)

	s := &Day07{}
	s.SetInput(lines)

	for b.Loop() {
		_ = s.SolvePart1()
	}
}

// 3. Benchmark SolvePart2 only
func BenchmarkDay07_SolvePart2(b *testing.B) {
	lines := loadRealInputDay07(b)

	s := &Day07{}
	s.SetInput(lines)

	for b.Loop() {
		_ = s.SolvePart2()
	}
}

// 4. Benchmark full pipeline
func BenchmarkDay07_FullPipeline(b *testing.B) {
	for b.Loop() {
		lines := loadRealInputDay07(b)
		s := &Day07{}
		s.SetInput(lines)
		_ = s.SolvePart1()
		_ = s.SolvePart2()
	}
}
