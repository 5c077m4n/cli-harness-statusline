package main

import (
	"bytes"
	"math"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"

	"github.com/fatih/color"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var testCfg *Config

func init() {
	color.NoColor = true
	testCfg = loadConfig()
}

func ptr[T any](v T) *T { return &v }

func runCmd(t testing.TB, dir, name string, args ...string) {
	t.Helper()

	cmd := exec.Command(name, args...)
	cmd.Dir = dir
	out, err := cmd.CombinedOutput()

	require.NoErrorf(t, err, "command %s %v failed in %s:\n%s", name, args, dir, string(out))
}

func TestTermSep(t *testing.T) {
	assert.Equal(t, " | ", termSep())
}

func TestModelSegment(t *testing.T) {
	tests := []struct {
		name string
		data Payload
		want string
	}{
		{
			name: "with display name",
			data: Payload{Model: ModelInfo{DisplayName: "gpt-4"}},
			want: testCfg.Icons.Model + " gpt-4",
		},
		{name: "empty display name", data: Payload{}, want: testCfg.Icons.Model + " unknown"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.want, modelSegment(testCfg, tt.data))
		})
	}
}

func TestFolderSegment(t *testing.T) {
	tests := []struct {
		name string
		data Payload
		want string
	}{
		{
			name: "current dir set",
			data: Payload{Workspace: WorkspaceInfo{CurrentDir: "/home/user/project"}},
			want: testCfg.Icons.Folder + " project",
		},
		{
			name: "falls back to cwd",
			data: Payload{Cwd: "/home/user/other"},
			want: testCfg.Icons.Folder + " other",
		},
		{name: "both empty maps to root", data: Payload{}, want: testCfg.Icons.Folder + " /"},
		{
			name: "root path",
			data: Payload{Workspace: WorkspaceInfo{CurrentDir: "/"}},
			want: testCfg.Icons.Folder + " /",
		},
		{
			name: "dot path maps to root",
			data: Payload{Workspace: WorkspaceInfo{CurrentDir: "."}},
			want: testCfg.Icons.Folder + " /",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.want, folderSegment(testCfg, tt.data))
		})
	}
}

func TestWorktreeSegment(t *testing.T) {
	tests := []struct {
		name string
		data Payload
		want string
	}{
		{
			name: "with worktree name",
			data: Payload{Worktree: WorktreeInfo{Name: "feature-branch"}},
			want: testCfg.Icons.Worktree + " feature-branch",
		},
		{name: "empty worktree name", data: Payload{}, want: ""},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.want, worktreeSegment(testCfg, tt.data))
		})
	}
}

func TestContextSegment(t *testing.T) {
	tests := []struct {
		name string
		data Payload
		want string
	}{
		{
			name: "0 percent",
			data: Payload{ContextWindow: ContextWindowInfo{UsedPercentage: ptr(0.0)}},
			want: testCfg.Icons.Context + " [----------] 0%",
		},
		{
			name: "50 percent",
			data: Payload{ContextWindow: ContextWindowInfo{UsedPercentage: ptr(50.0)}},
			want: testCfg.Icons.Context + " [#####-----] 50%",
		},
		{
			name: "100 percent",
			data: Payload{ContextWindow: ContextWindowInfo{UsedPercentage: ptr(100.0)}},
			want: testCfg.Icons.Context + " [##########] 100%",
		},
		{
			name: "negative percent",
			data: Payload{ContextWindow: ContextWindowInfo{UsedPercentage: ptr(-10.0)}},
			want: testCfg.Icons.Context + " [----------] -10%",
		},
		{
			name: "over 100 percent",
			data: Payload{ContextWindow: ContextWindowInfo{UsedPercentage: ptr(150.0)}},
			want: testCfg.Icons.Context + " [##########] 150%",
		},
		{
			name: "fractional percent floors",
			data: Payload{ContextWindow: ContextWindowInfo{UsedPercentage: ptr(59.9)}},
			want: testCfg.Icons.Context + " [#####-----] 59%",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.want, contextSegment(testCfg, tt.data))
		})
	}
}

func TestCostSegment(t *testing.T) {
	tests := []struct {
		name string
		data Payload
		want string
	}{
		{
			name: "with cost",
			data: Payload{Cost: CostInfo{TotalCostUSD: new(0.42)}},
			want: testCfg.Icons.Cost + " $0.42",
		},
		{
			name: "zero cost",
			data: Payload{Cost: CostInfo{TotalCostUSD: new(0.0)}},
			want: testCfg.Icons.Cost + " $0.00",
		},
		{name: "nil cost", data: Payload{}, want: ""},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.want, costSegment(testCfg, tt.data))
		})
	}
}

func TestVimSegment(t *testing.T) {
	tests := []struct {
		name string
		data Payload
		want string
	}{
		{
			name: "insert mode",
			data: Payload{Vim: VimInfo{Mode: "INSERT"}},
			want: testCfg.Icons.VimInsert + " ",
		},
		{
			name: "normal mode",
			data: Payload{Vim: VimInfo{Mode: "NORMAL"}},
			want: testCfg.Icons.VimNormal + " ",
		},
		{name: "empty mode", data: Payload{}, want: ""},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.want, vimSegment(testCfg, tt.data))
		})
	}
}

func TestAutorunSegment(t *testing.T) {
	tests := []struct {
		name string
		data Payload
		want string
	}{
		{name: "autorun on", data: Payload{Autorun: true}, want: testCfg.Icons.Autorun + " auto"},
		{name: "autorun off", data: Payload{}, want: ""},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.want, autorunSegment(testCfg, tt.data))
		})
	}
}

func TestMaxSegment(t *testing.T) {
	tests := []struct {
		name string
		data Payload
		want string
	}{
		{
			name: "max mode on",
			data: Payload{Model: ModelInfo{MaxMode: true}},
			want: testCfg.Icons.Max + " max",
		},
		{name: "max mode off", data: Payload{}, want: ""},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.want, maxSegment(testCfg, tt.data))
		})
	}
}

func TestGitSegment_NotAGitRepo(t *testing.T) {
	data := Payload{Cwd: t.TempDir(), Workspace: WorkspaceInfo{CurrentDir: t.TempDir()}}
	assert.Empty(t, gitSegment(testCfg, data))
}

func TestGitSegment_CleanRepo(t *testing.T) {
	tmpDir := t.TempDir()
	runCmd(t, tmpDir, "git", "init")
	runCmd(t, tmpDir, "git", "config", "user.email", "[EMAIL]")
	runCmd(t, tmpDir, "git", "config", "user.name", "Test")
	runCmd(t, tmpDir, "git", "commit", "--allow-empty", "-m", "initial")
	runCmd(t, tmpDir, "git", "branch", "-M", "main")

	data := Payload{Cwd: tmpDir, Workspace: WorkspaceInfo{CurrentDir: tmpDir}}
	assert.Equal(t, testCfg.Icons.Branch+" main", gitSegment(testCfg, data))
}

func TestGitSegment_DirtyRepo(t *testing.T) {
	tmpDir := t.TempDir()
	runCmd(t, tmpDir, "git", "init")
	runCmd(t, tmpDir, "git", "config", "user.email", "[EMAIL]")
	runCmd(t, tmpDir, "git", "config", "user.name", "Test")
	runCmd(t, tmpDir, "git", "commit", "--allow-empty", "-m", "initial")
	runCmd(t, tmpDir, "git", "branch", "-M", "main")
	require.NoError(t, os.WriteFile(filepath.Join(tmpDir, "dirty.txt"), []byte("dirty"), 0644))

	data := Payload{Cwd: tmpDir, Workspace: WorkspaceInfo{CurrentDir: tmpDir}}
	got := gitSegment(testCfg, data)
	assert.Contains(t, got, testCfg.Icons.Branch+" main")
	assert.True(t, strings.HasSuffix(got, "*"), "expected trailing * for dirty repo")
}

func TestGitInfo_NotAGitRepo(t *testing.T) {
	branch, dirty := gitInfo(t.TempDir())
	assert.Empty(t, branch)
	assert.False(t, dirty)
}

func TestGitInfo_CleanRepo(t *testing.T) {
	tmpDir := t.TempDir()
	runCmd(t, tmpDir, "git", "init")
	runCmd(t, tmpDir, "git", "config", "user.email", "[EMAIL]")
	runCmd(t, tmpDir, "git", "config", "user.name", "Test")
	runCmd(t, tmpDir, "git", "commit", "--allow-empty", "-m", "initial")
	runCmd(t, tmpDir, "git", "branch", "-M", "main")

	branch, dirty := gitInfo(tmpDir)
	assert.Equal(t, "main", branch)
	assert.False(t, dirty)
}

func TestGitInfo_DirtyRepo(t *testing.T) {
	tmpDir := t.TempDir()
	runCmd(t, tmpDir, "git", "init")
	runCmd(t, tmpDir, "git", "config", "user.email", "[EMAIL]")
	runCmd(t, tmpDir, "git", "config", "user.name", "Test")
	runCmd(t, tmpDir, "git", "commit", "--allow-empty", "-m", "initial")
	runCmd(t, tmpDir, "git", "branch", "-M", "main")
	require.NoError(t, os.WriteFile(filepath.Join(tmpDir, "dirty.txt"), []byte("dirty"), 0644))

	branch, dirty := gitInfo(tmpDir)
	assert.Equal(t, "main", branch)
	assert.True(t, dirty)
}

func TestGitInfo_DetachedHead(t *testing.T) {
	tmpDir := t.TempDir()
	runCmd(t, tmpDir, "git", "init")
	runCmd(t, tmpDir, "git", "config", "user.email", "[EMAIL]")
	runCmd(t, tmpDir, "git", "config", "user.name", "Test")
	runCmd(t, tmpDir, "git", "commit", "--allow-empty", "-m", "first")
	runCmd(t, tmpDir, "git", "checkout", "-b", "temp")
	runCmd(t, tmpDir, "git", "checkout", "HEAD~0")

	branch, dirty := gitInfo(tmpDir)
	assert.Empty(t, branch)
	assert.False(t, dirty)
}

func TestGitInfo_NoCommits(t *testing.T) {
	tmpDir := t.TempDir()
	runCmd(t, tmpDir, "git", "init")
	runCmd(t, tmpDir, "git", "config", "user.email", "[EMAIL]")
	runCmd(t, tmpDir, "git", "config", "user.name", "Test")

	branch, dirty := gitInfo(tmpDir)
	assert.NotEmpty(t, branch, "expected a branch name even with no commits")
	assert.False(t, dirty)
}

func TestGitSegment_FallbackToCwd(t *testing.T) {
	tmpDir := t.TempDir()
	runCmd(t, tmpDir, "git", "init")
	runCmd(t, tmpDir, "git", "config", "user.email", "[EMAIL]")
	runCmd(t, tmpDir, "git", "config", "user.name", "Test")
	runCmd(t, tmpDir, "git", "commit", "--allow-empty", "-m", "initial")
	runCmd(t, tmpDir, "git", "branch", "-M", "main")

	data := Payload{Cwd: tmpDir, Workspace: WorkspaceInfo{CurrentDir: ""}}
	assert.Equal(t, testCfg.Icons.Branch+" main", gitSegment(testCfg, data))
}

func TestReadPayload(t *testing.T) {
	oldStdin := os.Stdin
	t.Cleanup(func() { os.Stdin = oldStdin })

	r, w, err := os.Pipe()
	require.NoError(t, err)
	os.Stdin = r

	input := `{"cwd":"/test","autorun":true,"model":{"display_name":"gpt-4o","max_mode":true},"workspace":{"current_dir":"/home/project"},"worktree":{"name":"wt"},"context_window":{"used_percentage":42.5},"cost":{"total_cost_usd":1.23},"vim":{"mode":"INSERT"}}`
	_, err = w.Write([]byte(input))
	require.NoError(t, err)
	require.NoError(t, w.Close())

	data := readPayload()
	assert.Equal(t, "/test", data.Cwd)
	assert.True(t, data.Autorun)
	assert.Equal(t, "gpt-4o", data.Model.DisplayName)
	assert.True(t, data.Model.MaxMode)
	assert.Equal(t, "/home/project", data.Workspace.CurrentDir)
	assert.Equal(t, "wt", data.Worktree.Name)
	require.NotNil(t, data.ContextWindow.UsedPercentage)
	assert.Equal(t, 42.5, *data.ContextWindow.UsedPercentage)
	require.NotNil(t, data.Cost.TotalCostUSD)
	assert.Equal(t, 1.23, *data.Cost.TotalCostUSD)
	assert.Equal(t, "INSERT", data.Vim.Mode)
}

func TestReadPayload_EmptyInput(t *testing.T) {
	if os.Getenv("TEST_READ_PAYLOAD_EXIT") == "1" {
		r, w, _ := os.Pipe()
		os.Stdin = r
		_ = w.Close()
		readPayload()
		return
	}
	cmd := exec.Command(os.Args[0], "-test.run=TestReadPayload_EmptyInput")
	cmd.Env = append(os.Environ(), "TEST_READ_PAYLOAD_EXIT=1")
	err := cmd.Run()
	if e, ok := err.(*exec.ExitError); ok && !e.Success() {
		return
	}
	t.Fatal("expected readPayload() to os.Exit(1) on empty input")
}

func TestMainFullPayload(t *testing.T) {
	oldStdin := os.Stdin
	oldStdout := os.Stdout
	t.Cleanup(func() {
		os.Stdin = oldStdin
		os.Stdout = oldStdout
	})

	r, w, err := os.Pipe()
	require.NoError(t, err)
	os.Stdin = r

	outR, outW, err := os.Pipe()
	require.NoError(t, err)
	os.Stdout = outW

	input := `{"cwd":"/home/user/my-project","session_id":"abc123","session_name":"my-session","prompt_id":"550e8400-e29b-41d4-a716-446655440000","transcript_path":"/path/to/transcript.jsonl","version":"2.1.90","autorun":true,"fast_mode":true,"exceeds_200k_tokens":false,"model":{"id":"claude-sonnet-4-20250514","display_name":"claude-sonnet-4-20250514","max_mode":false},"workspace":{"current_dir":"/home/user/my-project","project_dir":"/home/user/my-project","added_dirs":["/other/lib"],"git_worktree":"feature-xyz","repo":{"host":"github.com","owner":"anthropics","name":"claude-code"}},"worktree":{"name":"my-feature","path":"/path/to/.claude/worktrees/my-feature","branch":"worktree-my-feature","original_cwd":"/home/user","original_branch":"main"},"context_window":{"total_input_tokens":15500,"total_output_tokens":1200,"context_window_size":200000,"used_percentage":35,"remaining_percentage":65,"current_usage":{"input_tokens":8500,"output_tokens":1200,"cache_creation_input_tokens":5000,"cache_read_input_tokens":2000}},"cost":{"total_cost_usd":0.05,"total_duration_ms":45000,"total_api_duration_ms":2300,"total_lines_added":156,"total_lines_removed":23},"vim":{"mode":"INSERT"},"effort":{"level":"high"},"thinking":{"enabled":true},"output_style":{"name":"default"},"rate_limits":{"five_hour":{"used_percentage":23.5,"resets_at":"2025-01-01T00:00:00Z"},"seven_day":{"used_percentage":41.2,"resets_at":"2025-01-01T00:00:00Z"},"spend_limit":{"used_percentage":62.8,"resets_at":"2025-01-01T00:00:00Z"}},"prompt_cache":{"warm":true,"caching_observed":true,"ttl":"1h","expires_at":"2025-01-01T00:00:00Z","requests":14,"misses":2,"expected_rebuilds":1,"hit_ratio":0.91,"cache_write_tokens":352000,"miss_recache_tokens":310200,"last_miss_at":"2025-01-01T00:00:00Z","recache_tokens_if_cold":45000},"agent":{"name":"security-reviewer"},"pr":{"number":1234,"url":"https://github.com/anthropics/claude-code/pull/1234","review_state":"pending"}}`
	_, err = w.Write([]byte(input))
	require.NoError(t, err)
	require.NoError(t, w.Close())

	color.NoColor = true
	main()
	require.NoError(t, outW.Close())

	var buf bytes.Buffer
	_, err = buf.ReadFrom(outR)
	require.NoError(t, err)
	output := buf.String()

	assert.Contains(t, output, testCfg.Icons.Model+" claude-sonnet-4-20250514")
	assert.Contains(t, output, testCfg.Icons.Folder+" my-project")
	assert.Contains(t, output, testCfg.Icons.Agent+" security-reviewer")
	assert.Contains(t, output, testCfg.Icons.Worktree+" my-feature")
	assert.Contains(t, output, testCfg.Icons.Session+" my-session")
	assert.Contains(t, output, testCfg.Icons.PR+" #1234")
	assert.Contains(t, output, testCfg.Icons.Context+" [###-------] 35%")
	assert.Contains(t, output, testCfg.Icons.Cost+" $0.05")
	assert.Contains(t, output, testCfg.Icons.Tokens+" 15k/1k")
	assert.Contains(t, output, testCfg.Icons.Cache+" ")
	assert.Contains(t, output, "warm")
	assert.Contains(t, output, "91%")
	assert.Contains(t, output, testCfg.Icons.VimInsert)
	assert.Contains(t, output, testCfg.Icons.Autorun+" auto")
	assert.Contains(t, output, testCfg.Icons.FastMode+" fast")
	assert.Contains(t, output, testCfg.Icons.Effort+" high")
	assert.Contains(t, output, testCfg.Icons.Thinking+" think")
	assert.Contains(t, output, testCfg.Icons.RateLimit+" 5h:24% 7d:41% $:63%")
}

func TestMainIntegration(t *testing.T) {
	oldStdin := os.Stdin
	oldStdout := os.Stdout
	t.Cleanup(func() {
		os.Stdin = oldStdin
		os.Stdout = oldStdout
	})

	r, w, err := os.Pipe()
	require.NoError(t, err)
	os.Stdin = r

	outR, outW, err := os.Pipe()
	require.NoError(t, err)
	os.Stdout = outW

	input := `{"cwd":"/home/user/my-project","autorun":true,"model":{"display_name":"claude-sonnet-4-20250514"},"workspace":{"current_dir":"/home/user/my-project"},"context_window":{"used_percentage":35},"cost":{"total_cost_usd":0.05},"vim":{"mode":"INSERT"}}`
	_, err = w.Write([]byte(input))
	require.NoError(t, err)
	require.NoError(t, w.Close())

	color.NoColor = false
	main()
	color.NoColor = true
	require.NoError(t, outW.Close())

	var buf bytes.Buffer
	_, err = buf.ReadFrom(outR)
	require.NoError(t, err)
	output := buf.String()

	assert.Contains(t, output, testCfg.Icons.Model+" claude-sonnet-4-20250514")
	assert.Contains(t, output, testCfg.Icons.Folder+" my-project")
	assert.Contains(t, output, testCfg.Icons.Context+" [###-------] 35%")
	assert.Contains(t, output, testCfg.Icons.Cost+" $0.05")
	assert.Contains(t, output, testCfg.Icons.VimInsert)
	assert.Contains(t, output, testCfg.Icons.Autorun+" auto")
	assert.NotContains(t, output, testCfg.Icons.Max+" max")
}

func TestContextSegmentColorThresholds(t *testing.T) {
	tests := []struct {
		name    string
		percent float64
		want    string
	}{
		{name: "under 60 green", percent: 0, want: testCfg.Icons.Context + " [----------] 0%"},
		{name: "under 60 green", percent: 59, want: testCfg.Icons.Context + " [#####-----] 59%"},
		{name: "at 60 yellow", percent: 60, want: testCfg.Icons.Context + " [######----] 60%"},
		{name: "at 84 yellow", percent: 84, want: testCfg.Icons.Context + " [########--] 84%"},
		{name: "at 85 red", percent: 85, want: testCfg.Icons.Context + " [########--] 85%"},
		{name: "at 100 red", percent: 100, want: testCfg.Icons.Context + " [##########] 100%"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			data := Payload{ContextWindow: ContextWindowInfo{UsedPercentage: ptr(tt.percent)}}
			assert.Contains(t, contextSegment(testCfg, data), tt.want)
		})
	}
}

func TestContextSegmentBarFilling(t *testing.T) {
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
			data := Payload{ContextWindow: ContextWindowInfo{UsedPercentage: ptr(tt.percent)}}
			got := contextSegment(testCfg, data)
			expectedBar := strings.Repeat("#", tt.filled) + strings.Repeat("-", 10-tt.filled)
			assert.Contains(t, got, "["+expectedBar+"]")
		})
	}
}

func TestWorktreeSegment_Empty(t *testing.T) {
	assert.Empty(t, worktreeSegment(testCfg, Payload{Worktree: WorktreeInfo{Name: ""}}))
}

func TestFolderSegment_CurrentDirPriority(t *testing.T) {
	data := Payload{Cwd: "/fallback", Workspace: WorkspaceInfo{CurrentDir: "/preferred"}}
	assert.Equal(t, testCfg.Icons.Folder+" preferred", folderSegment(testCfg, data))
}

func TestFastModeSegment(t *testing.T) {
	tests := []struct {
		name string
		data Payload
		want string
	}{
		{name: "fast mode on", data: Payload{FastMode: true}, want: testCfg.Icons.FastMode + " fast"},
		{name: "fast mode off", data: Payload{}, want: ""},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.want, fastModeSegment(testCfg, tt.data))
		})
	}
}

func TestEffortSegment(t *testing.T) {
	tests := []struct {
		name string
		data Payload
		want string
	}{
		{
			name: "with effort",
			data: Payload{Effort: &EffortInfo{Level: "high"}},
			want: testCfg.Icons.Effort + " high",
		},
		{name: "nil effort", data: Payload{}, want: ""},
		{name: "empty level", data: Payload{Effort: &EffortInfo{Level: ""}}, want: ""},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.want, effortSegment(testCfg, tt.data))
		})
	}
}

func TestThinkingSegment(t *testing.T) {
	tests := []struct {
		name string
		data Payload
		want string
	}{
		{
			name: "thinking enabled",
			data: Payload{Thinking: &ThinkingInfo{Enabled: true}},
			want: testCfg.Icons.Thinking + " think",
		},
		{
			name: "thinking disabled",
			data: Payload{Thinking: &ThinkingInfo{Enabled: false}},
			want: "",
		},
		{name: "nil thinking", data: Payload{}, want: ""},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.want, thinkingSegment(testCfg, tt.data))
		})
	}
}

func TestSessionSegment(t *testing.T) {
	tests := []struct {
		name string
		data Payload
		want string
	}{
		{
			name: "with session name",
			data: Payload{SessionName: "my-session"},
			want: testCfg.Icons.Session + " my-session",
		},
		{name: "empty session name", data: Payload{}, want: ""},
		{
			name: "truncates long name",
			data: Payload{SessionName: "abcdefghijklmnopqrstuvwxyz123456"},
			want: testCfg.Icons.Session + " abcdefghijklmnopqrstuvwx…",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.want, sessionSegment(testCfg, tt.data))
		})
	}
}

func TestAgentSegment(t *testing.T) {
	tests := []struct {
		name string
		data Payload
		want string
	}{
		{
			name: "with agent",
			data: Payload{Agent: &AgentInfo{Name: "reviewer"}},
			want: testCfg.Icons.Agent + " reviewer",
		},
		{name: "nil agent", data: Payload{}, want: ""},
		{name: "empty name", data: Payload{Agent: &AgentInfo{Name: ""}}, want: ""},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.want, agentSegment(testCfg, tt.data))
		})
	}
}

func TestPRSegment(t *testing.T) {
	tests := []struct {
		name string
		data Payload
		want string
	}{
		{
			name: "PR without review state",
			data: Payload{PR: &PRInfo{Number: 42, URL: "https://example.com/pr/42"}},
			want: testCfg.Icons.PR + " #42",
		},
		{name: "nil PR", data: Payload{}, want: ""},
		{
			name: "PR with approved review",
			data: Payload{PR: &PRInfo{Number: 1, ReviewState: ptr("approved")}},
			want: testCfg.Icons.PR + " #1 ✓",
		},
		{
			name: "PR with changes requested",
			data: Payload{PR: &PRInfo{Number: 2, ReviewState: ptr("changes_requested")}},
			want: testCfg.Icons.PR + " #2 ✗",
		},
		{
			name: "PR draft",
			data: Payload{PR: &PRInfo{Number: 3, ReviewState: ptr("draft")}},
			want: testCfg.Icons.PR + " #3 ○",
		},
		{
			name: "PR pending",
			data: Payload{PR: &PRInfo{Number: 4, ReviewState: ptr("pending")}},
			want: testCfg.Icons.PR + " #4 ●",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.want, prSegment(testCfg, tt.data))
		})
	}
}

func TestRateLimitSegment(t *testing.T) {
	tests := []struct {
		name string
		data Payload
		want string
	}{
		{name: "nil rate limits", data: Payload{}, want: ""},
		{
			name: "with five hour",
			data: Payload{
				RateLimits: &RateLimits{FiveHour: &RateLimitWindow{UsedPercentage: 23.5}},
			},
			want: testCfg.Icons.RateLimit + " 5h:24%",
		},
		{
			name: "with all windows",
			data: Payload{RateLimits: &RateLimits{
				FiveHour:   &RateLimitWindow{UsedPercentage: 10},
				SevenDay:   &RateLimitWindow{UsedPercentage: 50},
				SpendLimit: &RateLimitWindow{UsedPercentage: 75},
			}},
			want: testCfg.Icons.RateLimit + " 5h:10% 7d:50% $:75%",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.want, rateLimitSegment(testCfg, tt.data))
		})
	}
}

func TestTokenSegment(t *testing.T) {
	tests := []struct {
		name string
		data Payload
		want string
	}{
		{name: "zero tokens", data: Payload{}, want: ""},
		{
			name: "with tokens",
			data: Payload{
				ContextWindow: ContextWindowInfo{TotalInputTokens: 15500, TotalOutputTokens: 1200},
			},
			want: testCfg.Icons.Tokens + " 15k/1k",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.want, tokenSegment(testCfg, tt.data))
		})
	}
}

func TestCacheSegment(t *testing.T) {
	tests := []struct {
		name string
		data Payload
		want string
	}{
		{name: "nil prompt cache", data: Payload{}, want: ""},
		{
			name: "warm cache with hit ratio",
			data: Payload{PromptCache: &PromptCacheInfo{Warm: true, HitRatio: ptr(0.91)}},
			want: testCfg.Icons.Cache + " warm 91%",
		},
		{
			name: "cold cache",
			data: Payload{PromptCache: &PromptCacheInfo{Warm: false}},
			want: testCfg.Icons.Cache + " cold",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.want, cacheSegment(testCfg, tt.data))
		})
	}
}

func TestExceedsSegment(t *testing.T) {
	tests := []struct {
		name string
		data Payload
		want string
	}{
		{
			name: "exceeds 200k",
			data: Payload{Exceeds200k: true},
			want: testCfg.Icons.Exceeds + " >200k",
		},
		{name: "under 200k", data: Payload{}, want: ""},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.want, exceedsSegment(testCfg, tt.data))
		})
	}
}

func TestSegmentEdgeCases(t *testing.T) {
	t.Run("model with empty display name and max mode", func(t *testing.T) {
		assert.Equal(
			t,
			testCfg.Icons.Model+" unknown",
			modelSegment(testCfg, Payload{Model: ModelInfo{MaxMode: true}}),
		)
	})

	t.Run("cost with very large value", func(t *testing.T) {
		assert.Contains(
			t,
			costSegment(testCfg, Payload{Cost: CostInfo{TotalCostUSD: new(math.MaxFloat64)}}),
			testCfg.Icons.Cost,
		)
	})

	t.Run("vim with unknown mode uses normal icon", func(t *testing.T) {
		assert.Equal(
			t,
			testCfg.Icons.VimNormal+" ",
			vimSegment(testCfg, Payload{Vim: VimInfo{Mode: "SOMETHING"}}),
		)
	})
}
