#!/usr/bin/env bash
set -euo pipefail

BIN_NAME="cli-harness-statusline"
REPO="5c077m4n/${BIN_NAME}"

# --- Resolve install dirs ---
CURSOR_CONFIG="${XDG_CONFIG_HOME:-$HOME/.config}/cursor"
CLAUDE_CONFIG="$HOME/.claude"
GO_BIN="${GOBIN:-$(go env GOBIN 2>/dev/null || echo "$HOME/.local/share/go/bin")}"

# --- Step 1: Install the binary ---
if command -v go &>/dev/null; then
	echo "Go detected – installing via 'go install'..."
	go install "github.com/${REPO}@latest"
	BIN_PATH="$GO_BIN/$BIN_NAME"
else
	echo "Go not found – downloading pre-built binary from GitHub releases..."

	# Detect OS / arch
	ARCHIVE_OS=$(uname -s | tr '[:upper:]' '[:lower:]')
	ARCHIVE_ARCH=$(uname -m)
	case "$ARCHIVE_ARCH" in
	x86_64) ARCHIVE_ARCH="amd64" ;;
	aarch64) ARCHIVE_ARCH="arm64" ;;
	esac

	if [ "$ARCHIVE_OS" != "darwin" ] && [ "$ARCHIVE_OS" != "linux" ]; then
		echo "Unsupported OS: $ARCHIVE_OS"
		exit 1
	fi

	PLATFORM="${ARCHIVE_OS}_${ARCHIVE_ARCH}"

	# Fetch latest release tag from GitHub API
	echo "Fetching latest release info..."
	LATEST=$(curl -sL "https://api.github.com/repos/${REPO}/releases/latest" | jq -r '.tag_name')
	VERSION="${LATEST#v}"
	FILENAME="${BIN_NAME}_${VERSION}_${PLATFORM}.tar.gz"
	URL="https://github.com/${REPO}/releases/download/${LATEST}/${FILENAME}"

	LOCAL_BIN="$HOME/.local/bin"
	mkdir -p "$LOCAL_BIN"

	TMP_DIR=$(mktemp -d)
	trap 'rm -rf "$TMP_DIR"' EXIT

	echo "Downloading $FILENAME ..."
	curl -sL "$URL" -o "$TMP_DIR/$FILENAME"
	tar -xzf "$TMP_DIR/$FILENAME" -C "$TMP_DIR"
	mv "$TMP_DIR/$BIN_NAME" "$LOCAL_BIN/"
	chmod u+x "$LOCAL_BIN/$BIN_NAME"
	BIN_PATH="$LOCAL_BIN/$BIN_NAME"
fi

echo "Binary installed at: $BIN_PATH"

# --- Step 2: Configure the statusline ---
if ! command -v jq &>/dev/null; then
	echo "jq not found – skipping automatic configuration."
	echo "Manual config instructions: https://github.com/${REPO}#readme"
	exit 0
fi

CONFIG_ENTRY=$(
	cat <<JSON
{
  "statusLine": {
    "type": "command",
    "command": "$BIN_PATH"
  }
}
JSON
)

configured_count=0

if [ -f "$CURSOR_CONFIG/cli-config.json" ]; then
	echo "Updating $CURSOR_CONFIG/cli-config.json ..."
	# merge statusLine into the top-level object, overwriting if exists
	jq --argjson entry "$CONFIG_ENTRY" '.statusLine = $entry.statusLine' \
		"$CURSOR_CONFIG/cli-config.json" >"$CURSOR_CONFIG/cli-config.json.tmp" &&
		mv "$CURSOR_CONFIG/cli-config.json.tmp" "$CURSOR_CONFIG/cli-config.json"
	configured_count=$((configured_count + 1))
fi

if [ -f "$CLAUDE_CONFIG/settings.json" ]; then
	echo "Updating $CLAUDE_CONFIG/settings.json ..."
	jq --argjson entry "$CONFIG_ENTRY" '.statusLine = $entry.statusLine' \
		"$CLAUDE_CONFIG/settings.json" >"$CLAUDE_CONFIG/settings.json.tmp" &&
		mv "$CLAUDE_CONFIG/settings.json.tmp" "$CLAUDE_CONFIG/settings.json"
	configured_count=$((configured_count + 1))
fi

if [ "$configured_count" -eq 0 ]; then
	echo "No existing config files found at:"
	echo "  - $CURSOR_CONFIG/cli-config.json"
	echo "  - $CLAUDE_CONFIG/settings.json"
	echo ""
	echo "Create one manually. Example for $CURSOR_CONFIG/cli-config.json:"
	echo "$CONFIG_ENTRY" | jq .
fi

echo "Done."
