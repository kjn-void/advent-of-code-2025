package days

type Day01 struct {
	input []string
}

func init() { Register(1, func() Solution { return &Day01{} }) }

func (d *Day01) SetInput(lines []string) {
	d.input = lines
}

func (d *Day01) SolvePart1() string {
	return "Day 1 Part 1 not implemented"
}

func (d *Day01) SolvePart2() string {
	return "Day 1 Part 2 not implemented"
}
