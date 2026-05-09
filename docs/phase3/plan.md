# Phase 3 Completeness Plan

Date: 2026-05-10

Ranked by real-user impact from the Phase 3 audits.

## Picklist

1. A1 file input pathways: single upload, drag-drop, mobile picker, and multi-file.
2. A2 frontend format detection and filename normalization on upload/URL/import.
3. A3 URL input with raw GitHub conversion and CORS fallback guidance.
4. A4 batch audit with per-file progress and partial errors.
5. A6 clipboard read for manifest/HTML paste.
6. A7 samples for every supported format, not only three.
7. A8 draft input and session resume with one-click start fresh.
8. B9 JSON and CSV exports that work on real reports.
9. B10 copy summary/report/curl to clipboard with confirmation.
10. B11 downloadable versioned state file.
11. B12 shareable hash URL for small state with documented limits.
12. B13 print/PDF control with print-friendly CSS.
13. B14 API/curl-ready output generated from current input/backend.
14. C15 half-baked feature triage recorded in ADR 0063.
15. C16 finish kept HTML scraper and state persistence paths.
16. C17 hide/delete out-of-scope image/folder/embed/screenshot paths from claims.
17. C18 treat backend URL as the only setting and make it persist immediately.
18. C19 align README/API docs with reality and limitations.
19. D20 extract storage/import/export logic to a single frontend module.
20. D21 consolidate frontend API envelope validation.
21. D22 keep canonical frontend report types from Zod schemas.
22. D23 shared boundary schemas for report state and exported state.
23. E24 split App responsibilities into storage/export/sample helpers without engine changes.
24. E25 document dependency direction and keep UI imports one-way.
25. E27 define storage/export/API seams where interfaces are useful.
26. F28 delete no code unless truly unreferenced; audit found no production dead code.
27. F29 resolve hook-discovered multipart G120 finding.
28. G31 keep one user-facing error convention in client helpers.
29. G32 keep one state-management convention: React state plus versioned localStorage.
30. H35 eliminate unsafe frontend JSON casts.
31. H36 validate every frontend external boundary with Zod.
32. I38 saved draft/report/batch/backend survive reload.
33. I39 add localStorage state migrations.
34. I40 clear-state operation works.
35. I41 exported state round-trips through import.
36. J42 README features checklist matches tests.
37. J43 quickstart clarified so backend/frontend steps are honest.
38. J44 minimum inline help for URL, clipboard, batch, export, and backend URL.
39. J45 README limitations section.
40. K46 run private-window stranger test with real fixture.
41. K47 fix the top three stranger-test findings.

## Implementation Batches

1. ADRs 0060-0071.
2. Input completeness and state persistence.
3. Output completeness and export/import.
4. Type-safety, DRY, docs alignment.
5. Stranger test and top-three fixes.
6. Version bump, Pages build, release.

## Gate

The Phase 2 real-data fixture suite remains the floor. `make test`, `make lint`, `make build`, and `make smoke` must pass before push.
