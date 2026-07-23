# Session Handoff

Update this file at the end of every development session. Keep it factual and concise; Git history preserves older entries.

## Current State

- Active phase: **Phase 0 — Baseline**.
- Last completed task: repository governance documents created.
- Working tree expectation: governance documentation changes pending review.
- Reference repository: `helton-godoy/dialogbox`, commit `6989740`; local audit clone was placed at `/tmp/dialogbox-reference` and must not be treated as persistent.

## Verified Findings

- The current executable opens a static demonstration and does not connect the original `stdin/stdout` protocol.
- CMake registers only `tests/example_test.cpp`, a trivial GoogleTest.
- Qt Test sources under `tests/auto` are not part of the active CMake test target and reference parts of an obsolete layout.
- Local CMake configuration could not complete because Qt6 development files were unavailable in the host environment.

## Next Exact Step

Create Phase 0 characterization fixtures for `--help`, `--version`, and a minimal `stdin` pushbutton dialog based on the original behavior. Then wire production-linked tests into CMake without restoring the entire monolith blindly.

## Completion Checklist

- [ ] Update roadmap task status.
- [ ] Update compatibility evidence for behavior changes.
- [ ] Record commands and results below.
- [ ] State blockers and the next exact step.

## Validation Log

- `cmake -S . -B build-audit -DCMAKE_BUILD_TYPE=Debug`: blocked because `Qt6Config.cmake` was unavailable locally.
