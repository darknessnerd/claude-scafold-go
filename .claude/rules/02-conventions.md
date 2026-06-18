# Conventions

> Team: fill this in. Claude enforces these when writing or reviewing code.

## Naming

- Files: `snake_case`
- Packages: lowercase single word
- Exported types/funcs: `PascalCase`
- Unexported: `camelCase`
- Constants: `ALL_CAPS` only for truly global invariants

## File Layout

<!-- Describe expected package/module structure -->

## Error Handling

<!-- Wrap or sentinel? When to log vs return? -->

## Imports

<!-- Grouping order: stdlib / internal / external -->

## Forbidden Patterns

<!-- Things that look OK but aren't: e.g., no `panic` outside init, no global state -->
