package segments

import (
	"fmt"
	"log/slog"
	"math"
	"os/exec"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/5c077m4n/cursor-agent-statusline/config"
	"github.com/5c077m4n/cursor-agent-statusline/types"
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

const (
	osc8Start = "\033]8;;"
	osc8Sep   = "\033\\"
	osc8End   = "\033]8;;\033\\"
)

func usedPct(pct *float64) float64 {
	if pct == nil {
		return 0
	}
	return *pct
}

func TermSep() string {
	return colorGray.Sprint(" | ")
}

func gitInfo(directory string) (string, bool) {
	cmd := exec.Command("git", "-C", directory, "status", "--branch", "--porcelain")
	outputBytes, err := cmd.Output()
	if err != nil {
		slog.Warn(
			"gitInfo: git command failed",
			slog.String("dir", directory),
			slog.Any("error", err),
		)
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

func model(cfg *config.Config, data types.Payload) string {
	if cfg.Segments.Model.Disable {
		return ""
	}
	name := data.Model.DisplayName
	if name == "" {
		name = data.Model.ID
	}
	if name == "" {
		name = "unknown"
	}
	return colorBlue.Add(color.Bold).Sprintf("%s %s", cfg.Segments.Model.Icon, name)
}

func folder(cfg *config.Config, data types.Payload) string {
	if cfg.Segments.Folder.Disable {
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

	return colorDim.Sprintf("%s %s", cfg.Segments.Folder.Icon, folder)
}

func git(cfg *config.Config, data types.Payload) string {
	if cfg.Segments.Git.Disable {
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

	segment := colorMagenta.Sprintf("%s %s", cfg.Segments.Git.BranchIcon, branch)
	if dirty {
		segment += colorYellow.Sprint("*")
	}

	return segment
}

func worktree(cfg *config.Config, data types.Payload) string {
	if cfg.Segments.Worktree.Disable {
		return ""
	}
	if data.Worktree.Name == "" {
		return ""
	}
	return colorDim.Sprintf("%s %s", cfg.Segments.Worktree.Icon, data.Worktree.Name)
}

func context(cfg *config.Config, data types.Payload) string {
	if cfg.Segments.Context.Disable {
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
	return selectedColor.Sprintf("%s [%s] %.0f%%", cfg.Segments.Context.Icon, bar, percent)
}

func cost(cfg *config.Config, data types.Payload) string {
	if cfg.Segments.Cost.Disable {
		return ""
	}
	if data.Cost.TotalCostUSD == nil {
		return ""
	}
	segment := colorDim.Sprintf("%s $%.2f", cfg.Segments.Cost.Icon, *data.Cost.TotalCostUSD)

	if data.Cost.TotalDurationMs > 0 {
		totalSecs := int64(data.Cost.TotalDurationMs / 1000)
		mins := totalSecs / 60
		secs := totalSecs % 60
		segment += colorDim.Sprintf(" %s%dm%ds", cfg.Segments.Cost.DurationIcon, mins, secs)
	}

	if data.Cost.TotalLinesAdded > 0 || data.Cost.TotalLinesRemoved > 0 {
		segment += colorDim.Sprintf(
			" %s+%d-%d",
			cfg.Segments.Cost.LinesIcon,
			data.Cost.TotalLinesAdded,
			data.Cost.TotalLinesRemoved,
		)
	}

	return segment
}

func vim(cfg *config.Config, data types.Payload) string {
	if cfg.Segments.Vim.Disable {
		return ""
	}
	if data.Vim.Mode == "" {
		return ""
	}

	icon := cfg.Segments.Vim.IconNormal
	if data.Vim.Mode == "INSERT" {
		icon = cfg.Segments.Vim.IconInsert
	}

	return colorDim.Sprintf("%s ", icon)
}

func autorun(cfg *config.Config, data types.Payload) string {
	if cfg.Segments.Autorun.Disable {
		return ""
	}
	if !data.Autorun {
		return ""
	}
	return colorYellow.Sprintf("%s auto", cfg.Segments.Autorun.Icon)
}

func max(cfg *config.Config, data types.Payload) string {
	if cfg.Segments.Max.Disable {
		return ""
	}
	if !data.Model.MaxMode {
		return ""
	}
	return colorYellow.Sprintf("%s max", cfg.Segments.Max.Icon)
}

func fastMode(cfg *config.Config, data types.Payload) string {
	if cfg.Segments.FastMode.Disable {
		return ""
	}
	if !data.FastMode {
		return ""
	}
	return colorCyan.Sprintf("%s fast", cfg.Segments.FastMode.Icon)
}

func effort(cfg *config.Config, data types.Payload) string {
	if cfg.Segments.Effort.Disable {
		return ""
	}
	if data.Effort == nil || data.Effort.Level == "" {
		return ""
	}
	return colorDim.Sprintf("%s %s", cfg.Segments.Effort.Icon, data.Effort.Level)
}

func thinking(cfg *config.Config, data types.Payload) string {
	if cfg.Segments.Thinking.Disable {
		return ""
	}
	if data.Thinking == nil || !data.Thinking.Enabled {
		return ""
	}
	return colorDim.Sprintf("%s think", cfg.Segments.Thinking.Icon)
}

func session(cfg *config.Config, data types.Payload) string {
	if cfg.Segments.Session.Disable {
		return ""
	}
	name := data.SessionName
	if name == "" {
		return ""
	}
	if len(name) > 24 {
		name = name[:24] + "…"
	}
	return colorDim.Sprintf("%s %s", cfg.Segments.Session.Icon, name)
}

func pr(cfg *config.Config, data types.Payload) string {
	if cfg.Segments.PR.Disable {
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
	if data.PR.URL != "" {
		label = strings.Join([]string{osc8Start, data.PR.URL, osc8Sep, label, osc8End}, "")
	}
	return colorMagenta.Sprintf("%s %s", cfg.Segments.PR.Icon, label)
}

func agent(cfg *config.Config, data types.Payload) string {
	if cfg.Segments.Agent.Disable {
		return ""
	}
	if data.Agent == nil || data.Agent.Name == "" {
		return ""
	}
	return colorDim.Sprintf("%s %s", cfg.Segments.Agent.Icon, data.Agent.Name)
}

func rateLimit(cfg *config.Config, data types.Payload) string {
	if cfg.Segments.RateLimit.Disable {
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
	return colorDim.Sprintf("%s %s", cfg.Segments.RateLimit.Icon, strings.Join(parts, " "))
}

func token(cfg *config.Config, data types.Payload) string {
	if cfg.Segments.Token.Disable {
		return ""
	}
	if data.ContextWindow.TotalInputTokens == 0 && data.ContextWindow.TotalOutputTokens == 0 {
		return ""
	}
	in := data.ContextWindow.TotalInputTokens
	out := data.ContextWindow.TotalOutputTokens
	return colorDim.Sprintf("%s %dk/%dk", cfg.Segments.Token.TokensIcon, in/1000, out/1000)
}

func cache(cfg *config.Config, data types.Payload) string {
	if cfg.Segments.Cache.Disable {
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
	return colorDim.Sprintf("%s %s", cfg.Segments.Cache.Icon, label)
}

func exceeds(cfg *config.Config, data types.Payload) string {
	if cfg.Segments.Exceeds.Disable {
		return ""
	}
	if !data.Exceeds200k {
		return ""
	}
	return colorYellow.Sprintf("%s >200k", cfg.Segments.Exceeds.Icon)
}

func Render(cfg *config.Config, data types.Payload) string {
	segmentFuncs := [...]func(*config.Config, types.Payload) string{
		model,
		folder,
		git,
		worktree,
		pr,
		agent,
		session,
		context,
		cost,
		token,
		cache,
		vim,
		autorun,
		max,
		fastMode,
		effort,
		thinking,
		rateLimit,
		exceeds,
	}

	nonEmpty := make([]string, 0, len(segmentFuncs))
	for _, segFunc := range segmentFuncs {
		if s := segFunc(cfg, data); s != "" {
			nonEmpty = append(nonEmpty, s)
		}
	}
	return strings.Join(nonEmpty, TermSep())
}
