# VEDO EduTrack — Makefile (single entry point for all workflows)
#
# See ADR-IMPL.PROCESS.development-tooling §11 and M0.2 plan (T9).
# Windows: GNU Make via Git Bash or WSL.
# Environment: `.env` (root, gitignored) is auto-loaded if present; all
# compose variables have defaults in deploy/docker-compose.yml.
#
# Colored output: every target prints a header line and unified verdicts
# (PASS/FAIL/WARN/SKIP/INFO) via scripts/verdict.sh. Colors auto-disable when
# stdout is not a TTY or NO_COLOR is set (CI-safe).

# Shell selection. POSIX recipes require bash. On Windows, native GNU make
# honours SHELL set to the Git-for-Windows bash path (forward slashes are
# required — backslash-escaped wildcard results produce broken paths).
# (See ADR-IMPL.PROCESS.development-tooling §11.)
ifeq ($(OS),Windows_NT)
  SHELL := C:/Program Files/Git/bin/bash.exe
else
  SHELL := /bin/bash
endif

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

# ────────────────────────────────────────────────────────────────────────────
# Colored output helpers
# ────────────────────────────────────────────────────────────────────────────
# Headers: $(call header,<text>) — cyan "▶ text" (plain when non-TTY/NO_COLOR).
# NOTE: no leading @ inside the macros — recipes invoke them as `@$(call ...)`
# (or mid-line after `; \`) so the expansion is pure shell.
define header
	if [ -t 1 ] && [ -z "$(NO_COLOR)" ]; then printf "\033[1m▶\033[0m \033[36m%s\033[0m\n" "$(1)"; else printf "▶ %s\n" "$(1)"; fi
endef

# Verdicts: $(call verdict,<STATUS>,<message>) — unified colored result line.
# Implemented in scripts/verdict.sh so the gate runner shares the same format.
# Callers pass the message already quoted: $(call verdict,PASS,"text").
define verdict
	bash scripts/verdict.sh "$(1)" $(2)
endef

.PHONY: help up down dev build test test-e2e lint format gen dev-check check migrate migrate-down hooks ci gates gates-list gates-json clean

help: ## Print available targets
	@if [ -t 1 ] && [ -z "$(NO_COLOR)" ]; then \
		echo "VEDO EduTrack — available targets:"; \
		grep -E '^[a-zA-Z0-9_%-]+:.*## ' $(MAKEFILE_LIST) | while IFS=':' read -r name desc; do \
			desc="$${desc#*## }"; \
			printf '  \033[1;36m%-14s\033[0m %s\n' "$$name" "$$desc"; \
		done; \
	else \
		echo "VEDO EduTrack — available targets:"; \
		grep -E '^[a-zA-Z0-9_%-]+:.*## ' $(MAKEFILE_LIST) | while IFS=':' read -r name desc; do \
			desc="$${desc#*## }"; \
			printf '  %-14s %s\n' "$$name" "$$desc"; \
		done; \
	fi
	@echo ""
	@echo "See ADR-IMPL.PROCESS.development-tooling §11 for details."

up: ## Start the full dev stack (9 services)
	@$(call header,up: start dev stack)
	@if $(COMPOSE) up -d --wait; then \
		$(call verdict,PASS,"dev stack up (9 services)"); \
	else \
		$(call verdict,FAIL,"dev stack failed to start"); \
		exit 1; \
	fi

down: ## Stop and remove the dev stack (incl. volumes) — idempotent
	@$(call header,down: stop dev stack)
	@$(COMPOSE) down --volumes >/dev/null 2>&1; rc=$$?; \
	if [ $$rc -eq 0 ]; then \
		$(call verdict,PASS,"dev stack stopped and cleaned"); \
	else \
		$(call verdict,SKIP,"stack was not running (nothing to stop)"); \
	fi

dev: up ## Dev mode: stack up + hot-reload (air in backend container, Vite in frontend)
	@$(call header,dev: hot-reload mode)
	@$(call verdict,INFO,"hot reload runs inside containers")
	@$(call verdict,INFO,"backend  -> http://localhost:8080   (air)")
	@$(call verdict,INFO,"frontend -> http://localhost:5173   (vite)")
	@$(call verdict,INFO,"traefik  -> http://localhost:8082   (dashboard, dev only)")

build: ## Production build check (Go binary + SPA)
	@$(call header,build: Go binary + SPA)
	@(cd backend && go build -ldflags="-s -w -X vedo-edutrack/backend/internal/platform/config.Version=$(VERSION)" -o bin/vedo-edutrack ./cmd/vedo-edutrack); b=$$?; \
	(cd frontend && VITE_APP_VERSION=$(VERSION) pnpm build); f=$$?; \
	if [ $$b -eq 0 ] && [ $$f -eq 0 ]; then \
		$(call verdict,PASS,"backend binary + SPA build (version $(VERSION))"); \
	else \
		$(call verdict,FAIL,"build: backend=$$b frontend=$$f"); \
		exit 1; \
	fi

test: ## All tests (unit + frontend + e2e) — red scaffolds fail by design at M0.2
	@$(call header,test: unit + frontend + e2e)
	@(cd backend && go test ./... -count=1); t1=$$?; \
	(cd frontend && pnpm test); t2=$$?; \
	(cd tests/e2e/gui && pnpm install --ignore-workspace >/dev/null 2>&1 && pnpm test); t3=$$?; \
	[ $$t1 -eq 0 ] && $(call verdict,PASS,"backend unit tests") || $(call verdict,FAIL,"backend unit tests (exit $$t1)"); \
	[ $$t2 -eq 0 ] && $(call verdict,PASS,"frontend tests") || $(call verdict,FAIL,"frontend tests (exit $$t2)"); \
	[ $$t3 -eq 0 ] && $(call verdict,PASS,"e2e tests") || $(call verdict,FAIL,"e2e tests (exit $$t3)"); \
	if [ $$t1 -eq 0 ] && [ $$t2 -eq 0 ] && [ $$t3 -eq 0 ]; then \
		$(call verdict,PASS,"test: all suites"); \
	else \
		$(call verdict,FAIL,"test: see failures above"); \
		exit 1; \
	fi

test-e2e: ## Playwright E2E (tests/e2e/gui, M1–M10 Must-scenarios)
	@$(call header,test-e2e: Playwright)
	@if cd tests/e2e/gui && pnpm install --ignore-workspace >/dev/null 2>&1 && pnpm test; then \
		$(call verdict,PASS,"e2e tests"); \
	else \
		$(call verdict,FAIL,"e2e tests failed"); \
		exit 1; \
	fi

lint: ## Lint both ends
	@$(call header,lint: golangci-lint + biome)
	@(cd backend && golangci-lint run ./...); b=$$?; \
	(cd frontend && pnpm lint); f=$$?; \
	[ $$b -eq 0 ] && $(call verdict,PASS,"backend golangci-lint") || $(call verdict,FAIL,"backend golangci-lint (exit $$b)"); \
	[ $$f -eq 0 ] && $(call verdict,PASS,"frontend biome") || $(call verdict,FAIL,"frontend biome (exit $$f)"); \
	if [ $$b -eq 0 ] && [ $$f -eq 0 ]; then \
		$(call verdict,PASS,"lint: both ends"); \
	else \
		$(call verdict,FAIL,"lint: see failures above"); \
		exit 1; \
	fi

format: ## Auto-format both ends
	@$(call header,format: gofmt + biome)
	@(cd backend && gofmt -l -w .); b=$$?; \
	(cd frontend && pnpm format); f=$$?; \
	if [ $$b -eq 0 ] && [ $$f -eq 0 ]; then \
		$(call verdict,PASS,"format: both ends"); \
	else \
		$(call verdict,FAIL,"format: backend=$$b frontend=$$f"); \
		exit 1; \
	fi

gen: ## Regenerate code (OpenAPI + sqlc) — no-op until generators are wired (M0.3)
	@$(call header,gen: go generate)
	@if (cd backend && go generate ./...); then \
		$(call verdict,PASS,"generated code up to date"); \
	else \
		$(call verdict,FAIL,"code generation failed"); \
		exit 1; \
	fi

dev-check: ## Fast feedback loop (gate tier fast, T16)
	@$(call header,dev-check: gate tier fast)
	@if bash deploy/ci/run-gates.sh --tier fast; then \
		$(call verdict,PASS,"fast gates: zero blocking failures"); \
	else \
		$(call verdict,FAIL,"fast gates: blocking failures found"); \
		exit 1; \
	fi

check: ## Full delivery gate set for task handoff (gate tier delivery, T16)
	@$(call header,check: gate tier delivery)
	@if bash deploy/ci/run-gates.sh --tier delivery; then \
		$(call verdict,PASS,"delivery gates: zero blocking failures"); \
	else \
		$(call verdict,FAIL,"delivery gates: blocking failures found"); \
		exit 1; \
	fi

migrate: ## Apply DB migrations (wraps Atlas via the vedo-edutrack CLI, ADR-DES.API.cli-interface)
	@$(call header,migrate: Atlas up)
	@if (cd backend && go run ./cmd/vedo-edutrack migrate up); then \
		$(call verdict,PASS,"migrations applied"); \
	else \
		$(call verdict,FAIL,"migration failed"); \
		exit 1; \
	fi

migrate-down: ## Revert last migration
	@$(call header,migrate-down: Atlas down)
	@if (cd backend && go run ./cmd/vedo-edutrack migrate down); then \
		$(call verdict,PASS,"last migration reverted"); \
	else \
		$(call verdict,FAIL,"migration revert failed"); \
		exit 1; \
	fi

hooks: ## Install lefthook git hooks
	@$(call header,hooks: lefthook install)
	@if lefthook install; then \
		$(call verdict,PASS,"lefthook hooks installed"); \
	else \
		$(call verdict,FAIL,"lefthook install failed"); \
		exit 1; \
	fi

ci: ## Full local CI run (mirrors GitHub Actions; delegates to the gate runner, T16)
	@$(call header,ci: delivery gates (ci trigger))
	@if bash deploy/ci/run-gates.sh --tier delivery --trigger ci; then \
		$(call verdict,PASS,"ci: zero blocking failures"); \
	else \
		$(call verdict,FAIL,"ci: blocking failures found"); \
		exit 1; \
	fi

##@ Gates (quality gates, T16)

gates: ## Run all delivery gates
	@$(call header,gates: tier delivery)
	@if bash deploy/ci/run-gates.sh --tier delivery; then \
		$(call verdict,PASS,"delivery gates: zero blocking failures"); \
	else \
		$(call verdict,FAIL,"delivery gates: blocking failures found"); \
		exit 1; \
	fi

gates-list: ## List gates grouped by category with counts
	@bash deploy/ci/run-gates.sh --tier delivery --list | bash deploy/ci/gates-list-formatted.sh

gates-json: ## Run all delivery gates, machine-readable JSON output
	@bash deploy/ci/run-gates.sh --tier delivery --out-format json

gates-%: ## Run gates for a specific group (lint|typecheck|test|coverage|gen|db|validate|security|build)
	@$(call header,gates: $* group)
	@if bash deploy/ci/run-gates.sh --tier delivery --group $*; then \
		$(call verdict,PASS,"$* gates: zero blocking failures"); \
	else \
		$(call verdict,FAIL,"$* gates: blocking failures found"); \
		exit 1; \
	fi

clean: ## Full cleanup (stack + build artifacts) — idempotent
	@$(call header,clean: stack + artifacts)
	@$(COMPOSE) down --volumes >/dev/null 2>&1 || true
	@rm -rf backend/bin/ frontend/dist/
	@$(call verdict,PASS,"cleanup complete")
