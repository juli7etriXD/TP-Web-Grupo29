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
	docker compose exec -T app sh -c 'sqlc generate && go build ./... && for i in $$(seq 1 30); do curl -s -f http://localhost:8080/libros >/dev/null && break || sleep 0.5; done && hurl --test --variable "email=ana.$$(date +%s%N)@example.com" request.hurl && go test -v ./tests/...'
	