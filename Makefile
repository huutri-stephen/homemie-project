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
	go mod vendor

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

docker-build:
	docker build -t homemie .

docker-run:
	docker run -p 8080:8080 homemie

debug:
	docker run -it --rm -p 40000:40000 \
		-v $(pwd):/app \
		homemie \
		dlv debug --headless --listen=:40000 --api-version=2 --accept-multiclient ./app

minio-up:
	docker run -d \
	  -p 9000:9000 \
	  -p 9001:9001 \
	  --name minio \
	  -e "MINIO_ROOT_USER=admin" \
	  -e "MINIO_ROOT_PASSWORD=admin123" \
	  quay.io/minio/minio server /data --console-address ":9001"

minio-down:
	@echo "Stopping and removing MinIO container..."
	@docker stop minio || true
	@docker rm minio || true

.PHONY: run tidy lint build dev setup-env migrate-up migrate-down migrate-create test
