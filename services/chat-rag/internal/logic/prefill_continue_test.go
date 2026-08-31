package logic

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/zgsm-ai/chat-rag/internal/types"
)

func TestHasContinuationMarker(t *testing.T) {
	tests := []struct {
		name string
		req  *types.ChatCompletionRequest
		want bool
	}{
		{
			name: "顶层 continue_final_message=true 触发（vLLM/SGLang 方言）",
			req: &types.ChatCompletionRequest{
				LLMRequestParams: types.LLMRequestParams{
					Extra: map[string]any{"continue_final_message": true},
					Messages: []types.Message{
						{Role: "user", Content: "u"},
						{Role: "assistant", Content: "prefill"},
					},
				},
			},
			want: true,
		},
		{
			name: "顶层 continue_final_message=false 不触发",
			req: &types.ChatCompletionRequest{
				LLMRequestParams: types.LLMRequestParams{
					Extra: map[string]any{"continue_final_message": false},
				},
			},
			want: false,
		},
		{
			name: "顶层 continue_final_message 非布尔（字符串）不触发",
			req: &types.ChatCompletionRequest{
				LLMRequestParams: types.LLMRequestParams{
					Extra: map[string]any{"continue_final_message": "true"},
				},
			},
			want: false,
		},
		{
			name: "末尾 assistant 带 partial=true 触发（Kimi/Qwen 方言）",
			req: &types.ChatCompletionRequest{
				LLMRequestParams: types.LLMRequestParams{
					Messages: []types.Message{
						{Role: "user", Content: "u"},
						{Role: "assistant", Content: "prefill", Extra: map[string]any{"partial": true}},
					},
				},
			},
			want: true,
		},
		{
			name: "末尾 assistant 带 prefix=true 触发（DeepSeek 方言）",
			req: &types.ChatCompletionRequest{
				LLMRequestParams: types.LLMRequestParams{
					Messages: []types.Message{
						{Role: "user", Content: "u"},
						{Role: "assistant", Content: "prefill", Extra: map[string]any{"prefix": true}},
					},
				},
			},
			want: true,
		},
		{
			name: "末尾 assistant partial=false 不触发",
			req: &types.ChatCompletionRequest{
				LLMRequestParams: types.LLMRequestParams{
					Messages: []types.Message{
						{Role: "user", Content: "u"},
						{Role: "assistant", Content: "prefill", Extra: map[string]any{"partial": false}},
					},
				},
			},
			want: false,
		},
		{
			name: "标记在历史中段 assistant（末尾无标记）不触发",
			req: &types.ChatCompletionRequest{
				LLMRequestParams: types.LLMRequestParams{
					Messages: []types.Message{
						{Role: "user", Content: "u"},
						{Role: "assistant", Content: "mid", Extra: map[string]any{"partial": true}},
						{Role: "user", Content: "u2"},
					},
				},
			},
			want: false,
		},
		{
			name: "末尾是 user 且带 partial 不触发",
			req: &types.ChatCompletionRequest{
				LLMRequestParams: types.LLMRequestParams{
					Messages: []types.Message{
						{Role: "user", Content: "u", Extra: map[string]any{"partial": true}},
					},
				},
			},
			want: false,
		},
		{
			name: "无标记普通请求不触发",
			req: &types.ChatCompletionRequest{
				LLMRequestParams: types.LLMRequestParams{
					Extra: map[string]any{"tools": []any{}},
					Messages: []types.Message{
						{Role: "user", Content: "u"},
					},
				},
			},
			want: false,
		},
		{
			name: "空 messages + 顶层标记仍触发",
			req: &types.ChatCompletionRequest{
				LLMRequestParams: types.LLMRequestParams{
					Extra: map[string]any{"continue_final_message": true},
				},
			},
			want: true,
		},
		{
			name: "标记藏在 extra_body.Extra（嵌套）不触发——契约仅认顶层",
			req: &types.ChatCompletionRequest{
				LLMRequestParams: types.LLMRequestParams{
					ExtraBody: types.ExtraBody{
						Extra: map[string]any{"continue_final_message": true},
					},
				},
			},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.want, hasContinuationMarker(tt.req))
		})
	}
}

func TestResolveEffectivePromptMode(t *testing.T) {
	tests := []struct {
		name string
		req  *types.ChatCompletionRequest
		want types.PromptMode
	}{
		{
			name: "未设置模式 + 标记 → 升级 Raw",
			req: &types.ChatCompletionRequest{
				LLMRequestParams: types.LLMRequestParams{
					Extra: map[string]any{"continue_final_message": true},
				},
			},
			want: types.Raw,
		},
		{
			name: "未设置模式 + 无标记 → 保持空（默认链）",
			req: &types.ChatCompletionRequest{
				LLMRequestParams: types.LLMRequestParams{},
			},
			want: "",
		},
		{
			name: "显式 Raw + 标记 → 保持 Raw（不重复升级）",
			req: &types.ChatCompletionRequest{
				LLMRequestParams: types.LLMRequestParams{
					Extra:     map[string]any{"continue_final_message": true},
					ExtraBody: types.ExtraBody{PromptMode: types.Raw},
				},
			},
			want: types.Raw,
		},
		{
			name: "显式 Balanced + 标记 → 尊重显式（不升级）",
			req: &types.ChatCompletionRequest{
				LLMRequestParams: types.LLMRequestParams{
					Extra:     map[string]any{"continue_final_message": true},
					ExtraBody: types.ExtraBody{PromptMode: types.Balanced},
				},
			},
			want: types.Balanced,
		},
		{
			name: "显式 Raw 无标记 → 保持 Raw",
			req: &types.ChatCompletionRequest{
				LLMRequestParams: types.LLMRequestParams{
					ExtraBody: types.ExtraBody{PromptMode: types.Raw},
				},
			},
			want: types.Raw,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.want, resolveEffectivePromptMode(tt.req))
		})
	}
}
