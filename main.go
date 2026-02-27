package main

import (
	"errors"
	"flag"
	"fmt"
	"os"
	"strconv"
)

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

func run(args []string) (stdout string, stderr string, exitCode int) {
	fs := flag.NewFlagSet("calc", flag.ContinueOnError)
	fs.Usage = func() {
		fmt.Fprintln(fs.Output(), "Usage: calc <operation> <number1> <number2>")
		fmt.Fprintln(fs.Output())
		fmt.Fprintln(fs.Output(), "Operations:")
		fmt.Fprintln(fs.Output(), "  add  Add two numbers")
		fmt.Fprintln(fs.Output(), "  sub  Subtract the second number from the first")
		fmt.Fprintln(fs.Output(), "  mul  Multiply two numbers")
		fmt.Fprintln(fs.Output(), "  div  Divide the first number by the second")
	}

	if err := fs.Parse(args); err != nil {
		return "", "", 0
	}

	posArgs := fs.Args()
	if len(posArgs) != 3 {
		errMsg := "Error: exactly 3 arguments required: <operation> <number1> <number2>\n"
		fs.SetOutput(os.Stderr)
		fs.Usage()
		return "", errMsg, 1
	}

	operation := posArgs[0]
	a, err := strconv.ParseFloat(posArgs[1], 64)
	if err != nil {
		return "", fmt.Sprintf("Error: invalid number '%s'\n", posArgs[1]), 1
	}
	b, err := strconv.ParseFloat(posArgs[2], 64)
	if err != nil {
		return "", fmt.Sprintf("Error: invalid number '%s'\n", posArgs[2]), 1
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
			return "", fmt.Sprintf("Error: %s\n", divErr), 1
		}
		result = val
	default:
		return "", fmt.Sprintf("Error: unknown operation '%s'\n", operation), 1
	}

	return fmt.Sprintf("%g\n", result), "", 0
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
