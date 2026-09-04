package main

import (
	"bytes"
	"os"
	"testing"

	"github.com/5c077m4n/cli-harness-statusline/segments"
	"github.com/5c077m4n/cli-harness-statusline/types"
	"github.com/fatih/color"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

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

	data, err := types.NewPayLoad()
	require.NoError(t, err)

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

func TestReadPayload_FaultTolerant(t *testing.T) {
	oldStdin := os.Stdin
	t.Cleanup(func() { os.Stdin = oldStdin })

	r, w, err := os.Pipe()
	require.NoError(t, err)
	os.Stdin = r

	input := `{"cwd":"/safe","autorun":true,"model":null,"workspace":"oops","cost":{"total_cost_usd":0.99},"vim":123}`
	_, err = w.Write([]byte(input))
	require.NoError(t, err)
	require.NoError(t, w.Close())

	data, err := types.NewPayLoad()
	require.NoError(t, err)

	assert.Equal(t, "/safe", data.Cwd)
	assert.True(t, data.Autorun)
	assert.Equal(t, "", data.Model.DisplayName)
	assert.Equal(t, "", data.Workspace.CurrentDir)
	require.NotNil(t, data.Cost.TotalCostUSD)
	assert.Equal(t, 0.99, *data.Cost.TotalCostUSD)
	assert.Equal(t, "", data.Vim.Mode)
}

func TestReadPayload_EmptyInput(t *testing.T) {
	oldStdin := os.Stdin
	t.Cleanup(func() { os.Stdin = oldStdin })

	r, w, err := os.Pipe()
	require.NoError(t, err)
	os.Stdin = r
	require.NoError(t, w.Close())

	payload, err := types.NewPayLoad()
	assert.Nil(t, payload)
	assert.ErrorContains(t, err, "JSON unmarshal failed for")
	assert.ErrorContains(t, err, "unexpected EOF")
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

	input := `{"cwd":"/home/user/my-project","session_id":"abc123","session_name":"my-session","prompt_id":"550e8400-e29b-41d4-a716-446655440000","transcript_path":"/path/to/transcript.jsonl","version":"2.1.90","autorun":true,"fast_mode":true,"exceeds_200k_tokens":false,"model":{"id":"claude-sonnet-4-20250514","display_name":"claude-sonnet-4-20250514","max_mode":false},"workspace":{"current_dir":"/home/user/my-project","project_dir":"/home/user/my-project","added_dirs":["/other/lib"],"git_worktree":"feature-xyz","repo":{"host":"github.com","owner":"anthropics","name":"claude-code"}},"worktree":{"name":"my-feature","path":"/path/to/.claude/worktrees/my-feature","branch":"worktree-my-feature","original_cwd":"/home/user","original_branch":"main"},"context_window":{"total_input_tokens":15500,"total_output_tokens":1200,"context_window_size":200000,"used_percentage":35,"remaining_percentage":65,"current_usage":{"input_tokens":8500,"output_tokens":1200,"cache_creation_input_tokens":5000,"cache_read_input_tokens":2000}},"cost":{"total_cost_usd":0.05,"total_duration_ms":45000,"total_api_duration_ms":2300,"total_lines_added":156,"total_lines_removed":23},"vim":{"mode":"INSERT"},"effort":{"level":"high"},"thinking":{"enabled":true},"output_style":{"name":"default"},"rate_limits":{"five_hour":{"used_percentage":23.5,"resets_at":0},"seven_day":{"used_percentage":41.2,"resets_at":0},"spend_limit":{"used_percentage":62.8,"resets_at":0}},"prompt_cache":{"warm":true,"caching_observed":true,"ttl":"1h","expires_at":"2025-01-01T00:00:00Z","requests":14,"misses":2,"expected_rebuilds":1,"hit_ratio":0.91,"cache_write_tokens":352000,"miss_recache_tokens":310200,"last_miss_at":"2025-01-01T00:00:00Z","recache_tokens_if_cold":45000},"agent":{"name":"security-reviewer"},"pr":{"number":1234,"url":"https://github.com/anthropics/claude-code/pull/1234","review_state":"pending"}}`
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

	assert.Contains(t, output, segments.IconModel+" claude-sonnet-4-20250514")
	assert.Contains(t, output, segments.IconFolder+" my-project")
	assert.Contains(t, output, segments.IconAgent+" security-reviewer")
	assert.Contains(t, output, segments.IconWorktree+" my-feature")
	assert.Contains(t, output, segments.IconSession+" my-session")
	assert.Contains(
		t,
		output,
		segments.IconPR+" "+segments.OSC8Start+"https://github.com/anthropics/claude-code/pull/1234"+segments.OSC8Sep+"#1234",
	)
	assert.Contains(t, output, segments.IconContext+" [###-------] 35%")
	assert.Contains(t, output, segments.IconCost+" $0.05")
	assert.Contains(
		t,
		output,
		segments.IconTokenIn+"15k "+segments.IconTokenOut+"1k",
	)
	assert.Contains(t, output, segments.IconCache+" ")
	assert.Contains(t, output, "warm")
	assert.Contains(t, output, "91%")
	assert.Contains(t, output, segments.IconVimInsert)
	assert.Contains(t, output, segments.IconAutorun+" auto")
	assert.Contains(t, output, segments.IconFastMode+" fast")
	assert.Contains(t, output, segments.IconEffort+" high")
	assert.Contains(t, output, segments.IconThinking+" think")
	assert.Contains(t, output, segments.IconRateLimit+" 5h:24% 7d:41% $:63%")
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

	assert.Contains(t, output, segments.IconModel+" claude-sonnet-4-20250514")
	assert.Contains(t, output, segments.IconFolder+" my-project")
	assert.Contains(t, output, segments.IconContext+" [###-------] 35%")
	assert.Contains(t, output, segments.IconCost+" $0.05")
	assert.Contains(t, output, segments.IconVimInsert)
	assert.Contains(t, output, segments.IconAutorun+" auto")
	assert.NotContains(t, output, segments.IconMax+" max")
}
