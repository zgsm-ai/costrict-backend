package logic

import (
	"context"
	"net/http"
	"strings"
	"testing"
	"unicode/utf8"

	"github.com/golang/mock/gomock"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/stretchr/testify/assert"
	"github.com/zgsm-ai/chat-rag/internal/bootstrap"
	"github.com/zgsm-ai/chat-rag/internal/client"
	"github.com/zgsm-ai/chat-rag/internal/config"
	"github.com/zgsm-ai/chat-rag/internal/model"
	"github.com/zgsm-ai/chat-rag/internal/service/mocks"
	"github.com/zgsm-ai/chat-rag/internal/tokenizer"
	"github.com/zgsm-ai/chat-rag/internal/types"
)

// createTestContext creates a test context
func createTestContext() context.Context {
	return context.Background()
}

// createTestServiceContext creates a test ServiceContext
// tokenCounter parameter can be nil or *utils.TokenCounter type
func createTestServiceMock(t *testing.T) (*gomock.Controller, *mocks.MockLoggerInterface, *mocks.MockMetricsInterface) {
	ctrl := gomock.NewController(t)
	loggerMock := mocks.NewMockLoggerInterface(ctrl)
	metricsMock := mocks.NewMockMetricsInterface(ctrl)

	// Setup common mock expectations
	loggerMock.EXPECT().LogAsync(gomock.Any(), gomock.Any()).AnyTimes()
	loggerMock.EXPECT().SetMetricsService(gomock.Any()).AnyTimes()
	metricsMock.EXPECT().GetRegistry().Return(prometheus.NewRegistry()).AnyTimes()
	metricsMock.EXPECT().RecordChatLog(gomock.Any()).AnyTimes()

	return ctrl, loggerMock, metricsMock
}

func createTestServiceContext(t *testing.T, cfg *config.Config, tokenCounter interface{}) *bootstrap.ServiceContext {
	ctrl, loggerMock, metricsMock := createTestServiceMock(t)
	defer ctrl.Finish()

	svcCtx := &bootstrap.ServiceContext{
		Config: config.Config{
			LLM: cfg.LLM,
		},
		LoggerService:  loggerMock,
		MetricsService: metricsMock,
	}

	// If tokenCounter exists and type is correct, set it to ServiceContext
	if tc, ok := tokenCounter.(*tokenizer.TokenCounter); ok {
		svcCtx.TokenCounter = tc
	}

	return svcCtx
}

// createTestRequest creates a test ChatCompletionRequest
func createTestRequest(model string, messages []types.Message, stream bool) *types.ChatCompletionRequest {
	req := &types.ChatCompletionRequest{
		Model: model,
		LLMRequestParams: types.LLMRequestParams{
			Messages: messages,
		},
	}
	// Add stream to Extra map
	if req.Extra == nil {
		req.Extra = make(map[string]any)
	}
	req.Extra["stream"] = stream
	return req
}

// createTestIdentity creates a test Identity
func createTestIdentity() *model.Identity {
	return &model.Identity{
		ClientID:    "test-client",
		ProjectPath: "/test/path",
	}
}

// setupTestLogic combines all helper functions to create complete test logic
func setupTestLogic(t *testing.T, cfg *config.Config, tokenCounter interface{},
	model string, messages []types.Message, writer http.ResponseWriter) (*ChatCompletionLogic, *bootstrap.ServiceContext) {
	ctx := createTestContext()
	svcCtx := createTestServiceContext(t, cfg, tokenCounter)
	req := createTestRequest(model, messages, false)
	identity := createTestIdentity()
	headers := make(http.Header)

	// Set mock expectations
	if logger, ok := svcCtx.LoggerService.(*mocks.MockLoggerInterface); ok {
		logger.EXPECT().LogAsync(gomock.Any(), gomock.Any()).AnyTimes()
	}

	return NewChatCompletionLogic(ctx, svcCtx, req, writer, &headers, identity), svcCtx
}

func TestChatCompletionLogic_NewChatCompletionLogic(t *testing.T) {
	mockWriter := &mockResponseWriter{}
	cfg := &config.Config{}
	logic, svcCtx := setupTestLogic(t, cfg, nil, "test-model", []types.Message{
		{Role: "user", Content: "Hello"},
	}, mockWriter)

	assert.NotNil(t, logic)
	assert.Equal(t, createTestContext(), logic.ctx)
	assert.Equal(t, svcCtx, logic.svcCtx)
}

// TestLLMClientMock tests LLMClient mock
func TestLLMClientMock(t *testing.T) {
	mock := &client.LLMClient{}
	assert.NotNil(t, mock)
}

func TestChatCompletionLogic_countTokensInMessages_Fallback(t *testing.T) {
	cfg := &config.Config{}
	logic, _ := setupTestLogic(t, cfg, nil, "test-model", []types.Message{}, &mockResponseWriter{})

	messages := []types.Message{
		{Role: "user", Content: "Hello"},
		{Role: "assistant", Content: "Hi there!"},
	}

	count := logic.countTokensInMessages(messages)
	assert.Greater(t, count, 0) // Should return estimated token count
}

func TestChatCompletionLogic_ChatCompletion_ValidationErrors(t *testing.T) {
	tests := []struct {
		name     string
		config   *config.Config
		expected string
	}{
		{
			name:     "empty endpoint",
			config:   &config.Config{},
			expected: "NewLLMClient llmEndpoint cannot be empty",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			logic, _ := setupTestLogic(t, tt.config, nil, "test-model", []types.Message{}, &mockResponseWriter{})

			resp, err := logic.ChatCompletion()
			t.Log("==>", err)

			assert.Error(t, err)
			assert.Contains(t, err.Error(), tt.expected)
			assert.Nil(t, resp)
		})
	}

	// Test valid configuration
	t.Run("valid config", func(t *testing.T) {
		cfg := &config.Config{
			LLM: config.LLMConfig{
				Endpoint: "http://test-endpoint",
			},
		}
		_, svcCtx := setupTestLogic(t, cfg, nil, "test-model", []types.Message{}, &mockResponseWriter{})
		assert.NotNil(t, svcCtx)
		assert.Equal(t, "http://test-endpoint", svcCtx.Config.LLM.Endpoint)
	})
}

// Tests TokenCounter setup correctly
func TestChatCompletionLogic_WithTokenCounter(t *testing.T) {
	mockWriter := &mockResponseWriter{}
	cfg := &config.Config{}
	tokenCounter := &tokenizer.TokenCounter{}

	logic, svcCtx := setupTestLogic(t, cfg, tokenCounter, "test-model",
		[]types.Message{{Role: "user", Content: "Hello"}}, mockWriter)

	assert.NotNil(t, logic)
	assert.NotNil(t, svcCtx.TokenCounter)
	assert.Equal(t, tokenCounter, svcCtx.TokenCounter)
}

func TestChatCompletionLogic_ChatCompletion_BasicRequest(t *testing.T) {
	cfg := config.MustLoadConfig("../../etc/chat-api.yaml")

	// Initialize token counter
	tokenCounter, err := tokenizer.NewTokenCounter()
	if err != nil {
		tokenCounter = &tokenizer.TokenCounter{} // Fallback to basic counter
	}

	ctrl, _, _ := createTestServiceMock(t)
	defer ctrl.Finish()

	logic, _ := setupTestLogic(t, &cfg, tokenCounter,
		"gpt-3.5-turbo", []types.Message{
			{Role: "user", Content: "Hello, how are you?"},
		}, &mockResponseWriter{})

	// Test basic request
	resp, err := logic.ChatCompletion()

	// Verify response
	assert.Error(t, err)
	assert.Nil(t, resp)
}

// mockResponseWriter mocks http.ResponseWriter and http.Flusher for testing
type mockResponseWriter struct {
	data       []byte
	headers    http.Header
	statusCode int
	flushed    bool
}

func (m *mockResponseWriter) Header() http.Header {
	if m.headers == nil {
		m.headers = make(http.Header)
	}
	return m.headers
}

func (m *mockResponseWriter) Write(data []byte) (int, error) {
	m.data = append(m.data, data...)
	return len(data), nil
}

func (m *mockResponseWriter) WriteHeader(statusCode int) {
	m.statusCode = statusCode
}

// Flush implements http.Flusher interface
func (m *mockResponseWriter) Flush() {
	m.flushed = true
}

func TestChatCompletionLogic_ChatCompletion_StreamingRequest(t *testing.T) {
	// Load config
	cfg := config.MustLoadConfig("../../etc/chat-api.yaml")

	// Initialize token counter
	tokenCounter, _ := tokenizer.NewTokenCounter()

	// Setup mocks
	ctrl, loggerMock, metricsMock := createTestServiceMock(t)
	defer ctrl.Finish()

	// Prepare test data
	testModel := "gpt-3.5-turbo"
	testMessages := []types.Message{
		{Role: "user", Content: "Hello, how are you?"},
	}
	testWriter := &mockResponseWriter{}

	// Create service context
	svcCtx := &bootstrap.ServiceContext{
		Config:         cfg,
		LoggerService:  loggerMock,
		MetricsService: metricsMock,
		TokenCounter:   tokenCounter,
	}

	headers := make(http.Header)

	// Create logic instance
	logic := NewChatCompletionLogic(
		createTestContext(),
		svcCtx,
		createTestRequest(testModel, testMessages, true),
		testWriter,
		&headers,
		createTestIdentity(),
	)

	// Execute test
	err := logic.ChatCompletionStream()

	// Verify expected error
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "401 Authorization Required", "Expected 401 unauthorized error")

	// Verify response write attempt
	assert.Greater(t, len(testWriter.data), 0, "Expected response attempt data")
	assert.True(t, testWriter.flushed, "Expected response flush attempt")
}

func TestChatCompletionLogic_ErrorClassification_BadRequestDoesNotRetryButCanDegrade(t *testing.T) {
	logic, _ := setupTestLogic(t, &config.Config{}, nil, "auto", []types.Message{{Role: types.RoleUser, Content: "hello"}}, &mockResponseWriter{})
	err := &types.APIError{
		Code:       types.ErrCodeModelServiceUnavailable,
		Message:    "downstream bad request: missing prompt",
		Success:    false,
		StatusCode: http.StatusBadRequest,
		Type:       string(types.ErrServerModel),
	}

	assert.False(t, logic.isSameModelRetryableError(err))
	assert.True(t, logic.isDegradationSwitchableError(err))
}

func TestChatCompletionLogic_ErrorClassification_EmptyMessageBadRequestDoesNotDegrade(t *testing.T) {
	logic, _ := setupTestLogic(t, &config.Config{}, nil, "auto", []types.Message{{Role: types.RoleUser, Content: "hello"}}, &mockResponseWriter{})
	err := types.NewEmptyMessageContentError()

	assert.False(t, logic.isSameModelRetryableError(err))
	assert.False(t, logic.isDegradationSwitchableError(err))
}

func TestChatCompletionLogic_ErrorClassification_UnauthorizedDoesNotDegrade(t *testing.T) {
	logic, _ := setupTestLogic(t, &config.Config{}, nil, "auto", []types.Message{{Role: types.RoleUser, Content: "hello"}}, &mockResponseWriter{})
	err := &types.APIError{
		Code:       types.ErrCodeUnauthorized,
		Message:    types.ErrMsgUnauthorized,
		Success:    false,
		StatusCode: http.StatusUnauthorized,
		Type:       string(types.ErrServerModel),
	}

	assert.False(t, logic.isSameModelRetryableError(err))
	assert.False(t, logic.isDegradationSwitchableError(err))
}

func TestChatCompletionLogic_ErrorClassification_TooManyRequestsDoesNotRetryButCanDegrade(t *testing.T) {
	logic, _ := setupTestLogic(t, &config.Config{}, nil, "auto", []types.Message{{Role: types.RoleUser, Content: "hello"}}, &mockResponseWriter{})
	err := &types.APIError{
		Code:       types.ErrCodeTooManyRequests,
		Message:    types.ErrMsgTooManyRequests,
		Success:    false,
		StatusCode: http.StatusTooManyRequests,
		Type:       string(types.ErrServerModel),
	}

	assert.False(t, logic.isSameModelRetryableError(err))
	assert.True(t, logic.isDegradationSwitchableError(err))
}

func TestChatCompletionLogic_ErrorClassification_ServerErrorRetriesAndDegrades(t *testing.T) {
	logic, _ := setupTestLogic(t, &config.Config{}, nil, "auto", []types.Message{{Role: types.RoleUser, Content: "hello"}}, &mockResponseWriter{})
	err := &types.APIError{
		Code:       types.ErrCodeModelServiceUnavailable,
		Message:    types.ErrMsgModelServiceUnavailable,
		Success:    false,
		StatusCode: http.StatusInternalServerError,
		Type:       string(types.ErrServerModel),
	}

	assert.True(t, logic.isSameModelRetryableError(err))
	assert.True(t, logic.isDegradationSwitchableError(err))
}

func TestChatCompletionLogic_ErrorClassification_IdleTimeoutRetriesAndDegrades(t *testing.T) {
	logic, _ := setupTestLogic(t, &config.Config{}, nil, "auto", []types.Message{{Role: types.RoleUser, Content: "hello"}}, &mockResponseWriter{})
	err := types.NewStreamIdleTimeoutError()

	assert.True(t, logic.isSameModelRetryableError(err))
	assert.True(t, logic.isDegradationSwitchableError(err))
}

func TestTruncateBytes_UTF8Safe(t *testing.T) {
	assert.Equal(t, "hello", truncateBytes("hello", 100))
	assert.Equal(t, "", truncateBytes("", 10))

	s := strings.Repeat("é", 300)
	r := truncateBytes(s, 7)
	assert.LessOrEqual(t, len([]byte(r)), 7)
	assert.True(t, utf8.ValidString(r), "truncated string must be valid UTF-8")

	assert.Equal(t, "abc", truncateBytes("abc", 0), "max<=0 means no truncation")
}

func TestSanitizeDegradationErrorMessage_Redaction(t *testing.T) {
	result := sanitizeDegradationErrorMessage("authorization: Bearer abcdef123456")
	assert.NotContains(t, result, "abcdef123456")
	assert.Contains(t, result, "[redacted]")

	result = sanitizeDegradationErrorMessage("token=supersecretvalue")
	assert.NotContains(t, result, "supersecretvalue")
	assert.Contains(t, result, "[redacted]")

	result = sanitizeDegradationErrorMessage("invalid token count: 5")
	assert.Equal(t, "invalid token count: 5", result)
}

func TestAppendDegradationEvent_CapAndTruncation(t *testing.T) {
	logic, _ := setupTestLogic(t, &config.Config{}, nil, "auto", []types.Message{{Role: types.RoleUser, Content: "hi"}}, &mockResponseWriter{})
	logic.orderedModels = []string{"m1", "m2"}

	for i := 0; i < 100; i++ {
		logic.appendDegradationEvent(model.DegradationEvent{Event: model.DegradationEventAttemptStarted, Model: "m1", ModelIndex: 0})
	}

	assert.Equal(t, maxDegradationTraceEvents, len(logic.degradationTrace))
	assert.Equal(t, model.DegradationEventTraceTruncated, logic.degradationTrace[len(logic.degradationTrace)-1].Event)
	assert.True(t, logic.degradationTraceFull)
	assert.Equal(t, maxDegradationTraceEvents, logic.degradationTraceSeq)
	for i, ev := range logic.degradationTrace {
		assert.Equal(t, i+1, ev.Sequence, "sequences should be monotonically increasing")
	}
}

func TestDegradationModelIndex(t *testing.T) {
	logic, _ := setupTestLogic(t, &config.Config{}, nil, "auto", []types.Message{{Role: types.RoleUser, Content: "hi"}}, &mockResponseWriter{})
	logic.orderedModels = []string{"m1", "m2", "m3"}
	assert.Equal(t, 0, logic.degradationModelIndex("m1"))
	assert.Equal(t, 2, logic.degradationModelIndex("m3"))
	assert.Equal(t, -1, logic.degradationModelIndex("missing"))
}

func TestAppendDegradationEvent_NoopWithoutOrderedModels(t *testing.T) {
	logic, _ := setupTestLogic(t, &config.Config{}, nil, "auto", []types.Message{{Role: types.RoleUser, Content: "hi"}}, &mockResponseWriter{})
	// DO NOT set orderedModels — leave it empty/nil

	logic.appendDegradationEvent(model.DegradationEvent{Event: model.DegradationEventAttemptStarted, Model: "m1", ModelIndex: 0})
	assert.Equal(t, 0, len(logic.degradationTrace))
}
