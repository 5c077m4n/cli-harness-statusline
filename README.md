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

## Screenshot

![screenshot](./screenshot.png)

## Benchmarks (using [hyperfine](https://github.com/sharkdp/hyperfine))

See the latest results in
[CI run](https://github.com/5c077m4n/cursor-agent-statusline/actions/workflows/ci.yaml?query=branch%3Amaster)
under the Benchmarks step summary.
