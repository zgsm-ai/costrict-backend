package processor

import (
	"reflect"
	"testing"

	"github.com/zgsm-ai/chat-rag/internal/model"
	"github.com/zgsm-ai/chat-rag/internal/types"
)

func TestPromptMsgDoesNotModifySystemPromptByDefault(t *testing.T) {
	messages := []types.Message{
		{Role: types.RoleSystem, Content: "You are costrict"},
		{Role: types.RoleUser, Content: "Hello"},
	}

	promptMsg, err := NewPromptMsg(messages, false)
	if err != nil {
		t.Fatalf("NewPromptMsg() error = %v", err)
	}

	promptMsg.UpdateSystemMsg("modified")
	SetLanguage("zh-CN", promptMsg)

	if got := promptMsg.GetSystemMsg().Content; got != messages[0].Content {
		t.Fatalf("system content = %#v, want %#v", got, messages[0].Content)
	}
}

func TestPromptMsgModifiesSystemPromptWhenEnabled(t *testing.T) {
	messages := []types.Message{
		{Role: types.RoleSystem, Content: "You are costrict"},
		{Role: types.RoleUser, Content: "Hello"},
	}

	promptMsg, err := NewPromptMsg(messages, true)
	if err != nil {
		t.Fatalf("NewPromptMsg() error = %v", err)
	}

	promptMsg.UpdateSystemMsg("modified")

	contents, ok := promptMsg.GetSystemMsg().Content.([]model.Content)
	if !ok || len(contents) != 1 {
		t.Fatalf("system content = %#v, want one structured content item", promptMsg.GetSystemMsg().Content)
	}
	if contents[0].Text != "modified" {
		t.Fatalf("system text = %q, want %q", contents[0].Text, "modified")
	}
	wantCacheControl := map[string]interface{}{"type": "ephemeral"}
	if !reflect.DeepEqual(contents[0].CacheControl, wantCacheControl) {
		t.Fatalf("cache control = %#v, want %#v", contents[0].CacheControl, wantCacheControl)
	}
}
