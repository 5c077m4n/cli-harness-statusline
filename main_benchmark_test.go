package main

import (
	"os"
	"strings"
	"testing"
)

func BenchmarkFullAssembly(b *testing.B) {
	data := Payload{
		Cwd:           "/home/user/my-project",
		Autorun:       true,
		Model:         ModelInfo{DisplayName: "claude-sonnet-4-20250514", MaxMode: false},
		Workspace:     WorkspaceInfo{CurrentDir: "/home/user/my-project"},
		Worktree:      WorktreeInfo{Name: "feature-x"},
		ContextWindow: ContextWindowInfo{UsedPercentage: 42.5},
		Cost:          CostInfo{TotalCostUSD: ptr(0.05)},
		Vim:           VimInfo{Mode: "INSERT"},
	}

	for b.Loop() {
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
		_ = strings.Join(segments, "")
	}
}

func BenchmarkMain(b *testing.B) {
	input := `{"cwd":"/home/user/my-project","autorun":true,"model":{"display_name":"claude-sonnet-4-20250514"},"workspace":{"current_dir":"/home/user/my-project"},"context_window":{"used_percentage":35},"cost":{"total_cost_usd":0.05},"vim":{"mode":"INSERT"}}`

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
