package main

import (
	"encoding/json/v2"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/fatih/color"
)

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

func main() {
	color.NoColor = false
	cfg := loadConfig()

	data := readPayload()
	segments := []string{
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
	}
	var nonEmpty []string
	for _, s := range segments {
		if s != "" {
			nonEmpty = append(nonEmpty, s)
		}
	}
	fmt.Println(strings.Join(nonEmpty, termSep()))
}
