.PHONY: all test up down clean generate build

all: test

test:
	./test.sh

up:
	docker compose up -d --remove-orphans

down:
	docker compose down --volumes --remove-orphans
	# Se esta eliminando la base de datos, se debe eliminar el archivo de migraciones para que no se intente ejecutar nuevamente
	# rm -f db/migrations/*.sql
	# La base de datos debe ser persistente!
clean: down

generate:
	sqlc generate

build:
	go build ./...
