# 0042 - Inference Engine

## Status

Accepted

## Context

Filename-only detection misclassified `pyproject.toml`, pnpm lockfiles, and pasted GitHub HTML.

## Decision

Infer input kind and ecosystem from both filename and content. Return confidence and reason strings for ecosystem and parser decisions. Recognize package.json, package-lock, pnpm lock, go.mod, requirements.txt, pyproject.toml, and GitHub blob HTML.

## Consequences

The app can make a useful first guess and expose uncertainty instead of pretending unsupported inputs are clean.

## Alternatives Considered

Only adding more filename suffix checks was rejected because users paste content without reliable names.
