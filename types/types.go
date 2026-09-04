// Package types used in this porject
package types

import (
	"encoding/json/v2"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"os"
	"reflect"
	"strings"
)

var ErrSTDINRead = errors.New("STDIN read failed")

type (
	ModelInfo struct {
		ID          string `json:"id"`
		DisplayName string `json:"display_name"`
		MaxMode     bool   `json:"max_mode"`
	}
	RepoInfo struct {
		Host  string `json:"host"`
		Owner string `json:"owner"`
		Name  string `json:"name"`
	}
	WorkspaceInfo struct {
		CurrentDir  string    `json:"current_dir"`
		ProjectDir  string    `json:"project_dir"`
		AddedDirs   []string  `json:"added_dirs"`
		GitWorktree string    `json:"git_worktree"`
		Repo        *RepoInfo `json:"repo,omitzero"`
	}
	WorktreeInfo struct {
		Name           string `json:"name"`
		Path           string `json:"path"`
		Branch         string `json:"branch"`
		OriginalCwd    string `json:"original_cwd"`
		OriginalBranch string `json:"original_branch"`
	}
	ContextUsageInfo struct {
		InputTokens              int `json:"input_tokens"`
		OutputTokens             int `json:"output_tokens"`
		CacheCreationInputTokens int `json:"cache_creation_input_tokens"`
		CacheReadInputTokens     int `json:"cache_read_input_tokens"`
	}
	ContextWindowInfo struct {
		TotalInputTokens    int               `json:"total_input_tokens"`
		TotalOutputTokens   int               `json:"total_output_tokens"`
		ContextWindowSize   int               `json:"context_window_size"`
		UsedPercentage      *float64          `json:"used_percentage"`
		RemainingPercentage *float64          `json:"remaining_percentage"`
		CurrentUsage        *ContextUsageInfo `json:"current_usage"`
	}
	CostInfo struct {
		TotalCostUSD       *float64 `json:"total_cost_usd"`
		TotalDurationMs    float64  `json:"total_duration_ms"`
		TotalAPIDurationMs float64  `json:"total_api_duration_ms"`
		TotalLinesAdded    int      `json:"total_lines_added"`
		TotalLinesRemoved  int      `json:"total_lines_removed"`
	}
	VimInfo struct {
		Mode string `json:"mode"`
	}
	EffortInfo struct {
		Level string `json:"level"`
	}
	ThinkingInfo struct {
		Enabled bool `json:"enabled"`
	}
	RateLimitWindow struct {
		UsedPercentage float64 `json:"used_percentage"`
		ResetsAt       float64 `json:"resets_at"`
	}
	RateLimits struct {
		FiveHour   *RateLimitWindow `json:"five_hour,omitzero"`
		SevenDay   *RateLimitWindow `json:"seven_day,omitzero"`
		SpendLimit *RateLimitWindow `json:"spend_limit,omitzero"`
	}
	PromptCacheInfo struct {
		Warm                bool     `json:"warm"`
		CachingObserved     bool     `json:"caching_observed"`
		TTL                 string   `json:"ttl"`
		ExpiresAt           string   `json:"expires_at"`
		Requests            int      `json:"requests"`
		Misses              int      `json:"misses"`
		ExpectedRebuilds    int      `json:"expected_rebuilds"`
		HitRatio            *float64 `json:"hit_ratio"`
		CacheWriteTokens    int64    `json:"cache_write_tokens"`
		MissRecacheTokens   int64    `json:"miss_recache_tokens"`
		LastMissAt          string   `json:"last_miss_at"`
		RecacheTokensIfCold *float64 `json:"recache_tokens_if_cold"`
	}
	OutputStyleInfo struct {
		Name string `json:"name"`
	}
	AgentInfo struct {
		Name string `json:"name"`
	}
	PRInfo struct {
		Number      int     `json:"number"`
		URL         string  `json:"url"`
		ReviewState *string `json:"review_state,omitzero"`
		Kind        *string `json:"kind,omitzero"`
	}
	Payload struct {
		Cwd            string            `json:"cwd"`
		SessionID      string            `json:"session_id"`
		SessionName    string            `json:"session_name"`
		PromptID       string            `json:"prompt_id"`
		TranscriptPath string            `json:"transcript_path"`
		Version        string            `json:"version"`
		Autorun        bool              `json:"autorun"`
		FastMode       bool              `json:"fast_mode"`
		Exceeds200k    bool              `json:"exceeds_200k_tokens"`
		Model          ModelInfo         `json:"model"`
		Workspace      WorkspaceInfo     `json:"workspace"`
		Worktree       WorktreeInfo      `json:"worktree"`
		ContextWindow  ContextWindowInfo `json:"context_window"`
		Cost           CostInfo          `json:"cost"`
		Vim            VimInfo           `json:"vim"`
		Effort         *EffortInfo       `json:"effort,omitzero"`
		Thinking       *ThinkingInfo     `json:"thinking,omitzero"`
		OutputStyle    OutputStyleInfo   `json:"output_style"`
		RateLimits     *RateLimits       `json:"rate_limits,omitzero"`
		PromptCache    *PromptCacheInfo  `json:"prompt_cache,omitzero"`
		Agent          *AgentInfo        `json:"agent,omitzero"`
		PR             *PRInfo           `json:"pr,omitzero"`
	}
)

func (p *Payload) UnmarshalJSON(data []byte) error {
	var raw map[string]any
	if err := json.Unmarshal(data, &raw); err != nil {
		return err
	}

	pValue := reflect.ValueOf(p).Elem()
	pType := pValue.Type()

	for i := range pType.NumField() {
		field := pType.Field(i)
		tag := field.Tag.Get("json")
		if tag == "" {
			continue
		}
		key, _, _ := strings.Cut(tag, ",")

		rawValue, ok := raw[key]
		if !ok {
			continue
		}

		b, err := json.Marshal(rawValue)
		if err != nil {
			slog.Warn(
				"unmarshal payload: failed to marshal",
				slog.String("key", key),
				slog.Any("error", err),
			)
			continue
		}

		fieldValue := pValue.Field(i)
		if fieldValue.Kind() == reflect.Pointer {
			fieldValue.Set(reflect.New(fieldValue.Type().Elem()))
			if err := json.Unmarshal(b, fieldValue.Interface()); err != nil {
				slog.Warn(
					"unmarshal payload: failed to unmarshal",
					slog.String("key", key),
					slog.Any("error", err),
				)
			}
		} else {
			if err := json.Unmarshal(b, fieldValue.Addr().Interface()); err != nil {
				slog.Warn(
					"unmarshal payload: failed to unmarshal",
					slog.String("key", key),
					slog.Any("error", err),
				)
			}
		}
	}

	return nil
}

func NewPayLoad() (*Payload, error) {
	input, err := io.ReadAll(os.Stdin)
	if err != nil {
		return nil, errors.Join(ErrSTDINRead, err)
	}

	var data Payload
	if err := json.Unmarshal(input, &data); err != nil {
		return nil, errors.Join(fmt.Errorf("JSON unmarshal failed for '%s'", input), err)
	}

	return &data, nil
}
