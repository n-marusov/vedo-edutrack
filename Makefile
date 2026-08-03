# VEDO EduTrack — Makefile (single entry point for all workflows)
#
# See ADR-IMPL.PROCESS.development-tooling §11 and M0.2 plan (T9).
# Windows: GNU Make via Git Bash or WSL.
# Environment: `.env` (root, gitignored) is auto-loaded if present; all
# compose variables have defaults in deploy/docker-compose.yml.

SHELL := /bin/bash

COMPOSE_FILE := deploy/docker-compose.yml
COMPOSE := docker compose -f $(COMPOSE_FILE)

# Default DB connection (overridable via .env or environment).
DATABASE_URL ?= postgres://edutrack:edutrack@localhost:5432/edutrack?sslmode=disable
export DATABASE_URL

# Build version — injected into the binary and the SPA fallback
# (ADR-DES.INFRA.dynamic-config-injection). Derived from git; override with
# VERSION=x.y.z make build.
VERSION ?= $(shell git describe --tags --always --dirty 2>/dev/null || echo dev)

# Load .env into make variables, then export them so `docker compose`
# (a child process) sees them.
-include .env
export $(shell sed 's/=.*//' .env 2>/dev/null)

.PHONY: help up down dev build test test-e2e lint format gen dev-check check migrate migrate-down hooks ci clean

help: ## Print available targets
	@echo "VEDO EduTrack — available targets:"
	@grep -E '^[a-zA-Z0-9_-]+:.*## ' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*## "}; {printf "  %-14s %s\n", $$1, $$2}'
	@echo ""
	@echo "See ADR-IMPL.PROCESS.development-tooling §11 for details."

up: ## Start the full dev stack (9 services)
	$(COMPOSE) up -d --wait

down: ## Stop and remove the dev stack (incl. volumes) — idempotent
	$(COMPOSE) down --volumes || true

dev: up ## Dev mode: stack up + hot-reload (air in backend container, Vite in frontend)
	@echo "Dev mode active — hot reload runs inside containers."
	@echo "  backend  -> http://localhost:8080   (air)"
	@echo "  frontend -> http://localhost:5173   (vite)"
	@echo "  traefik  -> http://localhost:8082   (dashboard, dev only)"

build: ## Production build check (Go binary + SPA)
	cd backend && go build -ldflags="-s -w -X vedo-edutrack/backend/internal/platform/config.Version=$(VERSION)" -o bin/vedo-edutrack ./cmd/vedo-edutrack
	cd frontend && VITE_APP_VERSION=$(VERSION) pnpm build

test: ## All tests (unit + frontend + e2e) — red scaffolds fail by design at M0.2
	cd backend && go test ./... -count=1
	cd frontend && pnpm test
	$(MAKE) test-e2e

test-e2e: ## Playwright E2E (tests/e2e/gui, M1–M10 Must-scenarios)
	cd tests/e2e/gui && pnpm install && pnpm test

lint: ## Lint both ends
	cd backend && golangci-lint run ./...
	cd frontend && pnpm lint

format: ## Auto-format both ends
	cd backend && gofmt -l -w .
	cd frontend && pnpm format

gen: ## Regenerate code (OpenAPI + sqlc) — no-op until generators are wired (M0.3)
	cd backend && go generate ./...

dev-check: ## Fast feedback loop (gate tier fast, T16)
	bash deploy/ci/run-gates.sh --tier fast

check: ## Full delivery gate set for task handoff (gate tier delivery, T16)
	bash deploy/ci/run-gates.sh --tier delivery

migrate: ## Apply DB migrations (wraps Atlas via the vedo-edutrack CLI, ADR-DES.API.cli-interface)
	cd backend && go run ./cmd/vedo-edutrack migrate up

migrate-down: ## Revert last migration
	cd backend && go run ./cmd/vedo-edutrack migrate down

hooks: ## Install lefthook git hooks
	lefthook install

ci: ## Full local CI run (mirrors GitHub Actions; delegates to the gate runner, T16)
	bash deploy/ci/run-gates.sh --tier delivery --trigger ci

clean: ## Full cleanup (stack + build artifacts) — idempotent
	$(COMPOSE) down --volumes || true
	rm -rf backend/bin/ frontend/dist/
