# Conventions

> Team: fill in File Layout, Error Handling, Imports, Forbidden Patterns.
> Naming rules below are Go standard — override only if your team diverges.

## Naming

- Files: `snake_case`
- Packages: lowercase single word — name by what it contains, not what it does (`user` not `userservice`)
- Exported types/funcs: `PascalCase`
- Unexported: `camelCase`
- Constants: `ALL_CAPS` only for truly global invariants
- Interfaces: noun or noun phrase (`Repository`, `Notifier`), not `IRepository` or `RepositoryInterface`
- Constructors: `New<Type>` returning the interface, not the concrete type

> See `.claude/skills/solid-principles.md` for package naming rules derived from SRP.

## File Layout

<!-- Team: describe your actual module/package structure here -->
<!-- Scaffold default below — adjust if you use a different layout -->

```
cmd/
  main.go              ← entry point, wires dependencies
internal/
  domain/              ← interfaces + value types (no external imports)
  service/             ← business logic (imports domain only)
  repository/          ← data access (imports domain only)
  handler/             ← HTTP/gRPC (imports service only)
```

One file per logical concern. No `utils.go` or `helpers.go` — if it needs a file, it needs a package.

## Error Handling

<!-- Team: fill in your wrapping style and logging rules -->

Defaults until overridden:

- Wrap errors with context: `fmt.Errorf("user.Get: %w", err)` — include the call site
- Return errors to callers; log only at the top boundary (handler layer)
- Sentinel errors in `domain` package: `var ErrNotFound = errors.New("not found")`
- Never `panic` outside `init()` or package-level setup
- Never swallow errors with `_ = someFunc()`

## Imports

<!-- Team: confirm grouping order or override -->

Standard Go import grouping — enforced by `goimports`:

```go
import (
    // 1. stdlib
    "context"
    "fmt"

    // 2. internal packages
    "github.com/your-org/your-repo/internal/domain"

    // 3. external dependencies
    "github.com/some/library"
)
```

Blank line between each group. `goimports` or `gofumpt` handles this automatically.

## Forbidden Patterns

<!-- Team: add domain-specific forbidden patterns below -->

- No `panic` outside `init()` or package-level var initialization
- No global mutable state — pass dependencies via constructors
- No `interface{}` / `any` where a concrete type or typed interface suffices
- No `init()` functions that perform I/O or have side effects
- No `// nolint` on security-related linter warnings without team approval (see `04-security.md`)
- No type switches on concrete types to vary behavior — use interface methods instead (OCP violation)
- Constructors must accept interfaces, not concrete types (`func New(repo domain.UserReader)` not `func New(db *PostgresRepo)`)
