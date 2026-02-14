# Build stage
FROM golang:1.25-alpine AS builder

WORKDIR /src

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN go build -o /app/api ./cmd/api

# Runtime stage
FROM alpine:latest

WORKDIR /app

COPY --from=builder /app/api /app/api

# Required by internal/config/config.go; override at runtime (e.g. -e or env file).
# PORT must match the value used in EXPOSE / docker run -p.
ENV PORT=8080
ENV ENV=development
ENV LOG_LEVEL=info
ENV SERVICE_NAME=calculation_api
# DATABASE_URL has no default; must be set at runtime (secret).

EXPOSE ${PORT}

CMD ["/app/api"]
