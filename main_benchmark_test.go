package main

import (
	"os"
	"testing"
)

func BenchmarkMain(b *testing.B) {
	input := `{"cwd":"/home/user/my-project","session_id":"abc123","session_name":"my-session","prompt_id":"550e8400-e29b-41d4-a716-446655440000","transcript_path":"/path/to/transcript.jsonl","version":"2.1.90","autorun":true,"fast_mode":true,"exceeds_200k_tokens":false,"model":{"id":"claude-sonnet-4-20250514","display_name":"claude-sonnet-4-20250514","max_mode":false},"workspace":{"current_dir":"/home/user/my-project","project_dir":"/home/user/my-project","added_dirs":["/other/lib"],"git_worktree":"feature-xyz","repo":{"host":"github.com","owner":"anthropics","name":"claude-code"}},"worktree":{"name":"my-feature","path":"/path/to/.claude/worktrees/my-feature","branch":"worktree-my-feature","original_cwd":"/home/user","original_branch":"main"},"context_window":{"total_input_tokens":15500,"total_output_tokens":1200,"context_window_size":200000,"used_percentage":35,"remaining_percentage":65,"current_usage":{"input_tokens":8500,"output_tokens":1200,"cache_creation_input_tokens":5000,"cache_read_input_tokens":2000}},"cost":{"total_cost_usd":0.05,"total_duration_ms":45000,"total_api_duration_ms":2300,"total_lines_added":156,"total_lines_removed":23},"vim":{"mode":"INSERT"},"effort":{"level":"high"},"thinking":{"enabled":true},"output_style":{"name":"default"},"rate_limits":{"five_hour":{"used_percentage":23.5,"resets_at":9999999999},"seven_day":{"used_percentage":41.2,"resets_at":9999999999},"spend_limit":{"used_percentage":62.8,"resets_at":9999999999}},"prompt_cache":{"warm":true,"caching_observed":true,"ttl":"1h","expires_at":9999999999,"requests":14,"misses":2,"expected_rebuilds":1,"hit_ratio":0.91,"cache_write_tokens":352000,"miss_recache_tokens":310200,"last_miss_at":9999999999,"recache_tokens_if_cold":45000},"agent":{"name":"security-reviewer"},"pr":{"number":1234,"url":"https://github.com/anthropics/claude-code/pull/1234","review_state":"pending"}}`

	oldStdout := os.Stdout
	b.Cleanup(func() { os.Stdout = oldStdout })
	os.Stdout = nil

	for b.Loop() {
		r, w, _ := os.Pipe()
		_, _ = w.Write([]byte(input))
		_ = w.Close()
		os.Stdin = r
		main()
	}
}
