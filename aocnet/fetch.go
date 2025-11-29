package aocnet

import (
	"fmt"
	"io"
	"net/http"
	"os"
)

const year = 2025

func FetchInput(day int, session string) ([]string, error) {
	url := fmt.Sprintf("https://adventofcode.com/%d/day/%d/input", year, day)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Cookie", "session="+session)
	req.Header.Set("User-Agent", fmt.Sprintf("github.com/%s/aoc%d (Go client)", getUsername(), year))

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("failed to fetch input: status %d", resp.StatusCode)
	}

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	lines := []string{}
	line := ""
	for _, b := range data {
		if b == '\n' {
			lines = append(lines, line)
			line = ""
			continue
		}
		line += string(b)
	}
	if line != "" {
		lines = append(lines, line)
	}

	return lines, nil
}

func getUsername() string {
	user := os.Getenv("USER")
	if user == "" {
		user = os.Getenv("USERNAME")
	}
	if user == "" {
		user = "anonymous"
	}
	return user
}
