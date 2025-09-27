include .env
export

build-app:
	@go build -o ./.bin/app ./cmd/main.go

run:build-app
	@./.bin/app

new-migrate:
	@migrate create -ext sql -dir db/migrations -seq ${name}

migrate-up:
	@migrate -database ${DB_URL} -path db/migrations up

migrate-down:
	@migrate -database ${DB_URL} -path db/migrations down 1
