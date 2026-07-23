# Session Handoff

Update this file at the end of every development session. Keep it factual and concise; Git history preserves older entries.

## Current State

- Active phase: **Phase 0 — Baseline**.
- Last completed task: ordered protocol mutations and explicit `show` commands were characterized.
- Working tree expectation: clean after the protocol mutation characterization commit.
- Reference repository: `helton-godoy/dialogbox`, commit `6989740`; local audit clone was placed at `/tmp/dialogbox-reference` and must not be treated as persistent.

## Verified Findings

- The production target uses the preserved V1 CLI entry point and reads the original protocol from `stdin`.
- CMake links unit tests to production library code and registers compatibility and quality checks in CTest.
- Historical Qt Test sources under `tests/auto` are not active and require evaluation before reuse or removal.
- Host CMake configuration remains unavailable because Qt6 development files are not installed; Docker is the verified environment.

## Next Exact Step

Characterize `clear`, layout commands, and representative widget values.

## Completion Checklist

- [ ] Update roadmap task status.
- [ ] Update compatibility evidence for behavior changes.
- [ ] Record commands and results below.
- [ ] State blockers and the next exact step.

## Validation Log

- `cmake -S . -B build-audit -DCMAKE_BUILD_TYPE=Debug`: blocked because `Qt6Config.cmake` was unavailable locally.
- Added CTest characterizations for `--help`, `--version`, and a minimal checkbox `query`; execution remains blocked locally by the missing Qt6 development package.
- The CMake assertion script parsed and reported the expected failure when exercised with `/bin/false`.
- `make docker-build`: failed because the container could not reach Debian repositories; package indexes remained empty and `apt-get install` exited with code 100.
- A direct `apt-get update` in `debian:trixie-slim` also stalled on every `deb.debian.org` InRelease endpoint, confirming a Docker network/repository access problem rather than an individual package name.
- `docker build --network host --pull -t shantilly-dev:latest -f src/dev.Dockerfile .`: passed.
- Containerized `make build_internal`: passed with Qt 6.8.2; optional XKB and CUPS discovery emitted warnings.
- Containerized `ctest --output-on-failure`: 4/4 passed, including three compatibility tests and the production-linked `WidgetConfigsTest` cases.
- Updated README, installation, architecture, and implementation-status documents to match the verified CMake layout and current runtime state.
- Standardized the brand as `SHantilly` and technical Unix identifiers as `shantilly` under ADR 0004; removed obsolete generated/qmake test files.
- Containerized build and CTest after naming changes: passed, 5/5 tests, including the naming-convention quality gate.
- Characterized `-h`, `-v`, `--resizable`, `-r`, and invalid-option exit/stderr behavior in CTest.
- Clean container build in `/tmp/shantilly-build` and full CTest suite: passed, 10/10 tests.
- Characterized ordered `add`, `set`, `unset`, `enable`, `disable`, `show`, `hide`, `query`, and `remove` commands.
- Clean container build and full CTest suite after protocol expansion: passed, 11/11 tests.
