package days

type Solution interface {
	SetInput(lines []string)
	SolvePart1() string
	SolvePart2() string
}
