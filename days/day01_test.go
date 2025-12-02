package days

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// ------------------------
// Example test data
// ------------------------

var exampleDay01 = []string{
	"L68",
	"L30",
	"R48",
	"L5",
	"R60",
	"L55",
	"L1",
	"L99",
	"R14",
	"L82",
}

// ------------------------
// Unit tests
// ------------------------

func TestDay01Part1(t *testing.T) {
	s := &Day01{}
	s.SetInput(exampleDay01)

	got := s.SolvePart1()
	want := "3"

	if got != want {
		t.Fatalf("Day01 Part1 example failed: got %s, want %s", got, want)
	}
}

func TestDay01Part2(t *testing.T) {
	s := &Day01{}
	s.SetInput(exampleDay01)

	got := s.SolvePart2()
	want := "6"

	if got != want {
		t.Fatalf("Day01 Part2 example failed: got %s, want %s", got, want)
	}
}

// ------------------------
// Benchmarks
// ------------------------

// loadRealInput loads the actual AoC input from input/day01.txt
func loadRealInput(b *testing.B) []string {
	path := filepath.Join("..", "input", "day01.txt")

	data, err := os.ReadFile(path)
	if err != nil {
		b.Fatalf("Missing input file: %v", err)
	}
	raw := strings.Split(strings.TrimSpace(string(data)), "\n")
	return raw
}

// 1. Benchmark SetInput() — parsing overhead only
func BenchmarkDay01_SetInput(b *testing.B) {
	for b.Loop() {
		lines := loadRealInput(b)
		s := &Day01{}
		s.SetInput(lines)
	}
}

// 2. Benchmark SolvePart1 only
func BenchmarkDay01_SolvePart1(b *testing.B) {
	lines := loadRealInput(b)

	s := &Day01{}
	s.SetInput(lines)

	for b.Loop() {
		_ = s.SolvePart1()
	}
}

// 3. Benchmark SolvePart2 only
func BenchmarkDay01_SolvePart2(b *testing.B) {
	lines := loadRealInput(b)

	s := &Day01{}
	s.SetInput(lines)

	for b.Loop() {
		_ = s.SolvePart2()
	}
}

// 4. Benchmark full pipeline (SetInput + Part1 + Part2)
func BenchmarkDay01_FullPipeline(b *testing.B) {
	for b.Loop() {
		lines := loadRealInput(b)
		s := &Day01{}
		s.SetInput(lines)
		_ = s.SolvePart1()
		_ = s.SolvePart2()
	}
}
