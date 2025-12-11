package days

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

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
	var d Day11
	lines := splitLines(day11ExamplePart1)
	d.SetInput(lines)

	got := d.SolvePart1()
	want := "5"
	if got != want {
		t.Fatalf("Part1 example: got %s, want %s", got, want)
	}
}

func TestDay11_Part2_Example(t *testing.T) {
	var d Day11
	lines := splitLines(day11ExamplePart2)
	d.SetInput(lines)

	got := d.SolvePart2()
	want := "2"
	if got != want {
		t.Fatalf("Part2 example: got %s, want %s", got, want)
	}
}

// loadRealInputDay11 loads the actual AoC input from input/day11.txt
func loadRealInputDay11(b *testing.B) []string {
	path := filepath.Join("..", "input", "day11.txt")

	data, err := os.ReadFile(path)
	if err != nil {
		b.Fatalf("Missing input file: %v", err)
	}
	raw := strings.Split(strings.TrimRight(string(data), "\n"), "\n")
	return raw
}

func BenchmarkDay11_SetInput(b *testing.B) {
	lines := loadRealInputDay11(b)

	var d Day11

	for b.Loop() {
		d.SetInput(lines)
	}
}

func BenchmarkDay11_SolvePart1(b *testing.B) {
	lines := loadRealInputDay11(b)

	var d Day11
	d.SetInput(lines)

	for b.Loop() {
		_ = d.SolvePart1()
	}
}

func BenchmarkDay11_SolvePart2(b *testing.B) {
	lines := loadRealInputDay11(b)

	var d Day11
	d.SetInput(lines)

	for b.Loop() {
		_ = d.SolvePart2()
	}
}

func BenchmarkDay11_FullPipeline(b *testing.B) {
	lines := loadRealInputDay11(b)

	for b.Loop() {
		var d Day11
		d.SetInput(lines)
		_ = d.SolvePart1()
		_ = d.SolvePart2()
	}
}
