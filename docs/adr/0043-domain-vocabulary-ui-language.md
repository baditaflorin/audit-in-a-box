# 0043 - Domain Vocabulary And UI Language

## Status

Accepted

## Context

V1 surfaced implementation errors such as JSON unmarshal failures.

## Decision

User-facing errors and warnings use dependency-audit language: manifest, dependency, scanner, vulnerability database, license evidence, maintainer signal, and next step.

## Consequences

Technical details may still appear in debug/provenance, but the default message tells the user what happened and what to do.

## Alternatives Considered

Returning raw errors was rejected because it pushes backend details onto users.
