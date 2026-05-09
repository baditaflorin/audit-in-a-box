# Phase 2 Substance Real-Data Audit

Date: 2026-05-09

Mode: Mode C, unchanged from Phase 1.

The v1 backend was run against ten real-world inputs with `TOOL_TIMEOUT=5s`, `MAX_MAINTAINER_PACKAGES=2`, and `MAX_UPLOAD_BYTES=12000000` to make the pass repeatable on local hardware. Scanner timeouts are counted as product behavior because a stranger will see the same class of issue on a slow machine or cold scanner database.

## Fixture Results

| Fixture | Source | v1 behavior | Desired behavior | Failure type |
| --- | --- | --- | --- | --- |
| `express-package.json` | https://raw.githubusercontent.com/expressjs/express/master/package.json | 44 deps, 0 vulns, 44 license risks, score 100 critical. | Detect npm, report dependencies, avoid critical score based only on unknown license evidence. | Wrong-but-confident. |
| `vscode-package.json` | https://raw.githubusercontent.com/microsoft/vscode/main/package.json | 162 deps, 0 vulns, 162 license risks, score 100 critical. | Treat large npm manifests as normal and prioritize real risks. | Wrong-but-confident. |
| `docker-compose-go.mod` | https://raw.githubusercontent.com/docker/compose/main/go.mod | 139 deps, Trivy parse warning on numeric license confidence. | Parse Trivy JSON shape and preserve partial evidence. | Visible but technical. |
| `kubernetes-go.mod` | https://raw.githubusercontent.com/kubernetes/kubernetes/master/go.mod | 212 deps, same Trivy parse warning, license-risk overstatement. | Handle huge Go modules deterministically with confidence. | Wrong-but-confident plus technical warning. |
| `requests-requirements-dev.txt` | https://raw.githubusercontent.com/psf/requests/main/requirements-dev.txt | 6 deps, 6 license risks, low score. | Recognize Python dev requirements and avoid unknown-license noise. | Mildly wrong. |
| `homeassistant-requirements-all.txt` | https://raw.githubusercontent.com/home-assistant/core/dev/requirements_all.txt | 1139 deps, 1 vuln, 1139 license risks, critical score. | Prioritize confirmed risks, surface scale/progress metadata, do not drown the user. | Overwhelming and wrong-confident. |
| `poetry-pyproject.toml` | https://raw.githubusercontent.com/python-poetry/poetry/main/pyproject.toml | Unknown ecosystem, 0 deps, clean summary. | Treat pyproject as Python dependency input. | Silent wrongness. |
| `pnpm-lock.yaml` | https://raw.githubusercontent.com/pnpm/pnpm/main/pnpm-lock.yaml | Unknown ecosystem, Syft found 19 deps, critical score. | Recognize pnpm lockfiles and parse lock evidence. | Confusing and wrong-confident. |
| `express-github-blob.html` | https://github.com/expressjs/express/blob/master/package.json | Paste-HTML flow failed with package JSON parse error from page chrome. | Extract visible GitHub blob code or offer a raw-file next step. | Visible but not helpful. |
| `express-package-truncated.json` | Derived from Express package.json, first third only. | 422 with `unexpected end of JSON input`. | Return recoverable `manifest_truncated` with next step. | Visible but not actionable. |

## Top 5 Logic Gaps

1. Unknown license metadata is treated as confirmed license risk.
2. Manifest detection is filename-first and misses common real formats: `pyproject.toml`, pnpm lockfiles, package locks, and pasted GitHub blob pages.
3. Scanner output parsing is brittle when upstream schemas vary.
4. Errors are technical strings rather than dependency-audit guidance.
5. Large inputs produce giant undifferentiated risk reports instead of prioritized, confidence-aware output.

## Top 3 Intuition Failures

1. Clean popular projects can be reported as critical with zero vulnerabilities.
2. Unsupported inputs can be reported as clean.
3. Pasting a normal GitHub page fails even though the manifest is visibly on the page.

## Top 3 Feels-Stupid Moments

1. User must know raw GitHub URLs work better than normal GitHub pages.
2. User must know which dependency formats are supported.
3. User must interpret scanner warnings like `cannot unmarshal number`.

## Smart Means

- Infer ecosystem and manifest shape from content, not just filename.
- Separate confirmed risk from missing evidence, and attach confidence to every inference.
- Produce a prioritized report that says what to fix first and why.
- Explain failures in domain terms with a next step.
- Stay useful on huge, partial, and pasted inputs without pretending certainty.

## Success Metrics

- At least 7 of 10 fixtures produce useful reports with no manual correction.
- Zero fixtures produce a clean report when the app failed to understand the input.
- Identical input produces byte-identical stable report projections in tests.
- Every dependency, license risk, vulnerability, maintainer signal, and ecosystem inference has confidence metadata.
- Every recoverable failure includes what failed, why, and next step.
- Huge fixtures complete without backend crash and expose scale/progress metadata.

## Out Of Scope

- No new product surfaces.
- No polish work.
- No architecture mode change.
- No SaaS accounts, auth, team features, or automatic upgrade PRs.
- No legal-grade license advice.

## Phase 2 Result

The Phase 2 fixture test suite now runs these ten inputs as the regression contract.

| Fixture | Phase 2 behavior | Status |
| --- | --- | --- |
| `express-package.json` | npm package manifest inferred with 44 dependencies; missing license metadata is an evidence anomaly, not a critical risk. | Pass |
| `vscode-package.json` | large npm manifest inferred with 150+ dependencies and non-critical risk score when no vulnerability evidence is present. | Pass |
| `docker-compose-go.mod` | Go module parsed with 120+ dependencies; numeric Trivy confidence no longer breaks scanner parsing. | Pass |
| `kubernetes-go.mod` | huge Go module parsed deterministically with 200+ dependencies. | Pass |
| `requests-requirements-dev.txt` | Python requirements parsed with markers/comments normalized. | Pass |
| `homeassistant-requirements-all.txt` | 1000+ requirements parsed without crashing and surfaced as a large dependency surface. | Pass |
| `poetry-pyproject.toml` | pyproject/Poetry dependencies parsed as Python evidence instead of clean unknown input. | Pass |
| `pnpm-lock.yaml` | multi-document pnpm lockfile parsed across all YAML documents. | Pass |
| `express-github-blob.html` | GitHub blob HTML extractor pulls the visible `package.json` code table. | Pass |
| `express-package-truncated.json` | recoverable `manifest_truncated` error includes what/why/next step. | Pass |

Pass rate: 9/10 produce useful reports; 10/10 produce either a useful report or a domain-specific recoverable failure. The remaining non-report case is intentionally the truncated file.
