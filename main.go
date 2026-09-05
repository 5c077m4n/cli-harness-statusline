package main

import (
	"fmt"
	"log/slog"
	"os"

	"github.com/5c077m4n/cli-harness-statusline/config"
	"github.com/5c077m4n/cli-harness-statusline/logs"
	"github.com/5c077m4n/cli-harness-statusline/segments"
	"github.com/5c077m4n/cli-harness-statusline/types"
	"github.com/fatih/color"
)

func main() {
	if logFile, err := logs.Init(); logFile != nil && err == nil {
		defer func() { _ = logFile.Close() }()
	}
	color.NoColor = false

	cfg := config.Load()

	data, err := types.NewPayLoad()
	if err != nil {
		slog.Error("could not parse payload", slog.Any("error", err))
		os.Exit(1)
	}
	line := segments.Render(cfg, data)

	slog.Info(
		"statusline",
		slog.String("model", data.Model.DisplayName),
		slog.String("session", data.SessionName),
		slog.String("output", line),
	)

	for range cfg.Padding.Top {
		fmt.Println("")
	}
	fmt.Println(line)
	for range cfg.Padding.Bottom {
		fmt.Println("")
	}
}
