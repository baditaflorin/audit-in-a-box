# 0006 - WASM Modules

## Status

Accepted

## Context

Mode C uses native scanners in Docker, not browser-native analysis.

## Decision

Do not ship WASM modules in v1.

## Consequences

The initial JS bundle stays smaller. Scanner execution remains on the backend, where Trivy, Syft, Grype, DuckDB, and local LLM integrations are practical.

## Alternatives Considered

- DuckDB-WASM: rejected for v1 because DuckDB is only needed for backend rollups.
- Scanner WASM ports: rejected because the scanner database and vulnerability matching workflow remains native-first.
