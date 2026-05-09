# 0064 - DRY Consolidation Map

## Status

Accepted

## Context

Storage, export, and API boundary logic were starting to spread through UI files.

## Decision

Create frontend modules for state storage/import/export, samples, and report exports. Keep Zod schemas as canonical frontend types. Consolidate API envelope validation in the API client.

## Consequences

`App.tsx` and `ReportView.tsx` still render the product, but they stop owning serialization details.

## Alternatives Considered

A broader frontend architecture rewrite was rejected as polish/refactor beyond Phase 3.
