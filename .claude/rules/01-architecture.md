# Architecture

> Team: fill system overview, data flow, and external dependencies.
> Layer structure and boundaries below are scaffold defaults — override if your layout differs.

## System Overview

<!-- High-level: what this system does, who uses it, what problem it solves -->

## Components

Recommended Go layer layout. Rename or split as needed — but keep the import direction.

**Rule: ALL implementation packages MUST live under `/internal/`.** This is enforced by the Go compiler — external modules cannot import from `/internal/`. No implementation code belongs at the root or in non-internal paths.

| Component | Purpose | Location |
|-----------|---------|----------|
| `cmd/` | Entry point — wires all layers together. Only place that imports across all layers. | `cmd/main.go` |
| `internal/domain` | Interfaces + value types. Zero external imports. The layer every other layer depends on. | `internal/domain/` |
| `internal/service` | Business logic. Imports `domain` only. No DB, no HTTP, no infrastructure. | `internal/service/` |
| `internal/repository` | Data access. Implements `domain` interfaces. Imports `domain` only. | `internal/repository/` |
| `internal/handler` | HTTP/gRPC layer. Imports `service` only. No business logic. | `internal/handler/` |

> All packages in the table above are under `/internal/` — this is mandatory, not a suggestion. The Go compiler prevents any external module from importing them. New packages that do not need to be imported externally must also be placed under `/internal/`.

> See `.claude/skills/c4-architecture.md` for Mermaid diagram templates at all four C4 levels.

## Data Flow

<!-- Sequence: how a request enters the system and moves through layers -->
<!-- Example: HTTP request → handler → service → repository → DB → response -->

## Key Boundaries

These must never be crossed. Claude enforces them on every code review and new file.

- `domain` has no outbound imports — not `database/sql`, not `net/http`, not any infra package
- `service` imports `domain` only — never `repository/postgres`, never `database/sql` directly
- `handler` imports `service` only — never `repository`, never `domain` directly
- `cmd/main.go` is the only file allowed to import and wire all layers together
- One package = one responsibility (SRP). If a package can't be named in 5 words, split it.
- **All implementation packages MUST be under `/internal/`.** Creating packages outside `/internal/` (other than `cmd/`) is forbidden unless there is an explicit, documented reason to expose them to external modules.

> Dependency direction: `handler → service → domain ← repository`
> Swapping Postgres for another DB should touch only `internal/repository/`.

## Observability

These rules apply to every layer. Claude enforces them on every code review and new file.

- **Logging:** Define a `Logger` interface in `internal/domain/`. Business logic must never import `log`, `zap`, or `slog` directly. Pass the logger via constructor injection. The concrete implementation (e.g. `zap`, `slog`) is wired in `cmd/main.go` only.
- **Tracing:** `context.Context` carries trace context. Pass it as the first argument to every function that does I/O. Never store a `Context` in a struct field.
- **Health:** Every service must expose `GET /healthz` (liveness) and `GET /readyz` (readiness). These are wired in the handler layer (`internal/handler/`) — no business logic in health handlers.
- **Metrics:** Datadog MCP is available (see `.mcp.json`). For code-level metrics, define a `MetricsRecorder` interface in `internal/domain/` and inject it via constructor. Never import the Datadog SDK directly in `internal/service/` or `internal/domain/`. The concrete adapter lives in `internal/repository/` or a dedicated `internal/metrics/` package, wired in `cmd/main.go`.

## External Dependencies

<!-- Third-party services this system integrates with, and why each exists -->
<!-- Format: Name | Purpose | How connected (MCP / HTTP / SDK) -->

MCP connections available (see `.mcp.json`):
- GitHub — PR/issue/code search
- PostgreSQL — database schema and data
- Datadog — metrics, logs, monitors
