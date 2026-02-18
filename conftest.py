import os
import subprocess
import enable_type_union

enable_type_union.install()

# The test file contains a literal "{{ cli_command }}" template that was not replaced.
# Patch subprocess.run to replace it with the actual CLI command.
_CLI_COMMAND = os.environ.get("BASE_CLI_COMMAND", "./calc-go")
_original_run = subprocess.run
_original_popen = subprocess.Popen


def _patched_run(*args, **kwargs):
    args_list = list(args)
    if args_list and isinstance(args_list[0], (list, tuple)):
        cmd = list(args_list[0])
        if cmd and cmd[0] == "{{ cli_command }}":
            cmd[0] = _CLI_COMMAND
            args_list[0] = cmd
    elif args_list and isinstance(args_list[0], str) and "{{ cli_command }}" in args_list[0]:
        args_list[0] = args_list[0].replace("{{ cli_command }}", _CLI_COMMAND)
    return _original_run(*args_list, **kwargs)


class _PatchedPopen(_original_popen):
    def __init__(self, args, *pargs, **kwargs):
        if isinstance(args, (list, tuple)):
            args = list(args)
            if args and args[0] == "{{ cli_command }}":
                args[0] = _CLI_COMMAND
        elif isinstance(args, str) and "{{ cli_command }}" in args:
            args = args.replace("{{ cli_command }}", _CLI_COMMAND)
        super().__init__(args, *pargs, **kwargs)


subprocess.run = _patched_run
subprocess.Popen = _PatchedPopen
