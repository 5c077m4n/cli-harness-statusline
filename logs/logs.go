// Package logs setup
package logs

import (
	"errors"
	"fmt"
	"log/slog"
	"os"
	"path/filepath"
)

const maxSize = 10 << 20 // 10 MB

var ErrInitLogging = errors.New("could not init logs")

func Init() (*os.File, error) {
	stateDir := os.Getenv("XDG_STATE_HOME")
	if stateDir == "" {
		home, err := os.UserHomeDir()
		if err != nil {
			return nil, errors.Join(ErrInitLogging, err)
		}
		stateDir = filepath.Join(home, ".local", "state")
	}

	logDir := filepath.Join(stateDir, "cli-harness-statusline")
	if err := os.MkdirAll(logDir, 0755); err != nil {
		return nil, errors.Join(ErrInitLogging, err)
	}

	logPath := filepath.Join(logDir, "statusline.jsonl")

	if fileStats, err := os.Stat(logPath); err == nil && fileStats.Size() > maxSize {
		if err := os.Remove(logPath); err != nil {
			fmt.Fprintf(os.Stderr, "could not rotate large logfile @ %s: %v\n", logPath, err)
		}
	}

	logFile, err := os.OpenFile(logPath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return nil, errors.Join(ErrInitLogging, err)
	}

	slog.SetDefault(
		slog.New(slog.NewJSONHandler(logFile, &slog.HandlerOptions{Level: slog.LevelInfo})),
	)
	return logFile, nil
}
