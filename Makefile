SHELL := /bin/bash
GO ?= go
PNPM ?= pnpm
PNPM_VERSION ?= 10.20.0
DOCKER_COMPOSE ?= docker compose
GCLOUD ?= gcloud
CLOUD_BUILD_CONFIG ?= deploy/cloudbuild/cloudbuild.yaml

.PHONY: deps deps-backend deps-frontend lint lint-backend lint-frontend test test-backend test-frontend test-perf build build-backend build-frontend fmt fmt-backend fmt-frontend ci cloudbuild up up-detached stop down clean smoke-backend db-verify db-migrate-refactor ops-validate

deps: deps-backend deps-frontend

deps-backend:
	cd backend && $(GO) mod tidy

deps-frontend:
	cd frontend && { command -v corepack >/dev/null 2>&1 && corepack enable || true; }
	cd frontend && { command -v $(PNPM) >/dev/null 2>&1 && $(PNPM) install || npx pnpm@$(PNPM_VERSION) install; }

pnpm-install:
	cd frontend && { command -v corepack >/dev/null 2>&1 && corepack enable || true; }
	cd frontend && { command -v $(PNPM) >/dev/null 2>&1 && $(PNPM) install || npx pnpm@$(PNPM_VERSION) install; }

lint: lint-backend lint-frontend

lint-backend:
	cd backend && test -z "$$(gofmt -l .)" || (gofmt -l . && exit 1)
	cd backend && $(GO) vet ./...

lint-frontend:
	cd frontend && { \
		command -v $(PNPM) >/dev/null 2>&1 && $(PNPM) lint || \
		(command -v corepack >/dev/null 2>&1 && corepack pnpm lint) || \
		npx pnpm@$(PNPM_VERSION) lint; \
	}

test: test-backend test-frontend

test-backend:
	cd backend && $(GO) test ./...

test-frontend:
	cd frontend && { \
		command -v $(PNPM) >/dev/null 2>&1 && $(PNPM) test || \
		(command -v corepack >/dev/null 2>&1 && corepack pnpm test) || \
		npx pnpm@$(PNPM_VERSION) test; \
	}

test-perf:
	cd frontend && { \
		command -v $(PNPM) >/dev/null 2>&1 && $(PNPM) test:perf || \
		(command -v corepack >/dev/null 2>&1 && corepack pnpm test:perf) || \
		npx pnpm@$(PNPM_VERSION) test:perf; \
	}

smoke-backend:
	cd backend && SMOKE_BACKEND_SPAWN_LOCAL=1 ./scripts/api_smoke.sh

build: build-backend build-frontend

build-backend:
	cd backend && $(GO) build ./cmd/server

build-frontend:
	cd frontend && { \
		command -v $(PNPM) >/dev/null 2>&1 && $(PNPM) build || \
		(command -v corepack >/dev/null 2>&1 && corepack pnpm build) || \
		npx pnpm@$(PNPM_VERSION) build; \
	}

fmt: fmt-backend fmt-frontend

fmt-backend:
	cd backend && $(GO) fmt ./...

fmt-frontend:
	cd frontend && { \
		command -v $(PNPM) >/dev/null 2>&1 && $(PNPM) format || \
		(command -v corepack >/dev/null 2>&1 && corepack pnpm format) || \
		npx pnpm@$(PNPM_VERSION) format; \
	}

ci: lint test build

cloudbuild:
ifndef CLOUD_BUILD_SUBSTITUTIONS
	$(error CLOUD_BUILD_SUBSTITUTIONS is not set. Provide substitutions string, e.g. CLOUD_BUILD_SUBSTITUTIONS="_ENV=staging,_REGION=asia-northeast1,_ARTIFACT_REPO=personal-website,_BACKEND_SERVICE=personal-website-api,_FRONTEND_SERVICE=personal-website-frontend")
endif
	$(GCLOUD) builds submit --config $(CLOUD_BUILD_CONFIG) --substitutions $(CLOUD_BUILD_SUBSTITUTIONS)

up:
	$(DOCKER_COMPOSE) up --build

up-detached:
	$(DOCKER_COMPOSE) up --build -d

stop:
	$(DOCKER_COMPOSE) stop

down:
	$(DOCKER_COMPOSE) down --remove-orphans

clean:
	rm -rf backend/bin backend/coverage.out frontend/apps/*/dist

db-verify:
	./scripts/db/apply_migrations.sh

db-migrate-refactor:
	cd backend && $(GO) run ./cmd/tools/contentmodelrefactor --dsn "$${APP_DATABASE_DSN:?set APP_DATABASE_DSN}" --dry-run

ops-validate:
	bash ./scripts/ops/validate.sh
