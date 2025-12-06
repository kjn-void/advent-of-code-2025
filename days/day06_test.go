package days

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

var day06ExampleInput = []string{
	"123 328  51 64 ",
	" 45 64  387 23 ",
	"  6 98  215 314",
	"*   +   *   +  ",
}

func TestDay06Part1Example(t *testing.T) {
	s := &Day06{}
	s.SetInput(day06ExampleInput)

	got := s.SolvePart1()
	want := "4277556"

	if got != want {
		t.Fatalf("Day06 Part1 example failed: got %s, want %s", got, want)
	}
}

func TestDay06Part2Example(t *testing.T) {
	s := &Day06{}
	s.SetInput(day06ExampleInput)

	got := s.SolvePart2()
	want := "3263827"

	if got != want {
		t.Fatalf("Day06 Part2 example failed: got %s, want %s", got, want)
	}
}

// ------------------------
// Benchmarks
// ------------------------

func loadRealInputDay06(b *testing.B) []string {
	path := filepath.Join("..", "input", "day06.txt")

	data, err := os.ReadFile(path)
	if err != nil {
		b.Fatalf("Missing input file: %v", err)
	}
	return strings.Split(strings.TrimSpace(string(data)), "\n")
}

// 1. Benchmark SetInput() — parsing only
func BenchmarkDay06_SetInput(b *testing.B) {
	for b.Loop() {
		lines := loadRealInputDay06(b)
		s := &Day06{}
		s.SetInput(lines)
	}
}

// 2. Benchmark SolvePart1 only
func BenchmarkDay06_SolvePart1(b *testing.B) {
	lines := loadRealInputDay06(b)

	s := &Day06{}
	s.SetInput(lines)

	for b.Loop() {
		_ = s.SolvePart1()
	}
}

// 3. Benchmark SolvePart2 only
func BenchmarkDay06_SolvePart2(b *testing.B) {
	lines := loadRealInputDay06(b)

	s := &Day06{}
	s.SetInput(lines)

	for b.Loop() {
		_ = s.SolvePart2()
	}
}

// 4. Benchmark full pipeline (SetInput + Part1 + Part2)
func BenchmarkDay06_FullPipeline(b *testing.B) {
	for b.Loop() {
		lines := loadRealInputDay06(b)
		s := &Day06{}
		s.SetInput(lines)
		_ = s.SolvePart1()
		_ = s.SolvePart2()
	}
}
