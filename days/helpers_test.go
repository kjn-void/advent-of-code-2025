package days

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// loadRealInput reads input/dayNN.txt for benchmarks and returns its lines,
// preserving any meaningful blank lines inside the file.
func loadRealInput(b *testing.B, day int) []string {
	b.Helper()

	path := filepath.Join("..", "input", fmt.Sprintf("day%02d.txt", day))
	data, err := os.ReadFile(path)
	if err != nil {
		b.Fatalf("Missing input file: %v", err)
	}

	return strings.Split(strings.TrimRight(string(data), "\r\n"), "\n")
}

// benchmarkDay runs the standard Advent of Code benchmark suite for one day:
// parsing input, solving each part, and the full parse-plus-solve pipeline.
func benchmarkDay(b *testing.B, day int, newSolution func() Solution) {
	b.Helper()

	lines := loadRealInput(b, day)

	b.Run("SetInput", func(b *testing.B) {
		for b.Loop() {
			s := newSolution()
			s.SetInput(lines)
		}
	})

	b.Run("SolvePart1", func(b *testing.B) {
		s := newSolution()
		s.SetInput(lines)

		b.ResetTimer()
		for b.Loop() {
			_ = s.SolvePart1()
		}
	})

	b.Run("SolvePart2", func(b *testing.B) {
		s := newSolution()
		s.SetInput(lines)

		b.ResetTimer()
		for b.Loop() {
			_ = s.SolvePart2()
		}
	})

	b.Run("FullPipeline", func(b *testing.B) {
		for b.Loop() {
			s := newSolution()
			s.SetInput(lines)
			_ = s.SolvePart1()
			_ = s.SolvePart2()
		}
	})
}
