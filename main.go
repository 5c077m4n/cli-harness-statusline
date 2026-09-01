package main

import (
	"encoding/json/v2"
	"fmt"
	"io"
	"log/slog"
	"os"
	"path/filepath"
	"strings"

	"github.com/fatih/color"
)

func initLog() {
	stateDir := os.Getenv("XDG_STATE_HOME")
	if stateDir == "" {
		home, err := os.UserHomeDir()
		if err != nil {
			return
		}
		stateDir = filepath.Join(home, ".local", "state")
	}

	logDir := filepath.Join(stateDir, "cursor-agent-statusline")
	if err := os.MkdirAll(logDir, 0755); err != nil {
		return
	}

	f, err := os.OpenFile(
		filepath.Join(logDir, "statusline.log"),
		os.O_APPEND|os.O_CREATE|os.O_WRONLY,
		0644,
	)
	if err != nil {
		return
	}

	slog.SetDefault(slog.New(slog.NewJSONHandler(f, &slog.HandlerOptions{Level: slog.LevelInfo})))
}

func init() {
	initLog()
}

func readPayload() Payload {
	input, err := io.ReadAll(os.Stdin)
	if err != nil {
		os.Exit(1)
	}

	var data Payload
	if err := json.Unmarshal(input, &data); err != nil {
		os.Exit(1)
	}

	return data
}

func render(cfg *Config, data Payload) string {
	segments := []string{
		strings.Repeat("\n", cfg.Padding.Top),
		modelSegment(cfg, data),
		folderSegment(cfg, data),
		gitSegment(cfg, data),
		worktreeSegment(cfg, data),
		prSegment(cfg, data),
		agentSegment(cfg, data),
		sessionSegment(cfg, data),
		contextSegment(cfg, data),
		costSegment(cfg, data),
		tokenSegment(cfg, data),
		cacheSegment(cfg, data),
		vimSegment(cfg, data),
		autorunSegment(cfg, data),
		maxSegment(cfg, data),
		fastModeSegment(cfg, data),
		effortSegment(cfg, data),
		thinkingSegment(cfg, data),
		rateLimitSegment(cfg, data),
		exceedsSegment(cfg, data),
		strings.Repeat("\n", cfg.Padding.Bottom),
	}

	var nonEmpty []string
	for _, s := range segments {
		if s != "" {
			nonEmpty = append(nonEmpty, s)
		}
	}
	return strings.Join(nonEmpty, termSep())
}

func main() {
	color.NoColor = false
	cfg := loadConfig()

	data := readPayload()
	line := render(cfg, data)

	slog.Info(
		"statusline",
		slog.String("model", data.Model.DisplayName),
		slog.String("session", data.SessionName),
		slog.String("output", line),
	)

	fmt.Println(line)
}
