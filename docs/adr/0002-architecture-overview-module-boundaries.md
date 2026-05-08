# 0002 - Architecture Overview And Module Boundaries

## Status

Accepted

## Context

The app needs a static UI, a runtime analyzer, scanner integrations, maintainer-health lookup, report scoring, and local-summary generation.

## Decision

Use these boundaries:

- `frontend/`: Vite React TypeScript application.
- `docs/`: built GitHub Pages output.
- `cmd/server/`: API process entrypoint.
- `internal/api/`: HTTP routing, validation, request/response mapping.
- `internal/analysis/`: audit orchestration and risk scoring.
- `internal/sbom/`: manifest parsing and Syft/Grype/Trivy integration.
- `internal/maintainer/`: registry and GitHub public API lookups.
- `internal/duckdb/`: optional DuckDB CLI-backed rollups.
- `internal/llm/`: Ollama, command, and deterministic summary adapters.
- `internal/tools/`: command execution and tool availability.
- `pkg/version/`: build version and commit metadata.

## Consequences

Modules stay small and testable. Scanner, maintainer, LLM, and HTTP concerns can evolve independently.

## Alternatives Considered

- One large server package: rejected because it would mix external-process code, HTTP handlers, and scoring logic.
