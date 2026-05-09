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

## Phase 2 Parser Check

`go test ./internal/analysis ./internal/sbom` exercises all ten real-data fixtures, including the 26k-line pnpm lockfile and the 1000+ line Home Assistant requirements file. The fixture suite completes in under one second on the local machine when scanner CLIs are not invoked, which keeps the intelligence layer fast and deterministic.

Runtime audits can still take seconds because Trivy, Syft, Grype, DuckDB, and maintainer lookups are external processes/network calls. Phase 2 now marks scanner gaps as warnings/anomalies instead of treating partial evidence as certainty.
