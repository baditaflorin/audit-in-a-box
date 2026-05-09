# audit-in-a-box

Live site: https://baditaflorin.github.io/audit-in-a-box/

Repository: https://github.com/baditaflorin/audit-in-a-box

Support: https://www.paypal.com/paypalme/florinbadita

Static web UI plus local Docker analyzer for OSS dependency risk reports. Current release: v0.2.0.

## Quickstart

```bash
make install-hooks
make dev
make build
make test
make smoke
```

The GitHub Pages frontend is the public entrypoint. The analyzer runs as a local or hosted Docker backend so Trivy, Syft, Grype, DuckDB, and a local LLM can process user-provided manifests without putting secrets in the browser.

Supported v0.2 inputs include `package.json`, `package-lock.json`, `pnpm-lock.yaml`, `go.mod`, `pyproject.toml`, `requirements.txt`, and pasted GitHub blob HTML.

## Architecture

```mermaid
C4Container
title audit-in-a-box container view
Person(user, "Developer", "Drops package.json, go.mod, requirements.txt, or pasted HTML")
System_Boundary(pages, "GitHub Pages") {
  Container(frontend, "Static React UI", "TypeScript/Vite", "Uploads manifests, shows risk report, links to GitHub and PayPal")
}
System_Boundary(runtime, "Docker Backend") {
  Container(api, "Go API", "chi + slog + Prometheus", "Coordinates scans and risk scoring")
  ContainerDb(duckdb, "DuckDB", "CLI database", "Rollups and report-friendly tabular queries")
  Container(tools, "Scanner tools", "Trivy + Syft + Grype", "SBOM, vulnerability, and license evidence")
  Container(llm, "Local LLM", "Ollama or command", "Plain-English summary")
}
Rel(user, frontend, "Uses", "HTTPS")
Rel(frontend, api, "POST /api/v1/audits", "JSON/multipart")
Rel(api, tools, "Executes", "local process")
Rel(api, duckdb, "Queries", "local process")
Rel(api, llm, "Prompts", "local-only")
```

## Documentation

- Architecture: https://github.com/baditaflorin/audit-in-a-box/blob/main/docs/architecture.md
- API: https://github.com/baditaflorin/audit-in-a-box/blob/main/docs/api.md
- Deployment: https://github.com/baditaflorin/audit-in-a-box/blob/main/deploy/README.md
- ADRs: https://github.com/baditaflorin/audit-in-a-box/tree/main/docs/adr
- Postmortem: https://github.com/baditaflorin/audit-in-a-box/blob/main/docs/postmortem.md
- Phase 2 postmortem: https://github.com/baditaflorin/audit-in-a-box/blob/main/docs/postmortem-phase2-substance.md

## License

MIT
