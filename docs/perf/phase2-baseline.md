# Phase 2 Performance Baseline

Measurements from the v1 real-data audit were taken on local hardware with scanner timeout set to 5 seconds.

| Fixture | Elapsed |
| --- | ---: |
| express package.json | 10.18s |
| VS Code package.json | 8.21s |
| Docker Compose go.mod | 9.95s |
| Kubernetes go.mod | 8.78s |
| Requests dev requirements | 7.91s |
| Home Assistant requirements_all.txt | 9.54s |
| Poetry pyproject.toml | 6.93s |
| pnpm-lock.yaml | 9.48s |

The dominant cost is scanner process timeout/cold database behavior, not manifest parsing. Phase 2 should make this honest in the report: scanner partials must not look like final certainty.
