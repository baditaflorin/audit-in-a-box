# Phase 2 Substance Plan

This plan ranks substance items by impact on the real-data fixture set.

## Picklist

1. A1 parser fuzzing over real fixtures and synthetic edge cases.
2. A2 encoding and newline normalization.
3. A3 huge-input budget and test coverage.
4. A4 partial-input recovery for truncated manifests.
5. A5 adversarial JSON/HTML/YAML variants.
6. B6 auto-detect manifest structure from content.
7. B8 useful first guess from first input.
8. B9 format normalization by default.
9. C11 domain vocabulary in errors and report warnings.
10. C12 domain-aware validation for missing metadata vs confirmed risk.
11. C13 recognize package.json, package-lock, pnpm lock, go.mod, requirements, pyproject, and GitHub blob HTML.
12. C14 export/report metadata: source hash, schema version, confidence, parameters.
13. C15 dependency-domain conventions: lockfile parsing, requirements markers, Go indirect scope.
14. D16 confidence scores on ecosystem, dependency, vulnerability, license, maintainer, and report inference.
15. D17 suggested fixes for truncated/unsupported/pasted inputs.
16. D18 anomaly surfacing for unknown ecosystem, huge inputs, scanner timeouts, and partial scanner data.
17. D19 explain key inference decisions in report metadata.
18. F24 state taxonomy documentation.
19. F25 no stuck state policy for audit/request lifecycle.
20. F26 cancellation model via request context and frontend abort.
21. F27 concurrency safety for double-run behavior.
22. G28 real-data performance profile.
23. G31 cache scanner/tool availability and deterministic parse results where safe.
24. H32 actionable errors with what/why/now what.
25. H33 boundary validation for API inputs.
26. H34 recoverable vs fatal error taxonomy.
27. I35 deterministic stable report projection tests.
28. I37 debug surface with `?debug=1`.
29. I38 output provenance: source hash, schema version, app version, input kind, parameters.
30. J39 remember correction-like choices within session for backend URL and recent report state.

## Implementation Order

1. Fixtures and stable test harness.
2. Normalization and manifest inference.
3. Format parsers for pyproject, pnpm lock, package locks, and robust requirements.
4. License-risk and confidence model.
5. Actionable error taxonomy.
6. HTML extraction improvements.
7. Determinism, provenance, and performance documentation.
8. Frontend state/cancellation/debug wiring.
9. Postmortem and version bump.

## Implementation Status

Completed in v0.2.0. The fixture suite now enforces the selected substance items for all ten real-world inputs, with 9/10 useful report outputs and 10/10 graceful report-or-recoverable-error outcomes.
