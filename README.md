# Game Event Service

Backend-сервис для приёма и асинхронной обработки событий игровой платформы.

Сервис принимает игровые события по HTTP, сохраняет их в PostgreSQL, публикует в RabbitMQ и обрабатывает через пул worker'ов на goroutines и channels.

## Features

- REST API на Go
- PostgreSQL для хранения событий и игроков
- Redis для кеширования игроков
- RabbitMQ для асинхронной обработки событий
- Worker pool на goroutines и channels
- Retry с поддержкой context и timeout
- Structured logging через `log/slog`
- Request ID для трассировки HTTP-запросов
- Prometheus metrics
- Nginx reverse proxy
- Docker и Docker Compose
- Graceful shutdown
- Unit tests
- Race detector
- Go benchmarks

## Architecture

```text
                    ┌─────────────┐
                    │    Client   │
                    └──────┬──────┘
                           │
                           ▼
                    ┌─────────────┐
                    │    Nginx    │
                    └──────┬──────┘
                           │
                           ▼
                 ┌──────────────────┐
                 │    HTTP Handler  │
                 └────────┬─────────┘
                          │
              ┌───────────┴───────────┐
              ▼                       ▼
       ┌─────────────┐         ┌─────────────┐
       │   Service   │         │    Redis    │
       └──────┬──────┘         └─────────────┘
              │
       ┌──────┴──────┐
       ▼             ▼
┌────────────┐ ┌────────────┐
│ PostgreSQL │ │ RabbitMQ   │
└────────────┘ └─────┬──────┘
                     │
                     ▼
              ┌─────────────┐
              │   Consumer  │
              └──────┬──────┘
                     │
                     ▼
              ┌─────────────┐
              │ Worker Pool │
              │ goroutines  │
              │ + channels  │
              └─────────────┘
```

## Tech Stack

- Go
- PostgreSQL
- Redis
- RabbitMQ
- Nginx
- Docker
- Docker Compose
- Prometheus
- `log/slog`

## Project Structure

```text
game-event-service/
├── cmd/
│   └── server/
│       └── main.go
├── internal/
│   ├── cache/
│   ├── config/
│   ├── database/
│   ├── handler/
│   ├── logger/
│   ├── metrics/
│   ├── middleware/
│   ├── model/
│   ├── queue/
│   ├── repository/
│   ├── retry/
│   ├── service/
│   └── worker/
├── migrations/
├── .dockerignore
├── .env.example
├── .gitignore
├── Dockerfile
├── docker-compose.yml
├── go.mod
├── go.sum
└── nginx.conf
```

## API

### Healthcheck

```http
GET /health
```

Response:

```json
{
  "status": "ok"
}
```

### Create event

```http
POST /api/v1/events
Content-Type: application/json
```

Request:

```json
{
  "player_id": 1,
  "type": "purchase.created",
  "payload": "sword_01"
}
```

Response:

```json
{
  "player_id": 1,
  "type": "purchase.created",
  "payload": "sword_01"
}
```

### Get player

```http
GET /api/v1/players/{id}
```

The service uses Redis as a cache layer before querying PostgreSQL.

## Event Processing

Event creation follows this flow:

```text
HTTP request
     │
     ▼
Handler
     │
     ▼
Service
     │
     ├──────────────► PostgreSQL
     │
     ▼
RabbitMQ
     │
     ▼
Consumer
     │
     ▼
Worker Pool
     │
     ▼
Event processing
```

Publishing to RabbitMQ uses retry logic with:

- 3 attempts
- 200 ms delay between attempts
- 2 second timeout
- `context.Context` cancellation support

## Logging

The service uses structured JSON logging through Go's `log/slog`.

HTTP requests include:

- method
- path
- status
- duration
- request ID

The request ID is returned in the `X-Request-ID` response header and can also be supplied by the client.

## Metrics

Prometheus metrics are exposed at:

```http
GET /metrics
```

The service records HTTP request counters, HTTP errors, request duration and processed events.

Examples of exposed metrics:

```text
game_event_http_requests_total
game_event_http_errors_total
game_event_http_request_duration_seconds
game_event_events_processed_total
```

## Configuration

Copy the example environment file:

```bash
cp .env.example .env
```

Local development configuration is loaded from `.env`.

The main configuration values are:

```env
PORT=8080
APP_ENV=development
DATABASE_URL=...
REDIS_ADDR=localhost:6379
RABBITMQ_URL=...
```

The `.env` file is ignored by Git.

## Running with Docker Compose

Start the infrastructure and application:

```bash
docker compose up -d --build
```

Check running containers:

```bash
docker compose ps
```

The application is available through Nginx:

```text
http://localhost:8080
```

Stop the project:

```bash
docker compose down
```

To remove the PostgreSQL volume as well:

```bash
docker compose down -v
```

## Running Locally

Start PostgreSQL, Redis and RabbitMQ, configure `.env`, then run:

```bash
go run ./cmd/server
```

## Testing

Run all tests:

```bash
go test ./...
```

Run static analysis:

```bash
go vet ./...
```

Run tests with the race detector:

```bash
go test -race ./...
```

## Benchmarks

Run benchmarks for the retry package:

```bash
go test -bench=. ./internal/retry
```

Run worker benchmarks:

```bash
go test -bench=. ./internal/worker
```

Run HTTP handler benchmarks:

```bash
go test -bench=. ./internal/handler
```

Current benchmark results were measured locally on an Intel Core i5-5350U:

```text
retry.Do                 ~10 ns/op
EventWorker.Submit      ~525 ns/op
GET /health             ~3267 ns/op
```

These numbers are environment-dependent and are intended for local performance comparison rather than production capacity estimates.

## Graceful Shutdown

The service handles `SIGINT` and `SIGTERM`.

During shutdown:

1. The HTTP server stops accepting new requests.
2. Application context is cancelled.
3. Event consumers and workers stop.
4. Connected resources are closed.

## Development Checks

Before committing changes, the project can be checked with:

```bash
gofmt -w .
go test ./...
go vet ./...
go test -race ./...
```

## Project Status

The project is complete and includes:

- HTTP API
- PostgreSQL persistence
- Redis caching
- RabbitMQ messaging
- Concurrent event workers
- Retry and timeout handling
- Structured logging
- Prometheus metrics
- Nginx reverse proxy
- Docker Compose deployment
- Automated tests
- Race detection
- Performance benchmarks
- Graceful shutdown