SHELL := /bin/bash
.ONESHELL:

.PHONY: all test stop up test-api

all: test

test: test-api

stop:
	docker compose down --volumes --remove-orphans

up:
	docker compose up -d --build

test-api:
	@set -Eeuo pipefail
	trap 'docker compose down --volumes --remove-orphans || true' EXIT
	docker compose down --volumes --remove-orphans || true
	docker compose up -d --build postgres app
	docker compose exec -T app sh -lc 'sqlc generate && go build ./... && hurl --test --variable "email=ana.$$(date +%s%N)@example.com" request.hurl && go test -v ./tests/...'
	