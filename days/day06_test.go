package days

import "testing"

var day06ExampleInput = []string{
	"123 328  51 64 ",
	" 45 64  387 23 ",
	"  6 98  215 314",
	"*   +   *   +  ",
}

func TestDay06Part1(t *testing.T) {
	s := &day06{}
	s.SetInput(day06ExampleInput)

	got := s.SolvePart1()
	want := "4277556"

	if got != want {
		t.Fatalf("Day06 Part1 example failed: got %s, want %s", got, want)
	}
}

func TestDay06Part2(t *testing.T) {
	s := &day06{}
	s.SetInput(day06ExampleInput)

	got := s.SolvePart2()
	want := "3263827"

	if got != want {
		t.Fatalf("Day06 Part2 example failed: got %s, want %s", got, want)
	}
}

func BenchmarkDay06(b *testing.B) {
	benchmarkDay(b, 6, func() Solution { return &day06{} })
}
