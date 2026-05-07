package main

import (
	"aoc2025/aocnet"
	"bufio"
	"fmt"
	"os"
	"path/filepath"
)

// ReadLocalInput reads the file at path line by line and returns the input as
// a slice of strings, or an error if the file cannot be opened or scanned.
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

// FetchOrReadInput loads the puzzle input for day, preferring an online fetch
// when AOC_ONLINE=1 and AOC_SESSION is set, then falling back to the local
// input/dayXX.txt cache. It returns the input lines or the cache-read error.
func FetchOrReadInput(day int) ([]string, error) {
	session := os.Getenv("AOC_SESSION")
	online := os.Getenv("AOC_ONLINE") == "1"

	if online {
		if session == "" {
			fmt.Println("AOC_ONLINE=1 but AOC_SESSION is not set.")
		} else {
			lines, err := aocnet.FetchInput(day, session)
			if err == nil {
				if err := writeInputCache(day, lines); err != nil {
					fmt.Printf("Input cache write failed: %v\n", err)
				}
				return lines, nil
			}
			fmt.Printf("Network fetch failed: %v\n", err)
		}
	}

	// Fall back to cached file
	inputFile := filepath.Join("input", fmt.Sprintf("day%02d.txt", day))
	return ReadLocalInput(inputFile)
}

// ensureDir verifies that name exists as a directory, creating it when missing,
// and returns any filesystem error encountered along the way.
func ensureDir(name string) error {
	if _, err := os.Stat(name); os.IsNotExist(err) {
		return os.MkdirAll(name, 0755)
	} else if err != nil {
		return err
	}
	return nil
}

// writeInputCache writes lines to input/dayXX.txt for day, creating the input
// directory if needed, and returns the first write or filesystem error.
func writeInputCache(day int, lines []string) error {
	if err := ensureDir("input"); err != nil {
		return err
	}

	path := filepath.Join("input", fmt.Sprintf("day%02d.txt", day))
	f, err := os.Create(path)
	if err != nil {
		return err
	}
	defer f.Close()

	for _, l := range lines {
		if _, err := f.WriteString(l + "\n"); err != nil {
			return err
		}
	}

	return nil
}
