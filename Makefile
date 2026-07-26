include .env
DSN=mysql://$(MYSQL_USER):$(MYSQL_PASSWORD)@tcp(localhost:3306)/$(MYSQL_DATABASE)?multiStatements=true


export PROJECT_ROOT=$(shell pwd)

env-up:
	docker compose up 
env-down:
	docker compose down

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
