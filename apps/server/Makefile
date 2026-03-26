sqlc:
	sqlc generate

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