# Dialogbox Compatibility Matrix

The reference implementation is `helton-godoy/dialogbox`, commit `6989740`. Status must be supported by an automated test; code presence alone is insufficient.

Status values: **pass**, **partial**, **missing**, or **unverified**.

## Runtime Contract

| Capability | Dialogbox | SHantilly | Evidence | Status |
| --- | --- | --- | --- | --- |
| Commands read from `stdin` | Yes | Not connected by current `main()` | None | Missing |
| Events and values on `stdout` | Yes | Implementation retained, runtime unverified | None | Unverified |
| `--help` and `--version` | Yes | Not handled by current `main()` | None | Missing |
| `--resizable` | Yes | Constructor supports it; CLI does not | None | Partial |
| `--hidden` and explicit `show` | Yes | Not handled by current `main()` | None | Missing |
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
