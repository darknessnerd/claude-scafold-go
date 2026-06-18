# Testing

> Team: fill this in. Claude follows these when writing tests.

## Framework

<!-- e.g., Go standard `testing`, testify, etc. -->

## Coverage Expectations

<!-- Minimum % if enforced, or which paths must have tests -->

## Test Layout

<!-- Unit tests alongside source? Separate `_test` packages? Integration in `/test`? -->

## What Must Have Tests

- All exported functions
- All error paths
- <!-- add domain-specific requirements -->

## What Must NOT Be Mocked

<!-- e.g., "never mock the DB — use testcontainers" -->

## Running Tests

```bash
go test ./...
```
