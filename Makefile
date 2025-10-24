SHELL := /bin/bash
GO ?= go
PNPM ?= pnpm
DOCKER_COMPOSE ?= docker compose

.PHONY: deps deps-backend deps-frontend lint lint-backend lint-frontend test test-backend test-frontend build build-backend build-frontend up up-detached stop down clean smoke-backend

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
