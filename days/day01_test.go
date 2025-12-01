package days

import "testing"

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

func TestDay01Part1Example(t *testing.T) {
	s := &Day01{}
	s.SetInput(exampleDay01)

	got := s.SolvePart1()
	want := "3"

	if got != want {
		t.Fatalf("Day01 Part1 example failed: got %s, want %s", got, want)
	}
}

func TestDay01Part2Example(t *testing.T) {
	s := &Day01{}
	s.SetInput(exampleDay01)

	got := s.SolvePart2()
	want := "6"

	if got != want {
		t.Fatalf("Day01 Part2 example failed: got %s, want %s", got, want)
	}
}
