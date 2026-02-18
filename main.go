package main

import (
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"
)

const progName = "calc"

var validOps = []string{"add", "sub", "mul", "div"}

func add(a, b float64) float64 {
	return a + b
}

func sub(a, b float64) float64 {
	return a - b
}

func mul(a, b float64) float64 {
	return a * b
}

func div(a, b float64) (float64, error) {
	if b == 0 {
		return 0, fmt.Errorf("Cannot divide by zero")
	}
	return a / b, nil
}

// formatFloat formats a float64 to match Python's default str(float) behavior.
func formatFloat(f float64) string {
	if f == float64(int64(f)) && f >= -1e16 && f <= 1e16 {
		return strconv.FormatFloat(f, 'f', 1, 64)
	}
	return strconv.FormatFloat(f, 'g', -1, 64)
}

func usageLine() string {
	return fmt.Sprintf("usage: %s [-h] {add,sub,mul,div} a b", progName)
}

func printHelp(w io.Writer) {
	fmt.Fprintf(w, "%s\n", usageLine())
	fmt.Fprintln(w)
	fmt.Fprintln(w, "Simple CLI Calculator")
	fmt.Fprintln(w)
	fmt.Fprintln(w, "positional arguments:")
	fmt.Fprintln(w, "  {add,sub,mul,div}  Operation to perform")
	fmt.Fprintln(w, "  a                  First number")
	fmt.Fprintln(w, "  b                  Second number")
	fmt.Fprintln(w)
	fmt.Fprintln(w, "options:")
	fmt.Fprintln(w, "  -h, --help         show this help message and exit")
}

func argError(stderr io.Writer, msg string) int {
	fmt.Fprintln(stderr, usageLine())
	fmt.Fprintf(stderr, "%s: error: %s\n", progName, msg)
	return 2
}

func isValidOp(op string) bool {
	for _, v := range validOps {
		if op == v {
			return true
		}
	}
	return false
}

func run(args []string, stdout io.Writer, stderr io.Writer) int {
	// Separate flags from positional args
	var positional []string
	var extra []string
	helpRequested := false

	for _, arg := range args[1:] {
		if arg == "-h" || arg == "--help" {
			helpRequested = true
		} else {
			positional = append(positional, arg)
		}
	}

	if helpRequested {
		printHelp(stdout)
		return 0
	}

	// Check for missing required arguments
	switch len(positional) {
	case 0:
		return argError(stderr, "the following arguments are required: operation, a, b")
	case 1:
		// Have operation, missing a and b
		op := positional[0]
		if !isValidOp(op) {
			return argError(stderr, fmt.Sprintf("argument operation: invalid choice: '%s' (choose from add, sub, mul, div)", op))
		}
		return argError(stderr, "the following arguments are required: a, b")
	case 2:
		// Have operation and a, missing b
		op := positional[0]
		if !isValidOp(op) {
			return argError(stderr, fmt.Sprintf("argument operation: invalid choice: '%s' (choose from add, sub, mul, div)", op))
		}
		// Check if first operand is valid float
		_, err := strconv.ParseFloat(positional[1], 64)
		if err != nil {
			return argError(stderr, fmt.Sprintf("argument a: invalid float value: '%s'", positional[1]))
		}
		return argError(stderr, "the following arguments are required: b")
	case 3:
		// Correct number of positional args - process below
	default:
		// Too many arguments - collect extra ones
		extra = positional[3:]
		positional = positional[:3]
	}

	operation := positional[0]
	aStr := positional[1]
	bStr := positional[2]

	// Validate operation
	if !isValidOp(operation) {
		return argError(stderr, fmt.Sprintf("argument operation: invalid choice: '%s' (choose from add, sub, mul, div)", operation))
	}

	// Parse first operand
	a, err := strconv.ParseFloat(aStr, 64)
	if err != nil {
		return argError(stderr, fmt.Sprintf("argument a: invalid float value: '%s'", aStr))
	}

	// Parse second operand
	b, err := strconv.ParseFloat(bStr, 64)
	if err != nil {
		return argError(stderr, fmt.Sprintf("argument b: invalid float value: '%s'", bStr))
	}

	// Check for extra arguments after validation
	if len(extra) > 0 {
		return argError(stderr, fmt.Sprintf("unrecognized arguments: %s", strings.Join(extra, " ")))
	}

	// Execute operation
	var result float64
	switch operation {
	case "add":
		result = add(a, b)
	case "sub":
		result = sub(a, b)
	case "mul":
		result = mul(a, b)
	case "div":
		var divErr error
		result, divErr = div(a, b)
		if divErr != nil {
			fmt.Fprintf(stderr, "Error: %s\n", divErr.Error())
			return 1
		}
	}

	fmt.Fprintln(stdout, formatFloat(result))
	return 0
}

func main() {
	os.Exit(run(os.Args, os.Stdout, os.Stderr))
}
