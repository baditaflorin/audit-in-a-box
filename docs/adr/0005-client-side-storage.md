# 0005 - Client-Side Storage

## Status

Accepted

## Context

The frontend needs local preferences such as backend URL and a last successful report, but v1 does not require cross-device sync.

## Decision

Use `localStorage` for backend URL and the last successful report snapshot.

## Consequences

The UI can recover from reloads and failed fetches without adding auth or a database.

## Alternatives Considered

- IndexedDB: unnecessary for v1 report sizes.
- Server persistence: rejected because v1 should avoid stored user manifests.
