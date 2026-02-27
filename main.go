package main

import (
	"errors"
	"fmt"
	"os"
	"strconv"
	"strings"
)

// Arithmetic functions — pure, no I/O.

func add(a, b float64) float64 {
	return a + b
}

func subtract(a, b float64) float64 {
	return a - b
}

func multiply(a, b float64) float64 {
	return a * b
}

func divide(a, b float64) (float64, error) {
	if b == 0 {
		return 0, errors.New("Cannot divide by zero")
	}
	return a / b, nil
}

// formatResult formats a float64 to match Python's default float string
// representation: integer-valued floats include a trailing ".0" (e.g. "8.0"),
// while non-integer values use the minimum number of digits (e.g. "8.7").
func formatResult(v float64) string {
	s := strconv.FormatFloat(v, 'f', -1, 64)
	if !strings.Contains(s, ".") {
		s += ".0"
	}
	return s
}

const usageLine = "usage: calc [-h] {add,sub,mul,div} a b"

const helpText = `usage: calc [-h] {add,sub,mul,div} a b

Simple CLI Calculator

positional arguments:
  {add,sub,mul,div}  Operation to perform
  a                  First number
  b                  Second number

options:
  -h, --help         show this help message and exit
`

var validOps = map[string]bool{"add": true, "sub": true, "mul": true, "div": true}

// run encapsulates the CLI logic and returns stdout, stderr content, and an
// exit code. This design keeps main() thin and makes integration testing
// straightforward.
func run(args []string) (stdout string, stderr string, exitCode int) {
	// Check for help flags anywhere in args (matching argparse behavior)
	for _, arg := range args {
		if arg == "-h" || arg == "--help" {
			return helpText, "", 0
		}
	}

	// No arguments at all
	if len(args) == 0 {
		return "", usageLine + "\ncalc: error: the following arguments are required: operation, a, b\n", 2
	}

	operation := args[0]
	remaining := args[1:]

	// Unknown flags in operation position (argparse treats these as missing positional args)
	if strings.HasPrefix(operation, "-") {
		return "", usageLine + "\ncalc: error: the following arguments are required: operation, a, b\n", 2
	}

	// Validate operation
	if !validOps[operation] {
		return "", fmt.Sprintf("%s\ncalc: error: argument operation: invalid choice: '%s' (choose from add, sub, mul, div)\n", usageLine, operation), 2
	}

	// Check operand count
	if len(remaining) == 0 {
		return "", usageLine + "\ncalc: error: the following arguments are required: a, b\n", 2
	}

	if len(remaining) == 1 {
		// Try to parse the first operand; if invalid, report that error first
		// (matching argparse behavior which validates type before checking missing args)
		_, err := strconv.ParseFloat(remaining[0], 64)
		if err != nil {
			return "", fmt.Sprintf("%s\ncalc: error: argument a: invalid float value: '%s'\n", usageLine, remaining[0]), 2
		}
		return "", usageLine + "\ncalc: error: the following arguments are required: b\n", 2
	}

	if len(remaining) > 2 {
		extra := strings.Join(remaining[2:], " ")
		return "", fmt.Sprintf("%s\ncalc: error: unrecognized arguments: %s\n", usageLine, extra), 2
	}

	// Parse operands
	a, err := strconv.ParseFloat(remaining[0], 64)
	if err != nil {
		return "", fmt.Sprintf("%s\ncalc: error: argument a: invalid float value: '%s'\n", usageLine, remaining[0]), 2
	}

	b, err := strconv.ParseFloat(remaining[1], 64)
	if err != nil {
		return "", fmt.Sprintf("%s\ncalc: error: argument b: invalid float value: '%s'\n", usageLine, remaining[1]), 2
	}

	// Execute operation
	var result float64

	switch operation {
	case "add":
		result = add(a, b)
	case "sub":
		result = subtract(a, b)
	case "mul":
		result = multiply(a, b)
	case "div":
		val, divErr := divide(a, b)
		if divErr != nil {
			return "", fmt.Sprintf("Error: %s\n", divErr.Error()), 1
		}
		result = val
	}

	return formatResult(result) + "\n", "", 0
}

func main() {
	stdout, stderr, exitCode := run(os.Args[1:])

	if stderr != "" {
		fmt.Fprint(os.Stderr, stderr)
	}
	if stdout != "" {
		fmt.Print(stdout)
	}

	os.Exit(exitCode)
}
