# Advent of Code 2025

This repository contains a clean, extensible Go framework for solving Advent of Code 2025 using the following interface:

```go
type Aoc2025 interface {
    SetInput(lines []string)
    SolvePart1() string
    SolvePart2() string
}
```

Each day’s solution implements this interface and automatically registers itself.

## 📦 Project Structure

```
aoc2025/
│
├── README.md
├── go.mod
├── main.go
├── problems.yaml       # brief day titles and descriptions for verbose output
│
├── input/              # cached input files (auto-created)
│     └── (empty until downloaded)
│
├── aocnet/
│     └── fetch.go      # handles online input downloading
│
└── days/
      ├── interface.go
      ├── registry.go
      ├── day01.go
      └── ... up to day12.go
```

## 🚀 Running

By default, the CLI runs in brief mode: it prints only the day header and the two answers.

Run a single day:

    ./aoc2025 1

or

    go run . 1

Run multiple days:

    ./aoc2025 1 4 5

The CLI accepts any number of days between 1–12.

Show the brief problem description before solving each selected day:

    ./aoc2025 --verbose 1

or:

    ./aoc2025 -v 1 4 5

Descriptions are read from `problems.yaml`. If that file is unavailable, solving still works; the CLI prints a warning and continues without descriptions.

## 🌐 Automatic Input Download From adventofcode.com

This framework supports automatic downloading of puzzle input using your personal Advent of Code session cookie.

Advent of Code does not provide `OAuth` or API login.
Even if you sign in with Google, AoC internally uses a cookie named:

    session=YOUR_SESSION_TOKEN

If you include this cookie in HTTP requests, you can fetch your input programmatically.

## 🔑 How to Retrieve Your Session Token

1.	Log in to: https://adventofcode.com/
2.	Open your browser’s Developer Tools
    *   Safari: ⌥ Option + ⌘ Command + I
    *	Chrome: F12 → Application tab
    *	Firefox: F12 → Storage tab
3.	Go to Cookies → https://adventofcode.com
4.	Look for a cookie named: session
5.	Copy the value (a long hex-like string)

⚠️ This token is your personal authentication.
DO NOT commit it to Git or share it.

## 🧷 Enabling Automatic Download

Set two environment variables:

    export AOC_SESSION="your-session-token"
    export AOC_ONLINE=1

Now when you run:

    ./aoc2025 1

the program will:
1.	Attempt to download https://adventofcode.com/2025/day/1/input
2.	Save it to
`input/day01.txt`
3.	Use the downloaded input for solving

If downloading fails, it falls back to reading the file from disk.

## ⏱️ Benchmarks

First solve all days so input is stored locally, then run like ths

    cd days
    go test -bench=.

### Benchmark Summary — Apple M4 (darwin/arm64)

| Day | SetInput (µs) | SolvePart1 (µs) | SolvePart2 (µs) | FullPipeline (µs) |
| --- | ------------- | --------------- | --------------- | ----------------- |
| 01  | 76.44         | 12.29           | 419.30          | 508.37            |
| 02  | 9.37          | 0.62            | 7.49            | 17.76             |
| 03  | 34.80         | 23.89           | 39.50           | 98.72             |
| 04  | 11.69         | 144.72          | 368.18          | 691.84            |
| 05  | 57.36         | 5.02            | 0.05            | 62.48             |
| 06  | 10.55         | 65.64           | 90.15           | 168.72            |
| 07  | 12.27         | 10.01           | 10.35           | 32.51             |
| 08  | 6_333.86      | 8.28            | 14.64           | 6_378.43          |
| 09  | 27.92         | 96.99           | 2_948.83        | 3_053.64          |
| 10  | 153.21        | 97.51           | 2_507.93        | 2_899.36          |
| 11  | 73.23         | 9.28            | 151.00          | 237.73            |
| 12  | 150.41        | 4.94            | -               | 162.48            |
