package priority

import (
	"context"
	"reflect"
	"strings"
	"testing"
	"time"

	"github.com/zgsm-ai/chat-rag/internal/config"
	"github.com/zgsm-ai/chat-rag/internal/model"
	"github.com/zgsm-ai/chat-rag/internal/types"
)

func TestValidateConfigRejectsNegativeMinVipLevel(t *testing.T) {
	err := validateConfig(config.PriorityConfig{
		Candidates: []config.PriorityCandidate{
			{
				ModelName:   "vip-model",
				Enabled:     true,
				Priority:    100,
				Weight:      1,
				MinVipLevel: -1,
			},
		},
	})

	if err == nil {
		t.Fatal("expected negative minVipLevel to be rejected")
	}
	if !strings.Contains(err.Error(), "minVipLevel must be >= 0") {
		t.Fatalf("expected minVipLevel validation error, got %v", err)
	}
}

func testAutoRequest() *types.ChatCompletionRequest {
	req := &types.ChatCompletionRequest{
		Model: "auto",
	}
	req.Messages = []types.Message{
		{Role: types.RoleUser, Content: "hello"},
	}
	return req
}

func testIdentity(vip int, expire *time.Time) *model.Identity {
	return &model.Identity{
		UserName: "test-user",
		UserInfo: &model.UserInfo{
			Name:      "test-user",
			Vip:       vip,
			VipExpire: expire,
		},
	}
}

func contextWithIdentity(identity *model.Identity) context.Context {
	return context.WithValue(context.Background(), model.IdentityContextKey, identity)
}

func TestFallbackNotInCandidatesNewSucceedsRunSkips(t *testing.T) {
	strategy, err := New(config.PriorityConfig{
		Candidates: []config.PriorityCandidate{
			{ModelName: "normal-model", Enabled: true, Priority: 100, Weight: 1},
		},
		FallbackModelName: "missing-fallback",
	})
	if err != nil {
		t.Fatalf("New should succeed when fallback is not in candidates, got error: %v", err)
	}

	selected, _, ordered, err := strategy.Run(context.Background(), nil, nil, testAutoRequest())
	if err != nil {
		t.Fatalf("Run returned error: %v", err)
	}

	if selected != "normal-model" {
		t.Fatalf("selected = %q, want normal-model", selected)
	}
	want := []string{"normal-model"}
	if !reflect.DeepEqual(ordered, want) {
		t.Fatalf("ordered = %#v, want %#v", ordered, want)
	}
}

func TestPriorityRouterOrdinaryUserSkipsVipOnlyCandidates(t *testing.T) {
	strategy, err := New(config.PriorityConfig{
		Candidates: []config.PriorityCandidate{
			{ModelName: "vip-model", Enabled: true, Priority: 100, Weight: 1, MinVipLevel: 1},
			{ModelName: "normal-model", Enabled: true, Priority: 200, Weight: 1, MinVipLevel: 0},
		},
	})
	if err != nil {
		t.Fatalf("New returned error: %v", err)
	}

	selected, _, ordered, err := strategy.Run(context.Background(), nil, nil, testAutoRequest())
	if err != nil {
		t.Fatalf("Run returned error: %v", err)
	}

	if selected != "normal-model" {
		t.Fatalf("selected = %q, want normal-model", selected)
	}
	want := []string{"normal-model"}
	if !reflect.DeepEqual(ordered, want) {
		t.Fatalf("ordered = %#v, want %#v", ordered, want)
	}
}

func TestPriorityRouterVipUserCanUseVipOnlyCandidates(t *testing.T) {
	strategy, err := New(config.PriorityConfig{
		Candidates: []config.PriorityCandidate{
			{ModelName: "vip-model", Enabled: true, Priority: 100, Weight: 1, MinVipLevel: 1},
			{ModelName: "normal-model", Enabled: true, Priority: 200, Weight: 1, MinVipLevel: 0},
		},
	})
	if err != nil {
		t.Fatalf("New returned error: %v", err)
	}

	selected, _, ordered, err := strategy.Run(contextWithIdentity(testIdentity(1, nil)), nil, nil, testAutoRequest())
	if err != nil {
		t.Fatalf("Run returned error: %v", err)
	}

	if selected != "vip-model" {
		t.Fatalf("selected = %q, want vip-model", selected)
	}
	want := []string{"vip-model", "normal-model"}
	if !reflect.DeepEqual(ordered, want) {
		t.Fatalf("ordered = %#v, want %#v", ordered, want)
	}
}

func TestPriorityRouterExpiredVipUserSkipsVipOnlyCandidates(t *testing.T) {
	expired := time.Now().Add(-time.Minute)
	strategy, err := New(config.PriorityConfig{
		Candidates: []config.PriorityCandidate{
			{ModelName: "vip-model", Enabled: true, Priority: 100, Weight: 1, MinVipLevel: 1},
			{ModelName: "normal-model", Enabled: true, Priority: 200, Weight: 1, MinVipLevel: 0},
		},
	})
	if err != nil {
		t.Fatalf("New returned error: %v", err)
	}

	selected, _, ordered, err := strategy.Run(contextWithIdentity(testIdentity(1, &expired)), nil, nil, testAutoRequest())
	if err != nil {
		t.Fatalf("Run returned error: %v", err)
	}

	if selected != "normal-model" {
		t.Fatalf("selected = %q, want normal-model", selected)
	}
	want := []string{"normal-model"}
	if !reflect.DeepEqual(ordered, want) {
		t.Fatalf("ordered = %#v, want %#v", ordered, want)
	}
}

func TestPriorityRouterSkipsInvisibleFallback(t *testing.T) {
	strategy, err := New(config.PriorityConfig{
		Candidates: []config.PriorityCandidate{
			{ModelName: "normal-model", Enabled: true, Priority: 200, Weight: 1, MinVipLevel: 0},
			{ModelName: "vip-fallback", Enabled: true, Priority: 300, Weight: 1, MinVipLevel: 1},
		},
		FallbackModelName: "vip-fallback",
	})
	if err != nil {
		t.Fatalf("New returned error: %v", err)
	}

	selected, _, ordered, err := strategy.Run(context.Background(), nil, nil, testAutoRequest())
	if err != nil {
		t.Fatalf("Run returned error: %v", err)
	}

	if selected != "normal-model" {
		t.Fatalf("selected = %q, want normal-model", selected)
	}
	want := []string{"normal-model"}
	if !reflect.DeepEqual(ordered, want) {
		t.Fatalf("ordered = %#v, want %#v", ordered, want)
	}
}

func TestPriorityRouterAppendsVisibleFallbackWhenMissingFromOrder(t *testing.T) {
	strategy, err := New(config.PriorityConfig{
		Candidates: []config.PriorityCandidate{
			{ModelName: "primary", Enabled: true, Priority: 100, Weight: 1, MinVipLevel: 0},
			{ModelName: "fallback", Enabled: true, Priority: 200, Weight: 1, MinVipLevel: 0},
		},
		FallbackModelName: "fallback",
	})
	if err != nil {
		t.Fatalf("New returned error: %v", err)
	}

	selected, _, ordered, err := strategy.Run(context.Background(), nil, nil, testAutoRequest())
	if err != nil {
		t.Fatalf("Run returned error: %v", err)
	}

	if selected != "primary" {
		t.Fatalf("selected = %q, want primary", selected)
	}
	want := []string{"primary", "fallback"}
	if !reflect.DeepEqual(ordered, want) {
		t.Fatalf("ordered = %#v, want %#v", ordered, want)
	}
}

func TestAppendVisibleFallbackAppendsWhenMissingFromOrder(t *testing.T) {
	strategy, err := New(config.PriorityConfig{
		Candidates: []config.PriorityCandidate{
			{ModelName: "primary", Enabled: true, Priority: 100, Weight: 1, MinVipLevel: 0},
			{ModelName: "fallback", Enabled: true, Priority: 200, Weight: 1, MinVipLevel: 0},
		},
		FallbackModelName: "fallback",
	})
	if err != nil {
		t.Fatalf("New returned error: %v", err)
	}

	ordered := []string{"primary"}
	visible := []*ModelCandidate{
		{modelName: "primary", priority: 100, weight: 1, enabled: true, minVipLevel: 0},
		{modelName: "fallback", priority: 200, weight: 1, enabled: true, minVipLevel: 0},
	}

	result := strategy.appendVisibleFallback(context.Background(), ordered, visible)

	want := []string{"primary", "fallback"}
	if !reflect.DeepEqual(result, want) {
		t.Fatalf("result = %#v, want %#v", result, want)
	}
}

func TestPriorityRouterDoesNotDuplicateVisibleFallback(t *testing.T) {
	strategy, err := New(config.PriorityConfig{
		Candidates: []config.PriorityCandidate{
			{ModelName: "fallback", Enabled: true, Priority: 100, Weight: 1, MinVipLevel: 0},
			{ModelName: "secondary", Enabled: true, Priority: 200, Weight: 1, MinVipLevel: 0},
		},
		FallbackModelName: "fallback",
	})
	if err != nil {
		t.Fatalf("New returned error: %v", err)
	}

	selected, _, ordered, err := strategy.Run(context.Background(), nil, nil, testAutoRequest())
	if err != nil {
		t.Fatalf("Run returned error: %v", err)
	}

	if selected != "fallback" {
		t.Fatalf("selected = %q, want fallback", selected)
	}
	want := []string{"fallback", "secondary"}
	if !reflect.DeepEqual(ordered, want) {
		t.Fatalf("ordered = %#v, want %#v", ordered, want)
	}
}

func TestPriorityRouterReturnsErrorWhenNoCandidateVisible(t *testing.T) {
	strategy, err := New(config.PriorityConfig{
		Candidates: []config.PriorityCandidate{
			{ModelName: "vip-model", Enabled: true, Priority: 100, Weight: 1, MinVipLevel: 1},
		},
		FallbackModelName: "vip-model",
	})
	if err != nil {
		t.Fatalf("New returned error: %v", err)
	}

	selected, _, ordered, err := strategy.Run(context.Background(), nil, nil, testAutoRequest())
	if err == nil {
		t.Fatal("expected error when no candidates are visible")
	}
	if selected != "" {
		t.Fatalf("selected = %q, want empty", selected)
	}
	if ordered != nil {
		t.Fatalf("ordered = %#v, want nil", ordered)
	}
	if !strings.Contains(err.Error(), "no visible candidates") {
		t.Fatalf("error = %v, want no visible candidates", err)
	}
}

func TestOmittedMinVipLevelIsVisibleToOrdinaryUser(t *testing.T) {
	strategy, err := New(config.PriorityConfig{
		Candidates: []config.PriorityCandidate{
			{ModelName: "model-a", Enabled: true, Priority: 100, Weight: 1},
		},
	})
	if err != nil {
		t.Fatalf("New returned error: %v", err)
	}

	selected, _, ordered, err := strategy.Run(context.Background(), nil, nil, testAutoRequest())
	if err != nil {
		t.Fatalf("Run returned error: %v", err)
	}

	if selected != "model-a" {
		t.Fatalf("selected = %q, want model-a", selected)
	}
	want := []string{"model-a"}
	if !reflect.DeepEqual(ordered, want) {
		t.Fatalf("ordered = %#v, want %#v", ordered, want)
	}
}

func TestIsCandidateDeclaredReturnsTrueForDisabledCandidate(t *testing.T) {
	strategy, err := New(config.PriorityConfig{
		Candidates: []config.PriorityCandidate{
			{ModelName: "enabled-model", Enabled: true, Priority: 100, Weight: 1},
			{ModelName: "disabled-fallback", Enabled: false, Priority: 200, Weight: 1},
		},
	})
	if err != nil {
		t.Fatalf("New returned error: %v", err)
	}

	if !strategy.isCandidateDeclared("disabled-fallback") {
		t.Fatal("isCandidateDeclared(disabled-fallback) = false, want true")
	}
}
