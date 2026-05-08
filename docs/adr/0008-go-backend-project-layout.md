# 0008 - Go Backend Project Layout

## Status

Accepted

## Context

The backend is a runtime API in Mode C and should follow common Go layout conventions.

## Decision

Use `cmd/`, `internal/`, `pkg/`, `api/`, `configs/`, `scripts/`, and `test/`.

The server entrypoint is `cmd/server`.

## Consequences

Private implementation details stay under `internal/`. Public build metadata lives in `pkg/version`.

## Alternatives Considered

- Flat package layout: rejected because the backend has several independent concerns.
