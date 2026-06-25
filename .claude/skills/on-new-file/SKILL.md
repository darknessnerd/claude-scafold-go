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
