# 0044 - Confidence Model

## Status

Accepted

## Context

The app confused missing license metadata with confirmed risk.

## Decision

Add confidence scores and reasons to dependency, vulnerability, license, maintainer, ecosystem, and report-level inferences. Unknown evidence lowers confidence but does not automatically become high risk.

## Consequences

Reports become more honest. UI and exports can distinguish confirmed problems from incomplete evidence.

## Alternatives Considered

Binary risk flags were rejected because they create wrong-confident summaries on real data.
