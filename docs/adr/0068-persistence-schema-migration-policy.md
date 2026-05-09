# 0068 - Persistence Schema and Migration Policy

## Status

Accepted

## Context

Stored state needs to survive version bumps without unsafe casts.

## Decision

Use a versioned state schema in the frontend. Version 1 includes backend URL, draft input, HTML paste, URL input, last report, and batch reports. Unknown or invalid old data is ignored only after preserving the ability to export any valid report already parsed.

## Consequences

Reload and import behavior becomes deterministic and testable.

## Alternatives Considered

Silent `JSON.parse(...) as Type` was rejected as unsafe.
