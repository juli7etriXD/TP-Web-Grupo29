SHELL := /bin/bash
.ONESHELL:

.PHONY: all generate build run run-app up stop restart test

all: run-app

# Flujo completo de ejecución y limpieza automática (requiere solo Docker)
run-app:
	./run-app.sh

# Genera código de sqlc y compila los componentes .templ localmente
generate:
	sqlc generate
	templ generate

# Compila el proyecto en Go localmente
build: generate
	go build -v ./...

# Ejecuta el proyecto localmente (requiere PostgreSQL corriendo)
run: generate
	go run .

# Levanta los servicios (PostgreSQL y app) en segundo plano en Docker
up:
	docker compose up -d --build

# Detiene los contenedores y limpia volúmenes
stop:
	docker compose down --volumes --remove-orphans

# Reinicia los contenedores con reconstrucción
restart: stop up

# Corre las pruebas unitarias/integración dentro de Docker
test:
	@set -Eeuo pipefail
	trap 'docker compose down --volumes --remove-orphans || true' EXIT
	docker compose down --volumes --remove-orphans || true
	docker compose up -d --build postgres app
	docker compose exec -T app sh -c 'sqlc generate && templ generate && go build ./... && for i in $$(seq 1 30); do curl -s -f http://localhost:8080/ >/dev/null && break || sleep 0.5; done && go test -v ./tests/...'
