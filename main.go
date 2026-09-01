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

	data := readPayload()
	segments := []string{
		modelSegment(data),
		folderSegment(data),
		gitSegment(data),
		worktreeSegment(data),
		prSegment(data),
		agentSegment(data),
		sessionSegment(data),
		contextSegment(data),
		costSegment(data),
		tokenSegment(data),
		cacheSegment(data),
		vimSegment(data),
		autorunSegment(data),
		maxSegment(data),
		fastModeSegment(data),
		effortSegment(data),
		thinkingSegment(data),
		rateLimitSegment(data),
		exceedsSegment(data),
	}
	var nonEmpty []string
	for _, s := range segments {
		if s != "" {
			nonEmpty = append(nonEmpty, s)
		}
	}
	fmt.Println(strings.Join(nonEmpty, termSep()))
}
