package main

import (
	"os"
	"path/filepath"

	"github.com/knadh/koanf/parsers/json"
	"github.com/knadh/koanf/providers/file"
	"github.com/knadh/koanf/v2"
)

type TogglesConfig struct {
	Model     bool `json:"model"`
	Folder    bool `json:"folder"`
	Git       bool `json:"git"`
	Worktree  bool `json:"worktree"`
	PR        bool `json:"pr"`
	Agent     bool `json:"agent"`
	Session   bool `json:"session"`
	Context   bool `json:"context"`
	Cost      bool `json:"cost"`
	Token     bool `json:"token"`
	Cache     bool `json:"cache"`
	Vim       bool `json:"vim"`
	Autorun   bool `json:"autorun"`
	Max       bool `json:"max"`
	FastMode  bool `json:"fastMode"`
	Effort    bool `json:"effort"`
	Thinking  bool `json:"thinking"`
	RateLimit bool `json:"rateLimit"`
	Exceeds   bool `json:"exceeds"`
}

type IconsConfig struct {
	Model     string `json:"model"`
	Folder    string `json:"folder"`
	Branch    string `json:"branch"`
	Worktree  string `json:"worktree"`
	Context   string `json:"context"`
	Cost      string `json:"cost"`
	VimInsert string `json:"vimInsert"`
	VimNormal string `json:"vimNormal"`
	Autorun   string `json:"autorun"`
	Max       string `json:"max"`
	Session   string `json:"session"`
	FastMode  string `json:"fastMode"`
	Effort    string `json:"effort"`
	Thinking  string `json:"thinking"`
	Duration  string `json:"duration"`
	Lines     string `json:"lines"`
	PR        string `json:"pr"`
	Agent     string `json:"agent"`
	RateLimit string `json:"rateLimit"`
	Cache     string `json:"cache"`
	Tokens    string `json:"tokens"`
	Exceeds   string `json:"exceeds"`
}

type Config struct {
	Toggles TogglesConfig `json:"toggles"`
	Icons   IconsConfig   `json:"icons"`
}

func loadConfig() *Config {
	cfg := &Config{
		Toggles: TogglesConfig{
			Model:     true,
			Folder:    true,
			Git:       true,
			Worktree:  true,
			PR:        true,
			Agent:     true,
			Session:   true,
			Context:   true,
			Cost:      true,
			Token:     true,
			Cache:     true,
			Vim:       true,
			Autorun:   true,
			Max:       true,
			FastMode:  true,
			Effort:    true,
			Thinking:  true,
			RateLimit: true,
			Exceeds:   true,
		},
		Icons: IconsConfig{
			Model:     "󰚩",
			Folder:    "󰉋",
			Branch:    "󰊢",
			Worktree:  "󰙅",
			Context:   "󰧑",
			Cost:      "󰆛",
			VimInsert: "󰏫",
			VimNormal: "󰌌",
			Autorun:   "󰄙",
			Max:       "󰓅",
			Session:   "󰋩",
			FastMode:  "󰑤",
			Effort:    "󰗡",
			Thinking:  "󰜗",
			Duration:  "󰥔",
			Lines:     "󰆓",
			PR:        "󰐃",
			Agent:     "󰀈",
			RateLimit: "󰩯",
			Cache:     "󰚩",
			Tokens:    "󰹑",
			Exceeds:   "󰄘",
		},
	}

	home, err := os.UserHomeDir()
	if err != nil {
		return cfg
	}

	path := filepath.Join(home, ".config", "cursor-agent-statusline", "config.json")
	k := koanf.New(".")
	if err := k.Load(file.Provider(path), json.Parser()); err != nil {
		return cfg
	}

	if err := k.Unmarshal("", cfg); err != nil {
		return cfg
	}

	return cfg
}
