package main

import (
	"fmt"
	"os"
	"strconv"

	"aoc2025/days"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: ./aoc2025 <day> [<day> ...]")
		os.Exit(1)
	}

	for _, arg := range os.Args[1:] {
		day, err := strconv.Atoi(arg)
		if err != nil || day < 1 || day > 12 {
			fmt.Printf("Invalid day: %s\n", arg)
			continue
		}

		solver, ok := days.Get(day)
		if !ok {
			fmt.Printf("No solver for day %d\n", day)
			continue
		}

		lines, err := FetchOrReadInput(day)
		if err != nil {
			fmt.Printf("Error loading input for day %d: %v\n", day, err)
			continue
		}

		solver.SetInput(lines)

		fmt.Printf("🌟 Day %d 🌟\n", day)
		fmt.Printf("Part 1: %s\n", solver.SolvePart1())
		fmt.Printf("Part 2: %s\n", solver.SolvePart2())
		fmt.Println()
	}
}
