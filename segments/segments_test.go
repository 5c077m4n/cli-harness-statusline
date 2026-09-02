package segments

import (
	"math"
	"strings"
	"testing"

	"github.com/5c077m4n/cli-harness-statusline/config"
	"github.com/5c077m4n/cli-harness-statusline/types"
	"github.com/fatih/color"
	"github.com/stretchr/testify/assert"
)

var testCfg *config.Config

func init() {
	color.NoColor = true
	testCfg = config.Load()
}

func TestTermSep(t *testing.T) {
	assert.Equal(t, " | ", TermSep())
}

func TestModel(t *testing.T) {
	tests := []struct {
		name string
		data types.Payload
		want string
	}{
		{
			name: "with display name",
			data: types.Payload{Model: types.ModelInfo{DisplayName: "gpt-4"}},
			want: testCfg.Segments.Model.Icon + " gpt-4",
		},
		{
			name: "empty display name",
			data: types.Payload{},
			want: testCfg.Segments.Model.Icon + " unknown",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.want, model(testCfg, tt.data))
		})
	}
}

func TestFolder(t *testing.T) {
	tests := []struct {
		name string
		data types.Payload
		want string
	}{
		{
			name: "current dir set",
			data: types.Payload{Workspace: types.WorkspaceInfo{CurrentDir: "/home/user/project"}},
			want: testCfg.Segments.Folder.Icon + " project",
		},
		{
			name: "falls back to cwd",
			data: types.Payload{Cwd: "/home/user/other"},
			want: testCfg.Segments.Folder.Icon + " other",
		},
		{
			name: "both empty maps to root",
			data: types.Payload{},
			want: testCfg.Segments.Folder.Icon + " /",
		},
		{
			name: "root path",
			data: types.Payload{Workspace: types.WorkspaceInfo{CurrentDir: "/"}},
			want: testCfg.Segments.Folder.Icon + " /",
		},
		{
			name: "dot path maps to root",
			data: types.Payload{Workspace: types.WorkspaceInfo{CurrentDir: "."}},
			want: testCfg.Segments.Folder.Icon + " /",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.want, folder(testCfg, tt.data))
		})
	}
}

func TestWorktree(t *testing.T) {
	tests := []struct {
		name string
		data types.Payload
		want string
	}{
		{
			name: "with worktree name",
			data: types.Payload{Worktree: types.WorktreeInfo{Name: "feature-branch"}},
			want: testCfg.Segments.Worktree.Icon + " feature-branch",
		},
		{name: "empty worktree name", data: types.Payload{}, want: ""},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.want, worktree(testCfg, tt.data))
		})
	}
}

func TestContext(t *testing.T) {
	tests := []struct {
		name string
		data types.Payload
		want string
	}{
		{
			name: "0 percent",
			data: types.Payload{ContextWindow: types.ContextWindowInfo{UsedPercentage: new(0.0)}},
			want: testCfg.Segments.Context.Icon + " [----------] 0%",
		},
		{
			name: "50 percent",
			data: types.Payload{ContextWindow: types.ContextWindowInfo{UsedPercentage: new(50.0)}},
			want: testCfg.Segments.Context.Icon + " [#####-----] 50%",
		},
		{
			name: "100 percent",
			data: types.Payload{ContextWindow: types.ContextWindowInfo{UsedPercentage: new(100.0)}},
			want: testCfg.Segments.Context.Icon + " [##########] 100%",
		},
		{
			name: "negative percent",
			data: types.Payload{ContextWindow: types.ContextWindowInfo{UsedPercentage: new(-10.0)}},
			want: testCfg.Segments.Context.Icon + " [----------] -10%",
		},
		{
			name: "over 100 percent",
			data: types.Payload{ContextWindow: types.ContextWindowInfo{UsedPercentage: new(150.0)}},
			want: testCfg.Segments.Context.Icon + " [##########] 150%",
		},
		{
			name: "fractional percent floors",
			data: types.Payload{ContextWindow: types.ContextWindowInfo{UsedPercentage: new(59.9)}},
			want: testCfg.Segments.Context.Icon + " [#####-----] 59%",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.want, context(testCfg, tt.data))
		})
	}
}

func TestCost(t *testing.T) {
	tests := []struct {
		name string
		data types.Payload
		want string
	}{
		{
			name: "with cost",
			data: types.Payload{Cost: types.CostInfo{TotalCostUSD: new(0.42)}},
			want: testCfg.Segments.Cost.Icon + " $0.42",
		},
		{
			name: "zero cost",
			data: types.Payload{Cost: types.CostInfo{TotalCostUSD: new(0.0)}},
			want: testCfg.Segments.Cost.Icon + " $0.00",
		},
		{name: "nil cost", data: types.Payload{}, want: ""},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.want, cost(testCfg, tt.data))
		})
	}
}

func TestVim(t *testing.T) {
	tests := []struct {
		name string
		data types.Payload
		want string
	}{
		{
			name: "insert mode",
			data: types.Payload{Vim: types.VimInfo{Mode: "INSERT"}},
			want: testCfg.Segments.Vim.IconInsert + " ",
		},
		{
			name: "normal mode",
			data: types.Payload{Vim: types.VimInfo{Mode: "NORMAL"}},
			want: testCfg.Segments.Vim.IconNormal + " ",
		},
		{name: "empty mode", data: types.Payload{}, want: ""},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.want, vim(testCfg, tt.data))
		})
	}
}

func TestAutorun(t *testing.T) {
	tests := []struct {
		name string
		data types.Payload
		want string
	}{
		{
			name: "autorun on",
			data: types.Payload{Autorun: true},
			want: testCfg.Segments.Autorun.Icon + " auto",
		},
		{name: "autorun off", data: types.Payload{}, want: ""},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.want, autorun(testCfg, tt.data))
		})
	}
}

func TestMax(t *testing.T) {
	tests := []struct {
		name string
		data types.Payload
		want string
	}{
		{
			name: "max mode on",
			data: types.Payload{Model: types.ModelInfo{MaxMode: true}},
			want: testCfg.Segments.Max.Icon + " max",
		},
		{name: "max mode off", data: types.Payload{}, want: ""},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.want, max(testCfg, tt.data))
		})
	}
}

func TestGitNotAGitRepo(t *testing.T) {
	data := types.Payload{Cwd: t.TempDir(), Workspace: types.WorkspaceInfo{CurrentDir: t.TempDir()}}
	assert.Empty(t, git(testCfg, data))
}

func TestGitInfoNotAGitRepo(t *testing.T) {
	branch, dirty := gitInfo(t.TempDir())
	assert.Empty(t, branch)
	assert.False(t, dirty)
}

func TestContextColorThresholds(t *testing.T) {
	tests := []struct {
		name    string
		percent float64
		want    string
	}{
		{
			name:    "under 60 green",
			percent: 0,
			want:    testCfg.Segments.Context.Icon + " [----------] 0%",
		},
		{
			name:    "under 60 green",
			percent: 59,
			want:    testCfg.Segments.Context.Icon + " [#####-----] 59%",
		},
		{
			name:    "at 60 yellow",
			percent: 60,
			want:    testCfg.Segments.Context.Icon + " [######----] 60%",
		},
		{
			name:    "at 84 yellow",
			percent: 84,
			want:    testCfg.Segments.Context.Icon + " [########--] 84%",
		},
		{name: "at 85 red", percent: 85, want: testCfg.Segments.Context.Icon + " [########--] 85%"},
		{
			name:    "at 100 red",
			percent: 100,
			want:    testCfg.Segments.Context.Icon + " [##########] 100%",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			data := types.Payload{
				ContextWindow: types.ContextWindowInfo{UsedPercentage: new(tt.percent)},
			}
			assert.Contains(t, context(testCfg, data), tt.want)
		})
	}
}

func TestContextBarFilling(t *testing.T) {
	tests := []struct {
		percent float64
		filled  int
	}{
		{percent: 0, filled: 0},
		{percent: 5, filled: 0},
		{percent: 10, filled: 1},
		{percent: 15, filled: 1},
		{percent: 50, filled: 5},
		{percent: 99, filled: 9},
		{percent: 100, filled: 10},
		{percent: 150, filled: 10},
	}
	for _, tt := range tests {
		t.Run("", func(t *testing.T) {
			data := types.Payload{
				ContextWindow: types.ContextWindowInfo{UsedPercentage: new(tt.percent)},
			}
			got := context(testCfg, data)
			expectedBar := strings.Repeat("#", tt.filled) + strings.Repeat("-", 10-tt.filled)
			assert.Contains(t, got, "["+expectedBar+"]")
		})
	}
}

func TestWorktreeEmpty(t *testing.T) {
	assert.Empty(t, worktree(testCfg, types.Payload{Worktree: types.WorktreeInfo{Name: ""}}))
}

func TestFolderCurrentDirPriority(t *testing.T) {
	data := types.Payload{
		Cwd:       "/fallback",
		Workspace: types.WorkspaceInfo{CurrentDir: "/preferred"},
	}
	assert.Equal(t, testCfg.Segments.Folder.Icon+" preferred", folder(testCfg, data))
}

func TestFastMode(t *testing.T) {
	tests := []struct {
		name string
		data types.Payload
		want string
	}{
		{
			name: "fast mode on",
			data: types.Payload{FastMode: true},
			want: testCfg.Segments.FastMode.Icon + " fast",
		},
		{name: "fast mode off", data: types.Payload{}, want: ""},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.want, fastMode(testCfg, tt.data))
		})
	}
}

func TestEffort(t *testing.T) {
	tests := []struct {
		name string
		data types.Payload
		want string
	}{
		{
			name: "with effort",
			data: types.Payload{Effort: &types.EffortInfo{Level: "high"}},
			want: testCfg.Segments.Effort.Icon + " high",
		},
		{name: "nil effort", data: types.Payload{}, want: ""},
		{name: "empty level", data: types.Payload{Effort: &types.EffortInfo{Level: ""}}, want: ""},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.want, effort(testCfg, tt.data))
		})
	}
}

func TestThinking(t *testing.T) {
	tests := []struct {
		name string
		data types.Payload
		want string
	}{
		{
			name: "thinking enabled",
			data: types.Payload{Thinking: &types.ThinkingInfo{Enabled: true}},
			want: testCfg.Segments.Thinking.Icon + " think",
		},
		{
			name: "thinking disabled",
			data: types.Payload{Thinking: &types.ThinkingInfo{Enabled: false}},
			want: "",
		},
		{name: "nil thinking", data: types.Payload{}, want: ""},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.want, thinking(testCfg, tt.data))
		})
	}
}

func TestSession(t *testing.T) {
	tests := []struct {
		name string
		data types.Payload
		want string
	}{
		{
			name: "with session name",
			data: types.Payload{SessionName: "my-session"},
			want: testCfg.Segments.Session.Icon + " my-session",
		},
		{name: "empty session name", data: types.Payload{}, want: ""},
		{
			name: "truncates long name",
			data: types.Payload{SessionName: "abcdefghijklmnopqrstuvwxyz123456"},
			want: testCfg.Segments.Session.Icon + " abcdefghijklmnopqrstuvwx…",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.want, session(testCfg, tt.data))
		})
	}
}

func TestAgent(t *testing.T) {
	tests := []struct {
		name string
		data types.Payload
		want string
	}{
		{
			name: "with agent",
			data: types.Payload{Agent: &types.AgentInfo{Name: "reviewer"}},
			want: testCfg.Segments.Agent.Icon + " reviewer",
		},
		{name: "nil agent", data: types.Payload{}, want: ""},
		{name: "empty name", data: types.Payload{Agent: &types.AgentInfo{Name: ""}}, want: ""},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.want, agent(testCfg, tt.data))
		})
	}
}

func TestPR(t *testing.T) {
	tests := []struct {
		name string
		data types.Payload
		want string
	}{
		{
			name: "PR without review state",
			data: types.Payload{PR: &types.PRInfo{Number: 42, URL: "https://example.com/pr/42"}},
			want: testCfg.Segments.PR.Icon + " " + osc8Start + "https://example.com/pr/42" + osc8Sep + "#42" + osc8End,
		},
		{name: "nil PR", data: types.Payload{}, want: ""},
		{
			name: "PR with approved review",
			data: types.Payload{PR: &types.PRInfo{Number: 1, ReviewState: new("approved")}},
			want: testCfg.Segments.PR.Icon + " #1 ✓",
		},
		{
			name: "PR with changes requested",
			data: types.Payload{
				PR: &types.PRInfo{Number: 2, ReviewState: new("changes_requested")},
			},
			want: testCfg.Segments.PR.Icon + " #2 ✗",
		},
		{
			name: "PR draft",
			data: types.Payload{PR: &types.PRInfo{Number: 3, ReviewState: new("draft")}},
			want: testCfg.Segments.PR.Icon + " #3 ○",
		},
		{
			name: "PR pending",
			data: types.Payload{PR: &types.PRInfo{Number: 4, ReviewState: new("pending")}},
			want: testCfg.Segments.PR.Icon + " #4 ●",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.want, pr(testCfg, tt.data))
		})
	}
}

func TestRateLimit(t *testing.T) {
	tests := []struct {
		name string
		data types.Payload
		want string
	}{
		{name: "nil rate limits", data: types.Payload{}, want: ""},
		{
			name: "with five hour",
			data: types.Payload{
				RateLimits: &types.RateLimits{
					FiveHour: &types.RateLimitWindow{UsedPercentage: 23.5},
				},
			},
			want: testCfg.Segments.RateLimit.Icon + " 5h:24%",
		},
		{
			name: "with all windows",
			data: types.Payload{RateLimits: &types.RateLimits{
				FiveHour:   &types.RateLimitWindow{UsedPercentage: 10},
				SevenDay:   &types.RateLimitWindow{UsedPercentage: 50},
				SpendLimit: &types.RateLimitWindow{UsedPercentage: 75},
			}},
			want: testCfg.Segments.RateLimit.Icon + " 5h:10% 7d:50% $:75%",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.want, rateLimit(testCfg, tt.data))
		})
	}
}

func TestToken(t *testing.T) {
	tests := []struct {
		name string
		data types.Payload
		want string
	}{
		{name: "zero tokens", data: types.Payload{}, want: ""},
		{
			name: "with tokens",
			data: types.Payload{
				ContextWindow: types.ContextWindowInfo{
					TotalInputTokens:  15500,
					TotalOutputTokens: 1200,
				},
			},
			want: testCfg.Segments.Token.IconIn + "15k " + testCfg.Segments.Token.IconOut + "1k",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.want, token(testCfg, tt.data))
		})
	}
}

func TestCache(t *testing.T) {
	tests := []struct {
		name string
		data types.Payload
		want string
	}{
		{name: "nil prompt cache", data: types.Payload{}, want: ""},
		{
			name: "warm cache with hit ratio",
			data: types.Payload{
				PromptCache: &types.PromptCacheInfo{Warm: true, HitRatio: new(0.91)},
			},
			want: testCfg.Segments.Cache.Icon + " warm 91%",
		},
		{
			name: "cold cache",
			data: types.Payload{PromptCache: &types.PromptCacheInfo{Warm: false}},
			want: testCfg.Segments.Cache.Icon + " cold",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.want, cache(testCfg, tt.data))
		})
	}
}

func TestExceeds(t *testing.T) {
	tests := []struct {
		name string
		data types.Payload
		want string
	}{
		{
			name: "exceeds 200k",
			data: types.Payload{Exceeds200k: true},
			want: testCfg.Segments.Exceeds.Icon + " >200k",
		},
		{name: "under 200k", data: types.Payload{}, want: ""},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.want, exceeds(testCfg, tt.data))
		})
	}
}

func TestEdgeCases(t *testing.T) {
	t.Run("model with empty display name and max mode", func(t *testing.T) {
		assert.Equal(
			t,
			testCfg.Segments.Model.Icon+" unknown",
			model(testCfg, types.Payload{Model: types.ModelInfo{MaxMode: true}}),
		)
	})

	t.Run("cost with very large value", func(t *testing.T) {
		assert.Contains(
			t,
			cost(testCfg, types.Payload{Cost: types.CostInfo{TotalCostUSD: new(math.MaxFloat64)}}),
			testCfg.Segments.Cost.Icon,
		)
	})

	t.Run("vim with unknown mode uses normal icon", func(t *testing.T) {
		assert.Equal(
			t,
			testCfg.Segments.Vim.IconNormal+" ",
			vim(testCfg, types.Payload{Vim: types.VimInfo{Mode: "SOMETHING"}}),
		)
	})
}
