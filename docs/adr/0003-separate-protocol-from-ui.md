# ADR 0003: Separate Protocol Processing from the UI

- Status: Accepted
- Date: 2026-07-22

## Context

The original application combines parsing, command execution, widget management, and process I/O across a monolithic dialog class and parser thread. A previous incomplete refactoring introduced parallel abstractions without establishing a tested migration path.

## Decision

Refactor behind characterization tests into four boundaries: protocol I/O, tokenizer/parser, command execution, and UI builder/view. Dependencies flow toward interfaces; parsing must be testable without displaying a GUI, and widgets must not parse protocol text.

## Consequences

Migration proceeds incrementally rather than through a second implementation. Deprecated experimental code is excluded from active targets once its required behavior is represented in the chosen components and tests.
