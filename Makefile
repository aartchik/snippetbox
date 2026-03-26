include .env
export

export PROJECT_ROOT=$(shell pwd)

env-up:
	docker compose up 
env-down:
	docker compose down

