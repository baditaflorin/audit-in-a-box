# 0045 - State Taxonomy And State Machine

## Status

Accepted

## Context

Real audits can be slow, partial, cancelled, or recoverable.

## Decision

Use the backend and frontend states documented in `docs/phase2-substance/states.md`. Long-running frontend requests use `AbortController`; backend request context cancellation must stop subprocesses.

## Consequences

Double-clicks, cancellation, and backend failures have defined behavior.

## Alternatives Considered

Leaving request state implicit was rejected because it produces stuck or half-loaded states.
