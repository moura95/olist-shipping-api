include .env

migrate-up:
	migrate -database ${DB_SOURCE} -path db/migrations up

migrate-down:
	migrate -database ${DB_SOURCE} -path db/migrations down --all

migrate-create:
	@read -p "name of migration: " name; \
	migrate create -dir db/migrations -ext sql -seq $$name

down:
	docker compose -f deploy/docker-compose/docker-compose.yml down --volumes && docker volume prune -f

up:
	docker compose -f deploy/docker-compose/docker-compose.yml up -d
	sleep 5
	make migrate-up

sqlc:
	rm -rf internal/repository/
	sqlc generate

run:
	go run cmd/main.go

start:
	make up
	sleep 5
	make migrate-up
	go run cmd/main.go

restart:
	make down
	make up
	sleep 10
	make migrate-up
	go run cmd/main.go

swag:
	swag init -g cmd/main.go

# Test commands
test:
	go test -v ./...

test-unit:
	go test -v ./tests/server/service/... -run ".*Unit.*"

test-integration:
	go test -v ./tests/server/service/... -run ".*Integration.*"

test-repository:
	go test -v ./tests/server/repository/...

test-service:
	go test -v ./tests/server/service/...

.PHONY: migrate-up migrate-down migrate-create down up sqlc start run restart swag test test-unit test-integration test-repository test-service