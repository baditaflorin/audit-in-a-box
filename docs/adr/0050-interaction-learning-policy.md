# 0050 - Interaction Learning Policy

## Status

Accepted

## Context

V1 already stores backend URL and last report. Phase 2 should not add accounts or cross-device learning.

## Decision

Remember only local-session/browser choices: backend URL, last successful report, and user-selected sample/input state. Do not infer hidden per-user preferences beyond local storage.

## Consequences

The app feels less forgetful without creating privacy or sync concerns.

## Alternatives Considered

Server-side preferences and user accounts were rejected as out of scope.
