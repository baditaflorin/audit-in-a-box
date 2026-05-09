# Phase 3 Findings Synthesis

Date: 2026-05-10

## Top 5 Usability Gaps

1. A user can get a report but cannot copy, download, print, share, or re-import it.
2. Multi-file and drag-drop workflows are absent even though real repositories often have several manifests.
3. GitHub URLs are a natural input, but the UI only supports manual text/HTML paste.
4. Draft input is not restored after reload; only backend URL and last report persist.
5. Supported formats are not equally visible in the UI, so users may assume only three formats work.

## Top 5 Half-Baked Features

1. Last-report persistence: finish with versioned state import/export and migration.
2. HTML scraper: finish with clipboard and URL helpers.
3. Supported input list: finish with typed sample registry for all supported formats.
4. Debug surface: keep hidden behind `?debug=1`, document in Phase 3 docs.
5. Tool-status refresh: keep, but do not grow it into a settings panel; backend URL remains the only setting.

## Top 5 Codebase Pain Points

1. `App.tsx` mixes storage, input, mutation, samples, and layout.
2. `ReportView.tsx` is already large before export controls.
3. Frontend storage lacks schema versioning and migration.
4. API boundary casts bypass Zod in a few places.
5. Smoke tests do not cover real-user output paths.

## Top 5 Documentation/Reality Mismatches

1. README supported inputs are real in the parser but not first-class in the UI.
2. API docs show curl examples, but UI cannot generate curl for the current input.
3. The product implies a report artifact, but no export exists.
4. Pasted GitHub blob HTML is supported, but the UI does not tell users how to get there from a URL.
5. The quickstart says `make dev`, but that starts only the frontend and prints a backend instruction.

## Fully Usable Means

- A stranger can load one or more of their own manifests by upload, drag-drop, paste, clipboard, or a raw/GitHub URL.
- A stranger can run audits, see partial failures per file, and keep working.
- A stranger can take the report out as JSON, CSV, state file, clipboard text, print/PDF, or curl.
- A stranger can reload the page and recover the draft, last report, backend URL, and batch results.
- Every visible control performs the action its label promises on real user data.

## Phase 3 Success Metrics

- Input audit green/out-of-scope rows: 14/14.
- Output audit green/out-of-scope rows: 11/11.
- At least 20 catalog items implemented or explicitly closed by ADR.
- Zero production TODO/FIXME/XXX/HACK.
- Frontend localStorage reads validate and migrate through Zod schemas.
- Smoke test covers one import/output path in addition to the existing page-load checks.

## Out Of Scope

- No new scanner/engine behavior.
- No visual polish pass.
- No accounts, cloud sync, or backend persistence.
- No image/OCR manifest parsing.
- No embeddable report widgets.
