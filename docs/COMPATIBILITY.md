# Dialogbox Compatibility Matrix

The reference implementation is `helton-godoy/dialogbox`, commit `6989740`. Status must be supported by an automated test; code presence alone is insufficient.

Status values: **pass**, **partial**, **missing**, or **unverified**.

## Runtime Contract

| Capability | Dialogbox | SHantilly | Evidence | Status |
| --- | --- | --- | --- | --- |
| Commands read from `stdin` | Yes | Preserved CLI entry point selected for production | `compatibility_stdin_query` | Pass |
| Events and values on `stdout` | Yes | Query output verified; interactive events remain unverified | `compatibility_stdin_query` | Partial |
| `--help` and `--version` | Yes | Preserved CLI entry point selected for production | `compatibility_cli_help`, `compatibility_cli_version` | Pass |
| `--resizable` | Yes | Preserved CLI and constructor support it | None | Unverified |
| `--hidden` and explicit `show` | Yes | `--hidden` verified; explicit `show` remains unverified | `compatibility_stdin_query` | Partial |
| Accept/reject exit semantics | Yes | Legacy implementation retained | None | Unverified |
| Pipes/FIFOs and background input | Yes | Parser code is transitional | None | Unverified |

## Commands and Layout

| Area | Expected baseline | Status | Required evidence |
| --- | --- | --- | --- |
| `add`, `set`, `unset`, `remove`, `clear` | Original syntax and effects | Unverified | Protocol integration tests |
| `step`, `position`, `end` | Original nesting behavior | Unverified | Layout characterization tests |
| `query`, `print`, `show` | Original output and visibility | Unverified | Stdout/visibility tests |

## Widgets

| Widgets | Status | Notes |
| --- | --- | --- |
| Original set: checkbox, combobox, frame, groupbox, label, listbox, page, progressbar, pushbutton, radiobutton, separator, slider, tabs, textbox, textview | Unverified | Implementations exist, but active CTest does not exercise them |
| Calendar | Extension, unverified | Not part of original compatibility contract |
| Table | Extension, unverified | Not part of original compatibility contract |
| Chart | Extension, unverified | Not part of original compatibility contract |

## Update Rules

- Link each **pass** status to a test name or CI job.
- Add a row before changing observable protocol behavior.
- Never replace an original behavior silently; document intentional changes in an ADR.
- Update this matrix in the same PR that changes compatibility.
