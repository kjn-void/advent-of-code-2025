package days

import "testing"

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
	s := &day07{}
	s.SetInput(day07ExampleInput)

	got := s.SolvePart1()
	want := "21"

	if got != want {
		t.Fatalf("Day07 Part1 example failed: got %s, want %s", got, want)
	}
}

func TestDay07Part2(t *testing.T) {
	s := &day07{}
	s.SetInput(day07ExampleInput)

	got := s.SolvePart2()
	want := "40"

	if got != want {
		t.Fatalf("Day07 Part2 example failed: got %s, want %s", got, want)
	}
}

func BenchmarkDay07(b *testing.B) {
	benchmarkDay(b, 7, func() Solution { return &day07{} })
}
