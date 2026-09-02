# CLI Harness Statusline

A minimal and functional statusline for
[Cursor (CLI) Agnet](https://cursor.com/cli) &
[Claude Code](https://claude.com/product/claude-code)

## Installation + usage

### One-liner (auto-detects Go vs release)

```bash
curl -sL https://raw.githubusercontent.com/5c077m4n/cli-harness-statusline/refs/heads/master/install.sh | bash
```

> [!IMPORTANT]
> Always inspect scripts before piping them into your shell. You can review
> [install.sh](./install.sh) first.

## Configuration

The statusline is configured via `~/.config/cli-harness-statusline/config.json`.
All fields are optional — any omitted field falls back to its default.

```json
{
  "segments": {
    "model": { "disable": false, "icon": "󰚩" }, // Show current AI model name
    "folder": { "disable": false, "icon": "󰉋" }, // Show working directory name
    "git": { "disable": false, "branch": "󰊢" }, // Show git branch and dirty state
    "worktree": { "disable": false, "icon": "󰙅" }, // Show git worktree name
    "pr": { "disable": false, "icon": "󰐃" }, // Show PR number and review state
    "agent": { "disable": false, "icon": "󰀈" }, // Show active agent name
    "session": { "disable": false, "icon": "󰋩" }, // Show session name (truncated)
    "context": { "disable": false, "icon": "󰧑" }, // Show context usage bar and percentage
    "cost": { "disable": false, "duration": "󱎫", "lines": "󰆓" }, // Show total cost, duration, lines
    "token": { "disable": false, "tokens": "󰹑", "iconIn": "↑", "iconOut": "↓" }, // Show input/output token counts
    "cache": { "disable": false, "icon": "󰚩" }, // Show cache status and hit ratio
    "vim": {
      "disable": false,
      "icon": "󰌌",
      "iconInsert": "󰏫",
      "iconNormal": "󰌌"
    }, // Show vim INSERT/NORMAL mode
    "autorun": { "disable": false, "icon": "󰄙" }, // Show "auto" when autorun enabled
    "max": { "disable": false, "icon": "󰓅" }, // Show "max" when max mode enabled
    "fastMode": { "disable": false, "icon": "󰑤" }, // Show "fast" when fast mode enabled
    "effort": { "disable": false, "icon": "󰗡" }, // Show effort level
    "thinking": { "disable": false, "icon": "󰜗" }, // Show "think" when thinking enabled
    "rateLimit": { "disable": false, "icon": "󰩯" }, // Show rate limit usage percentages
    "exceeds": { "disable": false, "icon": "󰄘" } // Show 200k+ context warning marker
  },
  "padding": {
    "top": 1,
    "bottom": 1
  } // Number of blank lines above/below the statusline (default: 0)
}
```

## Screenshot

![screenshot](./screenshot.png)

## Benchmarks (using [hyperfine](https://github.com/sharkdp/hyperfine))

See the latest results in
[CI run](https://github.com/5c077m4n/cli-harness-statusline/actions/workflows/ci.yaml?query=branch%3Amaster)
under the Benchmarks step summary.
