package days

import "testing"

var exampleDay09 = []string{
	"7,1",
	"11,1",
	"11,7",
	"9,7",
	"9,5",
	"2,5",
	"2,3",
	"7,3",
}

func TestDay09ExamplePart1(t *testing.T) {
	d := &day09{}
	d.SetInput(exampleDay09)

	got := d.SolvePart1()
	want := "50"

	if got != want {
		t.Fatalf("Day09 Part1 got %s want %s", got, want)
	}
}

func TestDay09ExamplePart2(t *testing.T) {
	d := &day09{}
	d.SetInput(exampleDay09)

	got := d.SolvePart2()
	want := "24"

	if got != want {
		t.Fatalf("Day09 Part2 got %s want %s", got, want)
	}
}

func BenchmarkDay09(b *testing.B) {
	benchmarkDay(b, 9, func() Solution { return &day09{} })
}
