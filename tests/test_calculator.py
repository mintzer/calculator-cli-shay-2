"""Smoke tests for the calculator CLI."""

import subprocess
import sys

from calculator import add, subtract, multiply, divide, main


def test_add():
    """Test addition returns correct result."""
    assert add(1.0, 2.0) == 3.0


def test_subtract():
    """Test subtraction returns correct result."""
    assert subtract(5.0, 3.0) == 2.0


def test_multiply():
    """Test multiplication returns correct result."""
    assert multiply(4.0, 3.0) == 12.0


def test_divide():
    """Test division returns correct result."""
    assert divide(10.0, 2.0) == 5.0


def test_divide_by_zero():
    """Test division by zero raises ValueError."""
    import pytest
    with pytest.raises(ValueError, match="Cannot divide by zero"):
        divide(1.0, 0.0)


def test_cli_add(capsys, monkeypatch):
    """Test CLI produces correct output for addition."""
    monkeypatch.setattr(sys, "argv", ["calc", "add", "1", "2"])
    exit_code = main()
    captured = capsys.readouterr()
    assert exit_code == 0
    assert captured.out.strip() == "3.0"


def test_cli_division_by_zero(capsys, monkeypatch):
    """Test CLI returns exit code 1 on division by zero."""
    monkeypatch.setattr(sys, "argv", ["calc", "div", "1", "0"])
    exit_code = main()
    captured = capsys.readouterr()
    assert exit_code == 1
    assert "Cannot divide by zero" in captured.err
