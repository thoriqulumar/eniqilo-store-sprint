include .env
GIT_REF_TAG := $(shell git describe --always --tag)
GIT_REF_COMMIT := $(shell git rev-parse --short=7 HEAD)

VERSION ?= $(GIT_REF_TAG)

.PHONY: lint
lint: ## lint all source codes
	@golangci-lint run ./...

.PHONY: run-app
run-app:
	@echo "Running the application..."
	go run main.go

.PHONY: migrate-up
migrate-up:
	migrate -path database/migrations/ -database "postgresql://$(DB_USERNAME):$(DB_PASSWORD)@$(DB_HOST):$(DB_PORT)/$(DB_NAME)?sslmode=disable" -verbose up

.PHONY: migrate-down
migrate-down:
	migrate -path database/migrations/ -database "postgresql://$(DB_USERNAME):$(DB_PASSWORD)@$(DB_HOST):$(DB_PORT)/$(DB_NAME)?sslmode=disable" -verbose down
