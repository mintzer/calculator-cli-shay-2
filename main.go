package main

import (
	"errors"
	"fmt"
	"os"
	"strconv"
	"strings"
)

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

var validOps = map[string]bool{
	"add": true,
	"sub": true,
	"mul": true,
	"div": true,
}

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

func formatResult(f float64) string {
	s := strconv.FormatFloat(f, 'f', -1, 64)
	if !strings.Contains(s, ".") {
		s += ".0"
	}
	return s
}

func argparseError(msg string) (string, string, int) {
	return "", usageLine + "\ncalc: error: " + msg + "\n", 2
}

func run(args []string) (stdout string, stderr string, exitCode int) {
	// Check for -h/--help anywhere in args
	for _, arg := range args {
		if arg == "-h" || arg == "--help" {
			return helpText, "", 0
		}
	}

	// Separate positional args from flag-like args.
	// Negative numbers (e.g. -1, -5.3) are treated as positional.
	// Unknown flags (e.g. --unknown, -x) are ignored, matching Python argparse
	// behavior where missing positional args are reported first.
	var positional []string
	for _, arg := range args {
		if len(arg) > 1 && arg[0] == '-' {
			// Check if this is a negative number (digit or dot after -)
			if arg[1] >= '0' && arg[1] <= '9' || arg[1] == '.' {
				positional = append(positional, arg)
			}
			// Otherwise it's a flag — skip it
		} else {
			positional = append(positional, arg)
		}
	}

	// Validate based on number of positional arguments, mimicking Python argparse order
	switch len(positional) {
	case 0:
		return argparseError("the following arguments are required: operation, a, b")
	case 1:
		if !validOps[positional[0]] {
			return argparseError(fmt.Sprintf("argument operation: invalid choice: '%s' (choose from add, sub, mul, div)", positional[0]))
		}
		return argparseError("the following arguments are required: a, b")
	case 2:
		if !validOps[positional[0]] {
			return argparseError(fmt.Sprintf("argument operation: invalid choice: '%s' (choose from add, sub, mul, div)", positional[0]))
		}
		_, err := strconv.ParseFloat(positional[1], 64)
		if err != nil {
			return argparseError(fmt.Sprintf("argument a: invalid float value: '%s'", positional[1]))
		}
		return argparseError("the following arguments are required: b")
	case 3:
		// Correct number — continue processing below
	default:
		extra := strings.Join(positional[3:], " ")
		return argparseError(fmt.Sprintf("unrecognized arguments: %s", extra))
	}

	operation := positional[0]
	aStr := positional[1]
	bStr := positional[2]

	if !validOps[operation] {
		return argparseError(fmt.Sprintf("argument operation: invalid choice: '%s' (choose from add, sub, mul, div)", operation))
	}

	a, err := strconv.ParseFloat(aStr, 64)
	if err != nil {
		return argparseError(fmt.Sprintf("argument a: invalid float value: '%s'", aStr))
	}

	b, err := strconv.ParseFloat(bStr, 64)
	if err != nil {
		return argparseError(fmt.Sprintf("argument b: invalid float value: '%s'", bStr))
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
