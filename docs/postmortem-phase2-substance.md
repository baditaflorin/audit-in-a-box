# Phase 2 Substance Postmortem

Date: 2026-05-09

## Real-Data Pass Rate

Before: 0/10 high-confidence useful reports, with 2/10 partially useful reports and multiple wrong-confident results.

After: 9/10 fixtures produce useful reports with no manual correction. The tenth fixture is deliberately truncated and now returns a recoverable `manifest_truncated` error with what/why/next step. Deterministic fixture tests pass for every report-producing fixture.

| Fixture | Before | After |
| --- | --- | --- |
| `express-package.json` | Critical from unknown licenses. | npm inferred; unknown licenses are evidence notes. |
| `vscode-package.json` | Critical from unknown licenses. | large npm manifest handled without fake criticality. |
| `docker-compose-go.mod` | Trivy numeric confidence parse warning. | scanner JSON accepts string or numeric confidence. |
| `kubernetes-go.mod` | same scanner warning and noisy license risk. | deterministic Go module parsing. |
| `requests-requirements-dev.txt` | unknown-license noise. | normalized Python requirements. |
| `homeassistant-requirements-all.txt` | 1139 fake license risks. | large dependency surface anomaly, no fake license risks. |
| `poetry-pyproject.toml` | clean unknown input. | Python/Poetry dependencies parsed. |
| `pnpm-lock.yaml` | unknown input; only partial scanner evidence. | multi-document pnpm lock parsed. |
| `express-github-blob.html` | parsed page chrome and failed. | extracts GitHub blob code table. |
| `express-package-truncated.json` | raw JSON parser error. | recoverable domain error. |

## Gaps Closed

1. Unknown license metadata is no longer treated as confirmed risk; it is surfaced as incomplete evidence.
2. Manifest inference now covers `package.json`, `package-lock.json`, `pnpm-lock.yaml`, `go.mod`, `pyproject.toml`, requirements files, and GitHub blob HTML.
3. Trivy license confidence accepts schema variation.
4. API errors now carry code, message, why, next step, and recoverability.
5. Large inputs expose scale warnings/anomalies and deterministic confidence metadata.

## Smart Behaviors

- First useful guess appears from content and filename inference.
- Confidence is attached to ecosystem, dependencies, vulnerabilities, license risks, risk score, and provenance.
- Confirmed risk is separated from missing evidence.
- Pasted GitHub pages and real lockfiles work without manual raw-URL knowledge.
- `?debug=1` exposes the internal inference/provenance surface.

## Determinism

Pass. The fixture suite parses every report-producing input twice and compares the stable JSON dependency projection. Stable audit IDs and source hashes are derived from normalized content.

## Performance

Parser-level real-data fixture tests complete in under one second locally. Full runtime audits still depend on scanner process timeouts and maintainer HTTP calls; those external waits are now presented as partial evidence rather than silent certainty.

## Surprises

The pnpm fixture was a multi-document YAML stream, so a normal single-document parse silently saw only the first small lockfile. Streaming all documents closed that gap and made the fixture far more representative.

## Still Open

1. Runtime progress streaming or job polling for scanner-heavy audits.
2. Generated OpenAPI TypeScript client instead of hand-maintained Zod schemas.
3. Richer license evidence from package registries when scanners are absent.
4. Optional backend GitHub token for maintainer-health rate limits.
5. Per-package confidence rollups in exports beyond the current JSON report.

## Honest Take

It no longer feels like a toy on the tested inputs. The main engine now infers common real manifests, avoids fake certainty, and explains recoverable failures. It can still feel slow on full runtime scans because external scanner processes dominate latency, but the app is much more honest about that cost and about incomplete evidence.
