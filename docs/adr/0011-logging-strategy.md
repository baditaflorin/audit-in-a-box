# 0011 - Logging Strategy

## Status

Accepted

## Context

Mode C needs useful server logs without leaking manifests.

## Decision

Use Go `slog` JSON logs on stdout with `trace_id`, request path, status, duration, and high-level scanner outcomes. Avoid logging full manifest contents.

The frontend keeps production console output minimal and only logs unrecoverable diagnostics.

## Consequences

Docker logs are parseable and safe for normal operations.

## Alternatives Considered

- Text logs: rejected because JSON logs are easier to aggregate.
- Verbose request dumps: rejected because manifests may be sensitive.
