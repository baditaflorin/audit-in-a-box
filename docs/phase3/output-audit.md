# Phase 3 Output Pathway Audit

Date: 2026-05-10

| Exit path | Status before Phase 3 | Evidence | Decision |
| --- | --- | --- | --- |
| JSON report | Works internally only | Report is visible and persisted, but no explicit copy/download. | Finish |
| CSV export | Not built | Dependencies render in cards/tables only. | Finish for dependency inventory |
| Code/API export | Works partially | `docs/api.md` has curl examples; UI has none for the current backend/input. | Finish curl copy |
| Copy-to-clipboard | Not built | No clipboard write controls. | Finish |
| Downloadable state file | Not built | Last report is localStorage-only. | Finish |
| Import exported state | Not built | No import control. | Finish |
| Share link | Not built | No hash/state URL. | Finish for small states only |
| Print-friendly output | Works partially | Browser print can print page chrome; no print control or print CSS. | Finish minimum |
| Screenshot | Not built | Not claimed. Browser/system screenshot is sufficient. | Out of scope via ADR 0062 |
| Embed code | Not built | Not claimed and not appropriate for local-backend audit reports. | Out of scope via ADR 0062 |
| API/curl-ready | Works partially | Static docs only, not report/input-specific. | Finish |

## Summary

Green before Phase 3: 0.

Yellow before Phase 3: 3.

Red before Phase 3: 8.

Highest-impact gap: users can get a report but cannot reliably take it with them, rerun it, share it, or archive it.
