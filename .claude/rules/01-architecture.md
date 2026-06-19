# Architecture

> Team: fill system overview, data flow, and external dependencies.
> Layer structure and boundaries below are scaffold defaults — override if your layout differs.

## System Overview

<!-- High-level: what this system does, who uses it, what problem it solves -->

## Components

Recommended Go layer layout. Rename or split as needed — but keep the import direction.

| Component | Purpose | Location |
|-----------|---------|----------|
| `cmd/` | Entry point — wires all layers together. Only place that imports across all layers. | `cmd/main.go` |
| `internal/domain` | Interfaces + value types. Zero external imports. The layer every other layer depends on. | `internal/domain/` |
| `internal/service` | Business logic. Imports `domain` only. No DB, no HTTP, no infrastructure. | `internal/service/` |
| `internal/repository` | Data access. Implements `domain` interfaces. Imports `domain` only. | `internal/repository/` |
| `internal/handler` | HTTP/gRPC layer. Imports `service` only. No business logic. | `internal/handler/` |

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

> Dependency direction: `handler → service → domain ← repository`
> Swapping Postgres for another DB should touch only `internal/repository/`.

## External Dependencies

<!-- Third-party services this system integrates with, and why each exists -->
<!-- Format: Name | Purpose | How connected (MCP / HTTP / SDK) -->

MCP connections available (see `.mcp.json`):
- GitHub — PR/issue/code search
- PostgreSQL — database schema and data
- Datadog — metrics, logs, monitors
