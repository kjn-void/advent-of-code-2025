package days

import "testing"

var day02ExampleInput = []string{
	"11-22,95-115,998-1012,1188511880-1188511890,222220-222224," +
		"1698522-1698528,446443-446449,38593856-38593862,565653-565659," +
		"824824821-824824827,2121212118-2121212124",
}

func TestDay02Part1(t *testing.T) {
	s := &day02{}
	s.SetInput(day02ExampleInput)

	got := s.SolvePart1()
	want := "1227775554"

	if got != want {
		t.Fatalf("Day02 Part1 example failed: got %s, want %s", got, want)
	}
}

func TestDay02Part2(t *testing.T) {
	s := &day02{}
	s.SetInput(day02ExampleInput)

	got := s.SolvePart2()
	want := "4174379265"

	if got != want {
		t.Fatalf("Day02 Part2 example failed: got %s, want %s", got, want)
	}
}

func BenchmarkDay02(b *testing.B) {
	benchmarkDay(b, 2, func() Solution { return &day02{} })
}
