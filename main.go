package main

import (
	"encoding/json/v2"
	"fmt"
	"io"
	"log/slog"
	"os"

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

func readPayload() types.Payload {
	input, err := io.ReadAll(os.Stdin)
	if err != nil {
		slog.Error("readPayload: stdin read failed", slog.Any("error", err))
		os.Exit(1)
	}

	var data types.Payload
	if err := json.Unmarshal(input, &data); err != nil {
		slog.Error(
			"readPayload: json unmarshal failed",
			slog.Any("error", err),
			slog.Int("input_len", len(input)),
		)
		os.Exit(1)
	}

	return data
}

func main() {
	cfg := config.Load()

	data := readPayload()
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
