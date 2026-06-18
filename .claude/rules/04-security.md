# Security

> Claude must follow these rules without exception.

## Secrets

- Never hardcode credentials, tokens, or keys
- Never write to `.env` files
- Never log values from auth headers, passwords, or tokens
- All secrets via env vars or secret manager — see `.mcp.json` for pattern

## Auth Model

<!-- Describe: who authenticates, how tokens flow, where they're validated -->

## What Claude Must Never Do

- Commit `.env` or any file containing a secret
- Generate code that stores plaintext passwords
- Disable TLS verification
- Introduce SQL string concatenation (use parameterized queries)
- Add `// nolint` to security-related linter warnings without team approval

## Vulnerability Classes to Avoid

- SQL injection → parameterized queries always
- XSS → escape all user-controlled output
- SSRF → validate/allowlist outbound URLs
- Command injection → never `exec(userInput)`
