# 0003 - Frontend Framework And Build Tooling

## Status

Accepted

## Context

The frontend must be static, fast, accessible, and deployable to GitHub Pages from day one.

## Decision

Use Vite, React, TypeScript strict mode, Tailwind CSS, TanStack Query, Zod, and Lucide React.

The Vite base path is `/audit-in-a-box/`. Builds write directly to `docs/` so Pages can serve the exact generated files.

## Consequences

The UI is a static artifact, has typed API payloads, and can be previewed locally with the same path behavior as Pages.

## Alternatives Considered

- Next.js static export: more machinery than needed.
- Vanilla TypeScript: less ergonomic for the report UI and state management.
