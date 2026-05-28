.PHONY: run build test migrate-up migrate-down migrate-drop seed setup \
        docker-up docker-down docker-up-redis check-placeholders

run:
	go run ./cmd/api

build:
	CGO_ENABLED=0 go build -o bin/api ./cmd/api

test:
	go test ./... -v -count=1

migrate-up:
	go run ./cmd/migrate up

migrate-down:
	go run ./cmd/migrate down

migrate-drop:
	go run ./cmd/migrate drop

seed:
	go run ./cmd/seed

setup: migrate-up seed
	@echo "Setup complete — run: make run"

docker-up:
	docker compose up -d postgres

docker-down:
	docker compose down

docker-up-redis:
	docker compose --profile redis up -d

check-placeholders:
	pwsh -File scripts/check-placeholders.ps1
