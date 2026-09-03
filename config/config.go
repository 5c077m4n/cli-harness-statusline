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

type SegmentConfig struct {
	Disable bool `json:"disable"`
}
type SegmentsConfig struct {
	Model     SegmentConfig `json:"model"`
	Folder    SegmentConfig `json:"folder"`
	Git       SegmentConfig `json:"git"`
	Worktree  SegmentConfig `json:"worktree"`
	PR        SegmentConfig `json:"pr"`
	Agent     SegmentConfig `json:"agent"`
	Session   SegmentConfig `json:"session"`
	Context   SegmentConfig `json:"context"`
	Cost      SegmentConfig `json:"cost"`
	Token     SegmentConfig `json:"token"`
	Cache     SegmentConfig `json:"cache"`
	Vim       SegmentConfig `json:"vim"`
	Autorun   SegmentConfig `json:"autorun"`
	Max       SegmentConfig `json:"max"`
	FastMode  SegmentConfig `json:"fastMode"`
	Effort    SegmentConfig `json:"effort"`
	Thinking  SegmentConfig `json:"thinking"`
	RateLimit SegmentConfig `json:"rateLimit"`
	Exceeds   SegmentConfig `json:"exceeds"`
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
	cfg := &Config{}

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
