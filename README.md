# Loan Calculator

## Introduction

Loan Calculator is a microservice-based API that computes the **monthly repayment amount** for a loan given the principal, interest rate, and number of payments. It exposes the logic via **gRPC**, with an **HTTP REST API** provided by a gateway for easy integration. Calculation requests are persisted in **PostgreSQL**.

## Technologies

- **Language**: Go
- **RPC / API**: gRPC, gRPC-Gateway (REST)
- **Database**: PostgreSQL 18
- **ORM / SQL**: sqlc (type-safe SQL)
- **Validation**: go-playground/validator
- **Runtime**: Docker, Docker Compose

## Architecture

### Services

- **gateway**: HTTP/gRPC gateway. Acts as reverse proxy and exposes REST at `POST /v1/monthly-repayments/calculate`.
- **calculation_api** — Core gRPC service. Computes monthly repayments and stores each request in the database.

### Database

- **PostgreSQL** — Stores loan calculation requests.

## Folder structure

```bash
loan_calculator/
├── calculation_api/           # gRPC calculation service
│   ├── cmd/api/               # Application entrypoint (main.go)
│   ├── internal/
│   │   ├── config/            # Configuration and logger setup
│   │   ├── controllers/       # gRPC request handlers
│   │   ├── api_errors/        # API error mapping
│   │   ├── services/          # Business logic and DTOs
│   │   └── providers/
│   │       ├── database/      # sqlc, migrations, connection
│   │       ├── logger/        # Structured logging
│   │       └── proto/         # Generated gRPC/Protobuf code
│   ├── Dockerfile
│   ├── Makefile               # build, run, watch, sqlc, migrate_up, docker
│   ├── sqlc.yaml
│   └── example.env
├── gateway/                   # HTTP/gRPC gateway
│   ├── cmd/api/               # Application entrypoint
│   ├── internal/
│   │   ├── config/            # Configuration and logger
│   │   ├── server/            # HTTP server + gRPC gateway wiring
│   │   └── providers/
│   │       ├── logger/
│   │       └── proto/         # Generated gRPC + gateway code
│   ├── Dockerfile
│   ├── Makefile               # build, run, build_proto, docker
│   └── example.env
├── proto/                     # Shared API definition
│   ├── loan_calculator/v1/
│   │   └── calculation.proto  # CalculationService, request/response
│   └── google/api/            # google.api annotations for HTTP mapping
├── docker-compose.yaml        # PostgreSQL only (base)
├── docker-compose.services.yaml # calculation_api + gateway (includes base)
├── init-db.sh                 # Creates `loan` and `loan_test` databases
└── README.md
```

## How to run

### Prerequisites

- **Go** 1.25+
- **Docker** and **Docker Compose**
- **PostgreSQL** 18 (or use Docker)
- **migrate** 4.15.1
- Optional: **sqlc** and **protoc** (and plugins) for local code generation

### Run everything with Docker Compose (recommended)

From the **project root**:

```bash
# Start PostgreSQL, Calculation API, and Gateway (services file includes base compose)
docker compose -f docker-compose.services.yaml up -d
# Run migrations
cd calculation_api
make migrate_up
```

### Calling the API

**REST (via Gateway):**

```bash
curl -X POST http://localhost:8090/v1/monthly-repayments/calculate \
  -H "Content-Type: application/json" \
  -d '{"loan_amount": 100000, "interest_rate": 0.055, "number_of_payments": 36}'
```

Example response:

```json
{
  "id": "00000000-0000-0000-0000-000000000000",
  "monthlyRepaymentAmount": 3019.59
}
```

## License

See [LICENSE](LICENSE) in the repository.
