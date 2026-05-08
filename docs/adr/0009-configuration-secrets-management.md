# 0009 - Configuration And Secrets Management

## Status

Accepted

## Context

The frontend must never contain secrets. The backend needs environment-driven configuration.

## Decision

Use environment variables documented in `.env.example`. The backend reads env vars with `envconfig` and validates them at startup.

No API keys are required in v1. Optional local LLM settings point to local processes or local network endpoints.

## Consequences

Deployment is portable and secrets stay out of git. Users can run locally without provisioning accounts.

## Alternatives Considered

- Checked-in config files with secrets: rejected.
- Browser-held tokens: rejected.
