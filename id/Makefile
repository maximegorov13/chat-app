include .env

migrate-create:
	@read -p "Enter migration name: " name; \
	migrate create -ext sql -dir migrations/pg -seq ${name}

migrate-up:
	migrate -path migrations/pg -database $(POSTGRES_URL) up

migrate-down:
	migrate -path migrations/pg -database $(POSTGRES_URL) down
