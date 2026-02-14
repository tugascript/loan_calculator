build:
	@echo "Building..."
	@go build -o main cmd/api/main.go

clean:
	@echo "Cleaning..."
	@rm -f main

run:
	@go run cmd/api/main.go

watch:
	@if command -v air > /dev/null; then \
            air; \
            echo "Watching...";\
        else \
            read -p "Go's 'air' is not installed on your machine. Do you want to install it? [Y/n] " choice; \
            if [ "$$choice" != "n" ] && [ "$$choice" != "N" ]; then \
                go install github.com/air-verse/air@latest; \
                air; \
                echo "Watching...";\
            else \
                echo "You chose not to install air. Exiting..."; \
                exit 1; \
            fi; \
        fi

sqlc:
	@echo "Generating SQLC bindings..."
	sqlc generate

build_proto:
	@echo "Generating Go stubs..."
	protoc -I ../proto --go_out=./internal/providers/proto --go_opt=paths=source_relative --go-grpc_out=./internal/providers/proto --go-grpc_opt=paths=source_relative loan_calculator/v1/calculation.proto

docker_compose_run:
	@if docker compose up --build 2>/dev/null; then \
		: ; \
	else \
		echo "Falling back to Docker Compose V1"; \
		docker-compose up --build; \
	fi

docker_compose_down:
	@if docker compose down 2>/dev/null; then \
		: ; \
	else \
		echo "Falling back to Docker Compose V1"; \
		docker-compose down; \
	fi

migrate_up:
	migrate -path ./internal/providers/database/migrations -database "postgres://postgres:postgres@127.0.0.1:5432/loan?sslmode=disable" -verbose up