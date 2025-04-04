include .env

migrate-create:
	@read -p "Enter migration name: " name; \
	migrate create -ext sql -dir migrations/pg -seq ${name}

migrate-up:
	migrate -path migrations/pg -database $(POSTGRES_URL) up

migrate-down:
	migrate -path migrations/pg -database $(POSTGRES_URL) down

lint:
	golangci-lint run

lint-fix:
	golangci-lint run --fix

run:
	go run cmd/main.go

test:
	go test -v ./...

compose-up:
	docker compose up -d --build

compose-down:
	docker compose down
