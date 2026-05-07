package main

import (
	"bufio"
	"os"
	"strconv"
	"strings"
)

type problemDescription struct {
	Title       string
	Description string
}

// LoadProblemDescriptions reads the small problems.yaml file and returns day
// descriptions keyed by day number. It supports the simple shape used in this
// repository: day numbers at the top level with title and description fields.
func LoadProblemDescriptions(path string) (map[int]problemDescription, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	descriptions := make(map[int]problemDescription)
	currentDay := 0

	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}

		if strings.HasSuffix(line, ":") {
			day, err := strconv.Atoi(strings.TrimSuffix(line, ":"))
			if err != nil {
				currentDay = 0
				continue
			}
			currentDay = day
			descriptions[currentDay] = problemDescription{}
			continue
		}

		if currentDay == 0 {
			continue
		}

		key, value, ok := strings.Cut(line, ":")
		if !ok {
			continue
		}

		problem := descriptions[currentDay]
		switch strings.TrimSpace(key) {
		case "title":
			problem.Title = yamlString(value)
		case "description":
			problem.Description = yamlString(value)
		}
		descriptions[currentDay] = problem
	}

	return descriptions, scanner.Err()
}

func yamlString(value string) string {
	value = strings.TrimSpace(value)
	value = strings.Trim(value, "\"")
	return value
}
