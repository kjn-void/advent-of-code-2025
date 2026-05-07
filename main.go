package main

import (
	"fmt"
	"os"
	"strconv"

	"aoc2025/days"
)

func main() {
	if len(os.Args) < 2 {
		printUsage()
		os.Exit(1)
	}

	verbose, dayArgs := parseArgs(os.Args[1:])
	if len(dayArgs) == 0 {
		printUsage()
		os.Exit(1)
	}

	descriptions := map[int]problemDescription{}
	if verbose {
		loaded, err := LoadProblemDescriptions("problems.yaml")
		if err != nil {
			fmt.Printf("Warning: could not load problem descriptions: %v\n", err)
		} else {
			descriptions = loaded
		}
	}

	for _, arg := range dayArgs {
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
		if verbose {
			printProblemDescription(day, descriptions)
		}
		fmt.Printf("Part 1: %s\n", solver.SolvePart1())
		fmt.Printf("Part 2: %s\n", solver.SolvePart2())
		fmt.Println()
	}
}

func parseArgs(args []string) (bool, []string) {
	verbose := false
	dayArgs := make([]string, 0, len(args))

	for _, arg := range args {
		switch arg {
		case "-v", "--verbose":
			verbose = true
		case "-h", "--help":
			printUsage()
			os.Exit(0)
		default:
			dayArgs = append(dayArgs, arg)
		}
	}

	return verbose, dayArgs
}

func printProblemDescription(day int, descriptions map[int]problemDescription) {
	problem, ok := descriptions[day]
	if !ok {
		return
	}

	if problem.Title != "" {
		fmt.Printf("%s\n", problem.Title)
	}
	if problem.Description != "" {
		fmt.Printf("%s\n", problem.Description)
	}
}

func printUsage() {
	fmt.Println("Usage: ./aoc2025 [-v|--verbose] <day> [<day> ...]")
}
