# Contributing

Thanks for helping improve audit-in-a-box.

## Local setup

```bash
make install-hooks
make dev
make test
make smoke
```

Use Conventional Commits for commit messages, for example `feat: add license risk scoring`.

## Pull request checklist

- Tests were added or updated for changed behavior.
- `make lint`, `make test`, and `make smoke` pass locally.
- No secrets, tokens, private keys, or private hostnames were committed.
- ADRs were updated when architecture changed.
