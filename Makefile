include .env

export $(shell sed 's/=.*//' .env)

DB_URL=postgresql://$(DB_USER):$(DB_PASSWORD)@localhost:$(DB_PORT)/$(DB_NAME)?sslmode=$(DB_SSLMODE)

dev:
	docker-compose up --build -d

exec:
	docker-compose exec app sh

logs-app:
	docker-compose logs -f app

logs-db:
	docker-compose logs -f db

down:
	docker-compose down

migrate-new:
	migrate create -ext sql -dir internal/db/migrations -seq $(name)

migrate-up:
	migrate -path internal/db/migrations -database "$(DB_URL)" -verbose up

migrate-down:
	migrate -path internal/db/migrations -database "$(DB_URL)" -verbose down
  
sqlc:
	sqlc generate

print-db-url:
	@echo $(DB_URL)

format:
	go fmt ./...

# Swagger documentation
.PHONY: swagger
swagger:
	swag init -g cmd/server/main.go -o docs/ --parseDependency --parseInternal

.PHONY: swagger-serve
swagger-serve: swagger
	@echo "Swagger UI available at: http://localhost:8080/swagger/index.html"
	go run cmd/server/main.go

.PHONY: swagger-clean
swagger-clean:
	rm -rf docs/

