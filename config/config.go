// Package config setup
package config

import (
	"log/slog"
	"os"
	"path/filepath"

	"github.com/knadh/koanf/parsers/json"
	"github.com/knadh/koanf/providers/file"
	"github.com/knadh/koanf/v2"
)

type Segment struct {
	Disable bool   `json:"disable"`
	Icon    string `json:"icon"`
}

type VimSegment struct {
	Disable    bool   `json:"disable"`
	Icon       string `json:"icon"`
	IconInsert string `json:"iconInsert"`
	IconNormal string `json:"iconNormal"`
}

type GitSegment struct {
	Disable    bool   `json:"disable"`
	Icon       string `json:"icon"`
	BranchIcon string `json:"branch"`
}

type CostSegment struct {
	Disable      bool   `json:"disable"`
	Icon         string `json:"icon"`
	DurationIcon string `json:"duration"`
	LinesIcon    string `json:"lines"`
}

type TokenSegment struct {
	Disable    bool   `json:"disable"`
	Icon       string `json:"icon"`
	TokensIcon string `json:"tokens"`
	IconIn     string `json:"iconIn"`
	IconOut    string `json:"iconOut"`
}

type SegmentsConfig struct {
	Model     Segment      `json:"model"`
	Folder    Segment      `json:"folder"`
	Git       GitSegment   `json:"git"`
	Worktree  Segment      `json:"worktree"`
	PR        Segment      `json:"pr"`
	Agent     Segment      `json:"agent"`
	Session   Segment      `json:"session"`
	Context   Segment      `json:"context"`
	Cost      CostSegment  `json:"cost"`
	Token     TokenSegment `json:"token"`
	Cache     Segment      `json:"cache"`
	Vim       VimSegment   `json:"vim"`
	Autorun   Segment      `json:"autorun"`
	Max       Segment      `json:"max"`
	FastMode  Segment      `json:"fastMode"`
	Effort    Segment      `json:"effort"`
	Thinking  Segment      `json:"thinking"`
	RateLimit Segment      `json:"rateLimit"`
	Exceeds   Segment      `json:"exceeds"`
}

type Padding struct {
	Top    int `json:"top"`
	Bottom int `json:"bottom"`
}

type Config struct {
	Segments SegmentsConfig `json:"segments"`
	Padding  Padding        `json:"padding"`
}

func configDir() string {
	if dir := os.Getenv("XDG_CONFIG_HOME"); dir != "" {
		return filepath.Join(dir, "cli-harness-statusline")
	}
	home, err := os.UserHomeDir()
	if err != nil {
		slog.Error("configDir: UserHomeDir", slog.Any("error", err))
		return ""
	}
	return filepath.Join(home, ".config", "cli-harness-statusline")
}

func Load() *Config {
	cfg := &Config{
		Segments: SegmentsConfig{
			Model:     Segment{Icon: "󰚩"},
			Folder:    Segment{Icon: "󰉋"},
			Git:       GitSegment{BranchIcon: "󰊢"},
			Worktree:  Segment{Icon: "󰙅"},
			Context:   Segment{Icon: "󰧑"},
			Cost:      CostSegment{DurationIcon: "󱎫", LinesIcon: "󰆓"},
			Vim:       VimSegment{Icon: "󰌌", IconInsert: "󰏫", IconNormal: "󰌌"},
			Autorun:   Segment{Icon: "󰄙"},
			Max:       Segment{Icon: "󰓅"},
			Session:   Segment{Icon: "󰋩"},
			FastMode:  Segment{Icon: "󰑤"},
			Effort:    Segment{Icon: "󰗡"},
			Thinking:  Segment{Icon: "󰜗"},
			PR:        Segment{Icon: "󰐃"},
			Agent:     Segment{Icon: "󰀈"},
			RateLimit: Segment{Icon: "󰩯"},
			Cache:     Segment{Icon: "󰚩"},
			Token:     TokenSegment{TokensIcon: "󰹑", IconIn: "↓", IconOut: "↑"},
			Exceeds:   Segment{Icon: "󰄘"},
		},
	}

	dir := configDir()
	if dir == "" {
		return cfg
	}

	path := filepath.Join(dir, "config.json")
	k := koanf.New(".")
	if err := k.Load(file.Provider(path), json.Parser()); err != nil {
		slog.Warn(
			"loadConfig: config file not loaded",
			slog.String("path", path),
			slog.Any("error", err),
		)
		return cfg
	}

	if err := k.Unmarshal("", cfg); err != nil {
		slog.Warn(
			"loadConfig: config unmarshal failed",
			slog.String("path", path),
			slog.Any("error", err),
		)
		return cfg
	}

	return cfg
}
