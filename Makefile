SHELL := /usr/bin/env bash

GO_CACHE ?= /tmp/go-build
COMPOSE ?= docker compose

.PHONY: help up up-build down reset logs ps test test-nocache test-processor test-aggregator seed purge health smoke smoke-receive smoke-duplicate smoke-invalid smoke-last-activity

help:
	@echo "Targets principais:"
	@echo "  make up              Sobe os servicos"
	@echo "  make up-build        Rebuilda e sobe os servicos"
	@echo "  make down            Para os servicos"
	@echo "  make reset           Para e remove volumes"
	@echo "  make logs            Acompanha logs"
	@echo "  make test            Roda todos os testes Go"
	@echo "  make test-nocache    Roda todos os testes Go sem cache de resultado"
	@echo "  make seed            Popula raw-events"
	@echo "  make purge           Limpa processed-events"
	@echo "  make health          Consulta /health"
	@echo "  make smoke           Roda smoke receive + duplicate + invalid + last-activity"

up:
	$(COMPOSE) up -d

up-build:
	$(COMPOSE) up --build --force-recreate -d

down:
	$(COMPOSE) down

reset:
	$(COMPOSE) down -v

logs:
	$(COMPOSE) logs -f

ps:
	$(COMPOSE) ps

test: test-processor test-aggregator

test-nocache:
	cd services/processor && GOCACHE=$(GO_CACHE) go test -count=1 ./...
	cd services/aggregator && GOCACHE=$(GO_CACHE) go test -count=1 ./...

test-processor:
	cd services/processor && GOCACHE=$(GO_CACHE) go test ./...

test-aggregator:
	cd services/aggregator && GOCACHE=$(GO_CACHE) go test ./...

seed:
	./scripts/seed.sh

purge:
	./scripts/purge.sh

health:
	curl -fsS http://localhost:8080/health
	@echo

smoke: smoke-receive smoke-duplicate smoke-invalid smoke-last-activity

smoke-receive:
	./scripts/smoke-receive.sh

smoke-duplicate:
	./scripts/smoke-duplicate.sh

smoke-invalid:
	./scripts/smoke-invalid.sh

smoke-last-activity:
	./scripts/smoke-last-activity.sh
