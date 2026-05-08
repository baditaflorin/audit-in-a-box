# 0013 - Testing Strategy

## Status

Accepted

## Context

Risk scoring, parsers, API behavior, and the built Pages app need verification.

## Decision

Use Go unit tests for parsers and scoring, Vitest for frontend logic, and Playwright smoke tests against the built `docs/` directory.

Targets:

- `make test`
- `make test-integration`
- `make smoke`

## Consequences

Fast tests can run in local hooks and before pushes.

## Alternatives Considered

- GitHub Actions: rejected by project constraint.
- Only manual browser testing: rejected.
