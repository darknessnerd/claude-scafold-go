# Security

> Claude must follow these rules without exception.
> Team: fill in Auth Model. All other sections are enforced as written.

## Secrets

- Never hardcode credentials, tokens, or keys
- Never write to `.env` files (blocked by hook + deny list)
- Never log values from auth headers, passwords, or tokens
- All secrets via env vars or secret manager — see `.mcp.json` for the pattern
- If a column name contains `token`, `secret`, `password`, `key`, or `credential` — flag it in review

## Auth Model

<!-- Team: fill in who authenticates, how tokens flow, where they're validated -->
<!-- Example: JWT signed with RS256, validated in handler middleware, claims passed via context.Context -->

## What Claude Must Never Do

- Commit `.env` or any file containing a secret
- Generate code that stores plaintext passwords
- Disable TLS verification (`InsecureSkipVerify: true`)
- Introduce SQL string concatenation — parameterized queries only
- Add `// nolint` to security-related linter warnings without team approval
- Call `exec.Command` with user-supplied input — allowlist commands, never interpolate user data
- Write outbound HTTP calls without URL validation — validate or allowlist the host first

## Vulnerability Classes to Avoid

| Class | Rule |
|-------|------|
| SQL injection | Parameterized queries always — `db.QueryContext(ctx, query, args...)`, never `fmt.Sprintf` into SQL |
| XSS | Escape all user-controlled output; in Go templates use `html/template`, never `text/template` for HTML |
| SSRF | Validate or allowlist outbound URLs; never fetch a URL constructed from user input without host check |
| Command injection | Never `exec.Command(userInput)` — allowlist commands and args separately |
| Path traversal | Never use user input in file paths; use `filepath.Clean` and validate stays within allowed root |
| Open redirect | Validate redirect targets against an allowlist of trusted hosts |

## Dependency Security

- Run `go mod tidy` before commit — remove unused dependencies
- Pin dependencies to specific versions in `go.sum`
- Flag any dependency added without a clear justification in the PR description
