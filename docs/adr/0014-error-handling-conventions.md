# 0014 - Error Handling Conventions

## Status

Accepted

## Context

Scanner and registry failures are expected and should degrade cleanly.

## Decision

Return structured API errors with stable `code` and `message` fields. Internal Go code wraps errors with `%w`.

Provide `internal/utils.HandleErrorOrLogWithMessages(err, errMsg, successMsg)` for the standing convention.

## Consequences

Failures can be surfaced clearly in the UI and logs without panics.

## Alternatives Considered

- Panic and recover: rejected for routine scanner failures.
