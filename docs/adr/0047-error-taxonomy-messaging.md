# 0047 - Error Taxonomy And Messaging

## Status

Accepted

## Context

V1 returned raw parse errors for common user mistakes.

## Decision

Errors include `code`, `message`, `why`, `next_step`, and `recoverable`. Recoverable examples: `manifest_truncated`, `unsupported_manifest`, `html_manifest_not_found`, `scanner_partial`.

## Consequences

The frontend can show actionable guidance and preserve user input.

## Alternatives Considered

HTTP status plus raw message was rejected as insufficient.
