// Copyright 2026 TypeFox GmbH
// This program and the accompanying materials are made available under the
// terms of the MIT License, which is available in the project root.

package main

import (
	"bytes"
	"flag"
	"fmt"
	"math"
	"os"
	"os/exec"
	"regexp"
	"strconv"
)

var msPerResourcePattern = regexp.MustCompile(`([0-9]+(?:\.[0-9]+)?)\s+ms/resource`)
var mbPerSecondPattern = regexp.MustCompile(`([0-9]+(?:\.[0-9]+)?)\s+MB/s`)

const separator = "----------------------------------------------------------------"
const (
	defaultRuns  = 10
	defaultBench = "^BenchmarkWorkspaceCycle$"
	defaultPkg   = "./examples/statemachine"
)

// main runs a benchmark multiple times and reports min/avg/max for selected metrics.
// Repeating the benchmark helps reduce noise from unrelated system load while keeping the
// metrics stable for Fastbelt's multithreaded execution.
//
// The default benchmark is the workspace cycle from the statemachine example.
func main() {
	var runs int
	var bench string
	var pkg string
	var cpu int
	flag.IntVar(&runs, "runs", defaultRuns, "number of benchmark runs")
	flag.IntVar(&cpu, "cpu", 0, "CPU limit passed to go test -cpu (0 means no limit)")
	flag.StringVar(&bench, "bench", defaultBench, "benchmark regex passed to go test -bench")
	flag.StringVar(&pkg, "pkg", defaultPkg, "package passed to go test")
	flag.Parse()

	if runs < 1 {
		fatalf("invalid -runs=%d (must be >= 1)", runs)
	}
	if cpu < 0 {
		fatalf("invalid -cpu=%d (must be >= 0)", cpu)
	}

	msStats := newStats()
	mbStats := newStats()

	fmt.Printf("Running benchmark %d times...\n", runs)
	for i := 1; i <= runs; i++ {
		fmt.Printf("\n%s\n", separator)
		fmt.Printf("Run %d/%d\n", i, runs)
		fmt.Printf("%s\n", separator)
		output, err := runBenchmark(bench, pkg, cpu)
		if err != nil {
			fatalf("benchmark run %d failed: %v\n%s", i, err, output)
		}
		fmt.Print(output)

		msPerResource, mbPerSecond, err := parseBenchmarkMetrics(output)
		if err != nil {
			fatalf("benchmark run %d: %v", i, err)
		}
		msStats.observe(msPerResource)
		mbStats.observe(mbPerSecond)
	}

	fmt.Printf("\n%s\n", separator)
	fmt.Printf("Final result over %d runs\n", runs)
	fmt.Printf("%s\n", separator)
	fmt.Printf("MB/s:        min %.3f, avg %.3f, max %.3f\n", mbStats.min, mbStats.avg(), mbStats.max)
	fmt.Printf("ms/resource: min %.4f, avg %.4f, max %.4f\n", msStats.min, msStats.avg(), msStats.max)
}

func runBenchmark(bench, pkg string, cpu int) (string, error) {
	args := []string{"test", "-run", "^$", "-bench", bench, "-count=1"}
	if cpu > 0 {
		args = append(args, "-cpu", strconv.Itoa(cpu))
	}
	args = append(args, pkg)
	cmd := exec.Command("go", args...)
	var buffer bytes.Buffer
	cmd.Stdout = &buffer
	cmd.Stderr = &buffer
	err := cmd.Run()
	return buffer.String(), err
}

func parseBenchmarkMetrics(output string) (float64, float64, error) {
	msPerResourceMatch := msPerResourcePattern.FindStringSubmatch(output)
	if len(msPerResourceMatch) != 2 {
		return 0, 0, fmt.Errorf("could not find ms/resource in benchmark output")
	}
	msPerResource, err := strconv.ParseFloat(msPerResourceMatch[1], 64)
	if err != nil {
		return 0, 0, fmt.Errorf("invalid ms/resource value %q: %w", msPerResourceMatch[1], err)
	}

	mbPerSecondMatch := mbPerSecondPattern.FindStringSubmatch(output)
	if len(mbPerSecondMatch) != 2 {
		return 0, 0, fmt.Errorf("could not find MB/s in benchmark output")
	}
	mbPerSecond, err := strconv.ParseFloat(mbPerSecondMatch[1], 64)
	if err != nil {
		return 0, 0, fmt.Errorf("invalid MB/s value %q: %w", mbPerSecondMatch[1], err)
	}

	return msPerResource, mbPerSecond, nil
}

type stats struct {
	min   float64
	max   float64
	sum   float64
	count int
}

func newStats() stats {
	return stats{
		min: math.Inf(1),
		max: math.Inf(-1),
	}
}

func (s *stats) observe(value float64) {
	if value < s.min {
		s.min = value
	}
	if value > s.max {
		s.max = value
	}
	s.sum += value
	s.count++
}

func (s stats) avg() float64 {
	if s.count == 0 {
		return 0
	}
	return s.sum / float64(s.count)
}

func fatalf(format string, args ...any) {
	fmt.Fprintf(os.Stderr, format+"\n", args...)
	os.Exit(1)
}
