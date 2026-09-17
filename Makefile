-include .env

POSTGRES_DB ?= snippetbox
POSTGRES_USER ?= web
POSTGRES_PASSWORD ?= pass
POSTGRES_PORT ?= 5432
TEST_POSTGRES_DB ?= test_snippetbox
TEST_POSTGRES_USER ?= test_web
TEST_POSTGRES_PASSWORD ?= pass
TEST_POSTGRES_PORT ?= 5433
REDIS_PORT ?= 6380

POSTGRES_DSN ?= postgres://$(POSTGRES_USER):$(POSTGRES_PASSWORD)@127.0.0.1:$(POSTGRES_PORT)/$(POSTGRES_DB)?sslmode=disable
TEST_POSTGRES_DSN ?= postgres://$(TEST_POSTGRES_USER):$(TEST_POSTGRES_PASSWORD)@127.0.0.1:$(TEST_POSTGRES_PORT)/$(TEST_POSTGRES_DB)?sslmode=disable


export PROJECT_ROOT=$(shell pwd)

env-up:
	docker compose up 
env-down:
	docker compose down
test-env-up:
	docker compose --profile test up -d test-db redis
test-env-down:
	docker compose --profile test down
test:
	go test ./...
test-compose:
	docker compose --profile test run --rm test
test-e2e:
	./scripts/test-e2e.sh

migrate_up:
	migrate -path migrations_v2 -database "$(POSTGRES_DSN)" up

migrate_up_one:
	migrate -path migrations_v2 -database "$(POSTGRES_DSN)" up 1

migrate_down:
	migrate -path migrations_v2 -database "$(POSTGRES_DSN)" down
migrate_down_one:
	migrate -path migrations_v2 -database "$(POSTGRES_DSN)" down 1

migrate_create:
	migrate create -ext sql -dir migrations_v2 -seq $(name)
