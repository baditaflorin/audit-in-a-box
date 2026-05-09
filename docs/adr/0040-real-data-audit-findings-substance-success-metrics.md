# 0040 - Real-Data Audit Findings And Substance Success Metrics

## Status

Accepted

## Context

The v1 happy path works, but real dependency manifests exposed wrong-confident license scoring, unsupported common formats, brittle scanner parsing, and unhelpful errors.

## Decision

Use the ten fixtures in `test/fixtures/realdata/` as the Phase 2 grading rubric. Success means at least seven fixtures produce useful reports without manual correction, no misunderstood input is reported clean, and stable report projections are deterministic.

## Consequences

Every inference change must improve or preserve this fixture set. Regressions block release unless documented by a later ADR.

## Alternatives Considered

Continuing with synthetic demo manifests was rejected because it would not catch the observed failure modes.
