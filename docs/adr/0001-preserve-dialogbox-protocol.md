# ADR 0001: Preserve the Dialogbox Protocol

- Status: Accepted
- Date: 2026-07-22

## Context

SHantilly evolves `dialogbox`, whose defining contract is a text filter: commands enter through standard input, GUI events and values leave through standard output, and exit status communicates acceptance or rejection. The current refactoring no longer connects that behavior from the executable entry point.

## Decision

Treat the original command syntax, `stdin/stdout` behavior, CLI options, and exit semantics as the compatibility baseline. Refactoring must first characterize and preserve observable behavior. Extensions may add commands and widgets, but incompatible changes require an explicit protocol version and a new ADR.

## Consequences

Compatibility tests become release gates. Internal architecture may change freely when those tests remain green. Existing shell scripts remain valid, while SHantilly can evolve through additive features.
