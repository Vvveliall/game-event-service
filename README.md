# Game Event Service

Backend service for receiving and processing events from a game platform.

## Current stage

Block 1: HTTP API, configuration, healthcheck and graceful shutdown.

## Run

```bash
cp .env.example .env
set -a
source .env
set +a
go run ./cmd/server
```

## Endpoints

### Health

`GET /health`

### Create event

`POST /api/v1/events`

Example body:

```json
{
  "player_id": 1,
  "type": "purchase.created",
  "payload": "sword_01"
}
```
