package days

import "testing"

var day05ExampleInput = []string{
	"3-5",
	"10-14",
	"16-20",
	"12-18",
	"",
	"1",
	"5",
	"8",
	"11",
	"17",
	"32",
}

func TestDay05Part1Example(t *testing.T) {
	s := &day05{}
	s.SetInput(day05ExampleInput)

	got := s.SolvePart1()
	want := "3"

	if got != want {
		t.Fatalf("Day05 Part1 example failed: got %s, want %s", got, want)
	}
}

func TestDay05Part2Example(t *testing.T) {
	s := &day05{}
	s.SetInput(day05ExampleInput)

	got := s.SolvePart2()
	want := "14"

	if got != want {
		t.Fatalf("Day05 Part2 example failed: got %s, want %s", got, want)
	}
}

func BenchmarkDay05(b *testing.B) {
	benchmarkDay(b, 5, func() Solution { return &day05{} })
}
