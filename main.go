package main

import (
	"encoding/json/v2"
	"fmt"
	"io"
	"math"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/fatih/color"
)

const (
	iconModel     = "󰚩"
	iconFolder    = "󰉋"
	iconBranch    = "󰊢"
	iconWorktree  = "󰙅"
	iconContext   = "󰧑"
	iconCost      = "󰆛"
	iconVimInsert = "󰏫"
	iconVimNormal = "󰌌"
	iconAutorun   = "󰄙"
	iconMax       = "󰓅"
)

var (
	colorDim     = color.New(color.Faint)
	colorBlue    = color.New(color.FgCyan)
	colorGreen   = color.New(color.FgGreen)
	colorYellow  = color.New(color.FgYellow)
	colorRed     = color.New(color.FgRed)
	colorMagenta = color.New(color.FgMagenta)
	colorGray    = color.New(color.FgHiBlack)

	branchRegex = regexp.MustCompile(`^(?:No commits yet on )?([^\.\n]+)`)
)

type ModelInfo struct {
	DisplayName string `json:"display_name"`
	MaxMode     bool   `json:"max_mode"`
}

type WorkspaceInfo struct {
	CurrentDir string `json:"current_dir"`
}

type WorktreeInfo struct {
	Name string `json:"name"`
}

type ContextWindow struct {
	UsedPercentage float64 `json:"used_percentage"`
}

type CostInfo struct {
	TotalCostUSD *float64 `json:"total_cost_usd"`
}

type VimInfo struct {
	Mode string `json:"mode"`
}

type Payload struct {
	Cwd       string        `json:"cwd"`
	Autorun   bool          `json:"autorun"`
	Model     ModelInfo     `json:"model"`
	Workspace WorkspaceInfo `json:"workspace"`
	Worktree  WorktreeInfo  `json:"worktree"`
	Context   ContextWindow `json:"context_window"`
	Cost      CostInfo      `json:"cost"`
	Vim       VimInfo       `json:"vim"`
}

func termSep() string {
	return colorGray.Sprint(" | ")
}

func gitInfo(directory string) (string, bool) {
	cmd := exec.Command("git", "-C", directory, "status", "--porcelain", "-b")
	outputBytes, err := cmd.Output()
	if err != nil {
		return "", false
	}

	output := strings.TrimSpace(string(outputBytes))
	if output == "" {
		return "", false
	}

	lines := strings.SplitN(output, "\n", 2)
	header := lines[0]
	if !strings.HasPrefix(header, "## ") {
		return "", false
	}

	header = strings.TrimPrefix(header, "## ")
	if header == "HEAD (no branch)" {
		return "", false
	}

	if idx := strings.Index(header, "..."); idx != -1 {
		header = header[:idx]
	}

	match := branchRegex.FindStringSubmatch(header)
	if len(match) < 2 {
		return "", false
	}

	dirty := len(lines) > 1 && lines[1] != ""
	return match[1], dirty
}

func modelSegment(data Payload) string {
	name := data.Model.DisplayName
	if name == "" {
		name = "unknown"
	}
	return colorBlue.Add(color.Bold).Sprintf("%s %s", iconModel, name)
}

func folderSegment(data Payload) string {
	directory := data.Workspace.CurrentDir
	if directory == "" {
		directory = data.Cwd
	}
	if directory == "" {
		directory = "."
	}

	folder := filepath.Base(directory)
	if folder == "." || folder == "/" {
		folder = "/"
	}

	return termSep() + colorDim.Sprintf("%s %s", iconFolder, folder)
}

func gitSegment(data Payload) string {
	directory := data.Workspace.CurrentDir
	if directory == "" {
		directory = data.Cwd
	}
	if directory == "" {
		directory = "."
	}

	branch, dirty := gitInfo(directory)
	if branch == "" {
		return ""
	}

	segment := termSep() + colorMagenta.Sprintf("%s %s", iconBranch, branch)
	if dirty {
		segment += colorYellow.Sprint("*")
	}

	return segment
}

func worktreeSegment(data Payload) string {
	if data.Worktree.Name == "" {
		return ""
	}
	return termSep() + colorDim.Sprintf("%s %s", iconWorktree, data.Worktree.Name)
}

func contextSegment(data Payload) string {
	percent := math.Floor(data.Context.UsedPercentage)
	filled := int(math.Min(10, math.Max(0, math.Floor(percent/10))))

	selectedColor := colorGreen
	if percent >= 85 {
		selectedColor = colorRed
	} else if percent >= 60 {
		selectedColor = colorYellow
	}

	bar := strings.Repeat("#", filled) + strings.Repeat("-", 10-filled)
	return termSep() + selectedColor.Sprintf("%s [%s] %.0f%%", iconContext, bar, percent)
}

func costSegment(data Payload) string {
	if data.Cost.TotalCostUSD == nil {
		return ""
	}
	return termSep() + colorDim.Sprintf("%s $%.2f", iconCost, *data.Cost.TotalCostUSD)
}

func vimSegment(data Payload) string {
	if data.Vim.Mode == "" {
		return ""
	}

	icon := iconVimNormal
	if data.Vim.Mode == "INSERT" {
		icon = iconVimInsert
	}

	return termSep() + colorDim.Sprintf("%s ", icon)
}

func autorunSegment(data Payload) string {
	if !data.Autorun {
		return ""
	}
	return termSep() + colorYellow.Sprintf("%s auto", iconAutorun)
}

func maxSegment(data Payload) string {
	if !data.Model.MaxMode {
		return ""
	}
	return termSep() + colorYellow.Sprintf("%s max", iconMax)
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

func main() {
	data := readPayload()
	segments := []string{
		modelSegment(data),
		folderSegment(data),
		gitSegment(data),
		worktreeSegment(data),
		contextSegment(data),
		costSegment(data),
		vimSegment(data),
		autorunSegment(data),
		maxSegment(data),
	}
	fmt.Println(strings.Join(segments, ""))
}
