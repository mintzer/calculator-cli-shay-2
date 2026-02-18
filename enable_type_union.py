"""
Import hook to handle PEP 604 union type syntax (X | Y) on Python 3.9.
Injects 'from __future__ import annotations' into source files to defer
annotation evaluation, avoiding TypeError on 'str | None' etc.
"""
import importlib.machinery
import importlib.util
import sys


class _FutureAnnotationsLoader:
    """Loader that injects 'from __future__ import annotations' into source."""

    def __init__(self, origin, source_bytes):
        self._origin = origin
        self._source_bytes = source_bytes

    def create_module(self, spec):
        return None

    def exec_module(self, module):
        source = self._source_bytes.decode('utf-8', errors='replace')
        # Inject future annotations after any encoding declarations / shebang
        # but it must be the first statement (after docstrings)
        source = "from __future__ import annotations\n" + source
        code = compile(source, self._origin, 'exec')
        exec(code, module.__dict__)


class _FutureAnnotationsFinder:
    """Meta path finder that injects future annotations for specific modules."""

    def __init__(self):
        self._active = False

    def find_spec(self, fullname, path, target=None):
        if self._active:
            return None

        self._active = True
        try:
            # Find the module using the standard mechanism
            # Note: importlib.util.find_spec takes (name, package) not (name, path)
            # We pass None for package to search sys.path
            spec = importlib.util.find_spec(fullname)
            if spec is None:
                return None

            origin = spec.origin
            if origin is None or not origin.endswith('.py'):
                return None

            # Only intercept files in .mcode or test directories
            if '.mcode' not in origin and 'test' not in origin.lower():
                return None

            # Read the source
            try:
                with open(origin, 'rb') as f:
                    source_bytes = f.read()
            except (OSError, IOError):
                return None

            # Only inject if source contains union type syntax (X | None, etc.)
            source_text = source_bytes.decode('utf-8', errors='replace')
            if '| None' not in source_text and '|None' not in source_text:
                return None

            # Create a new spec with our custom loader
            loader = _FutureAnnotationsLoader(origin, source_bytes)
            new_spec = importlib.machinery.ModuleSpec(
                fullname,
                loader,
                origin=origin,
            )
            new_spec.submodule_search_locations = spec.submodule_search_locations
            return new_spec
        except Exception:
            return None
        finally:
            self._active = False


def install():
    """Install the import hook if running Python < 3.10."""
    if sys.version_info < (3, 10):
        sys.meta_path.insert(0, _FutureAnnotationsFinder())
