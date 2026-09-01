package main

import (
	"fmt"
	"math"
	"os/exec"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/fatih/color"
)

var (
	colorDim     = color.New(color.Faint)
	colorBlue    = color.New(color.FgCyan)
	colorGreen   = color.New(color.FgGreen)
	colorYellow  = color.New(color.FgYellow)
	colorRed     = color.New(color.FgRed)
	colorMagenta = color.New(color.FgMagenta)
	colorGray    = color.New(color.FgHiBlack)
	colorCyan    = color.New(color.FgCyan)

	branchRegex = regexp.MustCompile(`^(?:No commits yet on )?([^\.\n]+)`)
)

func usedPct(pct *float64) float64 {
	if pct == nil {
		return 0
	}
	return *pct
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

func modelSegment(cfg *Config, data Payload) string {
	if !cfg.Toggles.Model {
		return ""
	}
	name := data.Model.DisplayName
	if name == "" {
		name = data.Model.Id
	}
	if name == "" {
		name = "unknown"
	}
	return colorBlue.Add(color.Bold).Sprintf("%s %s", cfg.Icons.Model, name)
}

func folderSegment(cfg *Config, data Payload) string {
	if !cfg.Toggles.Folder {
		return ""
	}
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

	return colorDim.Sprintf("%s %s", cfg.Icons.Folder, folder)
}

func gitSegment(cfg *Config, data Payload) string {
	if !cfg.Toggles.Git {
		return ""
	}
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

	segment := colorMagenta.Sprintf("%s %s", cfg.Icons.Branch, branch)
	if dirty {
		segment += colorYellow.Sprint("*")
	}

	return segment
}

func worktreeSegment(cfg *Config, data Payload) string {
	if !cfg.Toggles.Worktree {
		return ""
	}
	if data.Worktree.Name == "" {
		return ""
	}
	return colorDim.Sprintf("%s %s", cfg.Icons.Worktree, data.Worktree.Name)
}

func contextSegment(cfg *Config, data Payload) string {
	if !cfg.Toggles.Context {
		return ""
	}
	pct := usedPct(data.ContextWindow.UsedPercentage)
	percent := math.Floor(pct)
	filled := int(math.Min(10, math.Max(0, math.Floor(percent/10))))

	selectedColor := colorGreen
	if percent >= 85 {
		selectedColor = colorRed
	} else if percent >= 60 {
		selectedColor = colorYellow
	}

	bar := strings.Repeat("#", filled) + strings.Repeat("-", 10-filled)
	return selectedColor.Sprintf("%s [%s] %.0f%%", cfg.Icons.Context, bar, percent)
}

func costSegment(cfg *Config, data Payload) string {
	if !cfg.Toggles.Cost {
		return ""
	}
	if data.Cost.TotalCostUSD == nil {
		return ""
	}
	segment := colorDim.Sprintf("%s $%.2f", cfg.Icons.Cost, *data.Cost.TotalCostUSD)

	if data.Cost.TotalDurationMs > 0 {
		totalSecs := data.Cost.TotalDurationMs / 1000
		mins := totalSecs / 60
		secs := totalSecs % 60
		segment += colorDim.Sprintf(" %s%dm%ds", cfg.Icons.Duration, mins, secs)
	}

	if data.Cost.TotalLinesAdded > 0 || data.Cost.TotalLinesRemoved > 0 {
		segment += colorDim.Sprintf(
			" %s+%d-%d",
			cfg.Icons.Lines,
			data.Cost.TotalLinesAdded,
			data.Cost.TotalLinesRemoved,
		)
	}

	return segment
}

func vimSegment(cfg *Config, data Payload) string {
	if !cfg.Toggles.Vim {
		return ""
	}
	if data.Vim.Mode == "" {
		return ""
	}

	icon := cfg.Icons.VimNormal
	if data.Vim.Mode == "INSERT" {
		icon = cfg.Icons.VimInsert
	}

	return colorDim.Sprintf("%s ", icon)
}

func autorunSegment(cfg *Config, data Payload) string {
	if !cfg.Toggles.Autorun {
		return ""
	}
	if !data.Autorun {
		return ""
	}
	return colorYellow.Sprintf("%s auto", cfg.Icons.Autorun)
}

func maxSegment(cfg *Config, data Payload) string {
	if !cfg.Toggles.Max {
		return ""
	}
	if !data.Model.MaxMode {
		return ""
	}
	return colorYellow.Sprintf("%s max", cfg.Icons.Max)
}

func fastModeSegment(cfg *Config, data Payload) string {
	if !cfg.Toggles.FastMode {
		return ""
	}
	if !data.FastMode {
		return ""
	}
	return colorCyan.Sprintf("%s fast", cfg.Icons.FastMode)
}

func effortSegment(cfg *Config, data Payload) string {
	if !cfg.Toggles.Effort {
		return ""
	}
	if data.Effort == nil || data.Effort.Level == "" {
		return ""
	}
	return colorDim.Sprintf("%s %s", cfg.Icons.Effort, data.Effort.Level)
}

func thinkingSegment(cfg *Config, data Payload) string {
	if !cfg.Toggles.Thinking {
		return ""
	}
	if data.Thinking == nil || !data.Thinking.Enabled {
		return ""
	}
	return colorDim.Sprintf("%s think", cfg.Icons.Thinking)
}

func sessionSegment(cfg *Config, data Payload) string {
	if !cfg.Toggles.Session {
		return ""
	}
	name := data.SessionName
	if name == "" {
		return ""
	}
	if len(name) > 24 {
		name = name[:24] + "…"
	}
	return colorDim.Sprintf("%s %s", cfg.Icons.Session, name)
}

func prSegment(cfg *Config, data Payload) string {
	if !cfg.Toggles.PR {
		return ""
	}
	if data.PR == nil || data.PR.Number == 0 {
		return ""
	}
	label := fmt.Sprintf("#%d", data.PR.Number)
	if data.PR.ReviewState != nil {
		switch *data.PR.ReviewState {
		case "approved":
			label += " " + colorGreen.Sprint("✓")
		case "changes_requested":
			label += " " + colorRed.Sprint("✗")
		case "draft":
			label += " " + colorDim.Sprint("○")
		default:
			label += " " + colorYellow.Sprint("●")
		}
	}
	return colorMagenta.Sprintf("%s %s", cfg.Icons.PR, label)
}

func agentSegment(cfg *Config, data Payload) string {
	if !cfg.Toggles.Agent {
		return ""
	}
	if data.Agent == nil || data.Agent.Name == "" {
		return ""
	}
	return colorDim.Sprintf("%s %s", cfg.Icons.Agent, data.Agent.Name)
}

func rateLimitSegment(cfg *Config, data Payload) string {
	if !cfg.Toggles.RateLimit {
		return ""
	}
	if data.RateLimits == nil {
		return ""
	}
	parts := []string{}
	if rl := data.RateLimits.FiveHour; rl != nil {
		parts = append(parts, fmt.Sprintf("5h:%.0f%%", rl.UsedPercentage))
	}
	if rl := data.RateLimits.SevenDay; rl != nil {
		parts = append(parts, fmt.Sprintf("7d:%.0f%%", rl.UsedPercentage))
	}
	if rl := data.RateLimits.SpendLimit; rl != nil {
		parts = append(parts, fmt.Sprintf("$:%.0f%%", rl.UsedPercentage))
	}
	if len(parts) == 0 {
		return ""
	}
	return colorDim.Sprintf("%s %s", cfg.Icons.RateLimit, strings.Join(parts, " "))
}

func tokenSegment(cfg *Config, data Payload) string {
	if !cfg.Toggles.Token {
		return ""
	}
	if data.ContextWindow.TotalInputTokens == 0 && data.ContextWindow.TotalOutputTokens == 0 {
		return ""
	}
	in := data.ContextWindow.TotalInputTokens
	out := data.ContextWindow.TotalOutputTokens
	return colorDim.Sprintf("%s %dk/%dk", cfg.Icons.Tokens, in/1000, out/1000)
}

func cacheSegment(cfg *Config, data Payload) string {
	if !cfg.Toggles.Cache {
		return ""
	}
	if data.PromptCache == nil {
		return ""
	}
	label := ""
	if data.PromptCache.Warm {
		label = colorGreen.Sprint("warm")
	} else {
		label = colorDim.Sprint("cold")
	}
	if data.PromptCache.HitRatio != nil {
		hit := math.Floor(*data.PromptCache.HitRatio * 100)
		label += fmt.Sprintf(" %.0f%%", hit)
	}
	return colorDim.Sprintf("%s %s", cfg.Icons.Cache, label)
}

func exceedsSegment(cfg *Config, data Payload) string {
	if !cfg.Toggles.Exceeds {
		return ""
	}
	if !data.Exceeds200k {
		return ""
	}
	return colorYellow.Sprintf("%s >200k", cfg.Icons.Exceeds)
}
