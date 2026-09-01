# Cursor (CLI) Agent Statusline

A minimal and functional statusline for
[Cursor (CLI) Agnet](https://cursor.com/cli)

## Installation + usage

### One-liner (auto-detects Go vs release)

```bash
curl -sL https://raw.githubusercontent.com/5c077m4n/cursor-agent-statusline/master/install.sh | bash
```

> [!IMPORTANT]
> Always inspect scripts before piping them into your shell. You can review
> [install.sh](./install.sh) first.

This will:

- Install the binary via `go install` if Go is found, or download a pre-built
  release otherwise
- Automatically configure the statusline in `~/.config/cursor/cli-config.json`
  and/or `~/.config/claude/settings.json` (whichever exist) using `jq`

### Manual steps (if you prefer to do it yourself)

```json
{
  "statusLine": {
    "type": "command",
    "command": "<value of $CURSOR_CONFIG - cursor does not expand env vars>/cursor-agent-statusline"
  }
}
```

> Replace `$CURSOR_CONFIG` above with the actual path (`~/.config/cursor` by
> default).

## Configuration

The statusline is configured via `~/.config/cursor-agent-statusline/config.json`.
All fields are optional — any omitted field falls back to its default.

```json
{
  "toggles": {
    "model": true, // Show current AI model name
    "folder": true, // Show working directory name
    "git": true, // Show git branch and dirty state
    "worktree": true, // Show git worktree name
    "pr": true, // Show PR number and review state
    "agent": true, // Show active agent name
    "session": true, // Show session name (truncated)
    "context": true, // Show context usage bar and percentage
    "cost": true, // Show total cost and duration
    "token": true, // Show input/output token counts
    "cache": true, // Show cache status and hit ratio
    "vim": true, // Show vim INSERT/NORMAL mode
    "autorun": true, // Show "auto" when autorun enabled
    "max": true, // Show "max" when max mode enabled
    "fastMode": true, // Show "fast" when fast mode enabled
    "effort": true, // Show effort level
    "thinking": true, // Show "think" when thinking enabled
    "rateLimit": true, // Show rate limit usage percentages
    "exceeds": true // Show 200k+ context warning marker
  },
  "icons": {
    "model": "󰚩",
    "folder": "󰉋",
    "branch": "󰊢",
    "worktree": "󰙅",
    "context": "󰧑",
    "cost": "󰆛",
    "vimInsert": "󰏫",
    "vimNormal": "󰌌",
    "autorun": "󰄙",
    "max": "󰓅",
    "session": "󰋩",
    "fastMode": "󰑤",
    "effort": "󰗡",
    "thinking": "󰜗",
    "duration": "󰥔",
    "lines": "󰆓",
    "pr": "󰐃",
    "agent": "󰀈",
    "rateLimit": "󰩯",
    "cache": "󰚩",
    "tokens": "󰹑",
    "exceeds": "󰄘"
  }
}
```

## Screenshot

![screenshot](./screenshot.png)

## Benchmarks (using [hyperfine](https://github.com/sharkdp/hyperfine))

See the latest results in
[CI run](https://github.com/5c077m4n/cursor-agent-statusline/actions/workflows/ci.yaml?query=branch%3Amaster)
under the Benchmarks step summary.
