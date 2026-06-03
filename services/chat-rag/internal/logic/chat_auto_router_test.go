package logic

import (
	"context"
	"fmt"
	"net/http"
	"strings"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/stretchr/testify/assert"
	"github.com/zgsm-ai/chat-rag/internal/bootstrap"
	"github.com/zgsm-ai/chat-rag/internal/config"
	"github.com/zgsm-ai/chat-rag/internal/model"
	"github.com/zgsm-ai/chat-rag/internal/service/mocks"
	"github.com/zgsm-ai/chat-rag/internal/types"
)

func setupAutoRouterTestSvcCtx(t *testing.T, routerCfg *config.RouterConfig) (*bootstrap.ServiceContext, *gomock.Controller) {
	ctrl := gomock.NewController(t)
	loggerMock := mocks.NewMockLoggerInterface(ctrl)
	metricsMock := mocks.NewMockMetricsInterface(ctrl)
	loggerMock.EXPECT().LogAsync(gomock.Any(), gomock.Any()).AnyTimes()
	metricsMock.EXPECT().GetRegistry().Return(prometheus.NewRegistry()).AnyTimes()

	svcCtx := &bootstrap.ServiceContext{
		Config: config.Config{
			FromNacos: config.FromNacos{
				Router: routerCfg,
			},
		},
		LoggerService:  loggerMock,
		MetricsService: metricsMock,
	}
	return svcCtx, ctrl
}

func autoRouterRequest() *types.ChatCompletionRequest {
	return &types.ChatCompletionRequest{
		Model: "auto",
		LLMRequestParams: types.LLMRequestParams{
			Messages: []types.Message{
				{Role: types.RoleUser, Content: "hello"},
			},
		},
	}
}

func autoRouterIdentityNoVip() *model.Identity {
	return &model.Identity{
		UserName: "test-user",
		UserInfo: &model.UserInfo{
			Vip:  0,
			Name: "test-user",
		},
	}
}

func autoRouterCtxWithIdentity(identity *model.Identity) context.Context {
	return context.WithValue(context.Background(), model.IdentityContextKey, identity)
}

func TestAutoRouter_ChatCompletion_NoVisibleCandidates_ReturnsError(t *testing.T) {
	routerCfg := &config.RouterConfig{
		Enabled:  true,
		Strategy: "priority",
		Priority: config.PriorityConfig{
			Candidates: []config.PriorityCandidate{
				{ModelName: "vip-only-model", Enabled: true, Priority: 100, Weight: 1, MinVipLevel: 5},
			},
		},
	}

	svcCtx, ctrl := setupAutoRouterTestSvcCtx(t, routerCfg)
	defer ctrl.Finish()

	req := autoRouterRequest()
	hdr := make(http.Header)
	identity := autoRouterIdentityNoVip()
	ctx := autoRouterCtxWithIdentity(identity)
	writer := &mockResponseWriter{}

	logic := NewChatCompletionLogic(ctx, svcCtx, req, writer, &hdr, identity)

	resp, err := logic.ChatCompletion()
	assert.Error(t, err, "ChatCompletion should return an error when no visible candidates")
	assert.Nil(t, resp, "response should be nil on auto router error")
	assert.True(t,
		strings.Contains(err.Error(), "no visible candidates") || strings.Contains(err.Error(), "auto router"),
		fmt.Sprintf("error should mention 'no visible candidates' or 'auto router', got: %v", err),
	)
}

func TestAutoRouter_ChatCompletion_NilRunner_ReturnsError(t *testing.T) {
	routerCfg := &config.RouterConfig{
		Enabled:  true,
		Strategy: "unknown-strategy",
		Priority: config.PriorityConfig{},
	}

	svcCtx, ctrl := setupAutoRouterTestSvcCtx(t, routerCfg)
	defer ctrl.Finish()

	req := autoRouterRequest()
	hdr := make(http.Header)
	identity := autoRouterIdentityNoVip()
	ctx := autoRouterCtxWithIdentity(identity)
	writer := &mockResponseWriter{}

	logic := NewChatCompletionLogic(ctx, svcCtx, req, writer, &hdr, identity)

	resp, err := logic.ChatCompletion()
	assert.Error(t, err, "ChatCompletion should return error when runner is nil for auto model")
	assert.Nil(t, resp)
	assert.True(t,
		strings.Contains(err.Error(), "unknown or has invalid configuration") ||
			strings.Contains(err.Error(), "auto router"),
		fmt.Sprintf("error should mention 'unknown or has invalid configuration', got: %v", err),
	)
}

func TestAutoRouter_ChatCompletion_InvalidPriorityConfig_NilRunner(t *testing.T) {
	routerCfg := &config.RouterConfig{
		Enabled:  true,
		Strategy: "priority",
		Priority: config.PriorityConfig{
			Candidates: []config.PriorityCandidate{},
		},
	}

	svcCtx, ctrl := setupAutoRouterTestSvcCtx(t, routerCfg)
	defer ctrl.Finish()

	req := autoRouterRequest()
	hdr := make(http.Header)
	identity := autoRouterIdentityNoVip()
	ctx := autoRouterCtxWithIdentity(identity)
	writer := &mockResponseWriter{}

	logic := NewChatCompletionLogic(ctx, svcCtx, req, writer, &hdr, identity)

	resp, err := logic.ChatCompletion()
	assert.Error(t, err, "ChatCompletion should return error when priority config is invalid (runner nil)")
	assert.Nil(t, resp)
	assert.True(t,
		strings.Contains(err.Error(), "unknown or has invalid configuration"),
		fmt.Sprintf("error should mention 'unknown or has invalid configuration', got: %v", err),
	)
}

func TestAutoRouter_ChatCompletionStream_NoVisibleCandidates_ReturnsError(t *testing.T) {
	routerCfg := &config.RouterConfig{
		Enabled:  true,
		Strategy: "priority",
		Priority: config.PriorityConfig{
			Candidates: []config.PriorityCandidate{
				{ModelName: "vip-only-model", Enabled: true, Priority: 100, Weight: 1, MinVipLevel: 5},
			},
		},
	}

	svcCtx, ctrl := setupAutoRouterTestSvcCtx(t, routerCfg)
	defer ctrl.Finish()

	req := autoRouterRequest()
	hdr := make(http.Header)
	identity := autoRouterIdentityNoVip()
	ctx := autoRouterCtxWithIdentity(identity)
	writer := &mockResponseWriter{}

	logic := NewChatCompletionLogic(ctx, svcCtx, req, writer, &hdr, identity)

	err := logic.ChatCompletionStream()
	assert.Error(t, err, "ChatCompletionStream should return an error when no visible candidates")
	assert.True(t,
		strings.Contains(err.Error(), "no visible candidates") || strings.Contains(err.Error(), "auto router"),
		fmt.Sprintf("error should mention 'no visible candidates' or 'auto router', got: %v", err),
	)
}

func TestAutoRouter_ChatCompletionStream_NilRunner_ReturnsError(t *testing.T) {
	routerCfg := &config.RouterConfig{
		Enabled:  true,
		Strategy: "unknown-strategy",
		Priority: config.PriorityConfig{},
	}

	svcCtx, ctrl := setupAutoRouterTestSvcCtx(t, routerCfg)
	defer ctrl.Finish()

	req := autoRouterRequest()
	hdr := make(http.Header)
	identity := autoRouterIdentityNoVip()
	ctx := autoRouterCtxWithIdentity(identity)
	writer := &mockResponseWriter{}

	logic := NewChatCompletionLogic(ctx, svcCtx, req, writer, &hdr, identity)

	err := logic.ChatCompletionStream()
	assert.Error(t, err, "ChatCompletionStream should return error when runner is nil for auto model")
	assert.True(t,
		strings.Contains(err.Error(), "unknown or has invalid configuration") ||
			strings.Contains(err.Error(), "auto router"),
		fmt.Sprintf("error should mention 'unknown or has invalid configuration', got: %v", err),
	)
}

func TestAutoRouter_ChatCompletionStream_InvalidPriorityConfig_NilRunner(t *testing.T) {
	routerCfg := &config.RouterConfig{
		Enabled:  true,
		Strategy: "priority",
		Priority: config.PriorityConfig{
			Candidates: []config.PriorityCandidate{},
		},
	}

	svcCtx, ctrl := setupAutoRouterTestSvcCtx(t, routerCfg)
	defer ctrl.Finish()

	req := autoRouterRequest()
	hdr := make(http.Header)
	identity := autoRouterIdentityNoVip()
	ctx := autoRouterCtxWithIdentity(identity)
	writer := &mockResponseWriter{}

	logic := NewChatCompletionLogic(ctx, svcCtx, req, writer, &hdr, identity)

	err := logic.ChatCompletionStream()
	assert.Error(t, err, "ChatCompletionStream should return error when priority config is invalid (runner nil)")
	assert.True(t,
		strings.Contains(err.Error(), "unknown or has invalid configuration"),
		fmt.Sprintf("error should mention 'unknown or has invalid configuration', got: %v", err),
	)
}

func TestAutoRouter_ChatCompletion_RunError_ReturnsOriginalError(t *testing.T) {
	strategy := &mockFailingStrategy{err: fmt.Errorf("no visible candidates configured for current user")}
	routerCfg := &config.RouterConfig{
		Enabled:  true,
		Strategy: "priority",
		Priority: config.PriorityConfig{
			Candidates: []config.PriorityCandidate{
				{ModelName: "model-a", Enabled: true, Priority: 100, Weight: 1},
			},
		},
	}

	svcCtx, ctrl := setupAutoRouterTestSvcCtx(t, routerCfg)
	defer ctrl.Finish()

	svcCtx.SetRouterStrategy(strategy)

	req := autoRouterRequest()
	hdr := make(http.Header)
	identity := autoRouterIdentityNoVip()
	ctx := autoRouterCtxWithIdentity(identity)
	writer := &mockResponseWriter{}

	logic := NewChatCompletionLogic(ctx, svcCtx, req, writer, &hdr, identity)

	resp, err := logic.ChatCompletion()
	assert.Error(t, err)
	assert.Nil(t, resp)
	assert.True(t,
		strings.Contains(err.Error(), "no visible candidates") || strings.Contains(err.Error(), "auto router"),
		fmt.Sprintf("error should mention 'no visible candidates' or 'auto router', got: %v", err),
	)
}

func TestAutoRouter_ChatCompletionStream_RunError_ReturnsOriginalError(t *testing.T) {
	strategy := &mockFailingStrategy{err: fmt.Errorf("no visible candidates configured for current user")}
	routerCfg := &config.RouterConfig{
		Enabled:  true,
		Strategy: "priority",
		Priority: config.PriorityConfig{
			Candidates: []config.PriorityCandidate{
				{ModelName: "model-a", Enabled: true, Priority: 100, Weight: 1},
			},
		},
	}

	svcCtx, ctrl := setupAutoRouterTestSvcCtx(t, routerCfg)
	defer ctrl.Finish()

	svcCtx.SetRouterStrategy(strategy)

	req := autoRouterRequest()
	hdr := make(http.Header)
	identity := autoRouterIdentityNoVip()
	ctx := autoRouterCtxWithIdentity(identity)
	writer := &mockResponseWriter{}

	logic := NewChatCompletionLogic(ctx, svcCtx, req, writer, &hdr, identity)

	err := logic.ChatCompletionStream()
	assert.Error(t, err)
	assert.True(t,
		strings.Contains(err.Error(), "no visible candidates") || strings.Contains(err.Error(), "auto router"),
		fmt.Sprintf("error should mention 'no visible candidates' or 'auto router', got: %v", err),
	)
}

func TestAutoRouter_ResolveAutoRouter_Success_UpdatesModel(t *testing.T) {
	strategy := &mockSuccessStrategy{
		selected: "model-a",
		ordered:  []string{"model-a", "model-b"},
	}
	routerCfg := &config.RouterConfig{
		Enabled:  true,
		Strategy: "priority",
		Priority: config.PriorityConfig{
			Candidates: []config.PriorityCandidate{
				{ModelName: "model-a", Enabled: true, Priority: 100, Weight: 1},
			},
		},
	}

	svcCtx, ctrl := setupAutoRouterTestSvcCtx(t, routerCfg)
	defer ctrl.Finish()

	svcCtx.SetRouterStrategy(strategy)

	req := autoRouterRequest()
	hdr := make(http.Header)
	identity := autoRouterIdentityNoVip()
	ctx := autoRouterCtxWithIdentity(identity)
	writer := &mockResponseWriter{}

	logic := NewChatCompletionLogic(ctx, svcCtx, req, writer, &hdr, identity)

	assert.Equal(t, "auto", logic.request.Model)
	assert.Nil(t, logic.orderedModels)

	err := logic.resolveAutoRouter()
	assert.NoError(t, err, "resolveAutoRouter should not error on successful routing")
	assert.Equal(t, "model-a", logic.request.Model,
		"model should have been updated by successful routing")
	assert.Equal(t, []string{"model-a", "model-b"}, logic.orderedModels)
}

func TestAutoRouter_ResolveAutoRouter_Bypass_NilRouter(t *testing.T) {
	svcCtx, ctrl := setupAutoRouterTestSvcCtx(t, nil)
	defer ctrl.Finish()

	req := autoRouterRequest()
	req.Model = "auto"
	hdr := make(http.Header)
	identity := autoRouterIdentityNoVip()
	ctx := autoRouterCtxWithIdentity(identity)
	writer := &mockResponseWriter{}

	logic := NewChatCompletionLogic(ctx, svcCtx, req, writer, &hdr, identity)

	assert.Equal(t, "auto", logic.request.Model)
	err := logic.resolveAutoRouter()
	assert.NoError(t, err, "resolveAutoRouter should return nil when Router config is nil")
	assert.Equal(t, "auto", logic.request.Model,
		"model should remain unchanged when Router config is nil")
	assert.Nil(t, logic.orderedModels,
		"orderedModels should remain nil when Router config is nil")
	assert.Empty(t, writer.Header().Get(types.HeaderSelectLLm),
		"x-select-llm header should not be set when Router config is nil")
}

func TestAutoRouter_ResolveAutoRouter_Bypass_RouterDisabled(t *testing.T) {
	svcCtx, ctrl := setupAutoRouterTestSvcCtx(t, &config.RouterConfig{
		Enabled:  false,
		Strategy: "priority",
	})
	defer ctrl.Finish()

	req := autoRouterRequest()
	req.Model = "auto"
	hdr := make(http.Header)
	identity := autoRouterIdentityNoVip()
	ctx := autoRouterCtxWithIdentity(identity)
	writer := &mockResponseWriter{}

	logic := NewChatCompletionLogic(ctx, svcCtx, req, writer, &hdr, identity)

	assert.Equal(t, "auto", logic.request.Model)
	err := logic.resolveAutoRouter()
	assert.NoError(t, err, "resolveAutoRouter should return nil when Router is disabled")
	assert.Equal(t, "auto", logic.request.Model,
		"model should remain unchanged when Router is disabled")
	assert.Nil(t, logic.orderedModels,
		"orderedModels should remain nil when Router is disabled")
	assert.Empty(t, writer.Header().Get(types.HeaderSelectLLm),
		"x-select-llm header should not be set when Router is disabled")
}

func TestAutoRouter_ResolveAutoRouter_Bypass_ModelNotAuto(t *testing.T) {
	routerCfg := &config.RouterConfig{
		Enabled:  true,
		Strategy: "priority",
		Priority: config.PriorityConfig{
			Candidates: []config.PriorityCandidate{
				{ModelName: "model-a", Enabled: true, Priority: 100, Weight: 1},
			},
		},
	}

	svcCtx, ctrl := setupAutoRouterTestSvcCtx(t, routerCfg)
	defer ctrl.Finish()

	req := autoRouterRequest()
	req.Model = "gpt-4"
	hdr := make(http.Header)
	identity := autoRouterIdentityNoVip()
	ctx := autoRouterCtxWithIdentity(identity)
	writer := &mockResponseWriter{}

	logic := NewChatCompletionLogic(ctx, svcCtx, req, writer, &hdr, identity)

	assert.Equal(t, "gpt-4", logic.request.Model)
	err := logic.resolveAutoRouter()
	assert.NoError(t, err, "resolveAutoRouter should return nil when model is not auto")
	assert.Equal(t, "gpt-4", logic.request.Model,
		"model should remain unchanged when model is not auto")
	assert.Nil(t, logic.orderedModels,
		"orderedModels should remain nil when model is not auto")
	assert.Empty(t, writer.Header().Get(types.HeaderSelectLLm),
		"x-select-llm header should not be set when model is not auto")
}

func TestAutoRouter_ResolveAutoRouter_EmptySelect_NoChange(t *testing.T) {
	strategy := &mockSuccessStrategy{
		selected: "",
		ordered:  nil,
	}
	routerCfg := &config.RouterConfig{
		Enabled:  true,
		Strategy: "priority",
		Priority: config.PriorityConfig{
			Candidates: []config.PriorityCandidate{
				{ModelName: "model-a", Enabled: true, Priority: 100, Weight: 1},
			},
		},
	}

	svcCtx, ctrl := setupAutoRouterTestSvcCtx(t, routerCfg)
	defer ctrl.Finish()

	svcCtx.SetRouterStrategy(strategy)

	req := autoRouterRequest()
	hdr := make(http.Header)
	identity := autoRouterIdentityNoVip()
	ctx := autoRouterCtxWithIdentity(identity)
	writer := &mockResponseWriter{}

	logic := NewChatCompletionLogic(ctx, svcCtx, req, writer, &hdr, identity)

	assert.Equal(t, "auto", logic.request.Model)
	assert.Nil(t, logic.orderedModels)

	err := logic.resolveAutoRouter()
	assert.NoError(t, err, "resolveAutoRouter should not error when strategy returns empty select")
	assert.Equal(t, "auto", logic.request.Model,
		"model should remain unchanged when selected is empty")
	assert.Nil(t, logic.orderedModels,
		"orderedModels should remain nil when selected is empty")
	assert.Empty(t, hdr.Get(types.HeaderOriginalModel),
		"x-original-model header should not be set when selected is empty")
	assert.Empty(t, writer.Header().Get(types.HeaderSelectLLm),
		"x-select-llm header should not be set when selected is empty")
}

type mockFailingStrategy struct {
	err error
}

func (m *mockFailingStrategy) Name() string { return "mock-failing" }
func (m *mockFailingStrategy) Run(_ context.Context, _ *bootstrap.ServiceContext, _ *http.Header, _ *types.ChatCompletionRequest) (string, string, []string, error) {
	return "", "", nil, m.err
}

type mockSuccessStrategy struct {
	selected string
	ordered  []string
}

func (m *mockSuccessStrategy) Name() string { return "mock-success" }
func (m *mockSuccessStrategy) Run(_ context.Context, _ *bootstrap.ServiceContext, _ *http.Header, _ *types.ChatCompletionRequest) (string, string, []string, error) {
	return m.selected, "", m.ordered, nil
}

func TestChatCompletionLogic_LogCompletionCopiesDegradationTrace(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	loggerMock := mocks.NewMockLoggerInterface(ctrl)
	var capturedLog *model.ChatLog
	loggerMock.EXPECT().LogAsync(gomock.Any(), gomock.Any()).Do(func(log *model.ChatLog, _ *http.Header) {
		capturedLog = log
	})

	svcCtx := &bootstrap.ServiceContext{LoggerService: loggerMock}
	req := autoRouterRequest()
	hdr := make(http.Header)
	identity := autoRouterIdentityNoVip()
	logic := NewChatCompletionLogic(context.Background(), svcCtx, req, &mockResponseWriter{}, &hdr, identity)
	logic.orderedModels = []string{"model-a", "model-b"}
	chatLog := logic.newChatLog(time.Now())

	logic.appendDegradationEvent(model.DegradationEvent{
		Event:      model.DegradationEventAttemptFailed,
		Model:      "model-a",
		ModelIndex: 0,
		Reason:     "model_call_failed",
	})

	logic.logCompletion(chatLog)

	assert.NotNil(t, capturedLog)
	assert.Len(t, capturedLog.DegradationTrace, 1)
	assert.Equal(t, model.DegradationEventAttemptFailed, capturedLog.DegradationTrace[0].Event)
	assert.Equal(t, "model-a", capturedLog.DegradationTrace[0].Model)
	assert.Equal(t, 1, capturedLog.DegradationTrace[0].Sequence)
	assert.Empty(t, capturedLog.Error)
}
