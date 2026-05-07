package days

import "testing"

var day04ExampleInput = []string{
	"..@@.@@@@.",
	"@@@.@.@.@@",
	"@@@@@.@.@@",
	"@.@@@@..@.",
	"@@.@@@@.@@",
	".@@@@@@@.@",
	".@.@.@.@@@",
	"@.@@@.@@@@",
	".@@@@@@@@.",
	"@.@.@@@.@.",
}

func TestDay04Part1(t *testing.T) {
	s := &day04{}
	s.SetInput(day04ExampleInput)

	got := s.SolvePart1()
	want := "13"

	if got != want {
		t.Fatalf("Day04 Part1 example failed: got %s, want %s", got, want)
	}
}

func TestDay04Part2(t *testing.T) {
	s := &day04{}
	s.SetInput(day04ExampleInput)

	got := s.SolvePart2()
	want := "43"

	if got != want {
		t.Fatalf("Day04 Part2 example failed: got %s, want %s", got, want)
	}
}

func BenchmarkDay04(b *testing.B) {
	benchmarkDay(b, 4, func() Solution { return &day04{} })
}
