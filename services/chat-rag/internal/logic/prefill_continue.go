package logic

import (
	"github.com/zgsm-ai/chat-rag/internal/types"
)

// hasContinuationMarker 判断请求是否携带 assistant 前缀续写的显式标记：
//   - 顶层 "continue_final_message": true（vLLM / SGLang 方言）
//   - 末尾 assistant 消息带 "partial" 或 "prefix": true
//     （Kimi / Qwen / DeepSeek 官方 API 方言）
//
// 此类请求必须走 raw 直通：默认 prompt 处理链会重建消息（只保留
// system/user），静默丢弃末尾 assistant prefill 并注入网关 system，
// 破坏 token 级续写。仅识别顶层标记（extra_body.Extra 嵌套形态不认，
// 该形态本就无法端到端生效）。
func hasContinuationMarker(req *types.ChatCompletionRequest) bool {
	if v, ok := req.Extra["continue_final_message"].(bool); ok && v {
		return true
	}
	n := len(req.Messages)
	if n == 0 || req.Messages[n-1].Role != types.RoleAssistant {
		return false
	}
	for _, key := range []string{"partial", "prefix"} {
		if v, ok := req.Messages[n-1].Extra[key].(bool); ok && v {
			return true
		}
	}
	return false
}

// resolveEffectivePromptMode 返回用于路由的生效 prompt 模式：
// 未显式设置模式且携带续写标记的请求升级为 Raw，使 prefill 原样
// 穿过 prompt 管线。显式设置的模式（含 Raw/Balanced/Cost/Performance/Auto）
// 一律尊重。请求结构本身不被改写——不向后端转发任何 extra_body 污染。
func resolveEffectivePromptMode(req *types.ChatCompletionRequest) types.PromptMode {
	mode := req.ExtraBody.PromptMode
	if mode == "" && hasContinuationMarker(req) {
		return types.Raw
	}
	return mode
}
