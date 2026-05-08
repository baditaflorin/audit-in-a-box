# 0015 - Deployment Topology

## Status

Accepted

## Context

Mode C uses a static frontend and a Docker backend.

## Decision

Serve the frontend from GitHub Pages. Run the backend through Docker Compose behind nginx on server host port `25342`.

The backend image is published to GHCR:

ghcr.io/baditaflorin/audit-in-a-box

## Consequences

The static and runtime surfaces are separated. The backend can be self-hosted or run locally.

## Alternatives Considered

- Backend serves frontend: rejected because the frontend lives on Pages.
- Single binary without Docker: rejected because scanner dependencies are easier to bundle in an image.
