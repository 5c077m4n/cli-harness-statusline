package main

import (
	"fmt"
	"log/slog"

	"github.com/5c077m4n/cli-harness-statusline/config"
	"github.com/5c077m4n/cli-harness-statusline/logs"
	"github.com/5c077m4n/cli-harness-statusline/segments"
	"github.com/5c077m4n/cli-harness-statusline/types"
	"github.com/fatih/color"
)

func init() {
	logs.Init()
	color.NoColor = false
}

func main() {
	cfg := config.Load()

	data, err := types.NewPayLoad()
	if err != nil {
		panic(err)
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
