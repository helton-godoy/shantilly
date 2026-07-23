# ADR 0002: Use CMake as the Primary Build System

- Status: Accepted
- Date: 2026-07-22

## Context

The repository contains both qmake and CMake configurations. Documentation and older tests reference paths that no longer exist, while CI and packaging already use CMake and vcpkg.

## Decision

CMake is the authoritative build and test system. New production targets and tests must be registered with CMake and executable through CTest. qmake files may remain temporarily for historical comparison but must not define the expected CI behavior.

## Consequences

Documentation, CI, local builds, coverage, and packaging share one dependency graph. Transitional qmake files should be removed after equivalent CMake coverage exists.
