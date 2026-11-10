# Code Standards

Standards inferred from the current state of the repository. Everything here is enforced or already followed in the code; keep new code in the same shape.

## Layout

- Single Go module (`go 1.27.3`), one package per concern: `main` (entry), `host/` (runtime/host detection with per-OS files), `segments/` (statusline segment rendering), `types/` (input parsing and shared types), `config/` (config loading).
- Configuration is a single JSON file (`~/.config/cli-harness-statusline/config.json`); every field is optional with a documented default.

## Pre-commit gates

Run automatically by `.githooks/` — all must pass before any push:

- `go fix ./...`
- `go test ./...`
- `go vet ./...`
- `golangci-lint run ./...` (enable `revive`; unused parameters allowed only when named with a leading `_`)
- `go build -o /dev/null .`
- `gitleaks git`

## Go

- `gofmt` formatting is the baseline; no comments unless they explain a non-obvious decision.
- Use `any` instead of `interface{}`.
- Wrap errors with `errors.Join` and a sentinel error declared outside the function (package-level `var Err... = errors.New(...)`); use `fmt.Errorf` only when dynamic context is needed (e.g. an input value in the message).
- Exported names should be self-explanatory; unexported helpers are fine without comments.

## Testing

- Tests live next to the code (`*_test.go`); shared helpers are extracted rather than duplicated.
- `go test ./...` must pass before every push (enforced by the pre-commit hook).

## Tasks

- `task test` runs the full suite verbosely with a 60s timeout; `task build` produces `bin/statusline`.
- Releases are tagged with `task bump-tag:patch` — never push tags by hand.
