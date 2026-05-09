# 0046 - Performance Budgets

## Status

Accepted

## Context

Scanner process time dominates real-data latency. Huge inputs need honest progress and scale metadata.

## Decision

Parsing and inference should complete in under 300ms for fixtures under 1MB. Scanner work can exceed that, but reports must mark scanner partials and include elapsed time. Operations over 5 seconds must be cancellable from the frontend.

## Consequences

Phase 2 focuses on truthful partial results and cancellation rather than hiding scanner latency.

## Alternatives Considered

Pretending synchronous audits are always fast was rejected.
