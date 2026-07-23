# Repository Guidelines

## Mission and Source of Truth

SHantilly modernizes `helton-godoy/dialogbox` while preserving its shell-facing behavior. Before changing code, read [the roadmap](docs/ROADMAP.md), [the compatibility matrix](docs/COMPATIBILITY.md), applicable [ADRs](docs/adr/), and [the session handoff](docs/SESSION_HANDOFF.md). Work only on the active roadmap phase unless the user explicitly changes priorities.

## Project Structure

- `src/code/shantilly/`: executable and transitional legacy code.
- `libs/SHantilly-ui/`: reusable Qt6 widgets and builders.
- `tests/`: CTest, Qt Test, integration, and compatibility tests.
- `examples/`: shell examples and behavioral demonstrations.
- `docs/`: user documentation, architecture, roadmap, and decisions.
- `packaging/`: distribution tooling.

Do not treat `legacy/v2_incomplete` as the target architecture. Use the original `dialogbox` behavior and the compatibility matrix as the baseline.

## Build and Validation

Use the Docker-backed targets when possible:

```bash
make build       # Configure and compile with CMake
make test        # Run CTest and integration checks
make lint        # Run Trunk and clang-format checks
make coverage    # Build and produce coverage data
```

Tests must exercise actual production targets. A trivial smoke test alone is not sufficient. Record environment limitations in `docs/SESSION_HANDOFF.md` instead of claiming validation.

## Coding and Testing Rules

Use C++17, Qt6 ownership conventions, four-space indentation, K&R braces, and the repository `.clang-format`. Classes use `PascalCase`, methods `camelCase`, and constants `UPPER_SNAKE_CASE`. Prefer explicit CMake source lists over recursive globs.

Add characterization tests before changing legacy behavior. Preserve command syntax, standard input/output protocol, exit codes, and CLI options unless an approved ADR explicitly changes them. Every migrated widget needs unit and compatibility coverage.

## Commits, Pull Requests, and Handoff

Use Conventional Commits, such as `refactor(parser): isolate command tokenizer`. PRs must identify the roadmap phase, linked issue, compatibility impact, tests executed, and documentation or ADR changes. Before ending work, update the roadmap status, compatibility matrix when behavior changed, and session handoff with the exact next step. Never mark a phase complete until all listed exit criteria pass in CI.
