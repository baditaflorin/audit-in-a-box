# 0012 - Metrics And Observability

## Status

Accepted

## Context

The backend must expose basic operational signals for Mode C.

## Decision

Expose `/metrics` with Prometheus metrics:

- HTTP request count and duration.
- Audit request count.
- Audit duration.
- Scanner failure count.
- Report risk score histogram.

The frontend ships no analytics by default.

## Consequences

Operators can monitor the backend while the GitHub Pages frontend remains privacy-preserving.

## Alternatives Considered

- Client analytics: rejected for v1.
- Logs only: rejected because RED metrics are cheap and useful.
