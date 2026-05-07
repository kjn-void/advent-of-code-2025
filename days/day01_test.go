package days

import "testing"

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
	s := &day01{}
	s.SetInput(exampleDay01)

	got := s.SolvePart1()
	want := "3"

	if got != want {
		t.Fatalf("Day01 Part1 example failed: got %s, want %s", got, want)
	}
}

func TestDay01Part2(t *testing.T) {
	s := &day01{}
	s.SetInput(exampleDay01)

	got := s.SolvePart2()
	want := "6"

	if got != want {
		t.Fatalf("Day01 Part2 example failed: got %s, want %s", got, want)
	}
}

func BenchmarkDay01(b *testing.B) {
	benchmarkDay(b, 1, func() Solution { return &day01{} })
}
