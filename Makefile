include .env
export

# ЛОКАЛЬНЫЙ ЗАПУСК ПРИЛОЖЕНИЯ

build-app:
	@go build -o ./.bin/app ./cmd/main.go

run:build-app
	@./.bin/app

# СОЗДАНИЕ И ЛОКАЛЬНЫЙ ЗАПУСК МИГРАЦИЙ

new-migrate:
	@migrate create -ext sql -dir db/migrations -seq ${name}

migrate-up:
	@migrate -database ${DB_URL} -path db/migrations up

migrate-down:
	@migrate -database ${DB_URL} -path db/migrations down 1

# DOCKER COMPOSE

ifeq ($(IN_MEM_STORAGE),true)
docker-up:
	@docker compose up -d --build app
else
docker-up:
	@docker compose up -d --build
endif

docker-down:
	@docker compose down

# ТЕСТЫ

tests:
	@go test -cover ./...