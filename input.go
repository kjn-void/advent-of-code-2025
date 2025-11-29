package main

import (
	"aoc2025/aocnet"
	"bufio"
	"fmt"
	"os"
	"path/filepath"
)

func ReadLocalInput(path string) ([]string, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	var lines []string
	sc := bufio.NewScanner(f)
	for sc.Scan() {
		lines = append(lines, sc.Text())
	}

	return lines, sc.Err()
}

func FetchOrReadInput(day int) ([]string, error) {
	session := os.Getenv("AOC_SESSION")
	online := os.Getenv("AOC_ONLINE") == "1"

	if online {
		if session == "" {
			fmt.Println("AOC_ONLINE=1 but AOC_SESSION is not set.")
		} else {
			lines, err := aocnet.FetchInput(day, session)
			if err == nil {
				writeInputCache(day, lines)
				return lines, nil
			}
			fmt.Printf("Network fetch failed: %v\n", err)
		}
	}

	// Fall back to cached file
	inputFile := filepath.Join("input", fmt.Sprintf("day%02d.txt", day))
	return ReadLocalInput(inputFile)
}

func ensureDir(name string) {
	if _, err := os.Stat(name); os.IsNotExist(err) {
		os.MkdirAll(name, 0755)
	}
}

func writeInputCache(day int, lines []string) error {
	ensureDir("input")

	path := filepath.Join("input", fmt.Sprintf("day%02d.txt", day))
	f, err := os.Create(path)
	if err != nil {
		return err
	}
	defer f.Close()

	for _, l := range lines {
		f.WriteString(l + "\n")
	}

	return nil
}
