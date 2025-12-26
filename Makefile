include .env
export

MIGRATIONS_PATH=./migrations

.PHONY: migrate-up migrate-down migrate-force run-app

migrate-up:
	docker run --rm -v $(shell pwd)/migrations:/migrations --network host migrate/migrate \
		-path=/migrations/ -database "$(POSTGRES_DSN)" up

migrate-down:
	docker run --rm -v $(shell pwd)/migrations:/migrations --network host migrate/migrate \
		-path=/migrations/ -database "$(POSTGRES_DSN)" down 1

migrate-force:
	docker run --rm -v $(shell pwd)/migrations:/migrations --network host migrate/migrate \
		-path=/migrations/ -database "$(POSTGRES_DSN)" force 1

run-app:
	go run cmd/main.go