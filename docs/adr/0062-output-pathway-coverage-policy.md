# 0062 - Output Pathway Coverage Policy

## Status

Accepted

## Context

Before Phase 3, reports were visible but not portable.

## Decision

Support copying summaries, report JSON, and curl commands; downloading report JSON, dependencies CSV, markdown summary, and full state; importing exported state; small hash share links; and browser print/PDF.

Screenshot and embed code stay out of scope because the browser and OS already handle screenshots, and embed widgets are not appropriate for local-backend security reports.

## Consequences

Reports become artifacts that can move into issues, docs, reviews, and archives.

## Alternatives Considered

Server-side saved reports were rejected because v1-v3 intentionally avoid accounts and backend persistence.
