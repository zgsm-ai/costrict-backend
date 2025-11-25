package types

const (
	// RoleSystem System role message
	RoleSystem = "system"

	// RoleUser User role message
	RoleUser = "user"

	// RoleAssistant AI assistant role message
	RoleAssistant = "assistant"
)

// PromptMode defines different types of chat
type PromptMode string

const (
	// Raw mode: No deep processing of user prompt, only necessary operations like compression
	Raw PromptMode = "raw"

	// Balanced mode: Considering both cost and performance, choosing a compromise approach
	// including rag and prompt compression
	Balanced PromptMode = "balanced"

	// Cost Cost-first mode: Minimizing LLM calls and context size to save cost
	Cost PromptMode = "cost"

	// Performance Performance-first mode: Maximizing output quality without considering cost
	Performance PromptMode = "performance"

	// Auto select mode: Default is balanced mode
	Auto PromptMode = "auto"

	// Strict mode: Strictly follow the workflow agent
	Strict PromptMode = "strict"
)

const (
	// Request Headers
	HeaderQuotaIdentity = "x-quota-identity"
	HeaderRequestId     = "x-request-id"
	HeaderCaller        = "x-caller"
	HeaderTaskId        = "zgsm-task-id"
	HeaderClientId      = "zgsm-client-id"
	HeaderClientIde     = "zgsm-client-ide"
	HeaderClientOS      = "X-Stainless-OS"
	HeaderLanguage      = "Accept-Language"
	HeaderAuthorization = "authorization"
	HeaderProjectPath   = "zgsm-project-path"
	HeaderClientVersion = "X-Costrict-Version"
	HeaderOriginalModel = "x-original-model"

	// Response Headers
	HeaderUserInput   = "x-user-input"
	HeaderSelectLLm   = "x-select-llm"
	HeaderOneAPIReqId = "x-oneapi-request-id"
)

// ResponseHeadersToForward defines the list of response headers that should be forwarded
var ResponseHeadersToForward = []string{
	HeaderUserInput,
	HeaderSelectLLm,
	HeaderOneAPIReqId,
}

// ToolStatus defines the status of the tool
type ToolStatus string

const (
	ToolStatusRunning ToolStatus = "running"
	ToolStatusSuccess ToolStatus = "success"
	ToolStatusFailed  ToolStatus = "failed"
)

// Redis key prefix for tool status
const ToolStatusRedisKeyPrefix = "tool_status:"

// Tool string filter
const StrFilterToolAnalyzing = "\n#### 💡 检索已完成，分析中"
const StrFilterToolSearchStart = "\n#### 🔍 "
const StrFilterToolSearchEnd = "工具检索中"

type ChatCompletionRequest struct {
	ChatLLMRequest               // Embedded ChatLLMRequest
	Stream         bool          `json:"stream,omitempty"`
	StreamOptions  StreamOptions `json:"stream_options,omitempty"`
}

type ExtraBody struct {
	PromptMode PromptMode `json:"prompt_mode,omitempty"`
	Mode       string     `json:"mode,omitempty"`
}

type ChatCompletionResponse struct {
	Id      string   `json:"id"`
	Object  string   `json:"object"`
	Created int64    `json:"created"`
	Model   string   `json:"model"`
	Choices []Choice `json:"choices"`
	Usage   Usage    `json:"usage"`
}

// LLMRequestParams contains parameters for LLM requests
type LLMRequestParams struct {
	Messages            []Message `json:"messages"`
	MaxTokens           *int      `json:"max_tokens,omitempty"`
	MaxCompletionTokens *int      `json:"max_completion_tokens,omitempty"`
	Temperature         *float64  `json:"temperature,omitempty"`
	ExtraBody           ExtraBody `json:"extra_body,omitempty"`
}

type ChatLLMRequest struct {
	Model            string `json:"model"`
	LLMRequestParams        // Embedded params
}

type ChatLLMRequestStream struct {
	ChatLLMRequest               // Embedded ChatLLMRequest
	Tools          []Function    `json:"tools,omitempty"`
	ToolChoice     string        `json:"tool_choice,omitempty"`
	Stream         bool          `json:"stream,omitempty"`
	StreamOptions  StreamOptions `json:"stream_options,omitempty"`
}

type Choice struct {
	Index        int     `json:"index"`
	Message      Message `json:"message,omitempty"`
	Delta        Delta   `json:"delta,omitempty"`
	FinishReason string  `json:"finish_reason,omitempty"`
}

type Message struct {
	Role    string `json:"role"`
	Content any    `json:"content"`
}

type Delta struct {
	Role             string `json:"role,omitempty"`
	Content          string `json:"content"`
	ReasoningContent string `json:"reasoning_content,omitempty"`
	ToolCalls        []any  `json:"tool_calls,omitempty"`
}

type StreamOptions struct {
	IncludeUsage bool `json:"include_usage"`
}

type Usage struct {
	PromptTokens     int `json:"prompt_tokens"`
	CompletionTokens int `json:"completion_tokens"`
	TotalTokens      int `json:"total_tokens"`
}

// FunctionCall is the structure of the function called by the LLM.
type Function struct {
	Type     string             `json:"type"`
	Function FunctionDefinition `json:"function"`
}

type FunctionDefinition struct {
	Name        string             `json:"name"`
	Description string             `json:"description"`
	Parameters  FunctionParameters `json:"parameters"`
}

type FunctionParameters struct {
	Type       string                     `json:"type"`
	Properties map[string]PropertyDetails `json:"properties"`
	Required   []string                   `json:"required"`
}

type PropertyDetails struct {
	Type        string      `json:"type"`
	Description string      `json:"description"`
	Default     interface{} `json:"default,omitempty"`
	Items       *Items      `json:"items,omitempty"` // 用于array类型
}

type Items struct {
	Type string `json:"type"`
}

// ToolStatusResponse defines tool status response structure
type ToolStatusResponse struct {
	Code    int            `json:"code"`
	Data    ToolStatusData `json:"data"`
	Message string         `json:"message"`
}

// ToolStatusData defines tool status data structure
type ToolStatusData struct {
	Tools map[string]ToolStatusDetail `json:"tools,omitempty"`
}

// ToolStatusDetail defines tool status detail structure
type ToolStatusDetail struct {
	Status string      `json:"status"`
	Result interface{} `json:"result,omitempty"`
}
