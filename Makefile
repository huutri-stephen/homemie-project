export PATH := $(PATH):$(shell go env GOPATH)/bin
include .env

MIGRATIONS_DIR = db/migrations
URL = postgres://$(HOMIE_APP_DB_USER):$(HOMIE_APP_DB_PASSWORD)@$(HOMIE_APP_DB_URL):$(HOMIE_APP_DB_PORT)/$(HOMIE_APP_DB_NAME)?sslmode=disable

.PHONY: models 

## Go modules
mod: ## Tidy go.mod & go.sum
	go mod tidy

## Migration commands
migrate-up:
	goose -dir $(MIGRATIONS_DIR) postgres "$(URL)" up

migrate-down:
	goose -dir $(MIGRATIONS_DIR) postgres "$(URL)" down

migrate-create:
	@read -p "Enter migration name: " name; \
	goose -dir $(MIGRATIONS_DIR) create $$name sql;

## Docker commands
docker-run: ## Build & start containers
	docker compose up --build -d

docker-remove: ## Stop & remove containers + volumes
	docker compose down -v

docker-stop: ## Stop containers
	docker compose stop

docker-start: ## Start stopped containers
	docker compose start

## MinIO commands
minio-up: ## Run MinIO server in Docker
	docker run -d \
	  -p 9000:9000 \
	  -p 9001:9001 \
	  --name minio \
	  -e "MINIO_ROOT_USER=admin" \
	  -e "MINIO_ROOT_PASSWORD=admin123" \
	  quay.io/minio/minio server /data --console-address ":9001"

minio-down: ## Stop & remove MinIO container
	@echo "Stopping and removing MinIO container..."
	@docker stop minio || true
	@docker rm minio || true

## Helpers
help: ## Show available commands
	@echo "Available commands:"
	@grep -E '^[a-zA-Z0-9_-]+:.*?## .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "  make %-15s %s\n", $$1, $$2}'

## SQLBoiler commands
models: ## Generate models using sqlboiler
	@(cd db && sqlboiler psql --no-tests)
