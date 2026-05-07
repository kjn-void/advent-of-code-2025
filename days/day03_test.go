package days

import "testing"

var day03ExampleInput = []string{
	"987654321111111",
	"811111111111119",
	"234234234234278",
	"818181911112111",
}

func TestDay03Part1(t *testing.T) {
	s := &day03{}
	s.SetInput(day03ExampleInput)

	got := s.SolvePart1()
	want := "357"

	if got != want {
		t.Fatalf("Day03 Part1 example failed: got %s, want %s", got, want)
	}
}

func TestDay03Part2(t *testing.T) {
	s := &day03{}
	s.SetInput(day03ExampleInput)

	got := s.SolvePart2()
	want := "3121910778619"

	if got != want {
		t.Fatalf("Day03 Part2 example failed: got %s, want %s", got, want)
	}
}

func BenchmarkDay03(b *testing.B) {
	benchmarkDay(b, 3, func() Solution { return &day03{} })
}
