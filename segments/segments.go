package segments

import (
	"fmt"
	"log/slog"
	"math"
	"path/filepath"
	"strings"

	"github.com/5c077m4n/cli-harness-statusline/config"
	"github.com/5c077m4n/cli-harness-statusline/types"
	"github.com/fatih/color"
	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
)

func usedPct(percent *float64) float64 {
	if percent == nil {
		return 0
	}
	return *percent
}

func truncate(s string, maxLen int) string {
	if maxLen <= 0 || len(s) <= maxLen {
		return s
	}
	prefix := maxLen / 2
	suffix := maxLen - prefix
	return s[:prefix] + "…" + s[len(s)-suffix:]
}

func TermSep() string {
	return colorGray.Sprint(" | ")
}

type gitStatus struct {
	branch   string
	detached bool
	dirty    bool
}

func gitInfo(directory string) gitStatus {
	repo, err := git.PlainOpenWithOptions(directory, &git.PlainOpenOptions{
		DetectDotGit:          true,
		EnableDotGitCommonDir: true,
	})
	if err != nil {
		slog.Warn(
			"gitInfo: failed to open git repository",
			slog.String("dir", directory),
			slog.Any("error", err),
		)
		return gitStatus{}
	}

	worktree, err := repo.Worktree()
	if err != nil {
		slog.Warn(
			"gitInfo: failed to get worktree",
			slog.String("dir", directory),
			slog.Any("error", err),
		)
		return gitStatus{}
	}

	status, err := worktree.Status()
	if err != nil {
		slog.Warn(
			"gitInfo: failed to get status",
			slog.String("dir", directory),
			slog.Any("error", err),
		)
		return gitStatus{}
	}
	dirty := !status.IsClean()

	head, err := repo.Head()
	if err != nil {
		slog.Warn(
			"gitInfo: failed to resolve HEAD",
			slog.String("dir", directory),
			slog.Any("error", err),
		)
		return gitStatus{}
	}

	if head.Name() != plumbing.HEAD {
		return gitStatus{branch: head.Name().Short(), dirty: dirty}
	}

	sha := head.Hash().String()[:7]
	if sha == "" {
		return gitStatus{}
	}
	return gitStatus{branch: sha, detached: true, dirty: dirty}
}

func model(cfg *config.Config, data *types.Payload) string {
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
	return colorBlue.Add(color.Bold).Sprintf("%s %s", IconModel, name)
}

func folder(cfg *config.Config, data *types.Payload) string {
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

	if !cfg.Segments.Folder.DisableTruncate {
		folder = truncate(folder, defaultTruncateLength)
	}

	return colorDim.Sprintf("%s %s", IconFolder, folder)
}

func gitBase(cfg *config.Config, data *types.Payload) string {
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

	info := gitInfo(directory)
	if info.branch == "" {
		return ""
	}

	branch := info.branch
	if !cfg.Segments.Git.DisableTruncate {
		branch = truncate(branch, defaultTruncateLength)
	}

	if info.detached {
		segment := colorYellow.Sprintf("%s @%s", IconGitBranch, branch)
		if info.dirty {
			segment += colorYellow.Sprint("*")
		}
		return segment
	}

	segment := colorMagenta.Sprintf("%s %s", IconGitBranch, branch)
	if info.dirty {
		segment += colorYellow.Sprint("*")
	}

	return segment
}

func gitWorktree(cfg *config.Config, data *types.Payload) string {
	if cfg.Segments.Worktree.Disable {
		return ""
	}
	if data.Worktree.Name == "" {
		return ""
	}
	return colorDim.Sprintf("%s %s", IconWorktree, data.Worktree.Name)
}

func context(cfg *config.Config, data *types.Payload) string {
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

	bar := strings.Repeat(IconBarFilled, filled) + strings.Repeat(IconBarEmpty, 10-filled)
	contextSegment := selectedColor.Sprintf(
		"%s [%s] %.0f%%",
		IconContext,
		bar,
		percent,
	)

	if !cfg.Segments.Exceeds.Disable && data.Exceeds200k {
		return contextSegment + colorRed.Sprintf(" > 200K(!)")
	}
	return contextSegment
}

func cost(cfg *config.Config, data *types.Payload) string {
	if cfg.Segments.Cost.Disable {
		return ""
	}
	if data.Cost.TotalCostUSD == nil {
		return ""
	}
	segment := colorDim.Sprintf("%s $%.2f", IconCost, *data.Cost.TotalCostUSD)

	if data.Cost.TotalDurationMs > 0 {
		totalSecs := int64(data.Cost.TotalDurationMs / 1000)
		segment += colorDim.Sprintf(" %s%s", IconDuration, formatDuration(totalSecs))
	}

	if data.Cost.TotalLinesAdded > 0 || data.Cost.TotalLinesRemoved > 0 {
		segment += colorDim.Sprintf(
			" %s+%d-%d",
			IconLines,
			data.Cost.TotalLinesAdded,
			data.Cost.TotalLinesRemoved,
		)
	}

	return segment
}

func formatDuration(totalSecs int64) string {
	if totalSecs >= secondsPerDay {
		days := totalSecs / secondsPerDay
		totalSecs %= secondsPerDay
		hours := totalSecs / secondsPerHour
		totalSecs %= secondsPerHour
		mins := totalSecs / secondsPerMinute
		secs := totalSecs % secondsPerMinute
		return fmt.Sprintf("%02dd%02dh%02dm%02ds", days, hours, mins, secs)
	}
	hours := totalSecs / secondsPerHour
	totalSecs %= secondsPerHour
	mins := totalSecs / secondsPerMinute
	secs := totalSecs % secondsPerMinute
	if hours > 0 {
		return fmt.Sprintf("%02dh%02dm%02ds", hours, mins, secs)
	}
	return fmt.Sprintf("%dm%ds", mins, secs)
}

func vim(cfg *config.Config, data *types.Payload) string {
	if cfg.Segments.Vim.Disable {
		return ""
	}
	if data.Vim.Mode == "" {
		return ""
	}

	icon := IconVimNormal
	if data.Vim.Mode == "INSERT" {
		icon = IconVimInsert
	}

	return colorDim.Sprintf("%s", icon)
}

func autorun(cfg *config.Config, data *types.Payload) string {
	if cfg.Segments.Autorun.Disable {
		return ""
	}
	if !data.Autorun {
		return ""
	}
	return colorYellow.Sprintf("%s auto", IconAutorun)
}

func max(cfg *config.Config, data *types.Payload) string {
	if cfg.Segments.Max.Disable {
		return ""
	}
	if !data.Model.MaxMode {
		return ""
	}
	return colorYellow.Sprintf("%s max", IconMax)
}

func fastMode(cfg *config.Config, data *types.Payload) string {
	if cfg.Segments.FastMode.Disable {
		return ""
	}
	if !data.FastMode {
		return ""
	}
	return colorCyan.Sprintf("%s fast", IconFastMode)
}

func effort(cfg *config.Config, data *types.Payload) string {
	if cfg.Segments.Effort.Disable {
		return ""
	}
	if data.Effort == nil || data.Effort.Level == "" {
		return ""
	}
	return colorDim.Sprintf("%s %s", IconEffort, data.Effort.Level)
}

func thinking(cfg *config.Config, data *types.Payload) string {
	if cfg.Segments.Thinking.Disable {
		return ""
	}
	if data.Thinking == nil || !data.Thinking.Enabled {
		return ""
	}
	return colorDim.Sprintf("%s think", IconThinking)
}

func session(cfg *config.Config, data *types.Payload) string {
	if cfg.Segments.Session.Disable {
		return ""
	}
	name := truncate(data.SessionName, 24)
	if name == "" {
		return ""
	}
	return colorDim.Sprintf("%s %s", IconSession, name)
}

func pr(cfg *config.Config, data *types.Payload) string {
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
		label = strings.Join([]string{OSC8Start, data.PR.URL, OSC8Sep, label, OSC8End}, "")
	}
	return colorMagenta.Sprintf("%s %s", IconPR, label)
}

func agent(cfg *config.Config, data *types.Payload) string {
	if cfg.Segments.Agent.Disable {
		return ""
	}
	if data.Agent == nil || data.Agent.Name == "" {
		return ""
	}
	return colorDim.Sprintf("%s %s", IconAgent, data.Agent.Name)
}

func rateLimit(cfg *config.Config, data *types.Payload) string {
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
	return colorDim.Sprintf("%s %s", IconRateLimit, strings.Join(parts, " "))
}

func token(cfg *config.Config, data *types.Payload) string {
	if cfg.Segments.Token.Disable {
		return ""
	}
	if data.ContextWindow.TotalInputTokens == 0 && data.ContextWindow.TotalOutputTokens == 0 {
		return ""
	}
	in := data.ContextWindow.TotalInputTokens
	out := data.ContextWindow.TotalOutputTokens
	return colorDim.Sprintf("%s%dk %s%dk", IconTokenIn, in/1000, IconTokenOut, out/1000)
}

func cache(cfg *config.Config, data *types.Payload) string {
	if cfg.Segments.Cache.Disable {
		return ""
	}
	if data.PromptCache == nil {
		return ""
	}

	var format *color.Color
	base := ""
	if data.PromptCache.Warm {
		format = colorOrange
		base = fmt.Sprintf("%s warm", IconCache)
	} else {
		format = colorLightBlue
		base = fmt.Sprintf("%s cold", IconCache)
	}
	if data.PromptCache.HitRatio != nil {
		hit := math.Floor(*data.PromptCache.HitRatio * 100)
		base += fmt.Sprintf(" %.0f%%", hit)
	}
	return format.Sprint(base)
}

func Render(cfg *config.Config, data *types.Payload) string {
	segmentFuncs := [...]func(*config.Config, *types.Payload) string{
		model,
		folder,
		gitBase,
		gitWorktree,
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
	}

	nonEmpty := make([]string, 0, len(segmentFuncs))
	for _, segFunc := range segmentFuncs {
		if s := segFunc(cfg, data); s != "" {
			nonEmpty = append(nonEmpty, s)
		}
	}
	return strings.Join(nonEmpty, TermSep())
}
