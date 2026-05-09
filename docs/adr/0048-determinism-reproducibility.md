# 0048 - Determinism And Reproducibility

## Status

Accepted

## Context

Reports need enough provenance to be trusted and rerun.

## Decision

Sort dependencies, vulnerabilities, license risks, warnings, and maintainer signals deterministically. Add input hash, schema version, normalized byte count, parser name, parser confidence, app version, commit, and scanner/tool status to report provenance. Fixture tests compare a stable projection that excludes wall-clock fields.

## Consequences

Same input yields the same substantive report and tests can detect accidental ordering drift.

## Alternatives Considered

Leaving map/slice ordering incidental was rejected.
