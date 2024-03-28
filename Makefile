.PHONY: install create-migration migrate

MIGRATION_NAME ?= $(shell bash -c 'read  -p "Enter Migration Name: " migrationName; echo $$migrationName')
SERVICE_NAME ?= $(shell bash -c 'read -p "Enter Service Name: " serviceName; echo $$serviceName')
DB_NAME ?= $(shell bash -c 'read  -p "Enter Db Name: " dbName; echo $$dbName')

install:
	docker-compose build
	go install goose


create-migration:
	goose -dir "./$(SERVICE_NAME)/db/migrations"  create $(MIGRATION_NAME) sql

migrate:
	goose -dir "./$(SERVICE_NAME)/db/migrations" postgres "host=localhost password=S3cret  user=citizix_user dbname=${DB_NAME} sslmode=disable" up-by-one
