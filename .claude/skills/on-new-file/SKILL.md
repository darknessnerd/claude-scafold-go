---
name: on-new-file
description: Post-creation checklist for new source files. Auto-triggers when Claude creates a new file. Verifies naming, package declaration, headers, secrets, and test file presence.
---

**Trigger:** Claude just created a new source file.

Apply automatically — no user prompt needed.

## Steps

1. Verify file name follows naming conventions in `.claude/rules/02-conventions.md`
2. Confirm package declaration matches the directory name
3. Add required file header if a team standard is defined in `rules/02-conventions.md`
4. Check no secrets or hardcoded values were introduced (tokens, passwords, IPs)
5. Confirm a corresponding `_test` file was created if the new file exports functions
6. Verify the file is under `internal/` or `cmd/` — no implementation code at the repo root or in non-internal paths (see `01-architecture.md` Key Boundaries)
7. Check import direction is respected:
   - `internal/handler/` → imports `internal/service/` only (never `internal/repository/`, never `internal/domain/`)
   - `internal/service/` → imports `internal/domain/` only (never `internal/repository/`, never `database/sql`)
   - `internal/repository/` → imports `internal/domain/` only
   - `internal/domain/` → zero external or internal imports
   - `cmd/main.go` → only file allowed to import across all layers
8. Confirm `internal/domain/` contains only value types (structs, enums, sentinel errors) — no interface definitions
