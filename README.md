# Distributor

A small Go service that listens to an Azure Service Bus queue, validates and records incoming messages, generates a simple PDF report, and uploads it to Azure Blob Storage. It uses PostgreSQL to persist received messages and to look up daily investment results used in the generated report.

## Overview

- Listens to an Azure Service Bus queue and consumes JSON messages.
- Validates basic fields and stores every received message in PostgreSQL.
- Looks up same-day results and computes a balance (gains - losses).
- Generates a 1‑page PDF with the balance.
- Uploads the PDF to Azure Blob Storage container `reports` under `/<ClientName>/<ReportName>-<timestamp>.pdf`.
- Uses a worker pool for concurrent processing and completes/dead‑letters messages accordingly.

## Architecture

- `cmd/distributor/main.go`: entrypoint; starts the message listener and graceful shutdown.
- `internal/messaging`: Azure Service Bus client and consumer loop.
- `internal/core`: main handler that validates, persists, queries, builds PDF, uploads.
- `internal/db`: PostgreSQL access via `pgxpool`.
- `internal/storage`: Azure Blob Storage upload via `azblob`.
- `internal/domain`: data types for incoming message and DB results.

## Prerequisites

- Go (recent version)
- Azure Service Bus (connection string + queue)
- Azure Blob Storage (connection string, container named `reports`)
- PostgreSQL 16 (local install or a single Docker container)

## Configuration

Provide configuration via environment variables. A `.env` file at the repo root is supported (the app will try to load it):

```
SERVICEBUS_CONNECTION_STRING=Endpoint=sb://...;SharedAccessKeyName=...;SharedAccessKey=...
SERVICEBUS_QUEUE=your-queue-name

BLOB_STORAGE_CONNECTION_STRING=DefaultEndpointsProtocol=https;AccountName=...;AccountKey=...;EndpointSuffix=core.windows.net

# When running the app locally (outside Docker), host is localhost
DATABASE_CONNECTION_STRING=postgres://admin:mypassword@localhost:5432/mydb?sslmode=disable
```

## Start PostgreSQL (Docker single container)

If you don’t have Postgres locally, you can start one with Docker and initialize the schema using the provided SQL:

```
docker volume create distributor_pgdata

docker run -d --name distributor-db \
  -e POSTGRES_USER=admin \
  -e POSTGRES_PASSWORD=mypassword \
  -e POSTGRES_DB=mydb \
  -p 127.0.0.1:5432:5432 \
  -v distributor_pgdata:/var/lib/postgresql/data \
  -v $(pwd)/internal/db/init.sql:/docker-entrypoint-initdb.d/init.sql:ro \
  postgres:16-alpine
```

This exposes Postgres on localhost:5432 and creates the `messages` and `investment_results` tables on first run.

## Run the Service (Go)

1) Create `.env` with the variables from the Configuration section (or export them in your shell).
2) Run the worker:

```
go run ./cmd/distributor
```

The process will start listening to the configured Service Bus queue and process messages as they arrive.

## Message Format

Incoming queue messages are JSON; minimal schema:

```json
{
  "clientName": "Acme Corp",
  "reportName": "Monthly_Sales_Sept",
  "email": "jane.doe@acme.com"
}
```

Processing Flow
- Persist the received message (with current timestamp) into `messages`.
- Query `investment_results` for the same email/report for the current date.
- Generate a PDF summarizing balance or a not‑found message.
- Upload PDF to Azure Blob Storage at `reports/<ClientName>/<ReportName>-<YYYYMMDD_HHmm>.pdf`.
- Complete the message on success; dead‑letter on validation/processing errors.

## Project Structure

```
.
├── cmd/
│   └── distributor/          # main entrypoint
├── internal/
│   ├── core/                 # message handler
│   ├── db/                   # DB connections and queries
│   ├── domain/               # domain models
│   ├── email/                # email validation helper
│   ├── messaging/            # service bus listener
│   └── storage/              # azure blob upload
├── internal/db/init.sql      # DB bootstrap DDL
└── README.md
```