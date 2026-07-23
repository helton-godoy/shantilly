# ADR 0004: Standardize Project Naming

- Status: Accepted
- Date: 2026-07-22

## Context

The repository uses `SHantilly` and `shantilly` interchangeably for the product, executable, packages, source identifiers, and documentation examples. Case-sensitive systems consequently expose broken commands and packaging paths.

## Decision

Use names according to their role:

- **SHantilly** is the product and brand name in prose, UI titles, documentation headings, and C++ types such as `SHantilly` and `SHantillyBuilder`.
- **shantilly** is the Unix executable, command, process, package, internal CMake project/target prefix, directory name, and logging-category prefix.
- **SHANTILLY_*** is the prefix for environment variables and preprocessor constants.
- C++ filenames that define a public type may match that type; other new files use lowercase snake case.

Do not introduce the variants `Shantilly`, `SHANTILLY`, or capitalized command examples. Existing repository URLs and published application IDs are compatibility identifiers and require a separate migration before being renamed.

## Consequences

Shell documentation consistently invokes `shantilly`, while users continue to see the SHantilly brand. Code review and CI can detect invalid variants without forcing breaking changes to published IDs.
