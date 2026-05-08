# 0010 - GitHub Pages Publishing

## Status

Accepted

## Context

The live URL is a first-class deliverable and should work from the first commit.

## Decision

Publish GitHub Pages from `main` branch `/docs`.

The Vite build writes hashed assets into `docs/assets/`, emits `docs/index.html`, and copies `docs/404.html` for SPA fallback behavior. The `.gitignore` explicitly keeps `docs/` committed.

Live URL:

https://baditaflorin.github.io/audit-in-a-box/

## Consequences

The repository contains built frontend assets. Rollback is a normal git revert of the publishing commit.

## Alternatives Considered

- `gh-pages` branch: rejected because it adds branch juggling without meaningful benefit.
- GitHub Actions deploy: rejected because this project intentionally uses local hooks and no Actions.
