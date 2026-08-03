# VEDO EduTrack — Makefile (single entry point for all workflows)
#
# See ADR-IMPL.PROCESS.development-tooling §11 and M0.2 plan (T9).
# Windows: GNU Make via Git Bash or WSL.
# Environment: `.env` (root, gitignored) is auto-loaded if present; all
# compose variables have defaults in deploy/docker-compose.yml.
#
# Colored output: every target prints a header line (emoji resolved by
# scripts/header.sh from an ASCII key) and unified verdicts
# (PASS/FAIL/WARN/SKIP/INFO) via scripts/verdict.sh — emoji + colors.
# Colors auto-disable when stdout is not a TTY or NO_COLOR is set (CI-safe);
# emojis are kept so verdicts stay scannable in plain logs too.
#
# IMPORTANT: keep recipe command lines ASCII-only — native GNU Make on
# Windows mangles non-ASCII characters in child-process arguments. All
# decorative glyphs (emoji, arrows) live inside scripts (header.sh,
# verdict.sh) and reach the terminal via child stdout, which passes
# through make unchanged.

# Strictness. NOTE: .ONESHELL / .SHELLFLAGS=-e are deliberately NOT used —
# recipes capture exit codes (`cmd; rc=$$?`) and would break under set -e.
.DEFAULT_GOAL := help
.DELETE_ON_ERROR:
MAKEFLAGS += --no-builtin-rules

# Tool paths (override via environment, e.g. GO=/usr/local/go/bin/go).
GO           ?= go
PNPM         ?= pnpm
GOFMT        ?= gofmt
GOLANGCI_LINT ?= golangci-lint
LEFTHOOK     ?= lefthook
DOCKER       ?= docker
GOBIN        ?= $(shell $(GO) env GOPATH)/bin

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
TEST_COMPOSE_FILE := deploy/docker-compose.test.yml

# Dev environment file (non-secret defaults, committed) — loaded for
# `make up/dev/...` and passed to compose. Override with ENV_FILE=<path> make up.
ENV_FILE ?= deploy/.env.dev
# Test environment file (non-secret defaults, committed) — passed to the test
# compose stack by `make test-up` / deploy/ci/e2e-run.sh.
TEST_ENV_FILE ?= deploy/.env.test

# Pass the env file to compose interpolation when it exists; otherwise let
# compose fall back to the process environment (exported below) and defaults.
ifneq ($(wildcard $(ENV_FILE)),)
COMPOSE_ENV := --env-file $(ENV_FILE)
endif
ifneq ($(wildcard $(TEST_ENV_FILE)),)
TEST_COMPOSE_ENV := --env-file $(TEST_ENV_FILE)
endif

COMPOSE := $(DOCKER) compose $(COMPOSE_ENV) -f $(COMPOSE_FILE)
# Test stack (ADR-IMPL.INFRA.dev-test-compose-separation): isolated compose
# project vedo-edutrack-test — no observability/traefik, separate volumes.
TEST_COMPOSE := $(DOCKER) compose $(TEST_COMPOSE_ENV) -f $(TEST_COMPOSE_FILE)

# Dev env vars exported below (make-level) would otherwise leak into the test
# stack and win over deploy/.env.test (shell env > --env-file). Unset them in
# test recipes so the test values (warn / 0 / :58080) take effect.
TEST_ENV_CLEAN := LOG_LEVEL OTEL_SAMPLING_RATIO PUBLIC_BASE_URL JWKS_URL

# Default DB connection (overridable via .env.dev or environment).
DATABASE_URL ?= postgres://edutrack:edutrack@localhost:5432/edutrack?sslmode=disable
export DATABASE_URL

# Build version — injected into the binary and the SPA fallback
# (ADR-DES.INFRA.dynamic-config-injection). Derived from git; override with
# VERSION=x.y.z make build.
VERSION ?= $(shell git describe --tags --always --dirty 2>/dev/null || echo dev)

# Load env files into make variables, then export them so `docker compose`
# (a child process) sees them. deploy/.env.dev (last) wins over legacy .env.
# NOTE: variable names are extracted with $(file <) + text functions, NOT
# $(shell sed ...) — GNU Make on Windows spawns sed directly at parse time
# (no shell metacharacters → CreateProcess fast-path) and fails when sed is
# not in the process PATH, breaking even `make help`. $(file <) spawns no
# child process, so it cannot fail this way on any platform.
-include .env
-include $(ENV_FILE)
ifneq ($(wildcard $(ENV_FILE)),)
export $(foreach _w,$(file < $(ENV_FILE)),$(if $(findstring =,$(_w)),$(firstword $(subst =, ,$(_w)))))
endif
ifneq ($(wildcard .env),)
export $(foreach _w,$(file < .env),$(if $(findstring =,$(_w)),$(firstword $(subst =, ,$(_w)))))
endif

# ────────────────────────────────────────────────────────────────────────────
# Colored output helpers
# ────────────────────────────────────────────────────────────────────────────
# Headers: $(call header,<icon-key>,<text>) — bold cyan "▶ <icon> <text>"
# line (plain when non-TTY/NO_COLOR). The icon is resolved from an ASCII key
# inside scripts/header.sh; <text> must stay ASCII-only (see top comment).
# NOTE: no leading @ inside the macros — recipes invoke them as `@$(call ...)`
# (or mid-line after `; \`) so the expansion is pure shell.
define header
	bash scripts/header.sh "$(1)" $(2)
endef

# Verdicts: $(call verdict,<STATUS>,<message>) — unified colored result line.
# Implemented in scripts/verdict.sh so the gate runner shares the same format.
# Callers pass the message already quoted: $(call verdict,PASS,"text").
define verdict
	bash scripts/verdict.sh "$(1)" $(2)
endef

.PHONY: help up down test-up test-down dev build build-frontend bench test test-e2e lint format gen dev-check check migrate migrate-down hooks ci gates gates-list gates-json docker-build docker-build-backend docker-build-local docker-build-frontend-nginx docker-build-all clean

help: ## Print available targets
	@if [ -t 1 ] && [ -z "$(NO_COLOR)" ]; then \
		printf '\033[1m%s\033[0m\n' "VEDO EduTrack - available targets:"; \
		grep -E '^[a-zA-Z0-9_%-]+:.*## ' $(firstword $(MAKEFILE_LIST)) | while IFS=':' read -r name desc; do \
			desc="$${desc#*## }"; \
			printf '  \033[1;36m%-14s\033[0m %s\n' "$$name" "$$desc"; \
		done; \
	else \
		echo "VEDO EduTrack - available targets:"; \
		grep -E '^[a-zA-Z0-9_%-]+:.*## ' $(firstword $(MAKEFILE_LIST)) | while IFS=':' read -r name desc; do \
			desc="$${desc#*## }"; \
			printf '  %-14s %s\n' "$$name" "$$desc"; \
		done; \
	fi
	@echo ""
	@echo "See ADR-IMPL.PROCESS.development-tooling, section 11, for details."

up: ## Start the full dev stack (9 services)
	@$(call header,up,"start dev stack")
	@if $(COMPOSE) up -d --wait; then \
		$(call verdict,PASS,"dev stack up (9 services)"); \
	else \
		$(call verdict,FAIL,"dev stack failed to start"); \
		exit 1; \
	fi

down: ## Stop and remove the dev stack (incl. volumes) — idempotent
	@$(call header,down,"stop dev stack")
	@$(COMPOSE) down --volumes >/dev/null 2>&1; rc=$$?; \
	if [ $$rc -eq 0 ]; then \
		$(call verdict,PASS,"dev stack stopped and cleaned"); \
	else \
		$(call verdict,SKIP,"stack was not running (nothing to stop)"); \
	fi

test-up: ## Start the TEST stack (postgres + backend + hub-mock + frontend; no observability)
	@$(call header,up,"start test stack")
	@unset $(TEST_ENV_CLEAN); \
	if $(TEST_COMPOSE) up -d --wait; then \
		$(call verdict,PASS,"test stack up (4 services)"); \
	else \
		$(call verdict,FAIL,"test stack failed to start"); \
		exit 1; \
	fi

test-down: ## Stop and remove the TEST stack (isolated project — never touches dev data)
	@$(call header,down,"stop test stack")
	@unset $(TEST_ENV_CLEAN); \
	$(TEST_COMPOSE) down --volumes >/dev/null 2>&1; rc=$$?; \
	if [ $$rc -eq 0 ]; then \
		$(call verdict,PASS,"test stack stopped and cleaned"); \
	else \
		$(call verdict,SKIP,"test stack was not running (nothing to stop)"); \
	fi

dev: up ## Dev mode: stack up + hot-reload (air in backend container, Vite in frontend)
	@$(call header,dev,"hot-reload mode")
	@$(call verdict,INFO,"hot reload runs inside containers")
	@$(call verdict,INFO,"backend  -> http://localhost:8080   (air)")
	@$(call verdict,INFO,"frontend -> http://localhost:5173   (vite)")
	@$(call verdict,INFO,"traefik  -> http://localhost:8082   (dashboard, dev only)")

build: build-frontend ## Production build check (Go binary with embedded SPA)
	@$(call header,build,"Go binary + SPA")
	@(cd backend && $(GO) build -ldflags="-s -w -X vedo-edutrack/backend/internal/platform/config.Version=$(VERSION)" -o bin/vedo-edutrack ./cmd/vedo-edutrack); b=$$?; \
	if [ $$b -eq 0 ]; then \
		$(call verdict,PASS,"backend binary with embedded SPA (version $(VERSION))"); \
	else \
		$(call verdict,FAIL,"backend build failed: $$b"); \
		exit 1; \
	fi

build-frontend: ## Build the SPA and copy dist into the backend embed path
	@$(call header,build,"frontend + embed")
	@(cd frontend && VITE_APP_VERSION=$(VERSION) $(PNPM) build); f=$$?; \
	if [ $$f -ne 0 ]; then \
		$(call verdict,FAIL,"frontend build failed: $$f"); \
		exit 1; \
	fi; \
	rm -rf backend/internal/platform/spa/frontend/dist && mkdir -p backend/internal/platform/spa/frontend && cp -r frontend/dist backend/internal/platform/spa/frontend/dist; c=$$?; \
	if [ $$c -eq 0 ]; then \
		$(call verdict,PASS,"SPA built and embedded ($(VERSION))"); \
	else \
		$(call verdict,FAIL,"embed copy failed: $$c"); \
		exit 1; \
	fi

test: ## All tests (unit + frontend + e2e) — red scaffolds fail by design at M0.2
	@$(call header,test,"unit + frontend + e2e")
	@(cd backend && $(GO) test ./... -count=1); t1=$$?; \
	(cd frontend && $(PNPM) test); t2=$$?; \
	(cd tests/e2e/gui && $(PNPM) install --ignore-workspace >/dev/null 2>&1 && $(PNPM) test); t3=$$?; \
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
	@$(call header,test-e2e,"Playwright E2E")
	@if cd tests/e2e/gui && $(PNPM) install --ignore-workspace >/dev/null 2>&1 && $(PNPM) test; then \
		$(call verdict,PASS,"e2e tests"); \
	else \
		$(call verdict,FAIL,"e2e tests failed"); \
		exit 1; \
	fi

bench: ## NFR-critical performance benchmarks (advisory perf-bench gate, T24)
	@$(call header,bench,"NFR-critical benchmarks")
	@(cd backend && \
		$(GO) test -run=^$$ -bench=. -benchtime=1x ./internal/modules/routeplanning/domain ./internal/modules/gapcoverage/domain ./internal/modules/resources/domain ./internal/modules/ontologyport/adapters/hub); b=$$?; \
	if [ $$b -eq 0 ]; then \
		$(call verdict,PASS,"benchmarks: zero failures"); \
	else \
		$(call verdict,FAIL,"benchmarks failed: $$b"); \
		exit 1; \
	fi

lint: ## Lint both ends
	@$(call header,lint,"golangci-lint + biome")
	@(cd backend && $(GOLANGCI_LINT) run ./...); b=$$?; \
	(cd frontend && $(PNPM) lint); f=$$?; \
	[ $$b -eq 0 ] && $(call verdict,PASS,"backend golangci-lint") || $(call verdict,FAIL,"backend golangci-lint (exit $$b)"); \
	[ $$f -eq 0 ] && $(call verdict,PASS,"frontend biome") || $(call verdict,FAIL,"frontend biome (exit $$f)"); \
	if [ $$b -eq 0 ] && [ $$f -eq 0 ]; then \
		$(call verdict,PASS,"lint: both ends"); \
	else \
		$(call verdict,FAIL,"lint: see failures above"); \
		exit 1; \
	fi

format: ## Auto-format both ends
	@$(call header,format,"gofmt + biome")
	@(cd backend && $(GOFMT) -l -w .); b=$$?; \
	(cd frontend && $(PNPM) format); f=$$?; \
	if [ $$b -eq 0 ] && [ $$f -eq 0 ]; then \
		$(call verdict,PASS,"format: both ends"); \
	else \
		$(call verdict,FAIL,"format: backend=$$b frontend=$$f"); \
		exit 1; \
	fi

gen: ## Regenerate code (OpenAPI via oapi-codegen + future sqlc)
	@$(call header,gen,"oapi-codegen")
	@(cd backend && \
	  OAPI=$$(command -v oapi-codegen || echo "$(GOBIN)/oapi-codegen"); \
	  if [ ! -x "$$OAPI" ]; then echo "oapi-codegen not installed; run: go install github.com/oapi-codegen/oapi-codegen/v2/cmd/oapi-codegen@v2.4.1"; exit 1; fi; \
	  $$OAPI -package api -generate types -o internal/api/types.gen.go api/openapi/v1.yaml; \
	  $$OAPI -package api -generate chi-server -exclude-operation-ids issueToken -o internal/api/server.gen.go api/openapi/v1.yaml); g=$$?; \
	if [ $$g -eq 0 ]; then \
		$(call verdict,PASS,"generated code up to date"); \
	else \
		$(call verdict,FAIL,"code generation failed"); \
		exit 1; \
	fi

dev-check: ## Fast feedback loop (gate tier fast, T16)
	@$(call header,dev-check,"gate tier fast")
	@if bash deploy/ci/run-gates.sh --tier fast; then \
		$(call verdict,PASS,"fast gates: zero blocking failures"); \
	else \
		$(call verdict,FAIL,"fast gates: blocking failures found"); \
		exit 1; \
	fi

check: ## Full delivery gate set for task handoff (gate tier delivery, T16)
	@$(call header,check,"gate tier delivery")
	@if bash deploy/ci/run-gates.sh --tier delivery; then \
		$(call verdict,PASS,"delivery gates: zero blocking failures"); \
	else \
		$(call verdict,FAIL,"delivery gates: blocking failures found"); \
		exit 1; \
	fi

migrate: ## Apply DB migrations (wraps Atlas via the vedo-edutrack CLI, ADR-DES.API.cli-interface)
	@$(call header,migrate,"Atlas up")
	@if (cd backend && $(GO) run ./cmd/vedo-edutrack migrate up); then \
		$(call verdict,PASS,"migrations applied"); \
	else \
		$(call verdict,FAIL,"migration failed"); \
		exit 1; \
	fi

migrate-down: ## Revert last migration
	@$(call header,migrate-down,"Atlas down")
	@if (cd backend && $(GO) run ./cmd/vedo-edutrack migrate down); then \
		$(call verdict,PASS,"last migration reverted"); \
	else \
		$(call verdict,FAIL,"migration revert failed"); \
		exit 1; \
	fi

hooks: ## Install lefthook git hooks
	@$(call header,hooks,"lefthook install")
	@if $(LEFTHOOK) install; then \
		$(call verdict,PASS,"lefthook hooks installed"); \
	else \
		$(call verdict,FAIL,"lefthook install failed"); \
		exit 1; \
	fi

ci: ## Full local CI run (mirrors GitHub Actions; delegates to the gate runner, T16)
	@$(call header,ci,"delivery gates (ci trigger)")
	@if bash deploy/ci/run-gates.sh --tier delivery --trigger ci; then \
		$(call verdict,PASS,"ci: zero blocking failures"); \
	else \
		$(call verdict,FAIL,"ci: blocking failures found"); \
		exit 1; \
	fi

##@ Docker — Production images

docker-build: docker-build-backend ## Build production images (unified backend + SPA binary)
	@$(call verdict,PASS,"production image built ($(VERSION))")

docker-build-all: docker-build-backend docker-build-frontend-nginx ## Build all production images (unified backend + nginx)

docker-build-backend: ## Build the unified production image (API + embedded SPA, multi-arch)
	@$(call header,docker,"backend $(VERSION)")
	@DOCKER_BUILDKIT=1 $(DOCKER) buildx build \
		--platform linux/amd64,linux/arm64 \
		--build-arg VERSION=$(VERSION) \
		-t vedo-edutrack:$(VERSION) \
		-f backend/Dockerfile .
	@$(call verdict,PASS,"unified image built ($(VERSION), amd64+arm64)")

docker-build-local: ## Build the unified image for the local arch (--load)
	@$(call header,docker,"backend $(VERSION) (local)")
	@DOCKER_BUILDKIT=1 $(DOCKER) buildx build \
		--build-arg VERSION=$(VERSION) \
		-t vedo-edutrack:$(VERSION) \
		-f backend/Dockerfile .
	@$(call verdict,PASS,"unified image built ($(VERSION), local arch)")

docker-build-frontend-nginx: ## Build production frontend image (nginx, SaaS/CDN)
	@$(call header,docker,"frontend-nginx $(VERSION)")
	@DOCKER_BUILDKIT=1 $(DOCKER) build \
		--build-arg VERSION=$(VERSION) \
		-t vedo-edutrack-nginx:$(VERSION) \
		-f frontend/Dockerfile .
	@$(call verdict,PASS,"frontend-nginx image ($(VERSION))")

##@ Gates (quality gates, T16)

gates: ## Run all delivery gates
	@$(call header,gates,"tier delivery")
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
	@$(call header,gates,"$* group")
	@if bash deploy/ci/run-gates.sh --tier delivery --group $*; then \
		$(call verdict,PASS,"$* gates: zero blocking failures"); \
	else \
		$(call verdict,FAIL,"$* gates: blocking failures found"); \
		exit 1; \
	fi

clean: ## Full cleanup (stack + build artifacts) — idempotent
	@$(call header,clean,"stack + artifacts")
	@$(COMPOSE) down --volumes >/dev/null 2>&1 || true
	@rm -rf backend/bin/ frontend/dist/
	@$(call verdict,PASS,"cleanup complete")
