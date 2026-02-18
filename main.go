package main

import (
	"fmt"
	"io"
	"os"
	"strconv"
)

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
// Whole numbers are printed as "8.0", non-whole numbers use minimal representation.
func formatFloat(f float64) string {
	if f == float64(int64(f)) && f >= -1e15 && f <= 1e15 {
		return strconv.FormatFloat(f, 'f', 1, 64)
	}
	return strconv.FormatFloat(f, 'g', -1, 64)
}

func run(args []string, stdout io.Writer, stderr io.Writer) int {
	if len(args) != 4 {
		fmt.Fprintf(stderr, "Usage: %s <operation> <a> <b>\n", args[0])
		return 1
	}

	operation := args[1]
	aStr := args[2]
	bStr := args[3]

	a, err := strconv.ParseFloat(aStr, 64)
	if err != nil {
		fmt.Fprintf(stderr, "Error: invalid number '%s'\n", aStr)
		return 1
	}

	b, err := strconv.ParseFloat(bStr, 64)
	if err != nil {
		fmt.Fprintf(stderr, "Error: invalid number '%s'\n", bStr)
		return 1
	}

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
	default:
		fmt.Fprintf(stderr, "Error: unknown operation '%s'\n", operation)
		return 1
	}

	fmt.Fprintln(stdout, formatFloat(result))
	return 0
}

func main() {
	os.Exit(run(os.Args, os.Stdout, os.Stderr))
}
