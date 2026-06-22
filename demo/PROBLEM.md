# Demo Problem: URL Shortener API

## Task

Build a minimal URL shortener HTTP service in Go.

### Requirements

1. `POST /shorten` — accepts `{"url": "https://example.com"}`, returns `{"short": "abc123"}`
2. `GET /:code` — redirects to the original URL (301)
3. `GET /healthz` — liveness probe

### Rules

- Store URLs in a PostgreSQL table `links (code TEXT PK, original_url TEXT NOT NULL, created_at TIMESTAMPTZ)`
- Code is 6-char alphanumeric, randomly generated
- If URL already shortened, return existing code (idempotent)
- Return `404` if code not found

### Acceptance Criteria

- [ ] All three endpoints work
- [ ] SQL injection not possible on any endpoint
- [ ] Unit tests for the shortening logic
- [ ] No business logic in the HTTP handler

---

## Two Solutions

| Branch | Claude config | Expected output |
|--------|--------------|-----------------|
| `demo/without-scaffold` | Vanilla Claude, no rules/hooks/skills | Works, but probably: flat structure, SQL concat, no tests |
| `demo/with-scaffold` | This scaffold active | Layer separation, parameterized SQL, table-driven tests, error wrapping |
