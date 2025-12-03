package days

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

var day03ExampleInput = []string{
	"987654321111111",
	"811111111111119",
	"234234234234278",
	"818181911112111",
}

func TestDay03Part1(t *testing.T) {
	s := &Day03{}
	s.SetInput(day03ExampleInput)

	got := s.SolvePart1()
	want := "357"

	if got != want {
		t.Fatalf("Day03 Part1 example failed: got %s, want %s", got, want)
	}
}

func TestDay03Part2(t *testing.T) {
	s := &Day03{}
	s.SetInput(day03ExampleInput)

	got := s.SolvePart2()
	want := "3121910778619"

	if got != want {
		t.Fatalf("Day03 Part2 example failed: got %s, want %s", got, want)
	}
}

// ------------------------
// Benchmarks
// ------------------------

func loadRealInputDay03(b *testing.B) []string {
	path := filepath.Join("..", "input", "day03.txt")

	data, err := os.ReadFile(path)
	if err != nil {
		b.Fatalf("Missing input file: %v", err)
	}
	return strings.Split(strings.TrimSpace(string(data)), "\n")
}

// 1. Benchmark SetInput() — parsing cost only
func BenchmarkDay03_SetInput(b *testing.B) {
	for b.Loop() {
		lines := loadRealInputDay03(b)
		s := &Day03{}
		s.SetInput(lines)
	}
}

// 2. Benchmark SolvePart1 only
func BenchmarkDay03_SolvePart1(b *testing.B) {
	lines := loadRealInputDay03(b)

	s := &Day03{}
	s.SetInput(lines)

	for b.Loop() {
		_ = s.SolvePart1()
	}
}

// 3. Benchmark SolvePart2 only
func BenchmarkDay03_SolvePart2(b *testing.B) {
	lines := loadRealInputDay03(b)

	s := &Day03{}
	s.SetInput(lines)

	for b.Loop() {
		_ = s.SolvePart2()
	}
}

// 4. Benchmark full pipeline (SetInput + Part1 + Part2)
func BenchmarkDay03_FullPipeline(b *testing.B) {
	for b.Loop() {
		lines := loadRealInputDay03(b)
		s := &Day03{}
		s.SetInput(lines)
		_ = s.SolvePart1()
		_ = s.SolvePart2()
	}
}
