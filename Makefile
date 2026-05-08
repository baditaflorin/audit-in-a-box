SHELL := /bin/bash
REPO := baditaflorin/audit-in-a-box
IMAGE := ghcr.io/$(REPO)
VERSION ?= 0.1.0
COMMIT := $(shell git rev-parse --short=12 HEAD 2>/dev/null || echo unknown)
DATE := $(shell date -u +"%Y-%m-%dT%H:%M:%SZ")
GO_LDFLAGS := -s -w -X github.com/baditaflorin/audit-in-a-box/pkg/version.Version=$(VERSION) -X github.com/baditaflorin/audit-in-a-box/pkg/version.Commit=$(COMMIT) -X github.com/baditaflorin/audit-in-a-box/pkg/version.Date=$(DATE)
GO := CGO_ENABLED=0 go
GO_PACKAGES := $(shell CGO_ENABLED=0 go list ./... | grep -v '/frontend/node_modules/')
GO_LINT_PATHS := ./cmd/... ./internal/... ./pkg/...

.PHONY: help install-hooks dev build data test test-integration smoke lint fmt pages-preview docker-build docker-push release compose-up compose-down clean hooks-pre-commit hooks-commit-msg hooks-pre-push hooks-post-merge hooks-post-checkout

help:
	@grep -E '^[a-zA-Z0-9_-]+:.*?## ' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "%-22s %s\n", $$1, $$2}'

install-hooks: ## wire local git hooks
	git config core.hooksPath .githooks
	chmod +x .githooks/*

dev: ## run frontend and backend locally
	@echo "Start backend with: go run ./cmd/server"
	npm --prefix frontend run dev

build: ## build backend and Pages-ready frontend into docs/
	$(GO) build -trimpath -ldflags="$(GO_LDFLAGS)" -o bin/audit-server ./cmd/server
	rm -rf docs/assets
	VITE_APP_VERSION=$(VERSION) npm --prefix frontend run build
	test -s docs/index.html

data: ## Mode C has no static data pipeline
	@echo "Mode C runtime audits do not generate committed static data artifacts."

test: ## run unit tests
	$(GO) test $(GO_PACKAGES)
	npm --prefix frontend run test

test-integration: ## run integration tests
	$(GO) test -tags=integration ./test/integration/...

smoke: ## build and run smoke tests
	./scripts/smoke.sh

lint: ## run linters and security checks
	$(GO) vet $(GO_PACKAGES)
	golangci-lint run --allow-parallel-runners $(GO_LINT_PATHS)
	npm --prefix frontend run lint
	npm --prefix frontend run format:check
	npm --prefix frontend run build
	CGO_ENABLED=0 govulncheck $(GO_PACKAGES)
	npm --prefix frontend audit --audit-level=high

fmt: ## autoformat source
	gofmt -w cmd internal pkg test
	npm --prefix frontend run format

pages-preview: ## serve docs/ locally as GitHub Pages would
	cd docs && python3 -m http.server 4173

docker-build: ## build linux/amd64 backend image
	docker buildx build --platform linux/amd64 --build-arg VERSION=$(VERSION) --build-arg COMMIT=$(COMMIT) --build-arg DATE=$(DATE) -t $(IMAGE):latest -t $(IMAGE):v$(VERSION) -t $(IMAGE):$(COMMIT) .

docker-push: ## push backend image to GHCR
	docker buildx build --platform linux/amd64 --push --build-arg VERSION=$(VERSION) --build-arg COMMIT=$(COMMIT) --build-arg DATE=$(DATE) -t $(IMAGE):latest -t $(IMAGE):v$(VERSION) -t $(IMAGE):$(COMMIT) .

release: build test smoke ## tag release and build image locally
	git tag v$(VERSION)
	$(MAKE) docker-build

compose-up: ## run production-like local stack
	test -f deploy/.env || cp deploy/.env.example deploy/.env
	docker compose -f deploy/docker-compose.yml up -d

compose-down: ## stop local stack
	docker compose -f deploy/docker-compose.yml down

clean: ## remove build artifacts
	rm -rf bin frontend/dist coverage audit-work

hooks-pre-commit:
	.githooks/pre-commit

hooks-commit-msg:
	.githooks/commit-msg .git/COMMIT_EDITMSG

hooks-pre-push:
	.githooks/pre-push

hooks-post-merge:
	.githooks/post-merge

hooks-post-checkout:
	.githooks/post-checkout
