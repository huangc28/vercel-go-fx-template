# Interfaces

Put cross-domain interfaces (and only minimal shared types) under `lib/interfaces/*` to avoid cyclic imports between domains.

Guideline:
- Domain A can depend on `lib/interfaces/foo`, not on Domain B's concrete implementation.
- Domain B can implement interfaces from `lib/interfaces/foo`, and the wiring happens via FX in the caller's `api/<domain>/core.go`.

