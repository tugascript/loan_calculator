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

COPY --from=builder /app/api /app/gateway

# Required by internal/config/config.go; override at runtime (e.g. -e or env file).
# PORT must match the value used in EXPOSE / docker run -p.
ENV PORT=8090
ENV ENV=development
ENV LOG_LEVEL=info
ENV SERVICE_NAME=gateway
ENV CALCULATION_DOMAIN=0.0.0.0:8080

EXPOSE ${PORT}

CMD ["/app/gateway"]