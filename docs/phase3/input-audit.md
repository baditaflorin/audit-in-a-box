# Phase 3 Input Pathway Audit

Date: 2026-05-10

Mode remains Mode C.

| Entry point | Status before Phase 3 | Evidence | Decision |
| --- | --- | --- | --- |
| File upload | Works partially | Single file picker reads text and fills filename/content. No multi-file, import state, drag-drop, or explicit parse preview. | Finish |
| Drag-drop file | Not built | No drag/drop handlers in `frontend/src/App.tsx`. | Finish |
| Paste manifest text | Works fully | Manifest content textarea posts to `/api/v1/audits`. | Keep |
| Paste rendered HTML | Works partially | HTML textarea + `/api/v1/scrape` exists, but no clipboard helper and no URL guidance. | Finish |
| Paste image | Not built | Product audits text manifests; OCR/image parsing is not in scope for Mode C v1-v3. | Out of scope via ADR 0061 |
| URL input | Not built | README claims pasted HTML, not URL fetching, but real users will bring GitHub URLs. | Finish with direct fetch plus CORS fallback guidance |
| Clipboard read | Not built | No `navigator.clipboard` use. | Finish |
| Mobile picker | Works partially | Native file input works on mobile Files; no explicit mobile camera/photo handling. Images are out of scope. | Document |
| Multi-file input | Not built | File input lacks `multiple`; audit mutation is single-report only. | Finish |
| Folder input | Not built | Browser folder picker is inconsistent and not necessary for manifest audit. | Out of scope via ADR 0061 |
| Sample/demo | Works partially | Three sample buttons exist, fewer than supported formats. | Finish |
| Deep links | Not built | No hash/state parsing. | Finish for small saved state |
| Imported state | Not built | No state-file import path. | Finish |
| Restored autosave | Works partially | Backend URL and last report persist; current draft input does not. | Finish |

## Summary

Green before Phase 3: 1.

Yellow before Phase 3: 5.

Red before Phase 3: 8.

Highest-impact gap: a stranger with several manifests or a GitHub URL has to manually copy text into the correct box and loses draft input on reload.
