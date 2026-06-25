# claude-baseline-go

![Claude Code](https://img.shields.io/badge/Claude_Code-compatible-6B48FF?logo=anthropic&logoColor=white)
![MCP](https://img.shields.io/badge/MCP-GitHub%20%7C%20Postgres%20%7C%20Datadog-0078D4?logo=amazonwebservices&logoColor=white)
![Go](https://img.shields.io/badge/Go-1.26-00ADD8?logo=go&logoColor=white)
![Hooks](https://img.shields.io/badge/hooks-bash-4EAA25?logo=gnubash&logoColor=white)
![Maintained](https://img.shields.io/badge/maintained-yes-brightgreen)

Team-shared Claude AI configuration for consistent, safe, context-aware behavior across the codebase.

---

## Agent Session Context

> **Claude: read this section at the start of every session.**

**What this repo is:** A scaffold template for wiring Claude Code into a Go project. It ships pre-configured hooks, permissions, MCP connections, rules, commands, and skills. Teams fork it and fill in the stubs.

**What is live vs stub:**

| Path | Status | Notes |
|------|--------|-------|
| `.claude/settings.json` | Live | Permissions + hooks wired. |
| `.claude/hooks/pre-bash.sh` | Live | Parses stdin JSON via `jq`. Blocks on `exit 2`. Requires `jq` installed. |
| `.claude/hooks/post-tool-use.sh` | Live | Parses stdin JSON via `jq`. Writes audit log + stderr alert on Edit/Write/Bash failures. |
| `.githooks/commit-msg` | Live | Bash hook — rejects non-Conventional Commit messages. Activate with `git config core.hooksPath .githooks`. |
| `.github/workflows/ci.yml` | Live | 7-gate CI pipeline: format → vet → lint → test → race → coverage → integration. |
| `.github/workflows/conventional-commits.yml` | Live | CI: validates PR title + all commit messages on every PR. |
| `.github/workflows/release.yml` | Live | Creates GitHub Release with changelog on every `v*` tag push. No extra tools needed. |
| `scripts/version-bump.sh` | Live | Auto-semver tagging from commit history. `--dry-run` flag available. |
| `.claude/rules/*.md` | Partial stubs | Scaffold defaults filled in (layer layout, interfaces, error handling, SOLID, CI gates, security vuln classes). Team-specific sections (system overview, auth model, coverage %) still need filling. |
| `.claude/commands/` | Live | `/review`, `/standup`, `/db-schema`, `/markitdown` are functional. |
| `.claude/skills/` | Live | `on-new-file`, `pre-commit-check`, `explain-error`, `c4-architecture`, `solid-principles`, `frontend-design` auto-trigger. |
| `.mcp.json` | Live config, disabled locally | All three MCP servers are disabled in `settings.local.json`. Confirm env vars are set before assuming MCP works. |
| `main.go` | Placeholder | GoLand demo code — replace with `cmd/main.go` wiring stub. |
| `cmd/main.go` | Stub | Wiring-only entry point. Fill in real dependencies. |
| `internal/` | Stub | Layer directories created (`domain`, `service`, `repository`, `handler`). Fill in business logic. |
| `go.mod` | Placeholder | Module named `github.com/your-org/your-repo`. Rename when forking. |

**Key invariants — maintain these when editing:**

- Never write secrets, tokens, or credentials to any file.
- Never modify `.env` files (blocked by hook and deny list).
- `CLAUDE.local.md` is gitignored — personal preferences live there, not in `CLAUDE.md`.
- The five rule files (`.claude/rules/*.md`) are stubs until the team fills them in. Treat missing content as "not defined yet", not as permission to invent conventions.
- MCP servers (GitHub, Postgres, Datadog) require env vars. If `settings.local.json` has `disabledMcpjsonServers`, those connections are off regardless of `.mcp.json`.
- Hook scripts must be executable: `chmod +x .claude/hooks/*.sh && chmod +x .githooks/commit-msg`.
- Git commit-msg hook must be activated: `git config core.hooksPath .githooks`.

**Known issues to fix before production use:**

1. `main.go` — remove GoLand TIP comments; replace with real wiring (see `cmd/main.go` for the scaffold stub).
2. Rules stubs — fill team-specific sections: system overview, auth model, coverage expectations.

---

## How It Works

Claude reads configuration files at startup. The structure splits concerns:

| File / Folder | Who edits | Committed | Purpose |
|---|---|---|---|
| `CLAUDE.md` | Team | Yes | What Claude knows about the project |
| `CLAUDE.local.md` | Individual | **No** | Personal overrides and notes |
| `.mcp.json` | Team | Yes | External tool connections (GitHub, DB, Datadog) |
| `.claude/settings.json` | Team | Yes | What Claude can and cannot run |
| `.claude/rules/` | Team | Yes | Detailed chapters of `CLAUDE.md` |
| `.claude/commands/` | Team | Yes | Prompts **you** trigger with `/command-name` |
| `.claude/skills/` | Team | Yes | Prompts **Claude** triggers automatically |
| `.claude/hooks/` | Team | Yes | Shell scripts that fire on Claude events |

---

## Setup

### 1. Copy `CLAUDE.local.md`

`CLAUDE.local.md` is the **personal** config. It is already gitignored.

```bash
# It's already there as an example — edit it directly
# Never commit it
```

### 2. Set environment variables for MCP

`.mcp.json` references env vars — never hardcoded secrets.

```bash
export GITHUB_TOKEN=...
export DATABASE_URL=postgres://...
export DD_API_KEY=...
export DD_APP_KEY=...
export DD_SITE=datadoghq.eu   # or datadoghq.com
```

If you have a `settings.local.json` with `disabledMcpjsonServers`, remove the entries you want active.

### 3. Fill in the rules

Each file in `.claude/rules/` has placeholder sections. Fill them in once:

```
.claude/rules/01-architecture.md  → system overview, components, data flow
.claude/rules/02-conventions.md   → naming, file layout, forbidden patterns
.claude/rules/03-testing.md       → framework, coverage, what must be tested
.claude/rules/04-security.md      → secret handling, auth model, vuln classes
.claude/rules/05-ci.md            → CI gate config, linter list, coverage thresholds
```

Claude reads these every session — keep them accurate. Until filled, Claude treats them as undefined.

### 4. Fix the hooks

```bash
chmod +x .claude/hooks/*.sh
chmod +x .githooks/commit-msg
chmod +x scripts/version-bump.sh
```

Both Claude hooks require `jq`. Install it if not present: `apt install jq` / `brew install jq`.

### 5. Enable the git commit-msg hook

```bash
git config core.hooksPath .githooks
```

This enforces [Conventional Commits](https://www.conventionalcommits.org/en/v1.0.0/) on every local commit. The same rules run in CI on PR titles and all commit messages.

---

## Minimal Fork Path

Start with the smallest viable structure and grow as needed. Do not add layers before you need them.

```
Stage 1 — no database yet
  cmd/main.go
  internal/
    domain/        (value types + errors only)
    service/       (business logic, defines its own interfaces)
    handler/       (HTTP, imports service only)

Stage 2 — add persistence
  internal/
    repository/    (add when you have a real DB; implements interface defined in service/)

Stage 3 — scale
  Split service/ into sub-packages by domain noun
  Add internal/platform/ for cross-cutting concerns (logging, tracing, health)
```

Rule: never add a layer "just in case." Add it when the next concrete feature requires it.

---

## Commands — You trigger these

Type `/command-name` in Claude to run a command.

| Command | What it does |
|---|---|
| `/review` | Review current diff or a file for bugs, conventions, security |
| `/standup` | Generate standup summary from yesterday's git log |
| `/db-schema` | Fetch and display DB schema via MCP postgres connection |
| `/markitdown <path>` | Convert file/URL to Markdown via markitdown, saved to `raw/` |

**To add a command:** create `.claude/commands/your-command.md`. Describe what Claude should do. Use `$ARGUMENTS` for user-provided input.

---

## Skills — Claude triggers these automatically

Skills are self-activating — Claude applies them without being asked.

| Skill | Trigger |
|---|---|
| `on-new-file` | Claude just created a source file |
| `pre-commit-check` | Claude is about to suggest a `git commit` |
| `explain-error` | A command exited non-zero |
| `c4-architecture` | Designing, diagramming, or documenting system architecture; filling `01-architecture.md` |
| `solid-principles` | Creating a package, designing an interface, or reviewing component structure |
| `frontend-design` | Building or reviewing web UI — design tokens (3-tier), typography constraints, WCAG AA accessibility, component states |
| `caveman` | User types `/caveman` — activates compressed response mode |

**To add a skill:** create `.claude/skills/your-skill/SKILL.md`. Start with a `**Trigger:**` line so Claude knows when to apply it.

### Versioned skills via GitHub Packages

Skills are distributed as versioned npm packages hosted on GitHub Packages.  
`@team/caveman-skill` is pulled automatically — no manual file copy needed.

**Architecture:**

```
@team/caveman-skill  (GitHub Packages, versioned)
        ↓  dependency of
@team/claude-skills  (GitHub Packages, versioned)
        ↓  devDependency of
consumer repo  →  npm install && npm run setup:claude
                  writes .claude/skills/caveman/ automatically
```

**One-time org setup — add `.npmrc` to every consumer repo:**

```ini
# .npmrc
@team:registry=https://npm.pkg.github.com
//npm.pkg.github.com/:_authToken=${GITHUB_TOKEN}
```

`GITHUB_TOKEN` is available automatically in GitHub Actions. Local dev: use a PAT with `read:packages` scope.

**Consumer repo `package.json`:**

```json
{
  "devDependencies": {
    "@team/claude-skills": "^1.0.0"
  },
  "scripts": {
    "setup:claude": "claude-skills install"
  }
}
```

**Install:**

```bash
npm install             # pulls @team/claude-skills + its @team/caveman-skill dep
npm run setup:claude    # copies all skills → .claude/skills/
```

**Upgrade caveman:**  
Bump `@team/caveman-skill` version in `packages/claude-skills/package.json`, publish, then consumer repos run `npm update @team/claude-skills && npm run setup:claude` and commit the lockfile.

**`.gitignore` in consumer repos:**
```
.claude/skills/caveman/   # generated — source of truth is @team/caveman-skill
```

Package source lives in `packages/claude-skills/`.

---

### Extending for another language

Fork this repo → rename → swap the Go-specific files:

| File | What to change |
|---|---|
| `.claude/settings.json` | Replace `go build/test/vet/fmt` with your toolchain |
| `.claude/skills/pre-commit-check/SKILL.md` | Replace `go vet ./...` and `go test ./...` |
| `.claude/rules/03-testing.md` | Replace `go test` with your test runner |
| `.githooks/commit-msg` | Allowed types list is language-agnostic — keep as-is or extend |
| `.github/workflows/release.yml` | Changelog grouping is commit-message-based — language-agnostic, keep as-is |
| `scripts/version-bump.sh` | Bump logic is commit-message-based — language-agnostic, keep as-is |
| `README.md` badge | Update Go version badge |
| `go.mod` / `main.go` | Remove or replace with your language entry point |

Everything else — MCP, hooks, commands, rules structure, caveman — is language-agnostic.

---

## Hooks — Shell scripts on events

Hooks run outside Claude, in the shell, on specific events.

| Hook | Event | What it does |
|---|---|---|
| `.claude/hooks/pre-bash.sh` | Before every Bash call | Blocks forbidden command patterns |
| `.claude/hooks/post-tool-use.sh` | After every tool call | Audit log + failure alerts |
| `.githooks/commit-msg` | On every `git commit` | Rejects non-Conventional Commit messages |

Make executable:
```bash
chmod +x .claude/hooks/*.sh
chmod +x .githooks/commit-msg
```

Activate the git hook:
```bash
git config core.hooksPath .githooks
```

Both Claude hooks parse payloads from stdin as JSON via `jq`.

---

## Settings — Allow / Deny

`.claude/settings.json` controls what Claude can run.

- **allow** — commands Claude runs without prompting you
- **deny** — commands Claude can never run, even if asked

Edit the lists to match your project's toolchain. The skeleton ships with safe defaults for a Go project.

---

## Rules vs CLAUDE.md

`CLAUDE.md` is the summary — short enough to read in 30 seconds.  
`.claude/rules/*.md` are the chapters — full detail Claude uses when writing code.

Both are always loaded. Keep `CLAUDE.md` as an index; put specifics in rules.

---

## Conventional Commits

All commits and PR titles must follow [Conventional Commits v1.0.0](https://www.conventionalcommits.org/en/v1.0.0/).

Format: `type(scope): description`

Allowed types: `feat` `fix` `docs` `style` `refactor` `perf` `test` `chore` `ci` `build` `revert`

Breaking changes: append `!` to the type (`feat!:`) or add `BREAKING CHANGE:` in the footer.

### Enforcement layers

| Layer | What checks | Setup needed |
|-------|-------------|-------------|
| `git commit` hook | Every local commit message | `git config core.hooksPath .githooks` |
| GitHub Actions | PR title + all commit messages in the PR | Runs automatically on push |

### Changelog

Generated automatically on each tag push by `.github/workflows/release.yml`. No local tooling needed.

Push a tag → GitHub Release is created with a changelog grouped by commit type (Features, Bug Fixes, etc.). Breaking changes are surfaced at the top.

### Version bumping

```bash
./scripts/version-bump.sh            # creates a new semver tag
./scripts/version-bump.sh --dry-run  # preview what tag would be created
```

Semver rules: `BREAKING CHANGE` → major bump, `feat` → minor bump, anything else → patch.

---

## Adding a New Team Member

1. Clone repo
2. Copy env var template (share out-of-band, never commit)
3. Edit `CLAUDE.local.md` with personal preferences
4. Run `chmod +x .claude/hooks/*.sh && chmod +x .githooks/commit-msg && chmod +x scripts/version-bump.sh`
5. Run `git config core.hooksPath .githooks` — activates Conventional Commits enforcement
6. Check `settings.local.json` — remove any `disabledMcpjsonServers` you need active
7. Start Claude — configuration is automatic
