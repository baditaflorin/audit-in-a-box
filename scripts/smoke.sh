#!/usr/bin/env bash
set -euo pipefail

ROOT="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
cd "$ROOT"

make build

go run ./cmd/server >/tmp/audit-in-a-box-server.log 2>&1 &
SERVER_PID=$!
cleanup() {
  kill "$SERVER_PID" >/dev/null 2>&1 || true
}
trap cleanup EXIT

for _ in {1..30}; do
  if curl -fsS http://localhost:8080/healthz >/dev/null; then
    break
  fi
  sleep 1
done

curl -fsS http://localhost:8080/readyz >/dev/null
curl -fsS http://localhost:8080/metrics >/dev/null

node <<'NODE'
const { chromium } = require('./frontend/node_modules/@playwright/test');
(async () => {
  const server = require('http').createServer((req, res) => {
    const fs = require('fs');
    const path = require('path');
    const pathname = req.url === '/' ? '/index.html' : req.url.replace('/audit-in-a-box', '');
    const file = path.join(process.cwd(), 'docs', pathname);
    fs.readFile(file, (err, data) => {
      if (err) {
        fs.readFile(path.join(process.cwd(), 'docs', 'index.html'), (fallbackErr, fallback) => {
          res.writeHead(fallbackErr ? 404 : 200, {'content-type': 'text/html'});
          res.end(fallbackErr ? 'not found' : fallback);
        });
        return;
      }
      res.writeHead(200);
      res.end(data);
    });
  }).listen(4173);
  const browser = await chromium.launch({ headless: true });
  const page = await browser.newPage();
  await page.goto('http://127.0.0.1:4173/audit-in-a-box/', { waitUntil: 'networkidle' });
  await page.getByText('audit-in-a-box').first().waitFor();
  await page.getByText('Star on GitHub').waitFor();
  await page.getByText('PayPal').waitFor();
  await page.getByText('Version').waitFor();
  await browser.close();
  server.close();
})().catch((err) => {
  console.error(err);
  process.exit(1);
});
NODE
