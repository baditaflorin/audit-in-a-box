# API

OpenAPI spec:

https://github.com/baditaflorin/audit-in-a-box/blob/main/api/openapi.yaml

Local backend URL:

http://localhost:25342

Development backend URL:

http://localhost:8080

## Health

```bash
curl http://localhost:8080/healthz
curl http://localhost:8080/readyz
```

## Tool Status

```bash
curl http://localhost:8080/api/v1/tools
```

## Audit A Manifest

```bash
curl -X POST http://localhost:8080/api/v1/audits \
  -H 'Content-Type: application/json' \
  -d '{"file_name":"package.json","content":"{\"dependencies\":{\"lodash\":\"4.17.20\"}}"}'
```

Supported v0.2 manifest inputs:

- `package.json`
- `package-lock.json`
- `pnpm-lock.yaml`
- `go.mod`
- `pyproject.toml`
- `requirements.txt`
- pasted GitHub blob HTML through `/api/v1/scrape`

Audit responses include `input`, `provenance`, per-item confidence metadata, `anomalies`, scanner warnings, and the plain-English summary.

Recoverable errors use this shape:

```json
{
  "error": {
    "code": "manifest_truncated",
    "message": "The manifest looks incomplete.",
    "why": "The parser reached the end of the file before the manifest structure was closed.",
    "next_step": "Paste or upload the full package.json file, then run the audit again.",
    "recoverable": true
  }
}
```

## Scrape Pasted HTML

```bash
curl -X POST http://localhost:8080/api/v1/scrape \
  -H 'Content-Type: application/json' \
  -d '{"html":"<pre>django==3.2.0\nrequests==2.25.1</pre>"}'
```
