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

const usageText = `Usage: calc <operation> <number1> <number2>

Operations:
  add  Add two numbers
  sub  Subtract the second number from the first
  mul  Multiply two numbers
  div  Divide the first number by the second
`

// run encapsulates the CLI logic and returns stdout, stderr content, and an
// exit code. This design keeps main() thin and makes integration testing
// straightforward.
func run(args []string) (stdout string, stderr string, exitCode int) {
	if len(args) == 1 && (args[0] == "-h" || args[0] == "--help") {
		return usageText, "", 0
	}

	if len(args) != 3 {
		return "", "Error: expected 3 arguments: <operation> <number1> <number2>\n" + usageText, 1
	}

	operation := args[0]
	aStr := args[1]
	bStr := args[2]

	a, err := strconv.ParseFloat(aStr, 64)
	if err != nil {
		return "", fmt.Sprintf("Error: invalid number %q\n", aStr), 1
	}

	b, err := strconv.ParseFloat(bStr, 64)
	if err != nil {
		return "", fmt.Sprintf("Error: invalid number %q\n", bStr), 1
	}

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
	default:
		return "", fmt.Sprintf("Error: unknown operation %q. Choose from: add, sub, mul, div\n", operation), 1
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
