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

Run a single day:

    ./aoc2025 1

or

    go run . 1

Run multiple days:

    ./aoc2025 1 4 5

The CLI accepts any number of days between 1–12.

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

## Benchmarks

First solve all days so input is stored locally, then run like ths

    cd days
    go test -bench=.

### Benchmark Summary — Apple M4 (darwin/arm64)

| Day | SetInput (µs) | SolvePart1 (µs) | SolvePart2 (µs) | FullPipeline (µs) |
| --- | ------------- | --------------- | --------------- | ----------------- |
| 01  | 71            | 12              | 417             | 499               |
| 02  | 6             | <1              | 8               | 15                |
| 03  | 32            | 11              | 40              | 79                |
| 04  | 9             | 119             | 268             | 520               |
