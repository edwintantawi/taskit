#!make
include $(PWD)/.env
export $(shell sed 's/=.*//' $(PWD)/.env)

POSTGRES_DSN=postgres://$(POSTGRES_USER):$(POSTGRES_PASSWORD)@$(POSTGRES_HOST):$(POSTGRES_PORT)/$(POSTGRES_DB)?sslmode=disable

## help: show help
.PHONY: help
help:
	@echo 'Usage:'
	@sed -n 's/^##//p' ${MAKEFILE_LIST} | column -t -s ':' |  sed -e 's/^/ /'

## dev: run development server
.PHONY: dev
dev:
	@APP_ENV=dev go run cmd/main.go

## mocks: generate or update mocks
.PHONY: mocks
mocks:
	@cd internal/domain && mockery --all

## test: run tests
.PHONY: test
test:
	@go test -v -cover -coverprofile coverage.out ./...

## coverage: run test coverage
.PHONY: cover
cover: test
	@go tool cover -html=coverage.out

## migrate/new name=$1: create a new database migration
.PHONY: migrate/new
migrate/new:
	go run -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@latest create -seq -ext=.sql -dir=migrations ${name}

## migrate/up: run all database migrations up
.PHONY: migrate/up
migrate/up:
	go run -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@latest -path migrations -database $(POSTGRES_DSN) -verbose up

## migrate/down: run all database migrations down
.PHONY: migrate/down
migrate/down:
	go run -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@latest -path migrations -database $(POSTGRES_DSN) -verbose down