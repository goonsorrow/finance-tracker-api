SHELL := /bin/sh
.DEFAULT_GOAL := help

COMPOSE := docker compose

ROOT := $(dir $(abspath $(lastword $(MAKEFILE_LIST))))
ENV_COMMON := $(ROOT)/.env.common
ENV_DOCKER := $(ROOT)/.env.docker
ENV_LOCAL  := $(ROOT)/.env.local
ENV_COMMON_EXAMPLE := $(ROOT)/.env.common.example
ENV_DOCKER_EXAMPLE := $(ROOT)/.env.docker.example
ENV_LOCAL_EXAMPLE  := $(ROOT)/.env.local.example

APP_MAIN := cmd/main.go
APP_NAME := finance-tracker


help:
	@echo "Targets:"
	@echo "  make init 			 - initialize environment files"
	@echo "  make docker         - db+migrate+app in Docker"
	@echo "  make run            - db+migrate in Docker, app locally"
	@echo "  make db-up          - start only db in Docker"
	@echo "  make migrate-docker - run migrations in Docker (one-off)"
	@echo "  make docker-down    - stop docker compose"
	@echo "  make db-logs        - follow db logs"
	@echo "  make build          - build local binary"
	@echo "  make test           - run tests"
	@echo "  make clean          - remove bin/"

init:
	@echo "Initializing environment files..."
	@test -f $(ENV_COMMON) || cp $(ENV_COMMON_EXAMPLE) $(ENV_COMMON)
	@test -f $(ENV_DOCKER) || cp $(ENV_DOCKER_EXAMPLE) $(ENV_DOCKER)
	@test -f $(ENV_LOCAL) || cp $(ENV_LOCAL_EXAMPLE) $(ENV_LOCAL)
	@echo "Done! Environment files created from examples."

docker: 
	@$(COMPOSE) --env-file $(ENV_COMMON) --env-file $(ENV_DOCKER) up --build

docker-down:
	@$(COMPOSE) --env-file $(ENV_COMMON) --env-file $(ENV_DOCKER) down

db-up:
	@$(COMPOSE) --env-file $(ENV_COMMON) --env-file $(ENV_DOCKER) up -d db
	@echo "Waiting for database to be ready..."
	@until docker exec $$(docker ps -q --filter name=finance-tracker-db) pg_isready -U postgres; do \
		sleep 1; \
	done
	@echo "Database is ready!"
	@docker exec $$(docker ps -q --filter name=finance-tracker-db) createdb finance_db -U postgres || true
	@echo "✅ Database finance_db ready!"


db-logs:
	@$(COMPOSE) --env-file $(ENV_COMMON) --env-file $(ENV_DOCKER) logs -f db

migrate-docker:
	@$(COMPOSE) --env-file $(ENV_COMMON) --env-file $(ENV_DOCKER) run --rm migrate

run: db-up
	sleep 5
	make migrate-docker
	@set -a; . $(ENV_COMMON); . $(ENV_LOCAL); set +a; \
	DB_HOST=localhost DB_PORT=$${POSTGRES_PORT:-5436} go run $(APP_MAIN)

build:
	@go build -o bin/$(APP_NAME) $(APP_MAIN)

test:
	@go test ./...

clean:
	@rm -rf bin/

docker-reset:
	$(COMPOSE) down --volumes
	$(COMPOSE) up --build

.PHONY: help init docker docker-down docker-reset db-up db-logs migrate-docker run build test clean