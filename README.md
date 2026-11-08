# Go Task Management API

A small REST API built with Go's standard library. It demonstrates HTTP routing, JSON handling, input validation, concurrency-safe in-memory storage, and automated tests.

## Endpoints
| Method | Endpoint | Purpose |
|---|---|---|
| GET | `/health` | Health check |
| GET | `/tasks` | List tasks |
| POST | `/tasks` | Create a task |
| DELETE | `/tasks/{id}` | Delete a task |

Example request:
```json
{ "title": "Learn Go" }
```

## Requirements
Go 1.22 or newer.

## Run locally
```bash
go test ./...
go run .
```

Open `http://localhost:8080/health`.

Create a task:
```bash
curl -X POST http://localhost:8080/tasks -H "Content-Type: application/json" -d "{\"title\":\"Learn Go\"}"
```

## Run with Docker
```bash
docker build -t go-task-api .
docker run --rm -p 8080:8080 go-task-api
```

## Important design note
This first version uses in-memory storage so the API is easy to understand and run. Data is lost when the process stops. PostgreSQL persistence, update endpoints, authentication, and structured logging are planned enhancements.

## Tests
```bash
go test -v ./...
```
