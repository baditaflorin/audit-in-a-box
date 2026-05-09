# Phase 3 Codebase Health Audit

Date: 2026-05-10

## Measurements Before Phase 3

| Metric | Count / finding |
| --- | --- |
| Frontend modules over 300 lines | 2: `frontend/src/App.tsx` (389), `frontend/src/features/audit/ReportView.tsx` (386) |
| Go modules over 300 lines | 0 |
| TODO/FIXME/XXX/HACK in source | 0 in source; one FIXME exists inside committed real-world fixture text and is not project debt |
| `any` in source | 19 raw matches, mostly JSON/API boundary code |
| `unknown` in source | 8 raw matches, mostly Zod/API boundary code |
| `@ts-ignore` | 0 |
| Commented-out blocks | 0 found |
| Dead production controls | 0 stubs found |

## DRY Findings

1. Report serialization exists only implicitly through `JSON.stringify(report)` in localStorage. No single source of truth for report export/import.
2. API JSON response parsing uses direct casts in `frontend/src/api/client.ts`; Zod validation exists for reports but not for tools/scrape wrapper envelopes.
3. Date/version/export formatting is ad hoc inside `ReportView`.

## SOLID Findings

1. `frontend/src/App.tsx` owns storage, backend settings, input handling, scraping, audit mutation, samples, and page layout. It has several reasons to change.
2. `frontend/src/features/audit/ReportView.tsx` owns presentation plus formatter logic and future export actions would make it larger.
3. `internal/analysis/service.go` orchestrates parsing, scanners, maintainer health, DuckDB, scoring, evidence, and summary. This is accepted for Phase 3 because engine changes are locked.

## Dead Code

No unreferenced production UI controls or dormant feature flags were found. Export/import code is absent rather than dead.

## Type-Safety Holes

1. `frontend/src/App.tsx` casts `JSON.parse(stored) as AuditReport` without Zod validation.
2. `frontend/src/App.tsx` casts GitHub commit API JSON directly.
3. `frontend/src/api/client.ts` casts tools and scrape envelopes before full boundary validation.
4. Go `any` usage is concentrated in JSON boundary helpers, TOML/YAML decoders, and report provenance parameters.

## Inconsistent Patterns

1. Storage keys live in multiple places (`client.ts`, `App.tsx`) with no migration strategy.
2. Frontend API errors use a helper, but direct fetches like the GitHub commit lookup do not share a wrapper.
3. Sample inputs are loaded from button labels rather than a typed sample registry.

## Test Coverage Holes

1. No frontend tests cover file upload, draft restoration, report export, or state import.
2. Smoke test checks only page load, links, version, commit, and form presence.
3. Phase 2 real-data fixtures cover parser behavior but not full browser input/output workflows.
