package main

import (
	"math"
	"strings"
	"testing"
)

// ---------------------------------------------------------------------------
// Unit tests for arithmetic functions
// ---------------------------------------------------------------------------

func TestAdd(t *testing.T) {
	tests := []struct {
		a, b, want float64
	}{
		{1, 1, 2},
		{-1, 1, 0},
		{0, 0, 0},
		{2.5, 3.5, 6},
		{-2.5, -3.5, -6},
		{1e10, 1e10, 2e10},
	}
	for _, tc := range tests {
		got := add(tc.a, tc.b)
		if got != tc.want {
			t.Errorf("add(%v, %v) = %v, want %v", tc.a, tc.b, got, tc.want)
		}
	}
}

func TestSubtract(t *testing.T) {
	tests := []struct {
		a, b, want float64
	}{
		{5, 3, 2},
		{3, 5, -2},
		{0, 0, 0},
		{-1, -1, 0},
		{10.5, 0.5, 10},
	}
	for _, tc := range tests {
		got := subtract(tc.a, tc.b)
		if got != tc.want {
			t.Errorf("subtract(%v, %v) = %v, want %v", tc.a, tc.b, got, tc.want)
		}
	}
}

func TestMultiply(t *testing.T) {
	tests := []struct {
		a, b, want float64
	}{
		{2, 3, 6},
		{-2, 3, -6},
		{0, 100, 0},
		{0.5, 4, 2},
		{-3, -4, 12},
	}
	for _, tc := range tests {
		got := multiply(tc.a, tc.b)
		if got != tc.want {
			t.Errorf("multiply(%v, %v) = %v, want %v", tc.a, tc.b, got, tc.want)
		}
	}
}

func TestDivide(t *testing.T) {
	tests := []struct {
		a, b    float64
		want    float64
		wantErr bool
	}{
		{20, 4, 5, false},
		{1, 3, 1.0 / 3.0, false},
		{-10, 2, -5, false},
		{0, 5, 0, false},
		{7, 2, 3.5, false},
		{1, 0, 0, true},
		{0, 0, 0, true},
	}
	for _, tc := range tests {
		got, err := divide(tc.a, tc.b)
		if tc.wantErr {
			if err == nil {
				t.Errorf("divide(%v, %v) expected error, got nil", tc.a, tc.b)
			}
		} else {
			if err != nil {
				t.Errorf("divide(%v, %v) unexpected error: %v", tc.a, tc.b, err)
			}
			if math.Abs(got-tc.want) > 1e-15 {
				t.Errorf("divide(%v, %v) = %v, want %v", tc.a, tc.b, got, tc.want)
			}
		}
	}
}

// ---------------------------------------------------------------------------
// Unit test for formatResult
// ---------------------------------------------------------------------------

func TestFormatResult(t *testing.T) {
	tests := []struct {
		input float64
		want  string
	}{
		{8, "8.0"},
		{0, "0.0"},
		{-3, "-3.0"},
		{2.5, "2.5"},
		{0.3333333333333333, "0.3333333333333333"},
		{100, "100.0"},
		{-0.5, "-0.5"},
	}
	for _, tc := range tests {
		got := formatResult(tc.input)
		if got != tc.want {
			t.Errorf("formatResult(%v) = %q, want %q", tc.input, got, tc.want)
		}
	}
}

// ---------------------------------------------------------------------------
// Integration tests via run()
// ---------------------------------------------------------------------------

func TestRunSuccess(t *testing.T) {
	tests := []struct {
		name       string
		args       []string
		wantStdout string
	}{
		{"add integers", []string{"add", "5", "3"}, "8.0\n"},
		{"sub integers", []string{"sub", "10", "4"}, "6.0\n"},
		{"mul integers", []string{"mul", "3", "7"}, "21.0\n"},
		{"div integers", []string{"div", "20", "4"}, "5.0\n"},
		{"div with remainder", []string{"div", "1", "3"}, "0.3333333333333333\n"},
		{"add floats", []string{"add", "1.5", "2.5"}, "4.0\n"},
		{"sub negative result", []string{"sub", "3", "5"}, "-2.0\n"},
		{"mul by zero", []string{"mul", "42", "0"}, "0.0\n"},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			stdout, stderr, exitCode := run(tc.args)
			if exitCode != 0 {
				t.Errorf("run(%v) exit code = %d, want 0; stderr = %q", tc.args, exitCode, stderr)
			}
			if stdout != tc.wantStdout {
				t.Errorf("run(%v) stdout = %q, want %q", tc.args, stdout, tc.wantStdout)
			}
			if stderr != "" {
				t.Errorf("run(%v) stderr = %q, want empty", tc.args, stderr)
			}
		})
	}
}

func TestRunDivideByZero(t *testing.T) {
	stdout, stderr, exitCode := run([]string{"div", "1", "0"})
	if exitCode != 1 {
		t.Errorf("exit code = %d, want 1", exitCode)
	}
	if !strings.Contains(stderr, "Error: Cannot divide by zero") {
		t.Errorf("stderr = %q, want it to contain 'Error: Cannot divide by zero'", stderr)
	}
	if stdout != "" {
		t.Errorf("stdout = %q, want empty", stdout)
	}
}

func TestRunInvalidOperation(t *testing.T) {
	stdout, stderr, exitCode := run([]string{"foo", "1", "2"})
	if exitCode != 2 {
		t.Errorf("exit code = %d, want 2", exitCode)
	}
	if !strings.Contains(stderr, "invalid choice: 'foo'") {
		t.Errorf("stderr = %q, want it to contain \"invalid choice: 'foo'\"", stderr)
	}
	if stdout != "" {
		t.Errorf("stdout = %q, want empty", stdout)
	}
}

func TestRunWrongArgCount(t *testing.T) {
	tests := []struct {
		name        string
		args        []string
		wantContain string
	}{
		{"no args", []string{}, "the following arguments are required: operation, a, b"},
		{"one arg", []string{"add"}, "the following arguments are required: a, b"},
		{"two args", []string{"add", "1"}, "the following arguments are required: b"},
		{"four args", []string{"add", "1", "2", "3"}, "unrecognized arguments: 3"},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			stdout, stderr, exitCode := run(tc.args)
			if exitCode != 2 {
				t.Errorf("run(%v) exit code = %d, want 2", tc.args, exitCode)
			}
			if !strings.Contains(stderr, tc.wantContain) {
				t.Errorf("run(%v) stderr = %q, want it to contain %q", tc.args, stderr, tc.wantContain)
			}
			if stdout != "" {
				t.Errorf("run(%v) stdout = %q, want empty", tc.args, stdout)
			}
		})
	}
}

func TestRunNonNumericOperands(t *testing.T) {
	tests := []struct {
		name string
		args []string
		bad  string
	}{
		{"first operand", []string{"add", "abc", "2"}, "argument a: invalid float value: 'abc'"},
		{"second operand", []string{"add", "2", "xyz"}, "argument b: invalid float value: 'xyz'"},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			stdout, stderr, exitCode := run(tc.args)
			if exitCode != 2 {
				t.Errorf("run(%v) exit code = %d, want 2", tc.args, exitCode)
			}
			if !strings.Contains(stderr, tc.bad) {
				t.Errorf("run(%v) stderr = %q, want it to contain %q", tc.args, stderr, tc.bad)
			}
			if stdout != "" {
				t.Errorf("run(%v) stdout = %q, want empty", tc.args, stdout)
			}
		})
	}
}

func TestRunHelp(t *testing.T) {
	stdout, stderr, exitCode := run([]string{"-h"})
	if exitCode != 0 {
		t.Errorf("exit code = %d, want 0", exitCode)
	}
	if !strings.Contains(stdout, "usage: calc") {
		t.Errorf("stdout = %q, want it to contain 'usage: calc'", stdout)
	}
	if !strings.Contains(stdout, "add") || !strings.Contains(stdout, "div") {
		t.Errorf("stdout = %q, want it to list operations", stdout)
	}
	if stderr != "" {
		t.Errorf("stderr = %q, want empty", stderr)
	}
}
