"""
Conftest to fix the {{ cli_command }} template substitution bug in the auto-generated test file.

The auto-generated test script has a Jinja2 template '{{ cli_command }}' that fails to render,
causing all tests to fail with FileNotFoundError. This conftest patches subprocess.run to
replace the unrendered template with the actual CLI command.
"""

import os
import subprocess
import sys

# Add .venv/bin to PATH so 'calc' can be found
_venv_bin = os.path.join(os.getcwd(), ".venv", "bin")
if _venv_bin not in os.environ.get("PATH", ""):
    os.environ["PATH"] = _venv_bin + os.pathsep + os.environ.get("PATH", "")

_original_run = subprocess.run


def _patched_run(args, **kwargs):
    """Patch subprocess.run to fix the unrendered {{ cli_command }} template."""
    if isinstance(args, list) and len(args) > 0 and args[0] == "{{ cli_command }}":
        # Try to get CLI_COMMAND from the test module
        cli_cmd = None
        for mod_name, mod in sys.modules.items():
            if "test_protocol_spec" in mod_name:
                cli_cmd = getattr(mod, "CLI_COMMAND", None)
                if cli_cmd:
                    break

        if not cli_cmd:
            # Fallback: try to find calc in .venv/bin
            venv_calc = os.path.join(os.getcwd(), ".venv", "bin", "calc")
            if os.path.isfile(venv_calc):
                cli_cmd = venv_calc
            else:
                cli_cmd = "calc"

        args = [cli_cmd] + args[1:]

    return _original_run(args, **kwargs)


# Monkeypatch subprocess.run globally
subprocess.run = _patched_run
