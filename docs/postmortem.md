# Postmortem

## What Was Built

audit-in-a-box now has a GitHub Pages frontend, a Dockerized Go backend, scanner orchestration for Trivy, Syft, and Grype, DuckDB rollup support, maintainer-health lookup, license risk scoring, local LLM summary adapters, docs, hooks, tests, and deployment artifacts.

## Was Mode C Correct?

Yes. The core feature depends on native scanner binaries, vulnerability databases, filesystem workspaces, DuckDB CLI execution, and optional local LLM access. Mode A would have made the most important parts fake or impractical. Mode B would only work for precomputed datasets, not arbitrary manifests supplied by a user at runtime.

## What Worked

- Keeping GitHub Pages as the public surface made the live URL simple.
- Calling battle-tested scanner CLIs kept the backend small.
- The local LLM adapter can use Ollama, any stdin/stdout command, or a deterministic fallback.

## What Did Not Work

- GitHub Pages cannot provide the runtime headers needed for heavy WASM scanner work.
- DuckDB CLI was not installed locally, so local smoke tests use the Go fallback unless DuckDB is installed or the Docker image is used.

## Surprises

- Scanner CLIs are available locally, which made local backend development practical.
- The Pages API required the nested JSON payload for source configuration.

## Accepted Tech Debt

- The frontend uses hand-written Zod schemas rather than generated OpenAPI types.
- Maintainer health uses public unauthenticated APIs and may hit rate limits.
- The first report flow is synchronous; large manifests may need job polling later.

## Next Improvements

1. Add a queued audit job API with progress events.
2. Generate TypeScript API types from `api/openapi.yaml`.
3. Add optional GitHub token support on the backend for higher maintainer-metadata rate limits.

## Estimate Versus Time

The bootstrap scope is closer to a multi-day product slice than a single setup task. The implemented v1 prioritizes an end-to-end working path and clear extension points over exhaustive ecosystem coverage.
