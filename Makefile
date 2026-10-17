.PHONY: all test

all: test

test:
	./test.sh

stop:
	docker compose down

up:
	docker compose -d up
	go run
