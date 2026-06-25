---
name: pre-commit-check
description: Pre-commit validation checklist for Go. Auto-triggers before any git commit suggestion. Verifies staged diff, no secrets, go vet, go test, and Conventional Commits message format.
---

**Trigger:** Claude is about to suggest a `git commit` command.

Apply automatically before proposing any commit.

## Steps

1. Run `git diff --staged` — confirm the staged diff matches what was discussed
2. Verify no `.env` or secret files are in the staged set
3. Run `go vet ./...` — do not proceed if errors
4. Run `go test ./...` — do not proceed if any tests are red
5. Verify the proposed commit message follows Conventional Commits: `type(scope): summary`
   - Allowed types: `feat` / `fix` / `docs` / `style` / `refactor` / `perf` / `test` / `chore` / `ci` / `build` / `revert`
   - Breaking change: append `!` after type (`feat!:`) or add `BREAKING CHANGE:` footer
   - Same regex enforced by `.githooks/commit-msg` and CI
6. If any check fails: report the finding, do **NOT** suggest the commit command
