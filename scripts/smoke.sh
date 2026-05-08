#!/usr/bin/env bash
set -euo pipefail

ROOT="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
cd "$ROOT"

make build

SERVER_ADDR=:18080 CGO_ENABLED=0 ./bin/audit-server >/tmp/audit-in-a-box-server.log 2>&1 &
SERVER_PID=$!
cleanup() {
  kill "$SERVER_PID" >/dev/null 2>&1 || true
}
trap cleanup EXIT

for _ in {1..30}; do
  if curl -fsS http://localhost:18080/healthz >/dev/null 2>&1; then
    READY=1
    break
  fi
  sleep 1
done
test "${READY:-0}" = "1"

curl -fsS http://localhost:18080/readyz >/dev/null
curl -fsS http://localhost:18080/metrics >/dev/null

node <<'NODE'
const { chromium } = require('./frontend/node_modules/@playwright/test');
(async () => {
  const server = require('http').createServer((req, res) => {
    const fs = require('fs');
    const path = require('path');
    const pathname = req.url === '/' ? '/index.html' : req.url.replace('/audit-in-a-box', '');
    const file = path.join(process.cwd(), 'docs', pathname);
    const contentType = file.endsWith('.js')
      ? 'application/javascript'
      : file.endsWith('.css')
        ? 'text/css'
        : file.endsWith('.webmanifest')
          ? 'application/manifest+json'
          : 'text/html';
    fs.readFile(file, (err, data) => {
      if (err) {
        fs.readFile(path.join(process.cwd(), 'docs', 'index.html'), (fallbackErr, fallback) => {
          res.writeHead(fallbackErr ? 404 : 200, {'content-type': 'text/html'});
          res.end(fallbackErr ? 'not found' : fallback);
        });
        return;
      }
      res.writeHead(200, {'content-type': contentType});
      res.end(data);
    });
  }).listen(4173);
  const browser = await chromium.launch({ headless: true });
  const page = await browser.newPage();
  await page.goto('http://127.0.0.1:4173/audit-in-a-box/', { waitUntil: 'networkidle' });
  await page.getByText('audit-in-a-box').first().waitFor();
  await page.getByText('Star on GitHub').waitFor();
  await page.getByText('PayPal').waitFor();
  await page.getByText(/^Version 0\.1\.0$/).waitFor();
  await browser.close();
  server.close();
})().catch((err) => {
  console.error(err);
  process.exit(1);
});
NODE
