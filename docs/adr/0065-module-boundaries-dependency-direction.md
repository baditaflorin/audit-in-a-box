# 0065 - Module Boundaries and Dependency Direction

## Status

Accepted

## Context

The app is small, but Phase 3 adds state and export paths that can tangle UI code quickly.

## Decision

Frontend dependency direction is:

`App.tsx` and feature views -> `features/audit/*` helpers -> `api/*` schemas/client -> browser primitives.

Backend engine boundaries remain unchanged for Phase 3 because no engine changes are allowed.

## Consequences

New completeness helpers are isolated and testable without changing scanner orchestration.

## Alternatives Considered

Introducing a global state library was rejected; React state plus versioned localStorage is sufficient.
