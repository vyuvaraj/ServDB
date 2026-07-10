# ServPool

```bash
docker run -p 8087:8087 ghcr.io/vyuvaraj/servpool:latest
```

ServPool is a database connection pooler and query routing proxy service of the Servverse ecosystem.

## Features
- **Connection Pooling**: Multiplexes and manages active pooled database connections to reduce startup overhead.
- **Read/Write Query Routing**: Automatically splits incoming operations (routing mutations to the Primary pool and read queries starting with `SELECT` to the Replica pool).
- **Multi-Database Dialects**: Dialect safety checks supporting PostgreSQL and MySQL parameter format checks.

## API Endpoints
- `POST /api/db/query` - Route and execute an SQL statement
- `GET /api/db/stats` - Fetch connection pooling and query statistics per pool

## Getting Started
To run the integration tests locally:
```bash
go test -v ./...
```

---

## Use Without Servverse (Standalone Quickstart)

`ServPool` can be used as a standalone database connection pooler and query router proxy:
1. Configure primary and replica database backend connection details via environment variables:
   ```bash
   export DB_PRIMARY_URL="postgres://user:pass@localhost:5432/primary?sslmode=disable"
   export DB_REPLICA_URL="postgres://user:pass@localhost:5432/replica?sslmode=disable"
   ```
2. Start `ServPool`:
   ```bash
   ./servpool --port 8087 --dialect postgres
   ```
3. Issue SQL statements via REST queries to the proxy endpoint:
   ```bash
   curl -X POST http://localhost:8087/api/db/query \
     -H "Content-Type: application/json" \
     -d '{"query": "SELECT * FROM users;"}'
   ```

