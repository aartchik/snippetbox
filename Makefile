-include .env

MYSQL_ROOT_PASSWORD ?= rootpass
MYSQL_DATABASE ?= snippetbox
MYSQL_USER ?= web
MYSQL_PASSWORD ?= pass
DB_PORT ?= 3308
TEST_DB_PORT ?= 3307
REDIS_PORT ?= 6380

DSN=mysql://$(MYSQL_USER):$(MYSQL_PASSWORD)@tcp(localhost:$(DB_PORT))/$(MYSQL_DATABASE)?multiStatements=true


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

migrate_up:
	migrate -path migrations -database "$(DSN)" up

migrate_up_one:
	migrate -path migrations -database "$(DSN)" up 1

migrate_down:
	 migrate -path migrations -database "$(DSN)" down
migrate_down_one:
	 migrate -path migrations -database "$(DSN)" down 1

migrate_create:
	migrate create -ext sql -dir migrations -seq $(name)
