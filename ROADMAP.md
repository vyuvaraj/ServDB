# ServDB Roadmap

This roadmap outlines the planned development phases for the ServDB database proxy service.

---

## Phase 1: Connection Proxying (In Progress)
- [x] **Connection pooling** — Shared pooling and connection reuse proxy. [June 29, 2026]
- [x] **Query routing** — Read replica routing, primary write routing. [June 29, 2026]
- [x] **Multi-database support** — Multi-dialect parser backend support. [June 29, 2026]
- [x] **Serv-lang integration** — Centralized client driver connection pool setup. [June 29, 2026]

- [x] **Slow query detection** — Slow query profiling telemetry. [June 29, 2026]
- [x] **Query analytics** — CPU cost and pattern aggregation. [June 29, 2026]
- [ ] **Query caching** — Invalidation caching via ServCache.
- [x] **Centralized migrations** — Centralized schema migration runner. [June 29, 2026]
