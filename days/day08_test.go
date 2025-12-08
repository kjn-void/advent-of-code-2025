package days

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

var exampleDay08 = []string{
	"162,817,812",
	"57,618,57",
	"906,360,560",
	"592,479,940",
	"352,342,300",
	"466,668,158",
	"542,29,236",
	"431,825,988",
	"739,650,466",
	"52,470,668",
	"216,146,977",
	"819,987,18",
	"117,168,530",
	"805,96,715",
	"346,949,466",
	"970,615,88",
	"941,993,340",
	"862,61,35",
	"984,92,344",
	"425,690,689",
}

func TestDay08ExamplePart1(t *testing.T) {
	d := &Day08{}
	d.SetInput(exampleDay08)

	// Example uses 10 shortest connections, not 1000
	sizes := runConnections(d.points, d.edges, 10)

	if len(sizes) < 3 {
		t.Fatalf("Expected at least 3 components, got %v", sizes)
	}

	got := sizes[0] * sizes[1] * sizes[2]
	want := 40

	if got != want {
		t.Fatalf("Day 08 Part 1 example: got %d, want %d", got, want)
	}
}

func TestDay08ExamplePart2(t *testing.T) {
	d := &Day08{}
	d.SetInput(exampleDay08)

	i, j := runUntilSingleCircuit(d.points, d.edges)
	xa := d.points[i].x
	xb := d.points[j].x
	got := xa * xb
	var want int64 = 25272

	if got != want {
		t.Fatalf("Day 08 Part 2 example: got %d, want %d", got, want)
	}
}

// ------------------------
// Benchmarks
// ------------------------

func loadRealInputDay08(b *testing.B) []string {
	path := filepath.Join("..", "input", "day08.txt")

	data, err := os.ReadFile(path)
	if err != nil {
		b.Fatalf("Missing input file: %v", err)
	}
	return strings.Split(strings.TrimRight(string(data), "\n"), "\n")
}

// 1. Benchmark SetInput() — parsing into vec3 only
func BenchmarkDay08_SetInput(b *testing.B) {
	for b.Loop() {
		lines := loadRealInputDay08(b)
		s := &Day08{}
		s.SetInput(lines)
	}
}

// 2. Benchmark SolvePart1 only
func BenchmarkDay08_SolvePart1(b *testing.B) {
	lines := loadRealInputDay08(b)

	s := &Day08{}
	s.SetInput(lines)

	for b.Loop() {
		_ = s.SolvePart1()
	}
}

// 3. Benchmark SolvePart2 only
func BenchmarkDay08_SolvePart2(b *testing.B) {
	lines := loadRealInputDay08(b)

	s := &Day08{}
	s.SetInput(lines)

	for b.Loop() {
		_ = s.SolvePart2()
	}
}

// 4. Benchmark full pipeline
func BenchmarkDay08_FullPipeline(b *testing.B) {
	for b.Loop() {
		lines := loadRealInputDay08(b)
		s := &Day08{}
		s.SetInput(lines)
		_ = s.SolvePart1()
		_ = s.SolvePart2()
	}
}
