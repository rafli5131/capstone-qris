.PHONY: run build tidy docker-up docker-down migrate swagger test-k6 test-k6-go test-k6-legacy

## Run locally (requires postgres + redis running)
run:
	cd go_sistem_baru && go run ./cmd/web/main.go

## Build binary
build:
	cd go_sistem_baru && CGO_ENABLED=0 go build -o bin/server ./cmd/web

## Tidy dependencies
tidy:
	cd go_sistem_baru && go mod tidy

## Generate Swagger docs (requires swag: go install github.com/swaggo/swag/cmd/swag@latest)
swagger:
	cd go_sistem_baru && swag init -g cmd/web/main.go -o docs

## Start all services via Docker Compose
docker-up:
	docker compose up --build -d

## Stop all services
docker-down:
	docker compose down -v

## Run only migrations against a running postgres
migrate-up:
	docker run --rm --network host \
		-v $(PWD)/go_sistem_baru/db/migrations:/migrations \
		migrate/migrate:v4.17.1 \
		-path=/migrations \
		-database="postgres://qris_user:qris_pass@localhost:5432/qris_db?sslmode=disable" \
		up

migrate-down:
	docker run --rm --network host \
		-v $(PWD)/go_sistem_baru/db/migrations:/migrations \
		migrate/migrate:v4.17.1 \
		-path=/migrations \
		-database="postgres://qris_user:qris_pass@localhost:5432/qris_db?sslmode=disable" \
		down -all

## Run Go vs Legacy Java comparison load test through Docker Compose
test-k6:
	docker compose --profile test run --rm k6-compare

test-k6-go:
	docker compose --profile test run --rm k6-go

test-k6-legacy:
	docker compose --profile test run --rm k6-legacy
