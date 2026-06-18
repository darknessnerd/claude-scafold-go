# claude-baseline-go

![Claude Code](https://img.shields.io/badge/Claude_Code-compatible-6B48FF?logo=anthropic&logoColor=white)
![MCP](https://img.shields.io/badge/MCP-GitHub%20%7C%20Postgres%20%7C%20Datadog-0078D4?logo=amazonwebservices&logoColor=white)
![Go](https://img.shields.io/badge/Go-1.26-00ADD8?logo=go&logoColor=white)
![Hooks](https://img.shields.io/badge/hooks-bash-4EAA25?logo=gnubash&logoColor=white)
![Maintained](https://img.shields.io/badge/maintained-yes-brightgreen)

Team-shared Claude AI configuration for consistent, safe, context-aware behavior across the codebase.

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

### 3. Fill in the rules

Each file in `.claude/rules/` has placeholder sections. Fill them in once:

```
.claude/rules/01-architecture.md  → system overview, components, data flow
.claude/rules/02-conventions.md   → naming, file layout, forbidden patterns
.claude/rules/03-testing.md       → framework, coverage, what must be tested
.claude/rules/04-security.md      → secret handling, auth model, vuln classes
```

Claude reads these every session — keep them accurate.

---

## Commands — You trigger these

Type `/command-name` in Claude to run a command.

| Command | What it does |
|---|---|
| `/review` | Review current diff or a file for bugs, conventions, security |
| `/standup` | Generate standup summary from yesterday's git log |
| `/db-schema` | Fetch and display DB schema via MCP postgres connection |

**To add a command:** create `.claude/commands/your-command.md`. Describe what Claude should do. Use `$ARGUMENTS` for user-provided input.

---

## Skills — Claude triggers these automatically

Skills are self-activating — Claude applies them without being asked.

| Skill | Trigger |
|---|---|
| `on-new-file` | Claude just created a source file |
| `pre-commit-check` | Claude is about to suggest a `git commit` |
| `explain-error` | A command exited non-zero |
| `caveman` | User types `/caveman` — activates compressed response mode |

**To add a skill:** create `.claude/skills/your-skill.md`. Start with a `**Trigger:**` line so Claude knows when to apply it.

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
                  writes .claude/skills/caveman.md automatically
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
.claude/skills/caveman.md   # generated — source of truth is @team/caveman-skill
```

Package source lives in `packages/claude-skills/`.

---

### Extending for another language

Fork this repo → rename → swap the Go-specific files:

| File | What to change |
|---|---|
| `.claude/settings.json` | Replace `go build/test/vet/fmt` with your toolchain |
| `.claude/skills/pre-commit-check.md` | Replace `go vet ./...` and `go test ./...` |
| `.claude/rules/03-testing.md` | Replace `go test` with your test runner |
| `README.md` badge | Update Go version badge |
| `go.mod` / `main.go` | Remove or replace with your language entry point |

Everything else — MCP, hooks, commands, rules structure, caveman — is language-agnostic.

---

## Hooks — Shell scripts on events

Hooks run outside Claude, in the shell, on specific events.

| Hook | Event | What it does |
|---|---|---|
| `pre-bash.sh` | Before every Bash call | Blocks forbidden command patterns |
| `post-tool-use.sh` | After every tool call | Audit log + failure alerts |

Hooks must be executable:
```bash
chmod +x .claude/hooks/*.sh
```

Wire them in `.claude/settings.json` under the `hooks` key (see Claude Code docs for event names).

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

## Adding a New Team Member

1. Clone repo
2. Copy env var template (share out-of-band, never commit)
3. Edit `CLAUDE.local.md` with personal preferences
4. Run `chmod +x .claude/hooks/*.sh`
5. Start Claude — configuration is automatic
