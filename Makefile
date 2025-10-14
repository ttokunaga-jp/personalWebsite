SHELL := /bin/bash
GO ?= go
PNPM ?= pnpm
DOCKER_COMPOSE ?= docker compose

.PHONY: deps deps-backend deps-frontend lint lint-backend lint-frontend test test-backend test-frontend build build-backend build-frontend up up-detached stop down clean

deps: deps-backend deps-frontend

deps-backend:
	cd backend && $(GO) mod tidy

deps-frontend:
	cd frontend && { command -v corepack >/dev/null 2>&1 && corepack enable || true; }
	cd frontend && { command -v $(PNPM) >/dev/null 2>&1 && $(PNPM) install || npx pnpm@8.15.4 install; }

lint: lint-backend lint-frontend

lint-backend:
	cd backend && test -z "$$(gofmt -l .)" || (gofmt -l . && exit 1)
	cd backend && $(GO) vet ./...

lint-frontend:
	cd frontend && $(PNPM) lint

test: test-backend test-frontend

test-backend:
	cd backend && $(GO) test ./...

test-frontend:
	cd frontend && $(PNPM) test

build: build-backend build-frontend

build-backend:
	cd backend && $(GO) build ./cmd/server

build-frontend:
	cd frontend && $(PNPM) build

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
