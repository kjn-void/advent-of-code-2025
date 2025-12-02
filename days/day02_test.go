package days

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

var day02ExampleInput = []string{
	"11-22,95-115,998-1012,1188511880-1188511890,222220-222224," +
		"1698522-1698528,446443-446449,38593856-38593862,565653-565659," +
		"824824821-824824827,2121212118-2121212124",
}

func TestDay02Part1(t *testing.T) {
	s := &Day02{}
	s.SetInput(day02ExampleInput)

	got := s.SolvePart1()
	want := "1227775554"

	if got != want {
		t.Fatalf("Day02 Part1 example failed: got %s, want %s", got, want)
	}
}

func TestDay02Part2(t *testing.T) {
	s := &Day02{}
	s.SetInput(day02ExampleInput)

	got := s.SolvePart2()
	want := "4174379265"

	if got != want {
		t.Fatalf("Day02 Part2 example failed: got %s, want %s", got, want)
	}
}

// ------------------------
// Benchmarks
// ------------------------

func loadRealInputDay02(b *testing.B) []string {
	path := filepath.Join("..", "input", "day02.txt")

	data, err := os.ReadFile(path)
	if err != nil {
		b.Fatalf("Missing input file: %v", err)
	}
	return strings.Split(strings.TrimSpace(string(data)), "\n")
}

// 1. Benchmark SetInput() — parsing only
func BenchmarkDay02_SetInput(b *testing.B) {
	for b.Loop() {
		lines := loadRealInputDay02(b)
		s := &Day02{}
		s.SetInput(lines)
	}
}

// 2. Benchmark SolvePart1 only
func BenchmarkDay02_SolvePart1(b *testing.B) {
	lines := loadRealInputDay02(b)

	s := &Day02{}
	s.SetInput(lines)

	for b.Loop() {
		_ = s.SolvePart1()
	}
}

// 3. Benchmark SolvePart2 only
func BenchmarkDay02_SolvePart2(b *testing.B) {
	lines := loadRealInputDay02(b)

	s := &Day02{}
	s.SetInput(lines)

	for b.Loop() {
		_ = s.SolvePart2()
	}
}

// 4. Benchmark full pipeline (SetInput + Part1 + Part2)
func BenchmarkDay02_FullPipeline(b *testing.B) {
	for b.Loop() {
		lines := loadRealInputDay02(b)
		s := &Day02{}
		s.SetInput(lines)
		_ = s.SolvePart1()
		_ = s.SolvePart2()
	}
}
