#!/usr/bin/env bash
# Hook: post-tool-use
# Event:  PostToolUse
# Fires:  after every tool execution (Bash, Edit, Write, Read, MCP calls, etc.)
# Effect: audit log + stderr alert on failures — does NOT block or undo the tool
#
# How Claude Code calls this:
#   $1 = tool name   (e.g. "Bash", "Edit", "Write")
#   $2 = exit code   (0 = success, non-zero = tool reported an error)

TOOL="$1"
EXIT_CODE="$2"

# Append one line per tool call to a local audit log.
# The log lives inside .claude/logs/ — add that path to .gitignore
# if you do not want it committed (it can grow large in active sessions).
LOGFILE=".claude/logs/tool-audit.log"
mkdir -p ".claude/logs"
echo "$(date -u +%Y-%m-%dT%H:%M:%SZ) TOOL=$TOOL EXIT=$EXIT_CODE" >> "$LOGFILE"

# Write failures to stderr so they surface in the Claude Code UI.
# Only alert on tools that mutate state — Read failures are expected and noisy.
if [[ "$EXIT_CODE" != "0" && "$TOOL" =~ ^(Edit|Write|Bash)$ ]]; then
  echo "⚠ Tool '$TOOL' exited $EXIT_CODE" >&2
fi

# Exit 0 always — this hook is observability only.
# Non-zero here would be treated as a hook crash, not a tool block.
exit 0
