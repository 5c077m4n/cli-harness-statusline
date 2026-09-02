// Package logs setup
package logs

import (
	"fmt"
	"log/slog"
	"os"
	"path/filepath"
)

const maxSize = 10 << 20 // 10 MB

func Init() {
	stateDir := os.Getenv("XDG_STATE_HOME")
	if stateDir == "" {
		home, err := os.UserHomeDir()
		if err != nil {
			fmt.Fprintf(os.Stderr, "initLog: UserHomeDir: %v\n", err)
			return
		}
		stateDir = filepath.Join(home, ".local", "state")
	}

	logDir := filepath.Join(stateDir, "cli-harness-statusline")
	if err := os.MkdirAll(logDir, 0755); err != nil {
		fmt.Fprintf(os.Stderr, "initLog: MkdirAll %s: %v\n", logDir, err)
		return
	}

	logPath := filepath.Join(logDir, "statusline.jsonl")

	if fileStats, err := os.Stat(logPath); err == nil && fileStats.Size() > maxSize {
		if err := os.Remove(logPath); err != nil {
			fmt.Fprintf(os.Stderr, "initLog: Remove %s: %v\n", logPath, err)
		}
	}

	logFile, err := os.OpenFile(logPath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		fmt.Fprintf(os.Stderr, "initLog: OpenFile %s: %v\n", logPath, err)
		return
	}

	slog.SetDefault(
		slog.New(slog.NewJSONHandler(logFile, &slog.HandlerOptions{Level: slog.LevelInfo})),
	)
}
