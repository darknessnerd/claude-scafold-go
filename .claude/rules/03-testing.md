# Testing

> Team: fill in Framework, Coverage Expectations, What Must NOT Be Mocked, and domain-specific requirements.
> Layout rules and must-have-tests list below are scaffold defaults.

## Framework

<!-- Team: specify your testing stack -->
<!-- Examples: standard `testing` only | testify/assert | gomock | testcontainers -->

## Coverage Expectations

<!-- Team: specify minimum coverage % if enforced, or which paths must have tests -->
<!-- Example: 80% overall, 100% on all exported functions in domain/ -->

## Test Layout

Tests live alongside source in the same package (white-box) unless testing the public API surface (black-box, `package foo_test`).

```
internal/
  service/
    user.go
    user_test.go        ← white-box: same package, tests unexported helpers too
  handler/
    user_test.go        ← black-box: package handler_test, tests public surface only
test/
  integration/          ← integration and e2e tests, separate from unit tests
```

One test file per source file. Integration tests in `/test/integration/`, tagged with `//go:build integration`.

## What Must Have Tests

- All exported functions in `domain/` and `service/`
- All error paths — every `if err != nil` branch that returns a non-nil error
- All HTTP handler routes — at minimum: happy path + 400 + 404/500
- Any function touching auth, tokens, or credentials

<!-- Team: add domain-specific requirements -->

## What Must NOT Be Mocked

<!-- Team: critical — specify what must use real implementations -->
<!-- Example: "never mock the DB — use testcontainers for repository tests" -->

## Running Tests

```bash
# unit tests only
go test ./...

# with race detector (required before merge)
go test -race ./...

# integration tests (requires running dependencies)
go test -tags=integration ./test/integration/...

# coverage report
go test -coverprofile=coverage.out ./... && go tool cover -html=coverage.out
```
