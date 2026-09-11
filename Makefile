.PHONY: all test up down clean generate build

all: test

test:
	./test.sh

up:
	docker compose up -d --remove-orphans

down:
	docker compose down --volumes --remove-orphans

clean: down

generate:
	sqlc generate

build:
	go build ./...
