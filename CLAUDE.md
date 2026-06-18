# Claude Configuration — Team Shared

Do NOT add personal preferences here. Use `CLAUDE.local.md` (gitignored).

## What Claude Is

AI pair-programmer for this repo. Reads code, writes code, runs safe shell commands.
Not an autonomous agent — confirm before destructive or irreversible ops.

## Architecture

> Fill in: system overview, major components, data flow, key boundaries.

See `.claude/rules/01-architecture.md` for full details.

## Conventions

> Fill in: language, framework, naming, file layout, import style.

See `.claude/rules/02-conventions.md` for full details.

## Testing

> Fill in: test framework, coverage expectations, what must have tests.

See `.claude/rules/03-testing.md` for full details.

## Security

> Fill in: secret handling, auth model, what never to log or commit.

See `.claude/rules/04-security.md` for full details.

## What Claude Can Touch

Controlled via `.claude/settings.json`. Summary:
- Read: anything
- Write/Edit: source files, configs (non-secret)
- Run: build, test, lint commands
- Never: force-push, drop tables, write `.env` files

## MCP Connections

Defined in `.mcp.json`: GitHub, database, Datadog.
Credentials come from env vars — never hardcoded.
