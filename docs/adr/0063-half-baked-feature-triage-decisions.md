# 0063 - Half-Baked Feature Triage Decisions

## Status

Accepted

## Context

The audit found no production stubs, but several partially-finished paths.

## Decision

Finish:

- Last-report persistence, by adding versioned state, import/export, migration, and clear state.
- HTML scraper, by adding clipboard and URL workflows.
- Supported input visibility, by expanding samples to all supported formats.
- Backend URL setting, by persisting immediately.
- Report portability, by adding export/copy/print/share controls.

Keep hidden:

- `?debug=1`, as a support surface rather than visible chrome.

Delete or keep unclaimed:

- Image input, folder input, screenshots, embed code. These remain documented out of scope and do not appear as product claims.

## Consequences

The UI has fewer mystery paths and more completed paths.

## Alternatives Considered

Adding a Settings page was rejected. The only real setting is backend URL, already on the main workflow.
