build:
	@echo "Building..."
	@go build -o main cmd/api/main.go

clean:
	@echo "Cleaning..."
	@rm -f main

run:
	@go run cmd/api/main.go

run_in_docker:
	@echo "Running in Docker..."
	@docker build -t gateway . && docker run -e PORT=8090 -e ENV=development -e LOG_LEVEL=info -e SERVICE_NAME=gateway -e CALCULATION_DOMAIN=0.0.0.0:8080 -p 8090:8090 gateway

build_proto:
	@echo "Generating Go stubs..."
	@protoc -I ../proto \
		--go_out=./internal/providers/proto \
		--go_opt=paths=source_relative \
		--go-grpc_out=./internal/providers/proto \
		--go-grpc_opt=paths=source_relative \
		--grpc-gateway_out=./internal/providers/proto \
		--grpc-gateway_opt=paths=source_relative loan_calculator/v1/calculation.proto