package main

import (
	"strings"
	"testing"
)

// --- Unit Tests for Arithmetic Functions ---

func TestAdd(t *testing.T) {
	tests := []struct {
		a, b, want float64
	}{
		{1, 1, 2},
		{-1, 1, 0},
		{0, 0, 0},
		{2.5, 3.5, 6},
		{-3, -7, -10},
	}
	for _, tc := range tests {
		got := add(tc.a, tc.b)
		if got != tc.want {
			t.Errorf("add(%g, %g) = %g, want %g", tc.a, tc.b, got, tc.want)
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
			t.Errorf("subtract(%g, %g) = %g, want %g", tc.a, tc.b, got, tc.want)
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
		{-4, -5, 20},
		{1.5, 2, 3},
	}
	for _, tc := range tests {
		got := multiply(tc.a, tc.b)
		if got != tc.want {
			t.Errorf("multiply(%g, %g) = %g, want %g", tc.a, tc.b, got, tc.want)
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
		{1, 0, 0, true},
		{0, 0, 0, true},
	}
	for _, tc := range tests {
		got, err := divide(tc.a, tc.b)
		if tc.wantErr {
			if err == nil {
				t.Errorf("divide(%g, %g) expected error, got nil", tc.a, tc.b)
			}
		} else {
			if err != nil {
				t.Errorf("divide(%g, %g) unexpected error: %v", tc.a, tc.b, err)
			}
			if got != tc.want {
				t.Errorf("divide(%g, %g) = %g, want %g", tc.a, tc.b, got, tc.want)
			}
		}
	}
}

// --- Integration Tests for CLI Behavior ---

func TestRunAdd(t *testing.T) {
	stdout, stderr, exitCode := run([]string{"add", "5", "3"})
	if exitCode != 0 {
		t.Errorf("expected exit code 0, got %d", exitCode)
	}
	if stderr != "" {
		t.Errorf("expected no stderr, got %q", stderr)
	}
	if strings.TrimSpace(stdout) != "8" {
		t.Errorf("expected stdout '8', got %q", strings.TrimSpace(stdout))
	}
}

func TestRunSub(t *testing.T) {
	stdout, stderr, exitCode := run([]string{"sub", "10", "4"})
	if exitCode != 0 {
		t.Errorf("expected exit code 0, got %d", exitCode)
	}
	if stderr != "" {
		t.Errorf("expected no stderr, got %q", stderr)
	}
	if strings.TrimSpace(stdout) != "6" {
		t.Errorf("expected stdout '6', got %q", strings.TrimSpace(stdout))
	}
}

func TestRunMul(t *testing.T) {
	stdout, stderr, exitCode := run([]string{"mul", "6", "7"})
	if exitCode != 0 {
		t.Errorf("expected exit code 0, got %d", exitCode)
	}
	if stderr != "" {
		t.Errorf("expected no stderr, got %q", stderr)
	}
	if strings.TrimSpace(stdout) != "42" {
		t.Errorf("expected stdout '42', got %q", strings.TrimSpace(stdout))
	}
}

func TestRunDiv(t *testing.T) {
	stdout, stderr, exitCode := run([]string{"div", "20", "4"})
	if exitCode != 0 {
		t.Errorf("expected exit code 0, got %d", exitCode)
	}
	if stderr != "" {
		t.Errorf("expected no stderr, got %q", stderr)
	}
	if strings.TrimSpace(stdout) != "5" {
		t.Errorf("expected stdout '5', got %q", strings.TrimSpace(stdout))
	}
}

func TestRunDivByZero(t *testing.T) {
	stdout, stderr, exitCode := run([]string{"div", "1", "0"})
	if exitCode != 1 {
		t.Errorf("expected exit code 1, got %d", exitCode)
	}
	if stdout != "" {
		t.Errorf("expected no stdout, got %q", stdout)
	}
	if !strings.Contains(stderr, "Error: Cannot divide by zero") {
		t.Errorf("expected stderr to contain 'Error: Cannot divide by zero', got %q", stderr)
	}
}

func TestRunInvalidOperation(t *testing.T) {
	_, stderr, exitCode := run([]string{"foo", "1", "2"})
	if exitCode != 1 {
		t.Errorf("expected exit code 1, got %d", exitCode)
	}
	if !strings.Contains(stderr, "unknown operation") {
		t.Errorf("expected stderr to contain 'unknown operation', got %q", stderr)
	}
}

func TestRunWrongArgCount(t *testing.T) {
	_, stderr, exitCode := run([]string{"add", "1"})
	if exitCode != 1 {
		t.Errorf("expected exit code 1, got %d", exitCode)
	}
	if stderr == "" {
		t.Error("expected stderr output for wrong arg count")
	}
}

func TestRunNoArgs(t *testing.T) {
	_, stderr, exitCode := run([]string{})
	if exitCode != 1 {
		t.Errorf("expected exit code 1, got %d", exitCode)
	}
	if stderr == "" {
		t.Error("expected stderr output for no args")
	}
}

func TestRunInvalidNumber(t *testing.T) {
	_, stderr, exitCode := run([]string{"add", "abc", "3"})
	if exitCode != 1 {
		t.Errorf("expected exit code 1, got %d", exitCode)
	}
	if !strings.Contains(stderr, "invalid number") {
		t.Errorf("expected stderr to contain 'invalid number', got %q", stderr)
	}
}

func TestRunFloatOutput(t *testing.T) {
	stdout, _, exitCode := run([]string{"div", "1", "3"})
	if exitCode != 0 {
		t.Errorf("expected exit code 0, got %d", exitCode)
	}
	if strings.TrimSpace(stdout) != "0.3333333333333333" {
		t.Errorf("expected stdout '0.3333333333333333', got %q", strings.TrimSpace(stdout))
	}
}
