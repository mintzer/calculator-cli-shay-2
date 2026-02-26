package main

import (
	"errors"
	"flag"
	"fmt"
	"os"
	"strconv"
	"strings"
)

const usageText = `Usage: calc <operation> <number1> <number2>

Operations:
  add  Add two numbers
  sub  Subtract the second number from the first
  mul  Multiply two numbers
  div  Divide the first number by the second
`

var errDivideByZero = errors.New("Cannot divide by zero")

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
		return 0, errDivideByZero
	}
	return a / b, nil
}

// formatResult formats a float64 to match Python's default float formatting,
// which always includes a decimal point (e.g., "8.0" rather than "8").
func formatResult(f float64) string {
	s := strconv.FormatFloat(f, 'f', -1, 64)
	if !strings.Contains(s, ".") {
		s += ".0"
	}
	return s
}

// run encapsulates the CLI logic. It returns stdout content, stderr content,
// and an exit code, making the core logic testable without calling os.Exit.
func run(args []string) (stdout string, stderr string, exitCode int) {
	fs := flag.NewFlagSet("calc", flag.ContinueOnError)

	// Suppress default flag error output; we handle messaging ourselves.
	var buf strings.Builder
	fs.SetOutput(&buf)
	fs.Usage = func() {
		fmt.Fprint(&buf, usageText)
	}

	if err := fs.Parse(args); err != nil {
		if err == flag.ErrHelp {
			return usageText, "", 0
		}
		return "", buf.String(), 1
	}

	posArgs := fs.Args()
	if len(posArgs) != 3 {
		return "", "Error: expected exactly 3 arguments: <operation> <number1> <number2>\n" + usageText, 1
	}

	operation := posArgs[0]
	aStr := posArgs[1]
	bStr := posArgs[2]

	a, err := strconv.ParseFloat(aStr, 64)
	if err != nil {
		return "", fmt.Sprintf("Error: invalid number '%s'\n", aStr), 1
	}

	b, err := strconv.ParseFloat(bStr, 64)
	if err != nil {
		return "", fmt.Sprintf("Error: invalid number '%s'\n", bStr), 1
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
		result, err = divide(a, b)
		if err != nil {
			return "", fmt.Sprintf("Error: %s\n", err.Error()), 1
		}
	default:
		return "", fmt.Sprintf("Error: unknown operation '%s'. Choose from: add, sub, mul, div\n", operation), 1
	}

	return formatResult(result) + "\n", "", 0
}

func main() {
	stdout, stderr, exitCode := run(os.Args[1:])
	if stdout != "" {
		fmt.Print(stdout)
	}
	if stderr != "" {
		fmt.Fprint(os.Stderr, stderr)
	}
	os.Exit(exitCode)
}
