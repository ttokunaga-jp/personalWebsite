SHELL := /bin/bash
GO ?= go
PNPM ?= pnpm
DOCKER_COMPOSE ?= docker compose
GCLOUD ?= gcloud
CLOUD_BUILD_CONFIG ?= deploy/cloudbuild/cloudbuild.yaml

.PHONY: deps deps-backend deps-frontend lint lint-backend lint-frontend test test-backend test-frontend test-perf build build-backend build-frontend fmt fmt-backend fmt-frontend ci cloudbuild up up-detached stop down clean smoke-backend db-verify db-migrate-refactor

deps: deps-backend deps-frontend

deps-backend:
	cd backend && $(GO) mod tidy

deps-frontend:
	cd frontend && { command -v corepack >/dev/null 2>&1 && corepack enable || true; }
	cd frontend && { command -v $(PNPM) >/dev/null 2>&1 && $(PNPM) install || npx pnpm@8.15.4 install; }

pnpm-install:
	cd frontend && { command -v corepack >/dev/null 2>&1 && corepack enable || true; }
	cd frontend && { command -v $(PNPM) >/dev/null 2>&1 && $(PNPM) install || npx pnpm@8.15.4 install; }

lint: lint-backend lint-frontend

lint-backend:
	cd backend && test -z "$$(gofmt -l .)" || (gofmt -l . && exit 1)
	cd backend && $(GO) vet ./...

lint-frontend:
	cd frontend && { \
		command -v $(PNPM) >/dev/null 2>&1 && $(PNPM) lint || \
		(command -v corepack >/dev/null 2>&1 && corepack pnpm lint) || \
		npx pnpm@8.15.4 lint; \
	}

test: test-backend test-frontend

test-backend:
	cd backend && $(GO) test ./...

test-frontend:
	cd frontend && { \
		command -v $(PNPM) >/dev/null 2>&1 && $(PNPM) test || \
		(command -v corepack >/dev/null 2>&1 && corepack pnpm test) || \
		npx pnpm@8.15.4 test; \
	}

test-perf:
	cd frontend && { \
		command -v $(PNPM) >/dev/null 2>&1 && $(PNPM) test:perf || \
		(command -v corepack >/dev/null 2>&1 && corepack pnpm test:perf) || \
		npx pnpm@8.15.4 test:perf; \
	}

smoke-backend:
	cd backend && ./scripts/api_smoke.sh

build: build-backend build-frontend

build-backend:
	cd backend && $(GO) build ./cmd/server

build-frontend:
	cd frontend && { \
		command -v $(PNPM) >/dev/null 2>&1 && $(PNPM) build || \
		(command -v corepack >/dev/null 2>&1 && corepack pnpm build) || \
		npx pnpm@8.15.4 build; \
	}

fmt: fmt-backend fmt-frontend

fmt-backend:
	cd backend && $(GO) fmt ./...

fmt-frontend:
	cd frontend && { \
		command -v $(PNPM) >/dev/null 2>&1 && $(PNPM) format || \
		(command -v corepack >/dev/null 2>&1 && corepack pnpm format) || \
		npx pnpm@8.15.4 format; \
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
