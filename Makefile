build:
	@go build -o bin/ecommerce.exe cmd/main.go

run: build
	@./bin/ecommerce

migration:
	@migrate create -ext sql -dir cmd/migrate/migrations $(filter-out $@,$(MAKECMDGOALS))

migrate-up:
	@go run cmd/migrate/main.go up

migrate-down:
	@go run cmd/migrate/main.go down