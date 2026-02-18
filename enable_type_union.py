"""
Backport PEP 604 (X | Y) union type syntax for Python 3.9.
This patches type.__or__ and type.__ror__ to support union type expressions
at runtime, matching Python 3.10+ behavior.
"""
import sys
import typing

if sys.version_info < (3, 10):
    def _type_or(self, other):
        return typing.Union[self, other]

    def _type_ror(self, other):
        return typing.Union[other, self]

    # Patch type to support X | Y syntax
    type.__or__ = _type_or
    type.__ror__ = _type_ror

    # Also patch NoneType for X | None
    NoneType = type(None)
    NoneType.__or__ = _type_or
    NoneType.__ror__ = _type_ror
