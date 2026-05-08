# Deployment

Frontend:

https://baditaflorin.github.io/audit-in-a-box/

Backend image:

ghcr.io/baditaflorin/audit-in-a-box:latest

Repository:

https://github.com/baditaflorin/audit-in-a-box

## Server Prerequisites

- Docker Engine with Compose v2
- DNS pointing your backend domain to the server
- Let's Encrypt certificates mounted at `/etc/letsencrypt/live/YOUR_DOMAIN/`
- Port `25342` open for HTTPS traffic

## First Deploy

```bash
cp deploy/.env.example deploy/.env
docker compose -f deploy/docker-compose.yml pull
docker compose -f deploy/docker-compose.yml up -d
```

Edit `deploy/nginx/nginx.conf` and replace `YOUR_DOMAIN` with the real certificate directory before production use.

## Pull And Restart

```bash
docker compose -f deploy/docker-compose.yml pull
docker compose -f deploy/docker-compose.yml up -d
```

## Rollback

Use an earlier image tag:

```bash
docker compose -f deploy/docker-compose.yml pull app
docker compose -f deploy/docker-compose.yml up -d app
```

For the frontend, revert the commit that changed `docs/` and push `main`.

## TLS

Use certbot or your preferred ACME client to create certificates under:

`/etc/letsencrypt/live/YOUR_DOMAIN/fullchain.pem`

`/etc/letsencrypt/live/YOUR_DOMAIN/privkey.pem`

## GitHub Pages

Pages is served from `main` branch `/docs`.

Manual republish:

```bash
make build
git add docs
git commit -m "chore: publish pages"
git push
```
