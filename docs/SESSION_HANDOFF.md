# Session Handoff

Update this file at the end of every development session. Keep it factual and concise; Git history preserves older entries.

## Current State

- Active phase: **Phase 0 — Baseline**.
- Last completed task: project naming was standardized under ADR 0004 and validated in the Qt6 container.
- Working tree expectation: naming, packaging, obsolete qmake cleanup, and documentation changes pending review.
- Reference repository: `helton-godoy/dialogbox`, commit `6989740`; local audit clone was placed at `/tmp/dialogbox-reference` and must not be treated as persistent.

## Verified Findings

- The current executable opens a static demonstration and does not connect the original `stdin/stdout` protocol.
- CMake registers only `tests/example_test.cpp`, a trivial GoogleTest.
- Qt Test sources under `tests/auto` are not part of the active CMake test target and reference parts of an obsolete layout.
- Local CMake configuration could not complete because Qt6 development files were unavailable in the host environment.

## Next Exact Step

Characterize the remaining CLI options and invalid-option exit behavior, then expand command coverage beyond the minimal `query` scenario.

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
