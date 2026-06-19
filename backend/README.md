# Leviosa Backend

Go modular monolith: a single binary (`cmd/app`) composed of internal domain packages (`internal/authuser`, `internal/booking`, `internal/catalog`, `internal/settings`, `internal/notification`, `internal/messaging`, `internal/common`), each following a hexagonal (ports & adapters) structure. See `CLAUDE.md` for the full architecture breakdown.

## Requirements

- Go 1.24.2
- Docker (for Postgres, Redis, RabbitMQ, and integration tests via testcontainers)

## Running locally

```bash
cp example.env development.env   # fill in required values
docker compose up                 # starts backend + Postgres, Redis, RabbitMQ
```

The API listens on `http://localhost:3500`, with a health check at `/health`.

To run the Go binary directly against already-running dependencies:

```bash
go run ./cmd/app
```

## Testing

Run from this directory (see `make test-help` for the full list):

```bash
make test-unit          # unit tests, all domains
make test-integration   # integration tests (spins up deps via testcontainers), all domains
make test-unit-<domain>        # e.g. make test-unit-authuser
make test-integration-<domain> # e.g. make test-integration-booking
make test-coverage      # HTML coverage report across all domains
```

## ENCX (field-level encryption)

Structs tagged for encryption are processed by [encx](https://github.com/hengadev/encx). See `make encx-help` for the full command list (`encx-validate`, `encx-generate`, `encx-clean`).

## Seeding data

```bash
go run ./cmd/seed
```

Configure seed data via `cmd/seed/seed_data.example.json`.
