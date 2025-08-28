# File: Makefile

APP_NAME = homemie
ENV_FILE = .env
MIGRATIONS_DIR = db/migrations
URL = postgres://$(DB_USER):$(DB_PASSWORD)@$(DB_URL):$(DB_PORT)/$(DB_NAME)?sslmode=disable
include $(ENV_FILE)
export $(shell sed 's/=.*//' $(ENV_FILE))


run:
	@echo "Running $(APP_NAME)..."
	go run cmd/main.go

mod:
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

migrate-create:
	@read -p "Enter migration name: " name; \
	migrate create -ext sql -dir $(MIGRATIONS_DIR) -seq $$name

migrate-up:
	migrate -path $(MIGRATIONS_DIR) -database "$(URL)" up

migrate-down:
	migrate -path $(MIGRATIONS_DIR) -database "$(URL)" down

migrate-version:
	migrate -path $(MIGRATIONS_DIR) -database "$(URL)" version

test:
	go test ./...

db-init:
	@echo "Checking if database '$(DB_NAME)' exists..."
	@if ! PGPASSWORD=$(DB_PASSWORD) psql -U $(DB_USER) -h $(DB_HOST) -p $(DB_PORT) -tAc "SELECT 1 FROM pg_database WHERE datname='$(DB_NAME)'" | grep -q 1; then \
		echo "Creating database $(DB_NAME)..."; \
		PGPASSWORD=$(DB_PASSWORD) createdb -U $(DB_USER) -h $(DB_HOST) -p $(DB_PORT) $(DB_NAME); \
	else \
		echo "Database $(DB_NAME) already exists."; \
	fi

.PHONY: run tidy lint build dev setup-env migrate-up migrate-down migrate-create test
