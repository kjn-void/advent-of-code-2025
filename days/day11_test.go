package days

import "testing"

const day11ExamplePart1 = `
aaa: you hhh
you: bbb ccc
bbb: ddd eee
ccc: ddd eee fff
ddd: ggg
eee: out
fff: out
ggg: out
hhh: ccc fff iii
iii: out
`

const day11ExamplePart2 = `
svr: aaa bbb
aaa: fft
fft: ccc
bbb: tty
tty: ccc
ccc: ddd eee
ddd: hub
hub: fff
eee: dac
dac: fff
fff: ggg hhh
ggg: out
hhh: out
`

func TestDay11_Part1_Example(t *testing.T) {
	var d day11
	lines := splitLines(day11ExamplePart1)
	d.SetInput(lines)

	got := d.SolvePart1()
	want := "5"
	if got != want {
		t.Fatalf("Part1 example: got %s, want %s", got, want)
	}
}

func TestDay11_Part2_Example(t *testing.T) {
	var d day11
	lines := splitLines(day11ExamplePart2)
	d.SetInput(lines)

	got := d.SolvePart2()
	want := "2"
	if got != want {
		t.Fatalf("Part2 example: got %s, want %s", got, want)
	}
}

func BenchmarkDay11(b *testing.B) {
	benchmarkDay(b, 11, func() Solution { return &day11{} })
}
