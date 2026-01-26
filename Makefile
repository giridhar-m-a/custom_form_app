include .env

export $(shell sed 's/=.*//' .env)

DB_URL=postgresql://$(DB_USER):$(DB_PASSWORD)@localhost:$(DB_PORT)/$(DB_NAME)?sslmode=$(DB_SSLMODE)

dev:
	docker compose up --build -d
	@echo "Waiting for database to be ready..."
	@sleep 5
	$(MAKE) migrate-up

exec-client:
	docker compose exec custom_form_client sh

exec-server:
	docker compose exec custom_form_server sh

logs-server:
	docker compose logs -f custom_form_server

logs-client:
	docker compose logs -f custom_form_client


logs-db:
	docker compose logs -f db

down:
	docker compose down

migrate-new:
	migrate create -ext sql -dir apps/server/internal/db/migrations -seq $(name)

migrate-up:
	migrate -path apps/server/internal/db/migrations -database "$(DB_URL)" -verbose up

migrate-down:
	migrate -path apps/server/internal/db/migrations -database "$(DB_URL)" -verbose down

migrate-force:
	migrate -path apps/server/internal/db/migrations -database "$(DB_URL)" force $(v)

migrate-steps:
	migrate -path apps/server/internal/db/migrations -database "$(DB_URL)" down $(s)

migrate-version:
	migrate -path apps/server/internal/db/migrations -database "$(DB_URL)" goto $(v)

migrate-all: migrate-down migrate-up
