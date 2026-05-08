# 0017 - Dependency Policy

## Status

Accepted

## Context

The project should avoid custom implementations where mature libraries exist.

## Decision

Use production-ready libraries:

- Go: chi, envconfig, validator, Prometheus client, goquery, x/mod, testify.
- Frontend: Vite, React, Tailwind CSS, TanStack Query, Zod, Lucide React, Vitest, Playwright.
- Scanner tools: Trivy, Syft, Grype, DuckDB CLI.

## Consequences

The implementation is smaller and easier to trust. Dependency updates must pass local hooks and vulnerability checks.

## Alternatives Considered

- Custom scanner/parsers for all ecosystems: rejected.
- Unpinned external scripts: rejected except documented Docker installer stages.
