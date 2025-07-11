# File: Makefile

APP_NAME = mihome
ENV_FILE = .env

run:
	@echo "Running $(APP_NAME)..."
	go run cmd/main.go

tidy:
	go mod tidy

lint:
	@echo "Running golangci-lint..."
	golangci-lint run

build:
	@echo "Building $(APP_NAME)..."
	go build -o bin/$(APP_NAME) cmd/main.go

dev:
	air

setup-env:
	@if [ ! -f $(ENV_FILE) ]; then cp .env.example $(ENV_FILE); fi

migrate-up:
	migrate -path db/migrations -database postgres://$$DB_USER:$$DB_PASSWORD@$$DB_HOST:$$DB_PORT/$$DB_NAME?sslmode=disable up

migrate-down:
	migrate -path db/migrations -database postgres://$$DB_USER:$$DB_PASSWORD@$$DB_HOST:$$DB_PORT/$$DB_NAME?sslmode=disable down

migrate-create:
	@read -p "Enter migration name: " name; \
	migrate create -ext sql -dir db/migrations -seq $$name

test:
	go test ./...

.PHONY: run tidy lint build dev setup-env migrate-up migrate-down migrate-create test
