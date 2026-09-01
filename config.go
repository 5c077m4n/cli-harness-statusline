package main

import (
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

func loadConfig() *Config {
	cfg := &Config{
		Segments: SegmentsConfig{
			Model:     Segment{Icon: "󰚩"},
			Folder:    Segment{Icon: "󰉋"},
			Git:       GitSegment{BranchIcon: "󰊢"},
			Worktree:  Segment{Icon: "󰙅"},
			Context:   Segment{Icon: "󰧑"},
			Cost:      CostSegment{DurationIcon: "󰥔", LinesIcon: "󰆓"},
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
			Token:     TokenSegment{TokensIcon: "󰹑"},
			Exceeds:   Segment{Icon: "󰄘"},
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
