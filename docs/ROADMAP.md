# SHantilly Modernization Roadmap

This roadmap controls the evolution from `dialogbox` to SHantilly. Only one phase may be active at a time. A phase is complete only when every exit criterion is demonstrated by automated tests.

## Status

| Phase | Goal | Status |
| --- | --- | --- |
| 0 | Establish governance and behavioral baseline | In progress |
| 1 | Restore the scriptable CLI contract | Pending |
| 2 | Separate protocol, parser, execution, and UI | Pending |
| 3 | Migrate widgets into `SHantilly-ui` | Pending |
| 4 | Add SHantilly extensions safely | Pending |
| 5 | Harden releases, packaging, and documentation | Pending |

## Phase 0 — Baseline

Scope: document the original contract, make the build deterministic, and run meaningful tests in CI.

- [x] Record mission, roadmap, compatibility matrix, ADRs, and handoff process.
- [x] Import representative `dialogbox` behavior as characterization fixtures.
- [x] Replace the trivial CTest target with production-linked tests.
- [ ] Remove stale build and test instructions.
- [ ] Establish a known-good baseline for supported platforms.

Exit criteria: CI builds production code and tests the CLI protocol, options, exit status, and representative widgets.

## Phase 1 — Functional Compatibility

Restore command input on `stdin`, event/value output on `stdout`, parser lifecycle, `--help`, `--version`, `--resizable`, and `--hidden`. New widgets are out of scope.

Exit criteria: selected original demos run unchanged and all Phase 1 matrix entries pass.

## Phase 2 — Architecture

Extract tokenizer, command model, parser, execution context, and builder interfaces without changing observable behavior. Remove deprecated V2 code from active targets.

Exit criteria: components have isolated tests and compatibility tests remain green.

## Phase 3 — Widget Migration

Move one widget at a time to `SHantilly-ui`, preserving ownership, focus, layout, reporting, and mutation semantics.

Exit criteria: every original widget is marked compatible and covered by tests.

## Phase 4 — Evolution

Stabilize themes, icons, calendar, table, and charts. Extensions must be backward compatible or introduced through a documented protocol version.

## Phase 5 — Release Readiness

Validate packages, installation, security checks, manuals, examples, migration notes, and supported-platform CI. Produce a release candidate only after the compatibility matrix is complete.

## Change Control

Changes to phase order, compatibility policy, or architecture require an ADR. Update [SESSION_HANDOFF.md](SESSION_HANDOFF.md) after every work session.
