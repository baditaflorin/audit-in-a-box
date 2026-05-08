# API

OpenAPI spec:

api/openapi.yaml

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

## Scrape Pasted HTML

```bash
curl -X POST http://localhost:8080/api/v1/scrape \
  -H 'Content-Type: application/json' \
  -d '{"html":"<pre>django==3.2.0\nrequests==2.25.1</pre>"}'
```
