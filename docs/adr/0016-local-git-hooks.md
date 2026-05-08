# 0016 - Local Git Hooks

## Status

Accepted

## Context

The project uses local checks instead of GitHub Actions.

## Decision

Use plain `.githooks/` wired by `make install-hooks`.

Hooks:

- `pre-commit`: formatting, linting, type checking, and Gitleaks.
- `commit-msg`: Conventional Commits validation.
- `pre-push`: tests, build, and smoke.
- `post-merge` and `post-checkout`: dependency/codegen refresh placeholders.

## Consequences

Checks are transparent and runnable manually through Make targets.

## Alternatives Considered

- Lefthook: rejected to keep the bootstrap dependency smaller.
