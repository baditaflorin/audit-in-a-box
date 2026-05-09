# 0049 - Inspectability And Debug Surface

## Status

Accepted

## Context

Power users need to understand why the app inferred an ecosystem or risk.

## Decision

Support `?debug=1` in the frontend to reveal report provenance, confidence reasons, warnings, and raw evidence metadata already present in the API response.

## Consequences

Support/debugging improves without changing the primary surface for normal users.

## Alternatives Considered

Only server logs were rejected because users of the Pages frontend cannot inspect backend internals.
