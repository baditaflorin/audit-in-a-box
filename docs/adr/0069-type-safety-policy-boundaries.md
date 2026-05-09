# 0069 - Type Safety Policy at Boundaries

## Status

Accepted

## Context

The audit found direct frontend casts around API and localStorage JSON. Go also uses `any` at JSON/YAML/TOML boundaries.

## Decision

Frontend external data must pass through Zod schemas before use. `unknown` is allowed only at boundary schemas. Go `any` remains allowed in JSON/TOML/YAML/provenance boundary code and is documented as intentional.

## Consequences

User state import/export and API envelopes become safer without forcing awkward Go abstractions around dynamic package-manager files.

## Alternatives Considered

Eliminating all Go `any` was rejected because package-manager lockfiles are dynamic maps at the boundary.
