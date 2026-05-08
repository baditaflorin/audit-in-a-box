# 0001 - Deployment Mode

## Status

Accepted

## Context

audit-in-a-box must accept user-supplied dependency manifests and produce SBOMs, vulnerability results, license risks, maintainer-health signals, and a plain-English risk summary. The requested toolchain includes Trivy, Syft, Grype, DuckDB, a paste-HTML scraper, and a local LLM.

GitHub Pages remains the preferred public surface, but Pages cannot execute native scanner binaries, maintain vulnerability databases, run DuckDB CLI queries, or call a local LLM process.

## Decision

Use Mode C: GitHub Pages frontend plus Docker backend.

The frontend is static and served from `main` branch `/docs` at:

https://baditaflorin.github.io/audit-in-a-box/

The backend is an API-only Go service packaged as a Docker image. It runs Trivy, Syft, Grype, DuckDB, and a local LLM adapter on a machine controlled by the user or operator.

## Consequences

- The public site remains static and cheap to host.
- Runtime scanning stays out of the browser and can use mature native tools.
- Users must run or point to a backend for real audits.
- CORS and backend URL configuration are required.

## Alternatives Considered

- Mode A: rejected because native scanners, vulnerability databases, DuckDB, and local LLM execution are not realistic as a complete browser-only v1.
- Mode B: rejected because audits depend on arbitrary user manifests at runtime, not only prebuilt static artifacts.
