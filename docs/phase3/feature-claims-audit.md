# Phase 3 Feature Claims Audit

Date: 2026-05-10

| Claim | Source | Status before Phase 3 | Evidence | Decision |
| --- | --- | --- | --- | --- |
| GitHub Pages frontend is public entrypoint | README | Shipped fully | Live Pages URL works. | Keep |
| Local/hosted Docker analyzer | README/deploy docs | Shipped fully | Dockerfile and compose exist; API runs locally. | Keep |
| Trivy, Syft, Grype, DuckDB, local LLM | README/in-app | Shipped partially | Back end orchestrates tools; DuckDB may fall back; local LLM optional. | Clarify limitations |
| User-provided manifests | README | Shipped partially | Single manifest works; batch and URL workflows absent. | Finish |
| Supported v0.2 inputs list | README/API docs | Shipped partially | Parser supports them; UI samples only cover 3 formats. | Finish UI samples |
| Pasted GitHub blob HTML | API docs | Shipped partially | Backend supports it; UI does not guide from URL/clipboard. | Finish |
| Confidence/provenance metadata | Phase 2 docs/UI | Shipped fully | Report schema and debug view carry metadata. | Keep |
| Plain-English summary | README/UI | Shipped fully | Fallback/LLM summary exists. | Keep |
| Version and commit on page | User request/UI | Shipped fully | Footer shows version and commit fallback/API lookup. | Keep |
| No analytics | privacy docs | Shipped fully | No analytics scripts found. | Keep |
| Export/report portability | Implied by audit/report product | Not shipped | No explicit export controls. | Finish |

Top mismatch: docs say multiple supported inputs, but first-time UI only makes three feel first-class.
