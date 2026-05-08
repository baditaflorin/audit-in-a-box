# 0004 - API Contract

## Status

Accepted

## Context

Mode C needs a browser-to-backend contract. The frontend must not hand-write opaque request shapes.

## Decision

Expose a small JSON and multipart REST API under `/api/v1`.

- `GET /api/v1/tools`: scanner and optional dependency availability.
- `POST /api/v1/scrape`: accepts pasted HTML and extracts manifest candidates.
- `POST /api/v1/audits`: accepts JSON or multipart manifest content and returns a complete risk report.

The OpenAPI document lives at `api/openapi.yaml`.

## Consequences

The API surface stays small enough for v1 and can later be used to generate a typed client.

## Alternatives Considered

- GraphQL: rejected because the workflow is command-like and report-shaped.
- WebSockets: rejected because v1 audits are short enough for a request/response flow.
