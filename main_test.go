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
	cases := []struct {
		a, b, want float64
	}{
		{1, 1, 2},
		{-1, 1, 0},
		{0, 0, 0},
		{-3, -7, -10},
		{1.5, 2.5, 4},
		{100, 200, 300},
	}
	for _, tc := range cases {
		got := add(tc.a, tc.b)
		if got != tc.want {
			t.Errorf("add(%v, %v) = %v; want %v", tc.a, tc.b, got, tc.want)
		}
	}
}

func TestSubtract(t *testing.T) {
	cases := []struct {
		a, b, want float64
	}{
		{5, 3, 2},
		{3, 5, -2},
		{0, 0, 0},
		{-1, -1, 0},
		{10.5, 0.5, 10},
	}
	for _, tc := range cases {
		got := subtract(tc.a, tc.b)
		if got != tc.want {
			t.Errorf("subtract(%v, %v) = %v; want %v", tc.a, tc.b, got, tc.want)
		}
	}
}

func TestMultiply(t *testing.T) {
	cases := []struct {
		a, b, want float64
	}{
		{2, 3, 6},
		{-2, 3, -6},
		{-2, -3, 6},
		{0, 100, 0},
		{1.5, 4, 6},
	}
	for _, tc := range cases {
		got := multiply(tc.a, tc.b)
		if got != tc.want {
			t.Errorf("multiply(%v, %v) = %v; want %v", tc.a, tc.b, got, tc.want)
		}
	}
}

func TestDivide(t *testing.T) {
	t.Run("normal division", func(t *testing.T) {
		cases := []struct {
			a, b, want float64
		}{
			{20, 4, 5},
			{9, 3, 3},
			{7, 2, 3.5},
			{-6, 3, -2},
			{0, 5, 0},
		}
		for _, tc := range cases {
			got, err := divide(tc.a, tc.b)
			if err != nil {
				t.Errorf("divide(%v, %v) unexpected error: %v", tc.a, tc.b, err)
				continue
			}
			if math.Abs(got-tc.want) > 1e-9 {
				t.Errorf("divide(%v, %v) = %v; want %v", tc.a, tc.b, got, tc.want)
			}
		}
	})

	t.Run("division by zero", func(t *testing.T) {
		_, err := divide(1, 0)
		if err == nil {
			t.Fatal("divide(1, 0) expected error, got nil")
		}
		if err.Error() != "Cannot divide by zero" {
			t.Errorf("divide(1, 0) error = %q; want %q", err.Error(), "Cannot divide by zero")
		}
	})
}

// ---------------------------------------------------------------------------
// Unit test for formatResult
// ---------------------------------------------------------------------------

func TestFormatResult(t *testing.T) {
	cases := []struct {
		val  float64
		want string
	}{
		{8, "8.0"},
		{-2, "-2.0"},
		{0, "0.0"},
		{8.7, "8.7"},
		{1.5, "1.5"},
		{100, "100.0"},
	}
	for _, tc := range cases {
		got := formatResult(tc.val)
		if got != tc.want {
			t.Errorf("formatResult(%v) = %q; want %q", tc.val, got, tc.want)
		}
	}
}

// ---------------------------------------------------------------------------
// Integration tests using the run() function
// ---------------------------------------------------------------------------

func TestRunSuccess(t *testing.T) {
	cases := []struct {
		name       string
		args       []string
		wantStdout string
	}{
		{"add integers", []string{"add", "5", "3"}, "8.0\n"},
		{"sub integers", []string{"sub", "10", "4"}, "6.0\n"},
		{"mul integers", []string{"mul", "3", "7"}, "21.0\n"},
		{"div integers", []string{"div", "20", "4"}, "5.0\n"},
		{"add floats", []string{"add", "1.5", "2.5"}, "4.0\n"},
		{"div non-integer result", []string{"div", "7", "2"}, "3.5\n"},
		{"sub negative result", []string{"sub", "3", "5"}, "-2.0\n"},
		{"mul by zero", []string{"mul", "42", "0"}, "0.0\n"},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			stdout, stderr, exitCode := run(tc.args)
			if exitCode != 0 {
				t.Errorf("run(%v) exitCode = %d; want 0 (stderr: %q)", tc.args, exitCode, stderr)
			}
			if stdout != tc.wantStdout {
				t.Errorf("run(%v) stdout = %q; want %q", tc.args, stdout, tc.wantStdout)
			}
			if stderr != "" {
				t.Errorf("run(%v) stderr = %q; want empty", tc.args, stderr)
			}
		})
	}
}

func TestRunDivisionByZero(t *testing.T) {
	stdout, stderr, exitCode := run([]string{"div", "1", "0"})
	if exitCode != 1 {
		t.Errorf("exitCode = %d; want 1", exitCode)
	}
	if !strings.Contains(stderr, "Cannot divide by zero") {
		t.Errorf("stderr = %q; want it to contain %q", stderr, "Cannot divide by zero")
	}
	if stdout != "" {
		t.Errorf("stdout = %q; want empty", stdout)
	}
}

func TestRunInvalidOperation(t *testing.T) {
	stdout, stderr, exitCode := run([]string{"foo", "1", "2"})
	if exitCode != 2 {
		t.Errorf("exitCode = %d; want 2", exitCode)
	}
	if !strings.Contains(stderr, "invalid choice") {
		t.Errorf("stderr = %q; want it to contain %q", stderr, "invalid choice")
	}
	if stdout != "" {
		t.Errorf("stdout = %q; want empty", stdout)
	}
}

func TestRunWrongArgCount(t *testing.T) {
	cases := []struct {
		name string
		args []string
	}{
		{"too few args", []string{"add", "1"}},
		{"no args", []string{}},
		{"too many args", []string{"add", "1", "2", "3"}},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			stdout, stderr, exitCode := run(tc.args)
			if exitCode != 2 {
				t.Errorf("run(%v) exitCode = %d; want 2", tc.args, exitCode)
			}
			if stderr == "" {
				t.Errorf("run(%v) stderr is empty; expected an error message", tc.args)
			}
			if stdout != "" {
				t.Errorf("run(%v) stdout = %q; want empty", tc.args, stdout)
			}
		})
	}
}

func TestRunNonNumericOperands(t *testing.T) {
	cases := []struct {
		name string
		args []string
	}{
		{"first arg non-numeric", []string{"add", "hello", "2"}},
		{"second arg non-numeric", []string{"add", "1", "world"}},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			stdout, stderr, exitCode := run(tc.args)
			if exitCode != 2 {
				t.Errorf("run(%v) exitCode = %d; want 2", tc.args, exitCode)
			}
			if !strings.Contains(stderr, "invalid float value") {
				t.Errorf("run(%v) stderr = %q; want it to contain %q", tc.args, stderr, "invalid float value")
			}
			if stdout != "" {
				t.Errorf("run(%v) stdout = %q; want empty", tc.args, stdout)
			}
		})
	}
}

func TestRunHelp(t *testing.T) {
	for _, flag := range []string{"-h", "--help"} {
		t.Run(flag, func(t *testing.T) {
			stdout, stderr, exitCode := run([]string{flag})
			if exitCode != 0 {
				t.Errorf("run(%q) exitCode = %d; want 0", flag, exitCode)
			}
			if !strings.Contains(stdout, "usage:") {
				t.Errorf("run(%q) stdout = %q; want it to contain usage text", flag, stdout)
			}
			if stderr != "" {
				t.Errorf("run(%q) stderr = %q; want empty", flag, stderr)
			}
		})
	}
}
