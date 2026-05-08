# Runbook

## Local Backend

```bash
go run ./cmd/server
```

Health:

```bash
curl http://localhost:8080/healthz
curl http://localhost:8080/readyz
```

Metrics:

```bash
curl http://localhost:8080/metrics
```

## Docker Backend

```bash
docker compose -f deploy/docker-compose.yml pull
docker compose -f deploy/docker-compose.yml up -d
```

The public HTTPS port is `25342`.

## Logs

```bash
docker compose -f deploy/docker-compose.yml logs -f app
```

Logs are JSON on stdout and include request path, status, duration, and trace id.

## Resource Sizing

Minimum:

- 2 CPU cores
- 3 GB RAM
- 4 GB free disk for scanner databases and temporary audit artifacts

Large manifests may need more disk and memory because scanners build package indexes.

## Backups

The backend is stateless. There is no application database to back up. Preserve only deployment config and any optional LLM model state managed outside this repository.

## Common Failures

- `trivy`, `syft`, or `grype` missing: use the Docker image or install the CLI locally.
- Browser CORS error: add the Pages origin to `ALLOWED_ORIGINS`.
- Mixed-content or localhost issue: use the Docker/nginx HTTPS endpoint or a browser that allows loopback calls from HTTPS pages.
