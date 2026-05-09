# 0066 - Error Handling Convention

## Status

Accepted

## Context

Backend errors already use domain error envelopes. Frontend errors must remain user-actionable.

## Decision

Backend JSON errors keep `{code,message,why,next_step,recoverable}`. Frontend boundary helpers parse that envelope and show a combined what/why/next message. UI-only errors follow the same shape in plain text: what happened, why, what to do next.

## Consequences

URL CORS failures, import failures, and clipboard permission failures become recoverable, not silent dead ends.

## Alternatives Considered

Raw thrown exceptions were rejected because they fail the stranger test.
