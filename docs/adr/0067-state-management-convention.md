# 0067 - State Management Convention

## Status

Accepted

## Context

The app stores backend URL and last report, but draft input and batch state were not coherent.

## Decision

Use React component state for live UI and versioned localStorage for durable user state. Persist backend URL, draft input, HTML paste, URL input, last report, and batch reports. Provide a clear-state action.

## Consequences

The app can recover from reloads without accounts or backend persistence.

## Alternatives Considered

IndexedDB was rejected because reports are small enough for localStorage in this workflow.
