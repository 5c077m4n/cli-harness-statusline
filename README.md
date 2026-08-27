# Cursor (CLI) Agent Statusline

A minimal and functional statusline for
[Cursor (CLI) Agnet](https://cursor.com/cli)

## Installation + usage

```bash
go install github.com/5c077m4n/cursor-agent-statusline@latest
```

and then add to your `~/.config/cursor/cli-config.json`:

```json
{
  "statusLine": {
    "type": "command",
    "command": "~/.local/share/go/bin/cursor-agent-statusline"
  }
}
```

## Screenshot

![screenshot](./screenshot.png)

## Benchmarks (using [hyperfine](https://github.com/sharkdp/hyperfine))

Using a full payload:

|  Mean  |   σ   | Median |  Min   |  Max   | User  | System | Runs |
| :----: | :---: | :----: | :----: | :----: | :---: | :----: | :--: |
| 11.1ms | 2.0ms | 10.7ms | 10.0ms | 30.6ms | 3.8ms | 5.6ms  | 267  |
