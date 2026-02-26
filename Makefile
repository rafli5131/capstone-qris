.PHONY: run build tidy docker-up docker-down migrate swagger test-k6

## Run locally (requires postgres + redis running)
run:
	go run ./cmd/web/main.go

## Build binary
build:
	CGO_ENABLED=0 go build -o bin/server ./cmd/web

## Tidy dependencies
tidy:
	go mod tidy

## Generate Swagger docs (requires swag: go install github.com/swaggo/swag/cmd/swag@latest)
swagger:
	swag init -g cmd/web/main.go -o docs

## Start all services via Docker Compose
docker-up:
	docker compose up --build -d

## Stop all services
docker-down:
	docker compose down -v

## Run only migrations against a running postgres
migrate-up:
	docker run --rm --network host \
		-v $(PWD)/db/migrations:/migrations \
		migrate/migrate:v4.17.1 \
		-path=/migrations \
		-database="postgres://qris_user:qris_pass@localhost:5432/qris_db?sslmode=disable" \
		up

migrate-down:
	docker run --rm --network host \
		-v $(PWD)/db/migrations:/migrations \
		migrate/migrate:v4.17.1 \
		-path=/migrations \
		-database="postgres://qris_user:qris_pass@localhost:5432/qris_db?sslmode=disable" \
		down -all

## Run k6 load test (requires k6 installed: https://k6.io/docs/getting-started/installation/)
test-k6:
	k6 run test/k6/load_test.js

## Run k6 with output to JSON
test-k6-json:
	k6 run --out json=results.json test/k6/load_test.js
