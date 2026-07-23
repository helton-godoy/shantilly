# Dialogbox Compatibility Matrix

The reference implementation is `helton-godoy/dialogbox`, commit `6989740`. Status must be supported by an automated test; code presence alone is insufficient.

Status values: **pass**, **partial**, **missing**, or **unverified**.

## Runtime Contract

| Capability | Dialogbox | SHantilly | Evidence | Status |
| --- | --- | --- | --- | --- |
| Commands read from `stdin` | Yes | Preserved CLI entry point selected for production | `compatibility_stdin_query` | Pass |
| Events and values on `stdout` | Yes | Query values, state transitions, and real pushbutton clicks verified | Protocol and interaction tests | Pass |
| `--help` and `--version` | Yes | Preserved CLI entry point selected for production | `compatibility_cli_help`, `compatibility_cli_version` | Pass |
| `--resizable` | Yes | Long and short options are accepted | `compatibility_cli_resizable`, `compatibility_cli_short_r` | Pass |
| `--hidden` and explicit `show` | Yes | Hidden startup and ordered `show` commands are accepted; visual state remains unverified | `compatibility_stdin_query`, `compatibility_protocol_mutations` | Partial |
| Invalid option handling | Exit 1 and diagnostic on `stderr` | Exact behavior preserved | `compatibility_cli_invalid_option` | Pass |
| Accept/reject exit semantics | Yes | Apply stays open; acceptance returns 1 and rejection returns 0 | `dialog_interaction_tests` | Pass |
| Pipes/FIFOs and background input | Yes | Parser code is transitional | None | Unverified |

## Commands and Layout

| Area | Expected baseline | Status | Required evidence |
| --- | --- | --- | --- |
| `add`, `set`, `unset`, `remove`, `clear` | Core mutations and dialog/list clearing verified | Pass | `compatibility_protocol_mutations`, `compatibility_clear_semantics` |
| `step`, `position`, `end` | Horizontal/vertical steps and container positioning preserve report order | Pass | `compatibility_layout_widgets` |
| `query`, `print`, `show` | `query` output and ordered `show` accepted; visibility and `print` pending | Partial | `compatibility_stdin_query`, `compatibility_protocol_mutations` |

## Widgets

| Widgets | Status | Notes |
| --- | --- | --- |
| Checkbox, combobox, groupbox, listbox, page, progressbar, radiobutton, slider, tabs, textbox | Partial | Values, item addressing, non-reportable progress, and legacy empty tab markers covered by compatibility fixtures |
| Frame, label, pushbutton, separator, textview | Partial | Non-reporting behavior, toggle transitions, nesting, and recursive frame removal covered by `compatibility_passive_widgets` |
| Calendar | Extension, unverified | Not part of original compatibility contract |
| Table | Extension, unverified | Not part of original compatibility contract |
| Chart | Extension, unverified | Not part of original compatibility contract |

## Update Rules

- Link each **pass** status to a test name or CI job.
- Add a row before changing observable protocol behavior.
- Never replace an original behavior silently; document intentional changes in an ADR.
- Update this matrix in the same PR that changes compatibility.
